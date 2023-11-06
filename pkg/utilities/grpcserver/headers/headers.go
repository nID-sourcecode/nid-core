// Package headers provides functionality for handling grpc headers
package headers

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/grpckeys"
)

// Error definitions
var (
	ErrUnableToDecodeHeaders = fmt.Errorf("could not decode headers")
	ErrHeaderNotFound        = fmt.Errorf("header not found")
)

// MetadataHelper metadatahelper is an interface for retrieving header values from
// context and metadata
type MetadataHelper interface {
	GetMetadataValue(md metadata.MD, headerName string) (string, error)
	CtxHasVal(ctx context.Context, headername string) bool
	GetValFromCtx(ctx context.Context, headerName string) (string, error)
	MetadataFromCtx(ctx context.Context) (metadata.MD, error)
	GetBasicAuth(ctx context.Context) (string, string, error)
	GetJWTToken(ctx context.Context) (string, error)
	GetIPFromCtx(ctx context.Context) (string, error)
	GetAcceptFromCtx(ctx context.Context) (string, error)
}

// GRPCMetadataHelper grpc metadata helper implementation
type GRPCMetadataHelper struct{}

// GetMetadataValue retrieve field from metadata
func (m *GRPCMetadataHelper) GetMetadataValue(md metadata.MD, headerName string) (string, error) {
	headerVal := md.Get(headerName)
	if len(headerVal) != 1 {
		return "", ErrHeaderNotFound
	}
	return headerVal[0], nil
}

// CtxHasVal verify if context has given headername present
func (m *GRPCMetadataHelper) CtxHasVal(ctx context.Context, headername string) bool {
	_, err := m.GetValFromCtx(ctx, headername)
	return err == nil
}

// GetValFromCtx retrieve headername from current context
func (m *GRPCMetadataHelper) GetValFromCtx(ctx context.Context, headerName string) (string, error) {
	md, err := m.MetadataFromCtx(ctx)
	if err != nil {
		return "", err
	}
	return m.GetMetadataValue(md, headerName)
}

// MetadataFromCtx retrieve metadata from context
func (m *GRPCMetadataHelper) MetadataFromCtx(ctx context.Context) (metadata.MD, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return md, ErrUnableToDecodeHeaders
	}
	return md, nil
}

// GetIPFromCtx retrieve IP from context
func (m *GRPCMetadataHelper) GetIPFromCtx(ctx context.Context) (string, error) {
	externalAddress, err := m.GetValFromCtx(ctx, grpckeys.EnvoyExternalAddress.String())
	if err == nil {
		return externalAddress, nil
	}

	forwardedFor, err := m.GetValFromCtx(ctx, grpckeys.ForwardedFor.String())
	if err != nil {
		return "", err
	}
	forwards := strings.Split(forwardedFor, ",")
	// The left most entry in the forwarded for header is the most clienty https://en.wikipedia.org/wiki/X-Forwarded-For
	return forwards[0], nil
}

// GetAcceptFromCtx retrieve accept header from context
func (m *GRPCMetadataHelper) GetAcceptFromCtx(ctx context.Context) (string, error) {
	return m.GetValFromCtx(ctx, grpckeys.AcceptKey.String())
}
