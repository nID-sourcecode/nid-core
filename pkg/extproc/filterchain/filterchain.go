// Package filterchain contains logic for chaining filters and applying them in order
package filterchain

import (
	"context"

	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

	"github.com/nID-sourcecode/nid-core/pkg/extproc/filter"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

var _ Chain = &DefaultChain{}

// Chain is responsible for chaining filters and applying them in order
type Chain interface {
	ProcessRequestHeaders(ctx context.Context, requestHeaders map[string]string, skipBody bool) (*ProcessingResponse, error)
	ProcessResponseHeaders(responseHeaders map[string]string)
	ProcessRequestBody(ctx context.Context, requestBody []byte) (*ProcessingResponse, error)
	ProcessResponseBody(ctx context.Context, requestBody []byte) (*ProcessingResponse, error)
}

// DefaultChain is the default Chain implementation
type DefaultChain struct {
	filters         []filter.Filter
	requestHeaders  map[string]string
	responseHeaders map[string]string
}

// Error definitions
var (
	ErrFailedToInitializeFilter   = errors.New("failed to initialise filter")
	ErrNonexistingFilterSpecified = errors.New("nonexisting filter specified")
)

// BuildDefaultFilterChain creates a new filter chain based on a list of selected filter names and a list of filter initializers.
func BuildDefaultFilterChain(selectedFilterNames []string, filterInitializers []filter.Initializer) (*DefaultChain, error) {
	filtersSelected := map[string]bool{}
	filtersApplied := map[string]bool{}
	for _, filterName := range selectedFilterNames {
		filtersSelected[filterName] = true
		filtersApplied[filterName] = false
	}

	chain := &DefaultChain{filters: []filter.Filter{}}
	for _, filterInitializer := range filterInitializers {
		filterName := filterInitializer.Name()
		filterSelected, ok := filtersSelected[filterName]
		if ok && filterSelected {
			thisFilter, err := filterInitializer.NewFilter()
			if err != nil {
				return nil, errors.Errorf("%w: %s", ErrFailedToInitializeFilter, filterName)
			}
			chain.filters = append(chain.filters, thisFilter)
			filtersApplied[filterName] = true
		}
	}

	for filterName, filterApplied := range filtersApplied {
		if !filterApplied {
			return nil, errors.Errorf("%w: %s", ErrNonexistingFilterSpecified, filterName)
		}
	}

	return chain, nil
}

// ProcessingResponse contains the results of processing a request or response in a filter chain.
type ProcessingResponse struct {
	Headers           map[string]string
	Body              []byte
	ImmediateResponse *ext_proc_pb.ImmediateResponse
}

// ProcessRequestHeaders processes request headers. If the body is skipped, the request is processed by the filter chain. Otherwise, the chain waits for the body.
func (c *DefaultChain) ProcessRequestHeaders(ctx context.Context, requestHeaders map[string]string, skipBody bool) (*ProcessingResponse, error) {
	c.requestHeaders = requestHeaders

	if skipBody {
		processingResponse, err := c.process(ctx, nil, ProcessTypeRequest)

		return processingResponse, errors.Wrap(err, "processing request")
	}

	return &ProcessingResponse{}, nil
}

// ProcessResponseHeaders saves the response headers to be processed when the body arrives. FIXME this needs to become configurable.
func (c *DefaultChain) ProcessResponseHeaders(responseHeaders map[string]string) {
	c.responseHeaders = responseHeaders
}

// ErrInvalidProcessingType is returned when an invalid processing type is used.
var ErrInvalidProcessingType = errors.New("invalid processing type")

func (c *DefaultChain) process(ctx context.Context, body []byte, processType ProcessType) (*ProcessingResponse, error) {
	var headers map[string]string
	filters := make([]filter.Filter, len(c.filters))

	// Invert the filters on a response
	switch processType {
	case ProcessTypeRequest:
		headers = c.requestHeaders
		copy(filters, c.filters)
	case ProcessTypeResponse:
		headers = c.responseHeaders
		for i, thisFilter := range c.filters { // invert filter chain
			filters[len(c.filters)-i-1] = thisFilter
		}
	default:
		return nil, ErrInvalidProcessingType
	}

	for _, thisFilter := range filters {
		log.Infof("applying filter %s", thisFilter.Name())

		var res *filter.ProcessingResponse
		var err error

		switch processType {
		case ProcessTypeRequest:
			res, err = thisFilter.OnHTTPRequest(ctx, body, headers)
		case ProcessTypeResponse:
			res, err = thisFilter.OnHTTPResponse(ctx, body, headers)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "processing filter %s", thisFilter.Name())
		}

		if res != nil {
			if res.ImmediateResponse != nil {
				return &ProcessingResponse{
					ImmediateResponse: res.ImmediateResponse,
				}, nil // Return response immediately
			}
			if res.NewHeaders != nil {
				headers = res.NewHeaders // Pass changed headers along to next filter in chain
			}
			if res.NewBody != nil {
				body = res.NewBody // Pass changed body along to next filter in chain
			}
		}
	}

	return &ProcessingResponse{
		Headers: headers,
		Body:    body,
	}, nil
}

// ProcessRequestBody passes the request body and headers through the filters.
func (c *DefaultChain) ProcessRequestBody(ctx context.Context, requestBody []byte) (*ProcessingResponse, error) {
	return c.processBody(ctx, requestBody, ProcessTypeRequest)
}

// ProcessResponseBody passes the response body and headers through the filters. This goes in inverse order.
func (c *DefaultChain) ProcessResponseBody(ctx context.Context, requestBody []byte) (*ProcessingResponse, error) {
	return c.processBody(ctx, requestBody, ProcessTypeResponse)
}

// ProcessType is the type of processing to be done
type ProcessType string

// Process types
const (
	ProcessTypeRequest  ProcessType = "REQUEST"
	ProcessTypeResponse ProcessType = "RESPONSE"
)

func (c *DefaultChain) processBody(ctx context.Context, body []byte, processType ProcessType) (*ProcessingResponse, error) {
	processingResponse, err := c.process(ctx, body, processType)

	return processingResponse, errors.Wrap(err, "processing")
}
