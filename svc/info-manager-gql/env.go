package main

import (
	"lab.weave.nl/nid/nid-core/pkg/environment"
)

// InfoManagerGQLConfig implements the used environment variables
type InfoManagerGQLConfig struct {
	environment.BaseConfig
	GqlPlaygroundEnabled bool   `envconfig:"GQL_PLAYGROUND_ENABLED"`
	InfoManagerURI       string `envconfig:"INFO_MANAGER_URI"`
}
