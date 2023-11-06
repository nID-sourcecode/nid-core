// This package contains implementations of the IdentityProvider interface.
// The IdentityProvider interface provides the identity of request by using metadata.
package identityprovider

import (
	"context"
	"regexp"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/identityprovider/config"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

// CertificateIdentityProvider returns the identity from the client certificate header.
type CertificateIdentityProvider struct {
	conf *config.CertificateIdentityProviderConfig
}

// NewCertificateIdentityProvider returns a new instance of CertificateIdentityProvider.
func NewCertificateIdentityProvider(conf *config.CertificateIdentityProviderConfig) *CertificateIdentityProvider {
	return &CertificateIdentityProvider{
		conf: conf,
	}
}

// GetIdentity returns the identity from the certificate.
func (p *CertificateIdentityProvider) GetIdentity(ctx context.Context, metadata *models.TokenRequestMetadata) (string, error) {
	if p.conf.ClientIDPattern == "" {
		return metadata.CertificateHeader, nil
	}

	regex, err := regexp.Compile(p.conf.ClientIDPattern)
	if err != nil {
		return "", errors.Wrap(contract.ErrInvalidArguments, "identity header pattern is not valid regex")
	}

	identity := regex.FindStringSubmatch(metadata.CertificateHeader)
	if len(identity) == 0 {
		return "", errors.Wrap(contract.ErrInvalidArguments, "identity header pattern does not match")
	}

	// Regex match index 0 is the full match, index 1 is the first capture group.
	if len(identity) == 1 {
		return identity[0], nil
	}

	// Return first capture group for remove prefix support.
	return identity[1], nil
}
