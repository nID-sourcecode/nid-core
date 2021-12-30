// Package proto contains the proto definitions for the auth package.
package proto

// CreateSessionRequest is an interface for auth requests that contain data to create a new session.
type CreateSessionRequest interface {
	GetClientId() string
	GetAudience() string
	GetRedirectUri() string
}
