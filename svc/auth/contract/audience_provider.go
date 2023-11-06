package contract

import (
	"context"

	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

// AudienceProvider provides the audience of the provided scopes.
type AudienceProvider interface {
	GetAudience(ctx context.Context, req *models.TokenClientFlowRequest, scopes []*models.Scope) ([]string, error)
}

// AudienceProviderType is an enum type of valid AudienceProviders.
type AudienceProviderType string

const (
	AudienceProviderTypeDatabase AudienceProviderType = "database"
	AudienceProviderTypeRequest  AudienceProviderType = "request"
)

func (t AudienceProviderType) Validate() error {
	switch t {
	case AudienceProviderTypeDatabase, AudienceProviderTypeRequest:
		return nil
	default:
		return ErrIncorrectEnvironmentConfig
	}
}
