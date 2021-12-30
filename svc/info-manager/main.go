//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install github.com/goadesign/goa/goagen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install lab.weave.nl/weave/generator/cmd/gen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin bin/gen -graphqlPath=../info-manager-gql/graphql
package main

import (
	"context"
	"time"

	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/servicebase"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage/s3"
	"lab.weave.nl/nid/nid-core/svc/info-manager/inforestarter"
	pb "lab.weave.nl/nid/nid-core/svc/info-manager/proto"
)

func initialise() (*InfoManagerRegistry, *InfoManagerConfig) {
	var conf *InfoManagerConfig
	err := envconfig.Init(&conf)
	if err != nil {
		log.WithError(err).Fatal("unable to load environment config")
	}

	db := initDB(conf, false)

	amountRetries := 15
	var storageClient objectstorage.Client

	for i := 0; i < amountRetries; i++ {
		storageClient, err = s3.NewClient(context.Background(), &conf.ObjectStorage.ClientConfig, conf.ObjectStorage.Bucket, nil)
		if err != nil {
			log.WithError(err).Warn("failed to create storage client, retrying in 1s")
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		log.WithError(err).Fatal("unable to create storage client")
	}

	infoRestarter, err := inforestarter.NewK8sRolloutInfoRestarter(conf.Namespace)
	if err != nil {
		log.WithError(err).Fatal("creating info restarter")
	}

	infoManagerServer := &InfoManagerServiceServer{
		db:            db,
		storageClient: storageClient,
		infoRestarter: infoRestarter,
	}

	registry := &InfoManagerRegistry{
		infoManagerClient: infoManagerServer,
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
		log.WithError(err).Fatal("error initialising grpc server")
	}
}

// InfoManagerRegistry implementation of grpc service registry
type InfoManagerRegistry struct {
	servicebase.Registry

	infoManagerClient *InfoManagerServiceServer
}

// RegisterServices register info-manager service
func (r InfoManagerRegistry) RegisterServices(grpcServer *grpc.Server) {
	pb.RegisterInfoManagerServer(grpcServer, r.infoManagerClient)
}
