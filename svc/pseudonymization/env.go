package main

import (
	"time"

	"lab.weave.nl/nid/nid-core/pkg/environment"
)

// PseudonymizationConfig contains the config for the pseudonymization server
type PseudonymizationConfig struct {
	environment.BaseConfig
	JWKURL        string        `envconfig:"JWKURL"`
	CacheDuration time.Duration `envconfig:"CACHE_DURATION"`
}
