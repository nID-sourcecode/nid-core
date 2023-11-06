// Package gqlclient provides an interface for executing gql queries
package gqlclient

import (
	"context"
	goerr "errors"
)

var (
	// ErrRemoteErrorResponse is returned on a remote error
	ErrRemoteErrorResponse = goerr.New("remote graphql error response")
	// ErrMethodNotSupported is returned on an unsupported method
	ErrMethodNotSupported = goerr.New("method not supported")
)

// Method represents the HTTP methods supported by this package
type Method string

const (
	// MethodGet represents the HTTP GET method
	MethodGet Method = "GET"
	// MethodPost represents the HTTP POST method
	MethodPost Method = "POST"
)

// Client graphql client interface
type Client interface {
	// Post executes gql query on gqlClient
	Post(ctx context.Context, req Request, resp interface{}) error // Defaults to POST
	// Get executes gql query on gqlClient
	Get(ctx context.Context, req Request, resp interface{}) error
	// Run executes gql query on gqlClient
	Run(ctx context.Context, req Request, resp interface{}, method Method) error
}

// NewClient creates GraphQL client for a new URL
func NewClient(url string) Client {
	return NewRestyGQLClientFactory().NewClient(url)
}

// Request is a graphQL request
type Request struct {
	Query     string
	Variables map[string]interface{}
	Headers   map[string]string
}

// NewRequest initialises a new request
func NewRequest(query string) Request {
	return Request{
		Query:     query,
		Variables: make(map[string]interface{}),
		Headers:   make(map[string]string),
	}
}

// ClientFactory creates new clients
type ClientFactory interface {
	NewClient(url string) Client
}
