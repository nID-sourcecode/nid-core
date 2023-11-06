package s3

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vrischmann/envconfig"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
	httputil "github.com/nID-sourcecode/nid-core/pkg/utilities/http"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/objectstorage"
)

// S3TestSuite tests the s3 package
type S3TestSuite struct {
	grpctesthelpers.ExtendedTestSuite
	httpMock     httputil.ClientMock
	client       *bucketClient
	mockS3Client *s3ClientMock
}

func (s *S3TestSuite) SetupTest() {
	s.httpMock = httputil.ClientMock{}
	s.mockS3Client = &s3ClientMock{}
	s.client = &bucketClient{
		client:     s.mockS3Client,
		bucketName: "test",
	}
}

func TestS3TestSuite(t *testing.T) {
	suite.Run(t, new(S3TestSuite))
}

type createRoundTripper func(*http.Request) (*http.Response, error)

func (f createRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// Source https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLocation.html
const sampleBucketLocationResponse = `
         <?xml version="1.0" encoding="UTF-8"?>
         <LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">Europe</LocationConstraint>
`

func (s *S3TestSuite) TestCreateClient() {
	rt := createRoundTripper(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(sampleBucketLocationResponse))),
		}, nil
	})
	config := &objectstorage.ClientConfig{
		// Use localhost to prevent accidentally calling services if the roundtripper doesn't work
		Host: "localhost",
		Port: 80,
	}

	client, err := NewClient(context.Background(), config, "test", &rt)
	s.Require().NoError(err)
	s.Require().NotNil(client)
}

func (s *S3TestSuite) TestCreateClientInvalidUrl() {
	rt := createRoundTripper(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(sampleBucketLocationResponse))),
		}, nil
	})
	config := &objectstorage.ClientConfig{
		// Use localhost to prevent accidentally calling services if the roundtripper doesn't work
		Host: "localhost:invalid",
		Port: 80,
	}

	client, err := NewClient(context.Background(), config, "test", &rt)
	s.True(errors.Is(err, objectstorage.ErrInternal), "Expected error to be objectstorage.ErrInternal but was: %s", err)
	s.Nil(client)
}

func (s *S3TestSuite) TestDefaultConfigValue() {
	var config *objectstorage.ClientConfig
	os.Setenv("ACCESS_KEY", "test")
	os.Setenv("ACCESS_SECRET", "test")
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "80")
	err := envconfig.Init(&config)
	s.NoError(err)
	s.Equal(true, config.Secure)
}

func (s *S3TestSuite) TestList() {
	channel := make(chan minio.ObjectInfo, 10)
	channel <- minio.ObjectInfo{Key: "N/iets"}
	channel <- minio.ObjectInfo{Key: "N/anders"}

	// Convert channel to read only
	ret := func(c chan minio.ObjectInfo) <-chan minio.ObjectInfo {
		return c
	}(channel)

	wg := sync.WaitGroup{}
	wg.Add(1)

	var objs []objectstorage.Object
	var err error
	go func() {
		s.mockS3Client.On("ListObjects", mock.Anything, "test", mock.Anything).Return(ret)
		objs, err = s.client.List(context.Background(), "N")
		wg.Done()
	}()

	close(channel)
	wg.Wait()

	s.Require().Nil(err)
	s.Require().Len(objs, 2)
	names := []string{objs[0].Key, objs[1].Key}
	s.Contains(names, "N/iets")
	s.Contains(names, "N/anders")
}

func (s *S3TestSuite) TestWrite() {
	key := "test/key/and.ext"
	reader := bytes.NewReader([]byte("Test"))

	ret := minio.UploadInfo{}
	s.mockS3Client.On("PutObject", mock.Anything, "test", key, reader, int64(4), mock.Anything).Return(ret, nil)
	s.mockS3Client.On("StatObject", mock.Anything, "test", key, mock.Anything).Return(minio.ObjectInfo{Key: key}, nil)

	err := s.client.Write(context.Background(), &objectstorage.Object{Key: key, Size: 4}, reader, true)
	s.Require().NoError(err)
	s.mockS3Client.AssertNumberOfCalls(s.T(), "PutObject", 1)
}

func (s *S3TestSuite) TestWriteValidateArguments() {
	reader := bytes.NewReader([]byte("test"))

	s.Run("NilObject", func() {
		err := s.client.Write(context.Background(), nil, reader, false)
		s.True(errors.Is(err, objectstorage.ErrInvalidArgument))
	})

	s.Run("EmptyKey", func() {
		err := s.client.Write(context.Background(), &objectstorage.Object{Size: 10}, reader, false)
		s.True(errors.Is(err, objectstorage.ErrInvalidArgument))
	})

	s.Run("SizeZero", func() {
		err := s.client.Write(context.Background(), &objectstorage.Object{Key: "iets"}, reader, false)
		s.True(errors.Is(err, objectstorage.ErrInvalidArgument))
	})
}

func (s *S3TestSuite) TestWriteBytes() {
	key := "test/key/and.ext"
	data := []byte("Test")

	ret := minio.UploadInfo{}
	s.mockS3Client.On("PutObject", mock.Anything, "test", key, bytes.NewReader(data), int64(4), mock.Anything).Return(ret, nil)
	s.mockS3Client.On("StatObject", mock.Anything, "test", key, mock.Anything).Return(minio.ObjectInfo{Key: key}, nil)

	err := s.client.WriteBytes(context.Background(), &objectstorage.Object{Key: key, Size: 4}, data, true)
	s.Require().NoError(err)
	s.mockS3Client.AssertNumberOfCalls(s.T(), "PutObject", 1)
}

