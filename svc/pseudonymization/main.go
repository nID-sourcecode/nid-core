// Package pseudonymization
package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/headers"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/metrics"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/servicebase"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/pseudonymization/keymanager"
	pseudoPB "github.com/nID-sourcecode/nid-core/svc/pseudonymization/proto"
)

const (
	durationDay = 24 * time.Hour
)

func main() {
	var config PseudonymizationConfig
	if err := envconfig.Init(&config); err != nil {
		log.WithError(err).Fatal("unable to read configuration from environment")
	}

	err := log.SetFormat(log.Format(config.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	log.WithField("cacheduration", config.CacheDuration).WithField("jwk_url", config.JWKURL).Info("Creating keymanager")
	keyManager := keymanager.NewKeyManager(config.JWKURL, config.CacheDuration, &keymanager.JWKSFetcher{})
	defer keyManager.Cleanup()

	// Init the prometheus scope
	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "pseudonymization")

	registry := &PseudonymServiceRegistry{
		pseudonymClient: &PseudonymizerServer{
			stats:          CreateStats(scope),
			KeyManager:     keyManager,
			metadataHelper: new(headers.GRPCMetadataHelper),
		},
	}

	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = config.Port
	grpcConfig.LogLevel = config.GetLogLevel()
	grpcConfig.LogFormatter = config.GetLogFormatter()
	err = grpcserver.InitWithConf(registry, &grpcConfig)
	if err != nil {
		log.WithError(err).Fatal("Error initialising grpc server")
	}
}

// PseudonymServiceRegistry implementation of grpc service registry
type PseudonymServiceRegistry struct {
	servicebase.Registry

	pseudonymClient *PseudonymizerServer
}

// RegisterServices register pseudonymizer server
func (a PseudonymServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	pseudoPB.RegisterPseudonymizerServer(grpcServer, a.pseudonymClient)
}
