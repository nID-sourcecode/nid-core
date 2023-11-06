// Package requestcontext contains logic for handling the various parts of http requests in an ext_proc service
package requestcontext

import (
	"context"
	"strings"

	ext_proc_filter "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"google.golang.org/grpc/metadata"

	"github.com/nID-sourcecode/nid-core/pkg/extproc/filter"
	"github.com/nID-sourcecode/nid-core/pkg/extproc/filterchain"
	"github.com/nID-sourcecode/nid-core/pkg/extproc/mutation"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	grpcerrors "github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

// RequestContext represents the context of a single http request. It handles the various messages that can be send to an ext_proc stream.
type RequestContext interface {
	OnRequestHeaders(ctx context.Context, req *ext_proc_pb.ProcessingRequest_RequestHeaders) (*ext_proc_pb.ProcessingResponse, error)
	OnRequestBody(ctx context.Context, req *ext_proc_pb.ProcessingRequest_RequestBody) (*ext_proc_pb.ProcessingResponse, error)
	OnResponseHeaders(ctx context.Context, req *ext_proc_pb.ProcessingRequest_ResponseHeaders) (*ext_proc_pb.ProcessingResponse, error)
	OnResponseBody(ctx context.Context, req *ext_proc_pb.ProcessingRequest_ResponseBody) (*ext_proc_pb.ProcessingResponse, error)
}

// Factory creates new request contexts
type Factory interface {
	NewRequestContext(ctx context.Context) (RequestContext, error)
}

// FilterChainRequestContext creats a request context that runs a filter chain
type FilterChainRequestContext struct {
	mutationCalculator      mutation.Calculator
	filterchain             filterchain.Chain
	originalRequestHeaders  map[string]string
	originalResponseHeaders map[string]string
}

// OnRequestHeaders processes request headers
func (d *FilterChainRequestContext) OnRequestHeaders(ctx context.Context, req *ext_proc_pb.ProcessingRequest_RequestHeaders) (*ext_proc_pb.ProcessingResponse, error) {
	requestHeaders := convertHTTPHeadersToStringMap(req.RequestHeaders)
	d.originalRequestHeaders = make(map[string]string)
	for k, v := range requestHeaders {
		d.originalRequestHeaders[k] = v
	}

	skipBody := requestHeaders[":method"] == "GET"

	headerProcessingResponse, err := d.filterchain.ProcessRequestHeaders(ctx, requestHeaders, skipBody)
	if err != nil {
		return nil, errors.Wrap(err, "processing request headers in filter chain")
	}
	if headerProcessingResponse.ImmediateResponse != nil {
		return &ext_proc_pb.ProcessingResponse{
			Response: &ext_proc_pb.ProcessingResponse_ImmediateResponse{ImmediateResponse: headerProcessingResponse.ImmediateResponse},
		}, nil
	}

	var modeOverride *ext_proc_filter.ProcessingMode
	if skipBody {
		modeOverride = &ext_proc_filter.ProcessingMode{
			RequestHeaderMode:   ext_proc_filter.ProcessingMode_SEND,
			ResponseHeaderMode:  ext_proc_filter.ProcessingMode_SEND,
			RequestBodyMode:     ext_proc_filter.ProcessingMode_NONE,
			ResponseBodyMode:    ext_proc_filter.ProcessingMode_BUFFERED,
			RequestTrailerMode:  ext_proc_filter.ProcessingMode_SKIP,
			ResponseTrailerMode: ext_proc_filter.ProcessingMode_SKIP,
		}
	}

	headerMutation := d.mutationCalculator.CalculateHeaderMutations(d.originalRequestHeaders, headerProcessingResponse.Headers)

	return &ext_proc_pb.ProcessingResponse{
		Response: &ext_proc_pb.ProcessingResponse_RequestHeaders{
			RequestHeaders: &ext_proc_pb.HeadersResponse{Response: &ext_proc_pb.CommonResponse{
				HeaderMutation: headerMutation,
			}},
		}, ModeOverride: modeOverride,
	}, nil
}

// OnRequestBody processes the request body
func (d *FilterChainRequestContext) OnRequestBody(ctx context.Context, req *ext_proc_pb.ProcessingRequest_RequestBody) (*ext_proc_pb.ProcessingResponse, error) {
	body := req.RequestBody.GetBody()
	originalBody := make([]byte, len(body))
	copy(originalBody, body)

	res, err := d.filterchain.ProcessRequestBody(ctx, body)
	if err != nil {
		return nil, errors.Wrap(err, "processing request body in filter chain")
	}

	if res.ImmediateResponse != nil {
		return &ext_proc_pb.ProcessingResponse{
			Response: &ext_proc_pb.ProcessingResponse_ImmediateResponse{ImmediateResponse: res.ImmediateResponse},
		}, nil
	}

	commonResponse := d.convertProcessingResponseToCommonResponse(d.originalRequestHeaders, originalBody, res)
	return &ext_proc_pb.ProcessingResponse{
		Response: &ext_proc_pb.ProcessingResponse_RequestBody{
			RequestBody: &ext_proc_pb.BodyResponse{Response: commonResponse},
		},
	}, nil
}

// OnResponseHeaders processes the response headers
func (d *FilterChainRequestContext) OnResponseHeaders(_ context.Context, req *ext_proc_pb.ProcessingRequest_ResponseHeaders) (*ext_proc_pb.ProcessingResponse, error) {
	responseHeaders := convertHTTPHeadersToStringMap(req.ResponseHeaders)
	d.originalResponseHeaders = make(map[string]string)
	for k, v := range responseHeaders {
		d.originalResponseHeaders[k] = v
	}

	d.filterchain.ProcessResponseHeaders(responseHeaders)

	return &ext_proc_pb.ProcessingResponse{Response: &ext_proc_pb.ProcessingResponse_ResponseHeaders{ResponseHeaders: &ext_proc_pb.HeadersResponse{}}}, nil
}

// OnResponseBody processes the response body
func (d *FilterChainRequestContext) OnResponseBody(ctx context.Context, req *ext_proc_pb.ProcessingRequest_ResponseBody) (*ext_proc_pb.ProcessingResponse, error) {
	body := req.ResponseBody.GetBody()
	originalBody := make([]byte, len(body))
	copy(originalBody, body)

	res, err := d.filterchain.ProcessResponseBody(ctx, body)
	if err != nil {
		return nil, errors.Wrap(err, "processing request body in filter chain")
	}

	if res.ImmediateResponse != nil {
		return &ext_proc_pb.ProcessingResponse{
			Response: &ext_proc_pb.ProcessingResponse_ImmediateResponse{ImmediateResponse: res.ImmediateResponse},
		}, nil
	}

	commonResponse := d.convertProcessingResponseToCommonResponse(d.originalRequestHeaders, originalBody, res)
	return &ext_proc_pb.ProcessingResponse{
		Response: &ext_proc_pb.ProcessingResponse_ResponseBody{
			ResponseBody: &ext_proc_pb.BodyResponse{Response: commonResponse},
		},
	}, nil
}

func (d *FilterChainRequestContext) convertProcessingResponseToCommonResponse(originalHeaders map[string]string, originalBody []byte, response *filterchain.ProcessingResponse) *ext_proc_pb.CommonResponse {
	return &ext_proc_pb.CommonResponse{
		HeaderMutation: d.mutationCalculator.CalculateHeaderMutations(originalHeaders, response.Headers),
		BodyMutation:   d.mutationCalculator.CalculateBodyMutation(originalBody, response.Body),
	}
}

// NewFilterChainRequestContextFactory creates a new filterchain request context factory
func NewFilterChainRequestContextFactory(filterInitializers []filter.Initializer, mutationCalculator mutation.Calculator) *FilterChainRequestContextFactory {
	return &FilterChainRequestContextFactory{
		filterInitializers: filterInitializers,
		mutationCalculator: mutationCalculator,
	}
}

// FilterChainRequestContextFactory is responsible for creating FilterChainRequestContexts
type FilterChainRequestContextFactory struct {
	filterInitializers []filter.Initializer
	mutationCalculator mutation.Calculator
}

// NewRequestContext creates a new FilterChainRequestContext
func (d *FilterChainRequestContextFactory) NewRequestContext(ctx context.Context) (RequestContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Error("received not-ok getting metadata")
		return nil, grpcerrors.ErrInternalServer()
	}
	log.Info("metadata: %+v", md)
	selectedFilterNamesHeaders, ok := md["x-selected-filters"]
	if !ok || len(selectedFilterNamesHeaders) == 0 {
		return nil, grpcerrors.ErrInvalidArgument("no filters metadata specified")
	}
	selectedFilterNamesCommaSeparated := selectedFilterNamesHeaders[0]
	selectedFilterNames := strings.Split(selectedFilterNamesCommaSeparated, ",")

	chain, err := filterchain.BuildDefaultFilterChain(selectedFilterNames, d.filterInitializers)
	if err != nil {
		log.WithError(err).Error("building filter chain")
		return nil, grpcerrors.ErrInternalServer()
	}

	return &FilterChainRequestContext{
		mutationCalculator: d.mutationCalculator,
		filterchain:        chain,
	}, nil
}

func convertHTTPHeadersToStringMap(headers *ext_proc_pb.HttpHeaders) map[string]string {
	headersMap := make(map[string]string)
	for _, header := range headers.GetHeaders().GetHeaders() {
		headersMap[header.Key] = header.Value
	}
	return headersMap
}
