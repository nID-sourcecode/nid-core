package main

import (
	"github.com/nID-sourcecode/nid-core/pkg/environment"
)

// AutoPseudoConfig implements the used environment variables
type AutoPseudoConfig struct {
	environment.BaseConfig
	Namespace string `envconfig:"NAMESPACE"`
	RSAPriv   string `envconfig:"RSA_PRIV"` // PEM encoded RSA private key used for encrypting pseudonyms
}
