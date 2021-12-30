// Package s3 provides functionality for communicating with a minio store
package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	// This package works with any S3 compatible storage
	"github.com/minio/minio-go/v7"
	miniocredentials "github.com/minio/minio-go/v7/pkg/credentials"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
)

const minioNoSuchKeyError = "NoSuchKey"

// client is a minimal interface of the functions we use from the minio package, useful for testing
type client interface {
	BucketExists(context.Context, string) (bool, error)
	StatObject(context.Context, string, string, minio.StatObjectOptions) (minio.ObjectInfo, error)

	PutObject(context.Context, string, string, io.Reader, int64, minio.PutObjectOptions) (minio.UploadInfo, error)

	ListObjects(context.Context, string, minio.ListObjectsOptions) <-chan minio.ObjectInfo

	GetObject(context.Context, string, string, minio.GetObjectOptions) (*minio.Object, error)

	Presign(context.Context, string, string, string, time.Duration, url.Values) (*url.URL, error)

	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
}

// bucketClient is a  S3 compilable storage client, it implements the objectstorage.Client interface
type bucketClient struct {
	client     client // The minio client is compilable with any S3 API compilable service
	bucketName string
}

// NewClient creates a new client for minio, returns an error if credentials are invalid, the bucket does not exist or a transport error occurs
// the configuration should be the credentials for the S3 compatible API of the provider
func NewClient(ctx context.Context, config *objectstorage.ClientConfig, bucketName string, transport http.RoundTripper) (objectstorage.Client, error) {
	client := bucketClient{
		bucketName: bucketName,
	}

	opts := &minio.Options{
		Creds:  miniocredentials.NewStaticV4(config.AccessKey, config.AccessSecret, ""),
		Secure: config.Secure,
	}
	if transport != nil {
		opts.Transport = transport
	}

	minioClient, err := minio.New(config.Host, opts)
	if err != nil {
		// Wrap error around objectstorage so we don't expose the minio errors
		return nil, errors.Wrapf(objectstorage.ErrInternal, "error trying to connect to minio: %s", err)
	}
	client.client = minioClient

	var buckedExists bool
	// Creating a client will return an error if connection cannot be made (when the sidecar has not started yet for example)
	// therefore attempt to get the bucket for at maximum maxRetries
	buckedExists, err = client.client.BucketExists(ctx, bucketName)
	if err != nil {
		// Wrap error around objectstorage so we don't expose the minio errors
		return nil, errors.Wrapf(objectstorage.ErrInternal, "unable to connect to minio: %s", err)
	}
	if !buckedExists {
		// Wrap error around objectstorage so we don't expose the minio errors
		return nil, errors.Wrap(objectstorage.ErrNotFound, "bucket does not exist")
	}

	return &bucketClient{
		client:     minioClient,
		bucketName: bucketName,
	}, nil
}

// Write writes the bytes to the bucket for the given object
func (c *bucketClient) Write(ctx context.Context, object *objectstorage.Object, data io.Reader, overwrite bool) error {
	if object == nil {
		return errors.Wrap(objectstorage.ErrInvalidArgument, "object is nil")
	}

	if object.Size == 0 {
		return errors.Wrap(objectstorage.ErrInvalidArgument, "object size is 0")
	}

	if object.Key == "" {
		return errors.Wrap(objectstorage.ErrInvalidArgument, "object key is empty")
	}

	if !overwrite {
		stat, err := c.Stat(ctx, object.Key)
		if err != nil && !errors.Is(err, objectstorage.ErrNotFound) {
			return errors.Wrap(objectstorage.ErrInternal, "unable to stat object")
		}

		if stat != nil {
			return errors.Wrapf(objectstorage.ErrAlreadyExists, "overwrite was false when writing: %s/%s", c.bucketName, object.Key)
		}
	}

	_, err := c.client.PutObject(ctx, c.bucketName, object.Key, data, object.Size, minio.PutObjectOptions{
		ContentType: object.ContentType,
	})
	if err != nil {
		return fmt.Errorf("unable to write file to bucket: %w", err)
	}
	return nil
}

