package grpcserver

import (
	"context"
	"runtime/debug"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/logfields"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/loggrpc"
)

// getRecoveryFunction returns a function that attempts to recover the panic and return the correct response.
// The problemDetails parameter specifies if the recoveryFunction should return problem details.
func getRecoveryFunction() grpc_recovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, p interface{}) error {
		//nolint:staticcheck //We should remvoe the usage of this deprecated function https://lab.weave.nl/weave/utilities/grpcserver/-/issues/6
		loggrpc.LoggerFromContext(ctx).WithField(logfields.RecoveredPanic, p).WithField(logfields.Stack, string(debug.Stack())).Error("unhandled panic recovered in panicRecoveryHandlerFunc")
		return status.Errorf(codes.Internal, "Internal server error")
	}
}
