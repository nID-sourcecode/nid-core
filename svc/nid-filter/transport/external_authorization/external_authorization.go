// Package externalauthorization transport layer for authv3 of envoyproxy.
package externalauthorization

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/nid-filter/contract"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// ExternalAuthz stores and executes the filters that must be run for the nid-filter requests.
type ExternalAuthz struct {
	filters []contract.AuthorizationRule
}

// New returns new instance of ExternalAuthz
func New(filters []contract.AuthorizationRule) *ExternalAuthz {
	return &ExternalAuthz{
		filters: filters,
	}
}

// Check implements the external authorization interface and runs Check of every filter
func (c *ExternalAuthz) Check(ctx context.Context, request *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	for i := 0; i < len(c.filters); i++ {
		filter := c.filters[i]
		err := filter.Check(ctx, request)

		if err == nil {
			continue
		}

		log.WithError(err).Errorf("checking request in filter %s", filter.Name())
		return returnUnauthenticated()
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

func returnUnauthenticated() (*authv3.CheckResponse, error) {
	return &authv3.CheckResponse{
		Status: &status.Status{
			Code:    int32(codes.Unauthenticated),
			Message: "unauthenticated",
			Details: nil,
		},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: typev3.StatusCode_Unauthorized,
				},
				Headers: nil,
				Body:    "unauthorized",
			},
		},
		DynamicMetadata: nil,
	}, nil
}
