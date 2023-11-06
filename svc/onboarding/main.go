// Package onboarding
package main

import (
	"github.com/nID-sourcecode/nid-core/pkg/gqlclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/headers"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/metrics"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/servicebase"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	onboardingPB "github.com/nID-sourcecode/nid-core/svc/onboarding/proto"
	pseudoPB "github.com/nID-sourcecode/nid-core/svc/pseudonymization/proto"
)

func initialise() (*OnboardingServiceRegistry, *OnboardingConfig) {
	var config *OnboardingConfig
	if err := envconfig.Init(&config); err != nil {
		log.WithError(err).Fatal("unable to load config from environment")
	}

	err := log.SetFormat(log.Format(config.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	connection, err := grpc.Dial(config.PseudonymizationURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithError(err).WithField("url", config.PseudonymizationURL).Fatal("unable to dial pseudonymization service")
	}

	// Init the prometheus scopes
	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "onboarding")

	datasourceServer := &DataSourceServiceServer{
		stats:                  CreateStats(scope),
		walletClient:           gqlclient.NewClient(config.WalletURL),
		pseudonimizationClient: pseudoPB.NewPseudonymizerClient(connection),
		metadataHelper:         new(headers.GRPCMetadataHelper),
	}

	registry := &OnboardingServiceRegistry{
		datasourceClient: datasourceServer,
	}

	return registry, config
}

func main() {
	registry, conf := initialise()

	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = conf.Port
	grpcConfig.LogLevel = conf.GetLogLevel()
	grpcConfig.LogFormatter = conf.GetLogFormatter()
	err := grpcserver.InitWithConf(registry, &grpcConfig)
	if err != nil {
		log.WithError(err).Fatal("Error initialising grpc server")
	}
}

// OnboardingServiceRegistry implementation of grpc service registry
type OnboardingServiceRegistry struct {
	servicebase.Registry

	datasourceClient *DataSourceServiceServer
}

// RegisterServices register dashboard server
func (a OnboardingServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	onboardingPB.RegisterDataSourceServiceServer(grpcServer, a.datasourceClient)
}
