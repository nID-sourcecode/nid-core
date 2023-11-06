// Package nid-filter
package main

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/nID-sourcecode/nid-core/pkg/extproc"
	"github.com/nID-sourcecode/nid-core/pkg/extproc/filter"
	"github.com/nID-sourcecode/nid-core/pkg/keyutil"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/dial"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/servicebase"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	authpb "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
	"github.com/nID-sourcecode/nid-core/svc/nid-filter/contract"
	"github.com/nID-sourcecode/nid-core/svc/nid-filter/filters/auditlog"
	"github.com/nID-sourcecode/nid-core/svc/nid-filter/filters/authswap"
	"github.com/nID-sourcecode/nid-core/svc/nid-filter/filters/autopseudo"
	"github.com/nID-sourcecode/nid-core/svc/nid-filter/filters/scopeverification"
	externalauthorization "github.com/nID-sourcecode/nid-core/svc/nid-filter/transport/external_authorization"
	walletPB "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initialise() (*NIDFilterServiceRegistry, *NIDFilterConfig) {
	var config *NIDFilterConfig
	if err := envconfig.Init(&config); err != nil {
		log.WithError(err).Fatal("unable to load config from environment")
	}

	err := log.SetFormat(log.Format(config.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	filters := []filter.Initializer{
		auditlog.NewFilterInitializer(log.GetLogger()),
	}
	authorizationRules := []contract.AuthorizationRule{
		scopeverification.New(),
	}

	if config.AutopseudoEnabled || config.AutobsnEnabled { //nolint:nestif
		if config.AutopseudoPriv == "" {
			log.Fatal("AUTOPSEUDO_PRIV is required if AUTOBSN_ENABLED or AUTOPSEUDO_ENABLED is true")
		}
		key, err := keyutil.ParseKeypair(config.AutopseudoPriv)
		if err != nil {
			log.Fatal(err)
		}

		if config.AutopseudoEnabled {
			autopseudoInitializer := autopseudo.New(&autopseudo.Config{
				Namespace:         config.Namespace,
				Key:               key,
				SubjectIdentifier: "$$nid:subject$$",
				FilterName:        "autopseudo",
			})

			authorizationRules = append(authorizationRules, autopseudoInitializer)
		}

		if config.AutobsnEnabled {
			if config.WalletURI == "" {
				log.Fatal("WALLET_URI is required if AUTOBSN_ENABLED is true")
			}
			connection, err := dial.Service(config.WalletURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.WithError(err).WithField("uri", config.WalletURI).Fatal("connecting to wallet")
			}
			walletClient := walletPB.NewWalletClient(connection)
			autoBSNInitializer := autopseudo.New(&autopseudo.Config{
				Namespace:         config.Namespace,
				Key:               key,
				TranslateToBSN:    true,
				WalletClient:      walletClient,
				SubjectIdentifier: "$$nid:bsn$$",
				FilterName:        "autobsn",
			})

			authorizationRules = append(authorizationRules, autoBSNInitializer)
		}
	}

	if config.AuthswapEnabled {
		if config.AuthURI == "" {
			log.Fatal("AUTH_URI is required if AUTHSWAP_ENABLED is true")
		}
		connection, err := dial.Service(config.AuthURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.WithError(err).WithField("uri", config.AuthURI).Fatal("connecting to auth")
		}
		authClient := authpb.NewAuthClient(connection)
		authswapInitializer := authswap.NewFilterInitializer(authClient)
		filters = append(filters, authswapInitializer)
	}

	externalProcessorServer := extproc.NewExternalProcessorServer(filters)
	registry := &NIDFilterServiceRegistry{
		externalProcessorService: externalProcessorServer,
		authorizationRules:       authorizationRules,
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

// NIDFilterServiceRegistry is an implementation of grpc service registry
type NIDFilterServiceRegistry struct {
	servicebase.Registry

	externalProcessorService ext_proc_pb.ExternalProcessorServer
	authorizationRules       []contract.AuthorizationRule
}

// RegisterServices registers the external processor server
func (r *NIDFilterServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	authv3.RegisterAuthorizationServer(grpcServer, externalauthorization.New(r.authorizationRules))

	ext_proc_pb.RegisterExternalProcessorServer(grpcServer, r.externalProcessorService)
}
