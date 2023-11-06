// This package contains implementations of the IdentityProvider interface.
// The IdentityProvider interface provides the identity of request by using metadata.
package identityprovider

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

// DatabaseIdentityProvider returns the identity from the database.
type DatabaseIdentityProvider struct {
	clientDB        *models.ClientDB
	passwordManager password.IManager
}

// NewDatabaseIdentityProvider returns a new instance of DatabaseIdentityProvider.
func NewDatabaseIdentityProvider(clientDB *models.ClientDB, passwordManager password.IManager) *DatabaseIdentityProvider {
	return &DatabaseIdentityProvider{
		clientDB:        clientDB,
		passwordManager: passwordManager,
	}
}

// GetIdentity returns the identity from the database.
func (p *DatabaseIdentityProvider) GetIdentity(ctx context.Context, metadata *models.TokenRequestMetadata) (string, error) {
	client, err := p.clientDB.GetClientByID(metadata.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.Wrapf(contract.ErrUnauthenticated, "invalid username")
		}

		return "", errors.Wrapf(contract.ErrInternalError, "unable to get client: %s", metadata.Username)
	}

	passwordMatches, err := p.passwordManager.ComparePassword(metadata.Password, client.Password)
	if err != nil {
		return "", errors.Wrapf(contract.ErrUnauthenticated, "unable to compare password")
	}
	if !passwordMatches {
		return "", errors.Wrapf(contract.ErrUnauthenticated, "incorrect password")
	}

	return metadata.Username, nil
}
