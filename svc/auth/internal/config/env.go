package config

import (
	"time"

	"github.com/nID-sourcecode/nid-core/pkg/environment"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/identityprovider/config"
	"github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc"
	"github.com/nID-sourcecode/nid-core/svc/auth/transport/http"
)

// AuthConfig implements the used environment variables
type AuthConfig struct {
	environment.BaseConfig
	ClusterHost         string `envconfig:"CLUSTER_HOST"`
	JWTPath             string `envconfig:"JWT_PATH"`
	JWKSURI             string `envconfig:"JWKS_URI"`
	Issuer              string `envconfig:"ISSUER"`
	AuthRequestURI      string `envconfig:"AUTH_REQUEST_URI"`
	PseudonymizationURI string `envconfig:"PSEUDONYMIZATION_URI"`
	// The OAuth 2.0 spec recommends a maximum lifetime of 10 minutes,
	// but in practice, most services set the expiration much shorter,
	// around 30-60 seconds
	AuthorizationCodeExpirationTime time.Duration `envconfig:"AUTHORIZATION_CODE_EXPIRATION_TIME"`
	// The authorization code itself can be of any length,
	// but the length of the codes should be documented.
	AuthorizationCodeLength             int                                      `envconfig:"AUTHORIZATION_CODE_LENGTH"`
	JWTExpirationHours                  int                                      `envconfig:"default=24,JWT_EXPIRE_HOURS"`
	JWTRefreshExpirationHours           int                                      `envconfig:"default=168,JWT_REFRESH_EXPIRE_HOURS"`
	WalletURI                           string                                   `envconfig:"WALLET_URI"`
	TestingClientPassword               string                                   `envconfig:"TESTING_CLIENT_PASSWORD"`
	PilotClientPassword                 string                                   `envconfig:"PILOT_CLIENT_PASSWORD"`
	CallbackMaxRetryAttempts            int                                      `envconfig:"default=10,CALLBACK_MAX_RETRY_ATTEMPTS"`
	AllowMultipleAudiences              bool                                     `envconfig:"default=false,ALLOW_MULTIPLE_AUDIENCES"`
	AudienceProvider                    contract.AudienceProviderType            `envconfig:"default=database,AUDIENCE_PROVIDER"`
	MarshalSingleAudienceOrScopeAsArray bool                                     `envconfig:"default=false,MARSHAL_SINGLE_AUDIENCE_OR_SCOPE_AS_ARRAY"`
	IdentityProvider                    contract.IdentityProviderType            `envconfig:"default=database,IDENTITY_PROVIDER"`
	CertificateIdentityProvider         config.CertificateIdentityProviderConfig `envconfig:"CERTIFICATE_IDENTITY_PROVIDER"`
	Transport                           TransportConfig                          `envconfig:"TRANSPORT"`
}

type TransportConfig struct {
	Http http.HttpConfig
	Grpc grpc.GrpcConfig
}

func (c *AuthConfig) Validate() error {
	err := c.AudienceProvider.Validate()
	if err != nil {
		return err
	}

	err = c.IdentityProvider.Validate()
	if err != nil {
		return err
	}

	err = c.CertificateIdentityProvider.Validate()
	if err != nil {
		return err
	}

	return nil
}
