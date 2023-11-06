package main

import (
	"github.com/nID-sourcecode/nid-core/pkg/environment"
)

// WalletConfig contains the configuration for the dashboard service
type WalletConfig struct {
	environment.BaseConfig
	JWTPath      string `envconfig:"JWT_PATH"`
	DefaultUsers string `envconfig:"DEFAULT_USERS"`
	AuthURI      string `envconfig:"AUTH_URI"`
	Postmark     struct {
		Account string `envconfig:"POSTMARK_ACCOUNT"`
		API     string `envconfig:"POSTMARK_API"`
	}
	Messagebird string `envconfig:"MESSAGEBIRD"`
}
