//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install github.com/goadesign/goa/goagen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install lab.weave.nl/weave/generator/cmd/gen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin bin/gen -graphqlPath=../auth-gql/graphql -authPath=../auth-gql/auth
package main

import (
	"github.com/nID-sourcecode/nid-core/pkg/gqlutil"
	"github.com/nID-sourcecode/nid-core/pkg/interceptor/xrequestid"
	"github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/pseudonym"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/headers"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/metrics"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/auth/app"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/audienceprovider"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/config"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/identityprovider"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/repository"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/retryhttp"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/stats"
	"github.com/nID-sourcecode/nid-core/svc/auth/transport/http"
	walletPB "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vrischmann/envconfig"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nID-sourcecode/nid-core/pkg/jwtconfig"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/servicebase"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/jwt/v3"
	grpcTransport "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc"
	pb "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
)

func initialiseRegistry(conf *config.AuthConfig, jwtClient *jwt.Client, app contract.App) *AuthServiceRegistry {
	err := log.SetFormat(log.Format(conf.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	log.WithField("level", conf.GetLogLevel()).Info("Setting log level")

	err = log.SetLevel(log.Level(conf.LogLevel))
	if err != nil {
		log.WithError(err).Fatal("unable to set log level")
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
		authClient:      grpcTransport.New(app, &headers.GRPCMetadataHelper{}, &conf.Transport.Grpc),
		wellKnownClient: wellKnownServer,
	}

	return registry
}

func getJWTClient(conf *config.AuthConfig) *jwt.Client {
	jwtKey, err := jwtconfig.Read(conf.JWTPath)
	if err != nil {
		log.WithError(err).Fatal("unable to load additional config from files")
	}

	jwtClientOpts := jwt.DefaultOpts()
	jwtClientOpts.HeaderOpts.KID = jwtKey.ID
	jwtClientOpts.MarshalSingleStringAsArray = conf.MarshalSingleAudienceOrScopeAsArray
	jwtClient := jwt.NewJWTClientWithOpts(jwtKey.PrivateKey, &jwtKey.PublicKey, jwtClientOpts)
	return jwtClient
}

func main() {
	var conf *config.AuthConfig
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to load environment config")
	}

	err := conf.Validate()
	if err != nil {
		log.WithError(err).Fatal("invalid config")
	}

	jwtClient := getJWTClient(conf)

	authApp := getAuthApp(conf, jwtClient)

	registry := initialiseRegistry(conf, jwtClient, authApp)
	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = conf.Transport.Grpc.Port
	grpcConfig.LogLevel = conf.GetLogLevel()
	grpcConfig.LogFormatter = conf.GetLogFormatter()
	grpcConfig.AdditionalInterceptors = []grpc.UnaryServerInterceptor{
		xrequestid.AddXRequestID,
		otelgrpc.UnaryServerInterceptor(),
	}
	go func() {
		err := grpcserver.InitWithConf(registry, &grpcConfig)
		if err != nil {
			log.WithError(err).Fatal("initialising grpc server")
		}
	}()

	httpServer := http.New(authApp, &conf.Transport.Http)
	err = httpServer.Run(conf.Transport.Http.Port)
	if err != nil {
		log.WithError(err).Fatal("running http server")
	}
}

func getAuthApp(conf *config.AuthConfig, jwtClient *jwt.Client) *app.App {
	db := repository.InitDB(conf)
	passwordManager := password.NewDefaultManager()

	walletConnection, err := grpc.Dial(conf.WalletURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithError(err).WithField("url", conf.WalletURI).Fatal("unable to dial wallet service")
	}

	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "auth")
	audienceProvider, err := getAudienceProvider(conf, db)
	if err != nil {
		log.WithError(err).Fatal("unable to get audience provider")
	}

	identityProvider, err := getIdentityProvider(conf, db, passwordManager)
	if err != nil {
		log.WithError(err).Fatal("unable to get identity provider")
	}

	authApp := app.New(conf,
		db,
		gqlutil.NewSchemaFetcher(gqlutil.DefaultGraphQLClient),
		stats.CreateStats(scope),
		retryhttp.NewCallbackHandler(conf.CallbackMaxRetryAttempts),
		passwordManager,
		jwtClient,
		pseudonym.NewPseudonymizer(conf.PseudonymizationURI),
		walletPB.NewWalletClient(walletConnection),
		audienceProvider,
		identityProvider,
	)
	return authApp
}

// AuthServiceRegistry implementation of grpc service registry
type AuthServiceRegistry struct {
	servicebase.Registry

	authClient      *grpcTransport.Server
	wellKnownClient *WellKnownServiceServer
}

// RegisterServices register auth service
func (a AuthServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	pb.RegisterAuthServer(grpcServer, a.authClient)
	pb.RegisterWellKnownServer(grpcServer, a.wellKnownClient)
}

func getAudienceProvider(conf *config.AuthConfig, db *repository.AuthDB) (contract.AudienceProvider, error) {
	switch conf.AudienceProvider {
	case contract.AudienceProviderTypeRequest:
		return &audienceprovider.RequestAudienceProvider{}, nil
	case contract.AudienceProviderTypeDatabase:
		return audienceprovider.NewDatabaseAudienceProvider(conf), nil
	default:
		return nil, contract.ErrInvalidAudienceProvider
	}
}

func getIdentityProvider(conf *config.AuthConfig, db *repository.AuthDB, passwordManager password.IManager) (contract.IdentityProvider, error) {
	switch conf.IdentityProvider {
	case contract.IdentityProviderTypeCertificate:
		return &identityprovider.CertificateIdentityProvider{}, nil
	case contract.IdentityProviderTypeDatabase:
		return identityprovider.NewDatabaseIdentityProvider(db.ClientDB, passwordManager), nil
	default:
		return nil, contract.ErrInvalidIdentityProvider
	}
}
