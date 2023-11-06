// Package filter contains interface definitions for http filters to be used in filterchains.
package filter

import (
	"context"

	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

// ProcessingResponse contains the information returned from a filter.
type ProcessingResponse struct {
	NewHeaders        map[string]string
	NewBody           []byte
	ImmediateResponse *ext_proc_pb.ImmediateResponse
}

// Initializer creates new filters.
type Initializer interface {
	NewFilter() (Filter, error)
	Name() string
}

// Filter processes HTTP requests and their corresponding responses.
type Filter interface {
	OnHTTPRequest(ctx context.Context, body []byte, headers map[string]string) (*ProcessingResponse, error)
	OnHTTPResponse(ctx context.Context, body []byte, headers map[string]string) (*ProcessingResponse, error)
	Name() string
}

// DefaultFilter is the base implementation for a filter. It implements the two HTTP methods but does nothing.
type DefaultFilter struct{}

// OnHTTPRequest processes an HTTP request. The default implementation is to do nothing.
func (d DefaultFilter) OnHTTPRequest(_ context.Context, _ []byte, _ map[string]string) (*ProcessingResponse, error) {
	return nil, nil
}

// OnHTTPResponse processes an HTTP response. The default implementation is to do nothing.
func (d DefaultFilter) OnHTTPResponse(_ context.Context, _ []byte, _ map[string]string) (*ProcessingResponse, error) {
	return nil, nil
}
