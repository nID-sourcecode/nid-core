// Package contract stores the errors and interfaces for the external-authorization-chain service.
package contract

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

// App defines the method for the logic for this handler
type App interface {
	CheckEndpoints(ctx context.Context, request *authv3.CheckRequest) error
}
