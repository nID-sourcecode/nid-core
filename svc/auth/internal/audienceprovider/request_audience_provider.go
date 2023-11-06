package audienceprovider

import (
	"context"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

// RequestAudienceProvider provides the audience of the request body.
type RequestAudienceProvider struct{}

// GetAudience returns the audience of the request body.
func (r *RequestAudienceProvider) GetAudience(_ context.Context, req *models.TokenClientFlowRequest, _ []*models.Scope) ([]string, error) {
	if req.Audience == "" {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "audience is required")
	}

	return []string{req.Audience}, nil
}
