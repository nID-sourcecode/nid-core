package main

import "github.com/nID-sourcecode/nid-core/pkg/environment"

// OnboardingConfig contains the configuration for the onboarding service
type OnboardingConfig struct {
	environment.BaseConfig
	WalletURL           string `envconfig:"WALLET_URL"`
	PseudonymizationURL string `envconfig:"PSEUDONYMIZATION_URL"`
}