func (s *S3TestSuite) TestWriteBytesInvalidArgument() {
	key := "something"
	data := []byte("four")

	err := s.client.WriteBytes(context.Background(), &objectstorage.Object{Key: key, Size: 1033}, data, false)
	s.True(errors.Is(err, objectstorage.ErrInvalidArgument), "Expected error to be objectstorage.ErrInvalidArgument but was: %s", err)
}

func (s *S3TestSuite) TestWriteOverwriteFalseExists() {
	key := "test"
	data := []byte("something")

	s.mockS3Client.On("StatObject", mock.Anything, "test", key, mock.Anything).Return(minio.ObjectInfo{Key: key}, nil)

	err := s.client.WriteBytes(context.Background(), &objectstorage.Object{Key: key, Size: int64(len(data))}, data, false)

	s.True(errors.Is(err, objectstorage.ErrAlreadyExists), "Expected error to be objectstorage.ErrAlreadyExists but was: %s", err)

	s.mockS3Client.AssertNumberOfCalls(s.T(), "StatObject", 1)
	s.mockS3Client.AssertNumberOfCalls(s.T(), "StatObject", 1)
	s.mockS3Client.AssertNumberOfCalls(s.T(), "PutObject", 0)
}

func (s *S3TestSuite) TestWriteOverwriteFalseNotExists() {
	key := "someKey"
	data := []byte("someData")

	s.mockS3Client.On("StatObject", mock.Anything, "test", key, mock.Anything).Return(minio.ObjectInfo{}, minio.ErrorResponse{Code: minioNoSuchKeyError})
	s.mockS3Client.On("PutObject", mock.Anything, "test", key, bytes.NewReader(data), int64(len(data)), mock.Anything).Return(minio.UploadInfo{Key: "test"}, nil)

	err := s.client.WriteBytes(context.Background(), &objectstorage.Object{Key: key, Size: int64(len(data))}, data, false)

	s.Require().NoError(err)

	s.mockS3Client.AssertNumberOfCalls(s.T(), "StatObject", 1)
	s.mockS3Client.AssertNumberOfCalls(s.T(), "PutObject", 1)
}

func (s *S3TestSuite) TestWriteOverwriteTrueExists() {
	key := "test"
	data := []byte("ja")

	s.mockS3Client.On("StatObject", mock.Anything, "test", key, mock.Anything).Return(minio.ObjectInfo{Key: key}, nil)
	s.mockS3Client.On("PutObject", mock.Anything, "test", key, bytes.NewReader(data), int64(len(data)), mock.Anything).Return(minio.UploadInfo{Key: "test"}, nil)

	err := s.client.WriteBytes(context.Background(), &objectstorage.Object{Key: key, Size: int64(len(data))}, data, true)

	s.Require().NoError(err)

	s.mockS3Client.AssertNumberOfCalls(s.T(), "StatObject", 0)
	s.mockS3Client.AssertNumberOfCalls(s.T(), "PutObject", 1)
}

func (s *S3TestSuite) TestReadErr() {
	key := "some/test/key.txt"

	expectedGetOptions := minio.GetObjectOptions{}
	s.mockS3Client.On("GetObject", mock.Anything, "test", key, expectedGetOptions).Return((*minio.Object)(nil), io.EOF)

	reader, err := s.client.Read(context.Background(), key)
	s.True(errors.Is(err, objectstorage.ErrInternal), "Expected error to be objectstorage.ErrInternal but was: %s", err)
	s.Nil(reader)
}

func (s *S3TestSuite) TestReadEmptyKey() {
	reader, err := s.client.Read(context.Background(), "")
	s.True(errors.Is(err, objectstorage.ErrInvalidArgument), "Expected error to be objectstorage.ErrInvalidArgument but was: %s", err)
	s.Nil(reader)
}

func (s *S3TestSuite) TestDeleteErr() {
	key := "key"
	expectedRemoveOptions := minio.RemoveObjectOptions{}
	s.mockS3Client.On("RemoveObject", mock.Anything, "test", key, expectedRemoveOptions).Return(io.EOF)

	err := s.client.Delete(context.Background(), key)

	s.True(errors.Is(err, objectstorage.ErrInternal), "Expected error to be objectstorage.ErrInternal but was: %s", err)
}

func (s *S3TestSuite) TestDelete() {
	key := "key"
	expectedRemoveOptions := minio.RemoveObjectOptions{}
	s.mockS3Client.On("RemoveObject", mock.Anything, "test", key, expectedRemoveOptions).Return(nil)

	err := s.client.Delete(context.Background(), key)
	s.NoError(err)
	s.mockS3Client.AssertNumberOfCalls(s.T(), "RemoveObject", 1)
}

func (s *S3TestSuite) TestDeleteEmptyKey() {
	key := ""

	err := s.client.Delete(context.Background(), key)
	s.True(errors.Is(err, objectstorage.ErrInvalidArgument), "Expected error to be objectstorage.ErrInvalidArgument but was: %s", err)
}

func (s *S3TestSuite) TestStatEmptyKey() {
	key := ""

	obj, err := s.client.Stat(context.Background(), key)
	s.True(errors.Is(err, objectstorage.ErrInvalidArgument), "Expected error to be objectstorage.ErrInvalidArgument but was: %s", err)
	s.Nil(obj)
}
