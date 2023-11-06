// Package xrequestid propagation of request id header
package xrequestid

import (
	"context"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const xRequestIDHeaderName = "x-request-id"

// AddXRequestID propagates x-request-id header for the grpc server/client
func AddXRequestID(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	header := md[xRequestIDHeaderName]

	if ok && len(header) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, xRequestIDHeaderName, header[0])
	}

	if len(header) > 1 {
		log.Warn("multiple x-request-id headers found")
	}

	return handler(ctx, req)
}
