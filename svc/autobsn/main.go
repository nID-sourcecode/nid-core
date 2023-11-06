package main

import (
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nID-sourcecode/nid-core/pkg/keyutil"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/dial"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/servicebase"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/autobsn/proto"
	walletPB "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
)

const (
	bearerScheme      = "Bearer "
	subjectIdentifier = "$$nid:bsn$$"
)

func main() {
	var conf AutoBSNConfig
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to load config from environment")
	}
	err := log.SetFormat(log.Format(conf.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}
	log.WithField("level", conf.GetLogLevel()).Info("Setting log level")

	err = log.SetLevel(log.Level(conf.LogLevel))
	if err != nil {
		log.WithError(err).Fatal("unable to set log level")
	}

	key, err := keyutil.ParseKeypair(conf.RSAPriv)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := dial.Service(conf.WalletURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithError(err).WithField("uri", conf.WalletURI).Fatal("connecting to wallet")
	}
	walletClient := walletPB.NewWalletClient(connection)

	registry := &AutoBSNServiceRegistry{
		autoBSNServer: NewAutoBSNServer(key, walletClient),
	}

	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = conf.Port
	grpcConfig.LogLevel = conf.GetLogLevel()
	grpcConfig.LogFormatter = conf.GetLogFormatter()
	err = grpcserver.InitWithConf(registry, &grpcConfig)
	if err != nil {
		log.WithError(err).Fatal("Error initialising grpc server")
	}
}

// AutoBSNServiceRegistry implementation of grpc service registry
type AutoBSNServiceRegistry struct {
	servicebase.Registry

	autoBSNServer *AutoBSNServer
}

// RegisterServices register autobsn server
func (a AutoBSNServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	proto.RegisterAutoBSNServer(grpcServer, a.autoBSNServer)
}
