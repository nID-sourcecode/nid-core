package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/environment"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/metrics"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/servicebase"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	pb "lab.weave.nl/nid/nid-core/svc/auditlog/proto"
)

func main() {
	conf := environment.BaseConfig{}
	err := envconfig.Init(&conf)
	if err != nil {
		log.WithError(err).Fatal("unable to load config")
	}

	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = conf.Port
	grpcConfig.LogLevel = conf.GetLogLevel()
	grpcConfig.LogFormatter = conf.GetLogFormatter()

	// Init the prometheus scope
	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "auditlog")

	auditLogServer := AuditLogServiceServer{
		logger: log.GetLogger(),
		stats:  CreateStats(scope),
	}
	err = auditLogServer.logger.SetLevel(log.Level(conf.GetLogLevel().String()))
	if err != nil {
		log.WithError(err).Fatal("can't set log level")
	}
	err = auditLogServer.logger.SetFormatter(conf.GetLogFormatter())
	if err != nil {
		log.WithError(err).Fatal("can't set log formatter")
	}
	registry := &AuditLogServiceRegistry{
		auditLogServer: auditLogServer,
	}

	err = grpcserver.InitWithConf(registry, grpcConfig)
	if err != nil {
		log.WithError(err).Fatal("unable to initialise grpc server")
	}
}

// AuditLogServiceRegistry implementation of grpc service registry
type AuditLogServiceRegistry struct {
	servicebase.Registry

	auditLogServer AuditLogServiceServer
}

// RegisterServices is used to register services
func (a AuditLogServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	pb.RegisterAuditlogServiceServer(grpcServer, &a.auditLogServer)
}
