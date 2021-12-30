package main

import (
	messagebirdUtils "github.com/messagebird/go-rest-api"
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
	postmarkUtils "lab.weave.nl/nid/nid-core/pkg/utilities/postmark/v2"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/messagebird"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/postmark"
	"lab.weave.nl/nid/nid-core/svc/wallet-rpc/gqlclient"
	pb "lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
)

func initialise() (*WalletServiceRegistry, *WalletConfig) {
	var config *WalletConfig
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

	pwManager := password.NewDefaultManager()

	db := initDB(config, false, pwManager)
	opts := jwt.DefaultOpts()
	opts.HeaderOpts.KID = jwtKey.ID
	opts.ClaimsOpts.Issuer = "wallet"
	opts.ClaimsOpts.Audience = []string{"auth", "wallet-gql", "wallet-rpc"}
	jwtClient := jwt.NewJWTClientWithOpts(jwtKey.PrivateKey, &jwtKey.PublicKey, opts)

	// Init the prometheus scope
	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "wallet-rpc")
	stats := CreateStats(scope)

	authorizationServer := &AuthorizationServer{
		db:             db,
		stats:          stats,
		metadataHelper: &headers.GRPCMetadataHelper{},
		jwtClient:      jwtClient,
		pwManager:      pwManager,
	}

	walletServer := &WalletServer{
		db:         db,
		stats:      stats,
		authClient: gqlclient.NewAuthClient("http://auth-gql.nid.svc.cluster.local/gql"),
	}

	verifierServer := &VerifierServer{
		stats:         stats,
		db:            db,
		emailVerifier: &postmark.Postmark{Client: postmarkUtils.NewClient(config.Postmark.API, config.Postmark.Account)},
		phoneVerifier: &messagebird.Messagebird{Client: messagebirdUtils.New(config.Messagebird)},
	}

	registry := &WalletServiceRegistry{
		authorizationClient: authorizationServer,
		walletClient:        walletServer,
		verifierClient:      verifierServer,
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

// WalletServiceRegistry implementation of grpc service registry
type WalletServiceRegistry struct {
	servicebase.Registry

	authorizationClient *AuthorizationServer
	walletClient        *WalletServer
	verifierClient      *VerifierServer
}

// RegisterServices register dashboard server
func (a WalletServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	pb.RegisterAuthorizationServer(grpcServer, a.authorizationClient)
	pb.RegisterWalletServer(grpcServer, a.walletClient)
	pb.RegisterVerificationServer(grpcServer, a.verifierClient)
}
