package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/jwtconfig"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/metrics"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/servicebase"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	pb "lab.weave.nl/nid/nid-core/svc/dashboard/proto"
)

func initialise() (*DashboardServiceRegistry, *DashBoardConfig) {
	var config *DashBoardConfig
	if err := envconfig.Init(&config); err != nil {
		log.WithError(err).Fatal("unable to load config from environment")
	}

	err := log.SetFormat(log.Format(config.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	jwtKey, err := jwtconfig.Read(config.JWTPath)
	if err != nil {
		log.WithError(err).Fatal("unable to load additional config from files")
	}

	db := initDB(config)

	dashboardServer, err := NewDashboardServiceServer(config)
	if err != nil {
		log.WithError(err).Fatal("unable to create dashboard service server")
	}

	opts := jwt.DefaultOpts()
	opts.HeaderOpts.KID = jwtKey.ID
	opts.ClaimsOpts.Issuer = "dashboard"
	opts.ClaimsOpts.Audience = []string{"dashboard"}
	jwtClient := jwt.NewJWTClientWithOpts(jwtKey.PrivateKey, &jwtKey.PublicKey, opts)

	// Init the prometheus scope
	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "dashboard")

	authorizationServer := &AuthorizationServiceServer{
		db:             db,
		stats:          CreateStats(scope),
		metadataHelper: &headers.GRPCMetadataHelper{},
		jwtClient:      jwtClient,
		pwManager:      password.NewDefaultManager(),
	}

	registry := &DashboardServiceRegistry{
		dashboardClient:     dashboardServer,
		authorizationClient: authorizationServer,
	}

	return registry, config
}

func main() {
	registry, conf := initialise()

	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = conf.Port
	grpcConfig.LogLevel = conf.GetLogLevel()
	grpcConfig.LogFormatter = conf.GetLogFormatter()
	err := grpcserver.InitWithConf(registry, grpcConfig)
	if err != nil {
		log.WithError(err).Fatal("Error initialising grpc server")
	}
}

// DashboardServiceRegistry implementation of grpc service registry
type DashboardServiceRegistry struct {
	servicebase.Registry

	dashboardClient     *DashboardServiceServer
	authorizationClient *AuthorizationServiceServer
}

// RegisterServices register dashboard server
func (a DashboardServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	pb.RegisterDashboardServer(grpcServer, a.dashboardClient)
	pb.RegisterAuthorizationServiceServer(grpcServer, a.authorizationClient)
}
