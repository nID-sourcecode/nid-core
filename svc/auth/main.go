//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install github.com/goadesign/goa/goagen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install lab.weave.nl/weave/generator/cmd/gen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin bin/gen -graphqlPath=../auth-gql/graphql -authPath=../auth-gql/auth
package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/gqlutil"
	"lab.weave.nl/nid/nid-core/pkg/jwtconfig"
	"lab.weave.nl/nid/nid-core/pkg/pseudonym"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/metrics"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/servicebase"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v3"
	log "lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	httpCallbackHandler "lab.weave.nl/nid/nid-core/svc/auth/internal/callbackhandler/retryhttp"
	pb "lab.weave.nl/nid/nid-core/svc/auth/proto"
	walletPB "lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
)

func initialise() (*AuthServiceRegistry, *AuthConfig) {
	// Init conf
	var conf *AuthConfig
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to load environment config")
	}
	err := log.SetFormat(log.Format(conf.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	jwtKey, err := jwtconfig.Read(conf.JWTPath)
	if err != nil {
		log.WithError(err).Fatal("unable to load additional config from files")
	}

	log.WithField("level", conf.GetLogLevel()).Info("Setting log level")

	err = log.SetLevel(log.Level(conf.LogLevel))
	if err != nil {
		log.WithError(err).Fatal("unable to set log level")
	}

	db := initDB(conf)

	walletConnection, err := grpc.Dial(conf.WalletURI, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).WithField("url", conf.WalletURI).Fatal("unable to dial wallet service")
	}

	pseudonymizer := pseudonym.NewPseudonymizer(conf.PseudonymizationURI)

	jwtClientOpts := jwt.DefaultOpts()
	jwtClientOpts.HeaderOpts.KID = jwtKey.ID
	jwtClient := jwt.NewJWTClientWithOpts(jwtKey.PrivateKey, &jwtKey.PublicKey, jwtClientOpts)
	gqlUtil := gqlutil.NewSchemaFetcher(gqlutil.DefaultGraphQLClient)
	metadataHelper := &headers.GRPCMetadataHelper{}
	passwordManager := password.NewDefaultManager()

	// Init the prometheus scope
	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "auth")

	authServer := &AuthServiceServer{
		db:              db,
		stats:           CreateStats(scope),
		pseudonymizer:   pseudonymizer,
		jwtClient:       jwtClient,
		schemaFetcher:   gqlUtil,
		walletClient:    walletPB.NewWalletClient(walletConnection),
		conf:            conf,
		metadataHelper:  metadataHelper,
		passwordManager: passwordManager,
		callbackhandler: httpCallbackHandler.NewCallbackHandler(conf.CallbackMaxRetryAttempts),
	}

	wellKnownServer := &WellKnownServiceServer{
		conf:      conf,
		jwtClient: jwtClient,
	}
	err = wellKnownServer.initWellKnown(pb.File_auth_proto)
	if err != nil {
		log.WithError(err).Fatal("initialising wellknown")
	}

	registry := &AuthServiceRegistry{
		authClient:      authServer,
		wellKnownClient: wellKnownServer,
	}

	return registry, conf
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

// AuthServiceRegistry implementation of grpc service registry
type AuthServiceRegistry struct {
	servicebase.Registry

	authClient      *AuthServiceServer
	wellKnownClient *WellKnownServiceServer
}

// RegisterServices register auth service
func (a AuthServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	pb.RegisterAuthServer(grpcServer, a.authClient)
	pb.RegisterWellKnownServer(grpcServer, a.wellKnownClient)
}
