package audienceprovider

import (
	"context"

	"github.com/nID-sourcecode/nid-core/pkg/sliceutil"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/config"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

// DatabaseAudienceProvider provides the audience of the provided scopes in the database.
type DatabaseAudienceProvider struct {
	conf *config.AuthConfig
}

// NewDatabaseAudienceProvider creates a new DatabaseAudienceProvider.
func NewDatabaseAudienceProvider(conf *config.AuthConfig) *DatabaseAudienceProvider {
	return &DatabaseAudienceProvider{
		conf: conf,
	}
}

// GetAudience returns the audience for the provided scopes in the database models.
func (r *DatabaseAudienceProvider) GetAudience(_ context.Context, req *models.TokenClientFlowRequest, scopes []*models.Scope) ([]string, error) {
	audiences := []string{}

	// Get the audiences from the scopes.
	for _, scope := range scopes {
		for _, audience := range scope.Audiences {
			audiences = append(audiences, audience.Audience)
		}
	}

	audiences = sliceutil.RemoveDuplicates(audiences)

	// Use the audience in the body if it is correct.
	if sliceutil.Contains(audiences, req.Audience) {
		audiences = []string{req.Audience}
	} else if len(audiences) > 1 && !r.conf.AllowMultipleAudiences {
		return nil, errors.Wrap(contract.ErrInvalidArguments, "Multiple audiences are not allowed, provide the audience in the request body")
	}

	return audiences, nil
}
