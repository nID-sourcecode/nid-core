// Package contract has the interfaces for external-authorization-chain
package contract

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

// AuthorizationRule interface for Authorization implementation
type AuthorizationRule interface {
	Name() string
	Check(ctx context.Context, request *authv3.CheckRequest) error
}
