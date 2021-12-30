// Package mock provides a mock for the metadata headers
package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

// GRPCMetadataHelperMock is a mock for MetadataHelper
type GRPCMetadataHelperMock struct {
	mock.Mock
}

// CtxHasVal verify if context has given headername present
func (m *GRPCMetadataHelperMock) CtxHasVal(ctx context.Context, headername string) bool {
	args := m.Called(ctx, headername)
	return args.Bool(0)
}

// GetValFromCtx retrieve headername from current context
func (m *GRPCMetadataHelperMock) GetValFromCtx(ctx context.Context, headerName string) (string, error) {
	args := m.Called(ctx, headerName)
	return args.String(0), args.Error(1)
}

// MetadataFromCtx retrieve metadata from context
func (m *GRPCMetadataHelperMock) MetadataFromCtx(ctx context.Context) (metadata.MD, error) {
	args := m.Called(ctx)
	return args.Get(0).(metadata.MD), args.Error(1)
}

// GetBasicAuth mocks the GetBasicAuth function
func (m *GRPCMetadataHelperMock) GetBasicAuth(ctx context.Context) (string, string, error) {
	args := m.Called(ctx)
	return args.String(0), args.String(1), args.Error(2)
}

// GetJWTToken function stub
func (m *GRPCMetadataHelperMock) GetJWTToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// GetMetadataValue function stub
func (m *GRPCMetadataHelperMock) GetMetadataValue(md metadata.MD, headerName string) (string, error) {
	args := m.Called(md, headerName)
	return args.String(0), args.Error(1)
}

// GetIPFromCtx function stub
func (m *GRPCMetadataHelperMock) GetIPFromCtx(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// GetAcceptFromCtx function stub
func (m *GRPCMetadataHelperMock) GetAcceptFromCtx(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
