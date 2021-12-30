package http

import "net/http"

// Client is a simple http client interface
type Client interface {
	Do(*http.Request) (*http.Response, error)
}
