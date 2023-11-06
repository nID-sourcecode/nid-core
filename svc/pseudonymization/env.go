package main

import (
	"time"

	"github.com/nID-sourcecode/nid-core/pkg/environment"
)

// PseudonymizationConfig contains the config for the pseudonymization server
type PseudonymizationConfig struct {
	environment.BaseConfig
	JWKURL        string        `envconfig:"JWKURL"`
	CacheDuration time.Duration `envconfig:"CACHE_DURATION"`
}
