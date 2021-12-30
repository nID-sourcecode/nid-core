package main

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vrischmann/envconfig"
	"github.com/xanzy/go-gitlab"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/metrics"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage/s3"
	"lab.weave.nl/nid/nid-core/svc/documentation/packages/git"
	documentationPB "lab.weave.nl/nid/nid-core/svc/documentation/proto"
)

func initialise() *DocumentationServiceRegistry {
	var config documentationConfig
	if err := envconfig.Init(&config); err != nil {
		log.WithError(err).Fatal("unable to load environment config")
	}

	err := log.SetFormat(log.Format(config.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	// Wait for istio proxy to start
	amountRetries := 15
	var storageClient objectstorage.Client
	for i := 0; i < amountRetries; i++ {
		storageClient, err = s3.NewClient(context.Background(), &config.ObjectStorage.ClientConfig, config.ObjectStorage.Bucket, nil)
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

	gitlabClient, err := gitlab.NewClient(config.GitlabAccessToken, gitlab.WithBaseURL(config.GitlabBaseURL))
	if err != nil {
		log.WithError(err).Fatal("failed to create git client")
	}

	// Init the prometheus scope
	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "documentation")

	registry := &DocumentationServiceRegistry{
		documentationClient: &DocumentationServiceServer{
			stats:         CreateStats(scope),
			conf:          &config,
			git:           git.NewGitClient(gitlabClient),
			storageClient: storageClient,
		},
	}

	return registry
}

func main() {
	registry := initialise()
	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.LogLevel = registry.documentationClient.conf.GetLogLevel()
	grpcConfig.LogFormatter = registry.documentationClient.conf.GetLogFormatter()
	grpcConfig.Port = registry.documentationClient.conf.Port

	if err := grpcserver.InitWithConf(registry, grpcConfig); err != nil {
		log.Fatalf("Service registry init default failed: %s", err)
	}
}

// DocumentationServiceRegistry implementation of documentation service registry
type DocumentationServiceRegistry struct {
	grpcserver.ServiceRegistry

	documentationClient *DocumentationServiceServer
}

// RegisterServices by default no services are registered
func (d DocumentationServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	documentationPB.RegisterDocumentationServer(grpcServer, d.documentationClient)
}
