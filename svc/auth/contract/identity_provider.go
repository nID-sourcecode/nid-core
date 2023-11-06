package contract

import (
	"context"

	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

// IdentityProvider provides the identity of request by using metadata.
type IdentityProvider interface {
	GetIdentity(context.Context, *models.TokenRequestMetadata) (string, error)
}

// IdentityProviderType is an enum type of valid IdentityProviders.
type IdentityProviderType string

const (
	// IdentityProviderTypeCertificate is the certificate identity provider.
	IdentityProviderTypeCertificate IdentityProviderType = "certificate"
	// IdentityProviderTypeDatabase is the database identity provider.
	IdentityProviderTypeDatabase IdentityProviderType = "database"
)

// Validate checks if the IdentityProviderType is valid.
func (t IdentityProviderType) Validate() error {
	switch t {
	case IdentityProviderTypeCertificate, IdentityProviderTypeDatabase:
		return nil
	default:
		return ErrIncorrectEnvironmentConfig
	}
}
