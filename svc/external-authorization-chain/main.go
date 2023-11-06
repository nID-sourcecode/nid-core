// Package main
package main

import (
	"os"

	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/internal"

	"gopkg.in/yaml.v3"

	"github.com/kelseyhightower/envconfig"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/app"
	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/transport/grpc"
)

func main() {
	var config internal.Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.WithError(err).Fatal("could not read environment for config")
	}

	err = log.SetLevel(log.Level(config.LogLevel))
	log.Infof("external-authorization-chain log level set to: %s", config.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("setting log level")
	}

	var grpcConfig internal.GRPCConfig
	err = envconfig.Process("", &grpcConfig)
	if err != nil {
		log.WithError(err).Fatal("tried filling the grpc config from the environment variables")
	}

	appConfig, err := getAppConfiguration()
	if err != nil {
		log.WithError(err).Fatal("tried reading the app configuration")
	}

	a, err := app.New(appConfig)
	if err != nil {
		log.WithError(err).Fatal("could not initialise app")
	}

	grpc.New(grpcConfig, a)
}

func getAppConfiguration() (internal.AppConfig, error) {
	endpointsFile, err := os.ReadFile("config/endpoints.yaml")
	if err != nil {
		log.WithError(err).Fatal("tried reading the endpoints file")
	}

	var appConfig internal.AppConfig
	err = yaml.Unmarshal(endpointsFile, &appConfig)
	if err != nil {
		log.WithError(err).Fatal("tried unmarshalling the endpoints yaml file for app config")
	}

	log.Debugf("unmarshalled yaml configmap: %+v \n DenyByDefault: %+v", appConfig, appConfig.DenyByDefault)

	if len(appConfig.Endpoints) == 0 {
		log.Fatal("No check endpoints given!")
	}
	return appConfig, err
}
