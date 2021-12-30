package main

import (
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/environment"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	pb "lab.weave.nl/nid/nid-core/svc/scopeverification/proto"
)

func main() {
	var conf *environment.BaseConfig
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to load config from environment")
	}

	err := log.SetFormat(log.Format(conf.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	registry := scopeVerificationServiceRegistry{
		scopeVerificationServer: NewScopeVerificationServer(),
	}

	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = conf.Port
	grpcConfig.LogLevel = conf.GetLogLevel()
	grpcConfig.LogFormatter = conf.GetLogFormatter()
	err = grpcserver.InitWithConf(registry, grpcConfig)
	if err != nil {
		log.WithError(err).Fatal("Error initialising grpc server")
	}
}

type scopeVerificationServiceRegistry struct {
	scopeVerificationServer *ScopeVerificationServer
}

// RegisterServices register dashboard server
func (a scopeVerificationServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	pb.RegisterScopeVerificationServer(grpcServer, a.scopeVerificationServer)
}
