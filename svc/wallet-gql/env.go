package main

import "github.com/nID-sourcecode/nid-core/pkg/environment"

// WalletConfig contains the configuration of the wallet service
type WalletConfig struct {
	environment.BaseConfig
	ClientID  string `envconfig:"CLIENT_ID"`
	UserScope string `envconfig:"USER_SCOPE"`
	Audience  string `envconfig:"AUDIENCE"`
	Postmark  struct {
		Account string `envconfig:"POSTMARK_ACCOUNT"`
		API     string `envconfig:"POSTMARK_API"`
	}
	Messagebird         string `envconfig:"MESSAGEBIRD"`
	PseudonymizationURL string `envconfig:"PSEUDONYMIZATION_URL"`
	AuthorizationURI    string `envconfig:"AUTHORIZATION_URI"`
}
