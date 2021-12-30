package s3

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
)

type s3ClientMock struct {
	mock.Mock
}

func (m *s3ClientMock) BucketExists(ctx context.Context, s2 string) (bool, error) {
	args := m.Called(ctx, s2)
	return args.Bool(0), args.Error(1)
}

func (m *s3ClientMock) StatObject(ctx context.Context, s2, s3 string, options minio.StatObjectOptions) (minio.ObjectInfo, error) {
	args := m.Called(ctx, s2, s3, options)
	return args.Get(0).(minio.ObjectInfo), args.Error(1)
}

// We implement an interface therefore we cannot make options a pointer
func (m *s3ClientMock) PutObject(ctx context.Context, s2, s3 string, reader io.Reader, i int64, options minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, s2, s3, reader, i, options)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

func (m *s3ClientMock) ListObjects(ctx context.Context, s2 string, options minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	args := m.Called(ctx, s2, options)
	return args.Get(0).(<-chan minio.ObjectInfo)
}

func (m *s3ClientMock) GetObject(ctx context.Context, s2, s3 string, options minio.GetObjectOptions) (*minio.Object, error) {
	args := m.Called(ctx, s2, s3, options)
	return args.Get(0).(*minio.Object), args.Error(1)
}

func (m *s3ClientMock) Presign(ctx context.Context, s2, s3, s4 string, duration time.Duration, values url.Values) (*url.URL, error) {
	args := m.Called(ctx, s2, s3, s4, duration, values)
	return args.Get(0).(*url.URL), args.Error(1)
}

func (m *s3ClientMock) RemoveObject(ctx context.Context, s2, s3 string, opts minio.RemoveObjectOptions) error {
	args := m.Called(ctx, s2, s3, opts)
	return args.Error(0)
}
