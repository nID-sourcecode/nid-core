package main

import (
	"time"

	"lab.weave.nl/nid/nid-core/pkg/environment"
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
	AuthorizationCodeLength  int    `envconfig:"AUTHORIZATION_CODE_LENGTH"`
	WalletURI                string `envconfig:"WALLET_URI"`
	TestingClientPassword    string `envconfig:"TESTING_CLIENT_PASSWORD"`
	PilotClientPassword      string `envconfig:"PILOT_CLIENT_PASSWORD"`
	CallbackMaxRetryAttempts int    `envconfig:"default=10,CALLBACK_MAX_RETRY_ATTEMPTS"`
}
