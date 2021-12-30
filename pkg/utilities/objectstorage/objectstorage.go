// Package objectstorage provides a generic interface for (cloud) object storage, currently an implementation using the S3 API standard is provided.
// This package supports multiple cloud providers like azure and gcloud
package objectstorage

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Object represents a blob storage object
type Object struct {
	Key         string
	Size        int64
	ContentType string
	// FIXME add support for versioning objects https://lab.weave.nl/weave/utilities/objectstorage/-/issues/1
}

// ClientConfig contains the configuration of a S3 compilable storage
type ClientConfig struct {
	AccessKey    string
	AccessSecret string
	Host         string
	Port         int
	Secure       bool `envconfig:"default=true"`
}

// ErrNotFound is returned if the object does not exist
var ErrNotFound = fmt.Errorf("not found")

// ErrInvalidArgument is returned when an invalid argument is provided
var ErrInvalidArgument = fmt.Errorf("invalid argument")

// ErrAlreadyExists is returned when an object is modified but it already exists (and should not be overwritten)
var ErrAlreadyExists = fmt.Errorf("object already exists")

// ErrInternal is returned when an internal or transport error occurs
var ErrInternal = fmt.Errorf("internal error")

// Client interface for the storage bucket client
type Client interface {
	// List lists all objects in the bucket with the given prefix, an error is returned if a transport error occurs
	List(ctx context.Context, prefix string) ([]Object, error)

	// Write writes some data to the object at the specified key, an error is returned if the object cannot be written to or a transport error occurs
	Write(ctx context.Context, object *Object, input io.Reader, overwrite bool) error

	// WriteBytes is a utility function for Write where a byte array is provided instead of an io.Reader
	// for performance reasons this should probably only be used for small objects
	WriteBytes(ctx context.Context, object *Object, data []byte, overwrite bool) error

	// Read reads the object from the specified key, an error is returned if the object does not exist or an transport error occurs
	Read(ctx context.Context, key string) (io.ReadCloser, error)

	// ReadBytes is a utility function around Read, instead of returning an io.Reader a byte array is used
	// for performance reasons this should only be used for small objects
	ReadBytes(context.Context, string) ([]byte, error)

	// Delete deletes an object
	Delete(context.Context, string) error

	// GetPresignedObjectURL creates a signed url for an object, this url can then be used to access this object for the duration provided
	// an error is returned if the url cannot be signed.
	// Validity can to be between 1 second and 7 days.
	GetPresignedObjectURL(ctx context.Context, key, method string, validity time.Duration) (string, error)

	// Stat stats an object and returns information about it, nil is returned if the object does not exist and an error is returned if a transport error occurs
	Stat(ctx context.Context, key string) (*Object, error)
}
