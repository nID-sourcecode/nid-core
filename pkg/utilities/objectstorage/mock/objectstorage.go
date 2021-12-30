// Package mock a mock for the objectstorage client
package mock

import (
	"context"
	"io"
	"time"

	"github.com/stretchr/testify/mock"

	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
)

// Client that mocks a writer
type Client struct {
	mock.Mock
}

// List mocks the List method
func (m *Client) List(ctx context.Context, prefix string) ([]objectstorage.Object, error) {
	args := m.MethodCalled("List", ctx, prefix)
	return args.Get(0).([]objectstorage.Object), args.Error(1)
}

// Write mocks the Write method
func (m *Client) Write(ctx context.Context, obj *objectstorage.Object, data io.Reader, overwrite bool) error {
	args := m.MethodCalled("Write", ctx, obj, data, overwrite)
	return args.Error(0)
}

// WriteBytes mocks the Write method
func (m *Client) WriteBytes(ctx context.Context, obj *objectstorage.Object, data []byte, overwrite bool) error {
	args := m.MethodCalled("WriteBytes", ctx, obj, data, overwrite)
	return args.Error(0)
}

// Read mocks the Read method
func (m *Client) Read(ctx context.Context, key string) (io.ReadCloser, error) {
	args := m.MethodCalled("Read", ctx, key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ReadBytes mocks the ReadBytes method
func (m *Client) ReadBytes(ctx context.Context, key string) ([]byte, error) {
	args := m.MethodCalled("ReadBytes", ctx, key)
	return args.Get(0).([]byte), args.Error(1)
}

// Delete mocks the delete method
func (m *Client) Delete(ctx context.Context, key string) error {
	args := m.MethodCalled("Delete", ctx, key)
	return args.Error(0)
}

// GetPresignedObjectURL mocks the GetPresignedObjectURL method
func (m *Client) GetPresignedObjectURL(ctx context.Context, key, method string, validity time.Duration) (string, error) {
	args := m.MethodCalled("GetPresignedObjectURL", ctx, key, method, validity)
	return args.String(0), args.Error(1)
}

// Stat mocks the Stat function
func (m *Client) Stat(ctx context.Context, key string) (*objectstorage.Object, error) {
	args := m.MethodCalled("Stat", ctx, key)
	return args.Get(0).(*objectstorage.Object), args.Error(1)
}
