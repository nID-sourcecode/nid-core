// Package extproc contains an implementation of envoy's ExternalProcessorServer that processes requests using a filter chain.
package extproc

import (
	"context"
	"io"

	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3alpha"
	"github.com/gofrs/uuid"

	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	"lab.weave.nl/nid/nid-core/pkg/extproc/mutation"
	"lab.weave.nl/nid/nid-core/pkg/extproc/requestcontext"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
)

// ExternalProcessorServer processes HTTP request/response streams
type ExternalProcessorServer struct {
	requestContextFactory requestcontext.Factory
}

// NewExternalProcessorServer creates a new external processor server
func NewExternalProcessorServer(filterInitializers []filter.Initializer) *ExternalProcessorServer {
	return &ExternalProcessorServer{
		requestContextFactory: requestcontext.NewFilterChainRequestContextFactory(
			filterInitializers,
			&mutation.DefaultCalculator{},
		),
	}
}

// Process processes a single HTTP request/response stream
func (s *ExternalProcessorServer) Process(stream ext_proc_pb.ExternalProcessor_ProcessServer) error {
	ctx := stream.Context()

	id := uuid.Must(uuid.NewV4()).String()
	requestContext, err := s.requestContextFactory.NewRequestContext(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "receiving request")
		}

		// handle in
		switch v := in.Request.(type) {
		case *ext_proc_pb.ProcessingRequest_RequestHeaders:
			if err := onRequestHeaders(ctx, id, stream, requestContext, v); err != nil {
				return err
			}
		case *ext_proc_pb.ProcessingRequest_RequestBody:
			if err := onRequestBody(ctx, id, stream, requestContext, v); err != nil {
				return err
			}
		case *ext_proc_pb.ProcessingRequest_ResponseHeaders:
			if err := onResponseHeaders(ctx, id, stream, requestContext, v); err != nil {
				return err
			}
		case *ext_proc_pb.ProcessingRequest_ResponseBody:
			if err := onResponseBody(ctx, id, stream, requestContext, v); err != nil {
				return err
			}
		}
	}
}

func onRequestHeaders(ctx context.Context, contextID string, stream ext_proc_pb.ExternalProcessor_ProcessServer, requestContext requestcontext.RequestContext, v *ext_proc_pb.ProcessingRequest_RequestHeaders) error {
	log.Infof("on request headers: %s", contextID)
	processingResponse, err := requestContext.OnRequestHeaders(ctx, v)
	if err != nil {
		log.WithError(err).Error("processing request headers")
		return grpcerrors.ErrInternalServer()
	}

	return sendResponse(contextID, "request headers", stream, processingResponse)
}

func onResponseHeaders(ctx context.Context, contextID string, stream ext_proc_pb.ExternalProcessor_ProcessServer, requestContext requestcontext.RequestContext, v *ext_proc_pb.ProcessingRequest_ResponseHeaders) error {
	log.Infof("on response headers: %s", contextID)
	processingResponse, err := requestContext.OnResponseHeaders(ctx, v)
	if err != nil {
		log.WithError(err).Error("processing response headers")
		return grpcerrors.ErrInternalServer()
	}

	return sendResponse(contextID, "response headers", stream, processingResponse)
}

func onRequestBody(ctx context.Context, contextID string, stream ext_proc_pb.ExternalProcessor_ProcessServer, requestContext requestcontext.RequestContext, v *ext_proc_pb.ProcessingRequest_RequestBody) error {
	log.Infof("on request body: %s", contextID)
	processingResponse, err := requestContext.OnRequestBody(ctx, v)
	if err != nil {
		log.WithError(err).Error("processing request body")
		return grpcerrors.ErrInternalServer()
	}

	return sendResponse(contextID, "request body", stream, processingResponse)
}

func onResponseBody(ctx context.Context, contextID string, stream ext_proc_pb.ExternalProcessor_ProcessServer, requestContext requestcontext.RequestContext, v *ext_proc_pb.ProcessingRequest_ResponseBody) error {
	log.Infof("on response body: %s", contextID)
	processingResponse, err := requestContext.OnResponseBody(ctx, v)
	if err != nil {
		log.WithError(err).Error("processing response body")
		return grpcerrors.ErrInternalServer()
	}

	return sendResponse(contextID, "response body", stream, processingResponse)
}

func sendResponse(contextID string, part string, stream ext_proc_pb.ExternalProcessor_ProcessServer, response *ext_proc_pb.ProcessingResponse) error {
	log.Infof("sending (%s) (%s): %+v", part, contextID, response)
	err := stream.Send(response)
	if err != nil {
		log.WithError(err).Errorf("sending %s response", part)
		return grpcerrors.ErrInternalServer()
	}

	return nil
}
