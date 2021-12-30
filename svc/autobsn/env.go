package main

import (
	"lab.weave.nl/nid/nid-core/pkg/environment"
)

// AutoBSNConfig implements the used environment variables
type AutoBSNConfig struct {
	environment.BaseConfig
	Namespace string `envconfig:"NAMESPACE"`
	RSAPriv   string `envconfig:"RSA_PRIV"` // PEM encoded RSA private key used for encrypting pseudonyms
	WalletURI string `envconfig:"WALLET_URI"`
}
