// Package internal has configuration for external-authorization-chain
package internal

// AppConfig configuration for App package
type AppConfig struct {
	DenyByDefault *DenyByDefault `yaml:"deny_by_default"`

	// Endpoints
	//
	// key: Server to send the request to
	//
	// Value: Endpoints endpoints of services which will be called by authz client (key).
	Endpoints map[string][]string `yaml:"endpoints"`
}

// DenyByDefault configuration for App package to configure the DenyByDefault feature.
type DenyByDefault struct {
	Enabled bool     `yaml:"enabled" default:"true"`
	Allow   []string `yaml:"allow"`
}

// GRPCConfig contains the configuration for the external authorization handler
type GRPCConfig struct {
	Port int ``
}

// Config configuration
type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}
