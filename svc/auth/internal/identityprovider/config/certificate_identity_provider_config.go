package config

// CertificateIdentityProviderConfig contains the configuration for the CertificateIdentityProvider.
type CertificateIdentityProviderConfig struct {
	// ClientIDPattern is a regex pattern that should match the identity header.
	// The first capture group will be used as the identity.
	// If no capture groups are present, the full match will be used as the identity.
	ClientIDPattern string `envconfig:"CLIENT_ID_PATTERN,optional"`
}

// Validate validates the CertificateIdentityProviderConfig.
func (c *CertificateIdentityProviderConfig) Validate() error {
	return nil
}
