// Package extauthz implements and runs the transport layer for authv3 of envoyproxy
package extauthz

import (
	"context"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/contract"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// ExternalAuthorization implements the authv3 of envoyproxy and checks each incoming requests.
type ExternalAuthorization struct {
	app contract.App
}

// New returns new instance of ExternalAuthorization
func New(app contract.App) *ExternalAuthorization {
	return &ExternalAuthorization{
		app: app,
	}
}

// Check implements the external authorization interface and runs Check of every filter
func (c *ExternalAuthorization) Check(ctx context.Context, request *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	if err := c.app.CheckEndpoints(ctx, request); err != nil {
		from, to := getSourceAndPathFromRequest(request)
		log.WithError(err).Errorf("tried checking the request from: %s to: %s", from, to)

		if errors.Is(err, contract.ErrTargetServiceHeaderNotFound) {
			return returnUnauthenticated("missing target header for service")
		}

		return returnUnauthenticated("unauthenticated")
	}

	return &authv3.CheckResponse{
		Status: &status.Status{
			Code:    int32(codes.OK),
			Message: "",
			Details: nil,
		},
		HttpResponse:    nil,
		DynamicMetadata: nil,
	}, nil
}

func getSourceAndPathFromRequest(request *authv3.CheckRequest) (string, string) {
	from := "nil"
	to := "nil"
	if request.GetAttributes() != nil {
		from = request.GetAttributes().GetSource().Address.String()
		httpRequest := request.GetAttributes().GetRequest().GetHttp()
		to = httpRequest.GetHost() + httpRequest.GetPath()
	}
	return from, to
}

func returnUnauthenticated(message string) (*authv3.CheckResponse, error) {
	return &authv3.CheckResponse{
		Status: &status.Status{
			Code:    int32(codes.Unauthenticated),
			Message: message,
			Details: nil,
		},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: typev3.StatusCode_Unauthorized,
				},
				Headers: nil,
				Body:    message,
			},
		},
		DynamicMetadata: nil,
	}, nil
}