// WriteBytes is a wrapper function around Write with []byte instead of an io.Reader as input
func (c *bucketClient) WriteBytes(ctx context.Context, object *objectstorage.Object, data []byte, overwrite bool) error {
	reader := bytes.NewReader(data)
	if int64(len(data)) != object.Size {
		return errors.Wrap(objectstorage.ErrInvalidArgument, "len(data) not equal to object.Size")
	}

	return c.Write(ctx, object, reader, overwrite)
}

// List lists all objects in the bucket with the given prefix
// an error is returned if the bucket does not exist or a transport error occurs
// note that the objects are returned in no particular (or deterministic) order, this purely dependent on the bucket
func (c *bucketClient) List(ctx context.Context, prefix string) ([]objectstorage.Object, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	objects := c.client.ListObjects(ctx, c.bucketName, minio.ListObjectsOptions{Prefix: prefix, Recursive: true, UseV1: true})

	resp := []objectstorage.Object{}
	for obj := range objects {
		object := obj
		if obj.Err != nil {
			return nil, errors.Wrapf(objectstorage.ErrInternal, "unable to list objects: %s", obj.Err)
		}
		resp = append(resp, *objInfoToobjectstorageObject(&object))
	}

	return resp, nil
}

// Read reads a file from the bucket, the reader should be closed
func (c *bucketClient) Read(ctx context.Context, key string) (io.ReadCloser, error) {
	if key == "" {
		return nil, errors.Wrap(objectstorage.ErrInvalidArgument, "key is empty")
	}

	obj, err := c.client.GetObject(ctx, c.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		errorResponse := minio.ToErrorResponse(err)
		if errorResponse.Code == minioNoSuchKeyError {
			return nil, objectstorage.ErrNotFound
		}
		return nil, errors.Wrapf(objectstorage.ErrInternal, "unable to get object from bucket: %s", err)
	}

	return obj, nil
}

// ReadBytes is a utility function to read an object to a bytearray
func (c *bucketClient) ReadBytes(ctx context.Context, key string) ([]byte, error) {
	reader, err := c.Read(ctx, key)
	if err != nil {
		return nil, err
	}

	ret, err := ioutil.ReadAll(reader)

	// Even if we have an error reading, still try to close the reader first before we return
	closerErr := reader.Close()
	if closerErr != nil {
		return nil, errors.Wrapf(objectstorage.ErrInternal, "unable to close object: %s", closerErr)
	}

	if err != nil {
		if minio.ToErrorResponse(err).Code == minioNoSuchKeyError {
			return nil, objectstorage.ErrNotFound
		}
		return nil, errors.Wrapf(objectstorage.ErrInternal, "unable to read object: %s", err)
	}

	return ret, nil
}

// GetPresignedObjectURL returns an url that can be used to access the object.
// Validity sets for how long the duration is valid (should be between 1sec and 7 days)
// Subject specifies the object
func (c *bucketClient) GetPresignedObjectURL(ctx context.Context, subject, method string, validity time.Duration) (string, error) {
	url, err := c.client.Presign(ctx, method, c.bucketName, subject, validity, nil)
	if err != nil {
		return "", fmt.Errorf("unable to sign url: %w", err)
	}
	return url.String(), nil
}

// Stat returns information about an object
func (c *bucketClient) Stat(ctx context.Context, key string) (*objectstorage.Object, error) {
	if key == "" {
		return nil, errors.Wrap(objectstorage.ErrInvalidArgument, "key is empty")
	}

	stat, err := c.client.StatObject(ctx, c.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == minioNoSuchKeyError {
			return nil, objectstorage.ErrNotFound
		}
		return nil, errors.Wrapf(objectstorage.ErrInternal, "unable to stat object: %s", err)
	}
	if stat.Err != nil {
		return nil, errors.Wrapf(objectstorage.ErrInternal, "unable to stat object: %s", stat.Err)
	}

	return objInfoToobjectstorageObject(&stat), nil
}

// Delete deletes an object
func (c *bucketClient) Delete(ctx context.Context, key string) error {
	if key == "" {
		return errors.Wrap(objectstorage.ErrInvalidArgument, "key is empty")
	}
	err := c.client.RemoveObject(ctx, c.bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		return objectstorage.ErrInternal
	}
	return nil
}

func objInfoToobjectstorageObject(in *minio.ObjectInfo) *objectstorage.Object {
	return &objectstorage.Object{
		Key:         in.Key,
		Size:        in.Size,
		ContentType: in.ContentType,
	}
}
