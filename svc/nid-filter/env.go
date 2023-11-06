package main

import (
	"github.com/nID-sourcecode/nid-core/pkg/environment"
)

// NIDFilterConfig contains the configuration for the nid-filter service
type NIDFilterConfig struct {
	environment.BaseConfig
	AutopseudoEnabled bool   `envconfig:"default=true,AUTOPSEUDO_ENABLED"`
	AutopseudoPriv    string `envconfig:"AUTOPSEUDO_PRIV,optional"`
	AutobsnEnabled    bool   `envconfig:"default=true,AUTOBSN_ENABLED"`
	WalletURI         string `envconfig:"WALLET_URI,optional"`
	AuthURI           string `envconfig:"AUTH_URI,optional"`
	AuthswapEnabled   bool   `envconfig:"default=true,AUTHSWAP_ENABLED"`
}
