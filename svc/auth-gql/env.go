package main

import (
	"github.com/nID-sourcecode/nid-core/pkg/environment"
)

// AuthGQLConfig implements the used environment variables
type AuthGQLConfig struct {
	environment.BaseConfig
	GqlPlaygroundEnabled bool `envconfig:"GQL_PLAYGROUND_ENABLED"`
}
