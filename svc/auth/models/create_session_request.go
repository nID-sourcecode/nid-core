package models

// CreateSessionRequest is an interface for auth requests that contain data to create a new session.
type CreateSessionRequest interface {
	GetClientID() string
	GetAudience() string
	GetRedirectURI() string
}
