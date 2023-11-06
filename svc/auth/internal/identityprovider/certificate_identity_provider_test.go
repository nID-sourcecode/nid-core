package identityprovider

import (
	"context"
	"testing"

	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/identityprovider/config"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	"github.com/stretchr/testify/assert"
)

func TestCertificateIdentityProviderParseCertificate(t *testing.T) {
	tests := []struct {
		Scenario string
		Conf     *config.CertificateIdentityProviderConfig
		Metadata *models.TokenRequestMetadata
		Expected string
		Error    error
	}{
		{
			Scenario: "pattern is empty",
			Conf: &config.CertificateIdentityProviderConfig{
				ClientIDPattern: "",
			},
			Metadata: &models.TokenRequestMetadata{
				CertificateHeader: "X=123 Y=hello Z=testvalue",
			},
			Expected: "X=123 Y=hello Z=testvalue",
			Error:    nil,
		},
		{
			Scenario: "pattern is not valid regex",
			Conf: &config.CertificateIdentityProviderConfig{
				ClientIDPattern: "[",
			},
			Metadata: &models.TokenRequestMetadata{
				CertificateHeader: "X=123 Y=hello Z=testvalue",
			},
			Expected: "",
			Error:    contract.ErrInvalidArguments,
		},
		{
			Scenario: "pattern does not match",
			Conf: &config.CertificateIdentityProviderConfig{
				ClientIDPattern: `Z=(\w+)`,
			},
			Metadata: &models.TokenRequestMetadata{
				CertificateHeader: "no-match",
			},
			Expected: "",
			Error:    contract.ErrInvalidArguments,
		},
		{
			Scenario: "pattern matches, full pattern",
			Conf: &config.CertificateIdentityProviderConfig{
				ClientIDPattern: `Z=\w+`,
			},
			Metadata: &models.TokenRequestMetadata{
				CertificateHeader: "X=123 Y=hello Z=testvalue",
			},
			Expected: "Z=testvalue",
			Error:    nil,
		},
		{
			Scenario: "pattern matches, first capture group",
			Conf: &config.CertificateIdentityProviderConfig{
				ClientIDPattern: `Z=(\w+)`,
			},
			Metadata: &models.TokenRequestMetadata{
				CertificateHeader: "X=123 Y=hello Z=testvalue",
			},
			Expected: "testvalue",
			Error:    nil,
		},
		{
			Scenario: "pattern matches multiple times, full pattern",
			Conf: &config.CertificateIdentityProviderConfig{
				ClientIDPattern: `Z=\w+`,
			},
			Metadata: &models.TokenRequestMetadata{
				CertificateHeader: "X=123 Z=testvalue Y=hello Z=anothervalue",
			},
			Expected: "Z=testvalue",
			Error:    nil,
		},

		{
			Scenario: "pattern matches multiple times, first capture group",
			Conf: &config.CertificateIdentityProviderConfig{
				ClientIDPattern: `Z=(\w+)`,
			},
			Metadata: &models.TokenRequestMetadata{
				CertificateHeader: "X=123 Z=testvalue Y=hello Z=anothervalue",
			},
			Expected: "testvalue",
			Error:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			provider := NewCertificateIdentityProvider(test.Conf)
			result, err := provider.GetIdentity(context.Background(), test.Metadata)

			assert.Equal(t, test.Expected, result)
			assert.ErrorIs(t, err, test.Error)
		})
	}
}
