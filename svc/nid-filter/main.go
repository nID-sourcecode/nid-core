package main

import (
	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3alpha"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/extproc"
	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/dial"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/servicebase"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	authpb "lab.weave.nl/nid/nid-core/svc/auth/proto"
	"lab.weave.nl/nid/nid-core/svc/autopseudo/keyutil"
	"lab.weave.nl/nid/nid-core/svc/nid-filter/filters/auditlog"
	"lab.weave.nl/nid/nid-core/svc/nid-filter/filters/authswap"
	"lab.weave.nl/nid/nid-core/svc/nid-filter/filters/autopseudo"
	"lab.weave.nl/nid/nid-core/svc/nid-filter/filters/scopeverification"
	walletPB "lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
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
		scopeverification.NewScopeVerificationFilterInitializer(),
	}

	if config.AutopseudoEnabled || config.AutobsnEnabled {
		if config.AutopseudoPriv == "" {
			log.Fatal("AUTOPSEUDO_PRIV is required if AUTOBSN_ENABLED or AUTOPSEUDO_ENABLED is true")
		}
		key, err := keyutil.ParseKeypair(config.AutopseudoPriv)
		if err != nil {
			log.Fatal(err)
		}

		if config.AutopseudoEnabled {
			autopseudoInitializer := autopseudo.NewFilterInitializer(&autopseudo.Config{
				Namespace:         config.Namespace,
				Key:               key,
				SubjectIdentifier: "$$nid:subject$$",
				FilterName:        "autopseudo",
			})

			filters = append(filters, autopseudoInitializer)
		}

		if config.AutobsnEnabled {
			if config.WalletURI == "" {
				log.Fatal("WALLET_URI is required if AUTOBSN_ENABLED is true")
			}
			connection, err := dial.Service(config.WalletURI, grpc.WithInsecure())
			if err != nil {
				log.WithError(err).WithField("uri", config.WalletURI).Fatal("connecting to wallet")
			}
			walletClient := walletPB.NewWalletClient(connection)
			autobsnInitializer := autopseudo.NewFilterInitializer(&autopseudo.Config{
				Namespace:         config.Namespace,
				Key:               key,
				TranslateToBSN:    true,
				WalletClient:      walletClient,
				SubjectIdentifier: "$$nid:bsn$$",
				FilterName:        "autobsn",
			})

			filters = append(filters, autobsnInitializer)
		}
	}

	if config.AuthswapEnabled {
		if config.AuthURI == "" {
			log.Fatal("AUTH_URI is required if AUTHSWAP_ENABLED is true")
		}
		connection, err := dial.Service(config.AuthURI, grpc.WithInsecure())
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

// NIDFilterServiceRegistry is an implementation of grpc service registry
type NIDFilterServiceRegistry struct {
	servicebase.Registry

	externalProcessorService ext_proc_pb.ExternalProcessorServer
}

// RegisterServices registers the external processor server
func (r NIDFilterServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	ext_proc_pb.RegisterExternalProcessorServer(grpcServer, r.externalProcessorService)
}
