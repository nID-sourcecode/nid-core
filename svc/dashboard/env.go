package main

import (
	"lab.weave.nl/nid/nid-core/pkg/environment"
)

// DashBoardConfig contains the configuration for the dashboard service
type DashBoardConfig struct {
	environment.BaseConfig
	JWTPath             string `envconfig:"JWT_PATH"`
	RegistrySecretPath  string `envconfig:"REGISTRY_SECRET_PATH"`
	AuthorizationURI    string `envconfig:"AUTHORIZATION_URI"`
	AutopseudoImage     string `envconfig:"AUTOPSEUDO_IMAGE"`
	ImagePullPolicy     string `envconfig:"IMAGE_PULL_POLICY"`
	AuthorizationIssuer string `envconfig:"AUTHORIZATION_ISSUER"`
	BaseDomain          string `envconfig:"BASE_DOMAIN"`
	DefaultUser         string `envconfig:"DEFAULT_USER"`
	DefaultUserPass     string `envconfig:"DEFAULT_USER_PASS"`
	PilotUser           string `envconfig:"PILOT_USER"`
	PilotUserPass       string `envconfig:"PILOT_USER_PASS"`
}
