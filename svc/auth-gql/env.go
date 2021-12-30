package main

import (
	"lab.weave.nl/nid/nid-core/pkg/environment"
)

// AuthGQLConfig implements the used environment variables
type AuthGQLConfig struct {
	environment.BaseConfig
	GqlPlaygroundEnabled bool `envconfig:"GQL_PLAYGROUND_ENABLED"`
}
