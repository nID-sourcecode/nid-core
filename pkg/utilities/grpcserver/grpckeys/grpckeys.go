// Package grpckeys provides the keys for the grpc metadata
package grpckeys

type reqHeader string

func (r reqHeader) String() string {
	return string(r)
}

// ContextTag is a tag that can be set on the context
type ContextTag string

func (c ContextTag) String() string {
	return string(c)
}

// constant header keys
const (
	// DefaultXRequestIDKey default header for the request id
	DefaultXRequestIDKey reqHeader = "x-request-id"

	EnvoyPathKey         ContextTag = "x-envoy-original-path"
	EnvoyExternalAddress ContextTag = "x-envoy-external-address"
	ForwardedFor         ContextTag = "x-forwarded-for"
	RequestIDKey         ContextTag = "x-request-id"
	AuthorizationKey     ContextTag = "authorization"
	AcceptKey            ContextTag = "accept"
	UserIDKey            ContextTag = "user_id"
	APIKeyIDKey          ContextTag = "api_key_id"
	AccountIDKey         ContextTag = "account_id"
	IPAddress            ContextTag = "ip_address"
)
