// //go:generate env GO111MODULE=on GOBIN=$PWD/bin go install github.com/goadesign/goa/goagen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install lab.weave.nl/weave/generator/cmd/gen
// //go:generate env GO111MODULE=on GOBIN=$PWD/bin bin/goagen -d lab.weave.nl/nid/nid-core/services/gqlciz/design gen --pkg-path=lab.weave.nl/weave/generator
//go:generate env GO111MODULE=on GOBIN=$PWD/bin bin/gen lab.weave.nl/nid/nid-core/services/luarunner
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/environment"
	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/dial"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	authProto "lab.weave.nl/nid/nid-core/svc/auth/proto"
	"lab.weave.nl/nid/nid-core/svc/luarunner/models"
)

func main() {
	var config environment.BaseConfig
	if err := envconfig.Init(&config); err != nil {
		log.WithError(err).Fatal("unable to load configuration from environment")
	}

	err := log.SetFormat(log.Format(config.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	authDB := database.MustConnectCustomWithCustomLogger(&database.DBConfig{
		Host:           config.PGHost,
		User:           config.PGUser,
		Pass:           config.PGPass,
		Port:           config.PGPort,
		RetryOnFailure: true,
		TestMode:       database.TestModeOff,
		DBName:         "auth",
		LogMode:        false,
		AutoMigrate:    false,
		Extensions:     []string{"uuid-ossp"},
	}, models.GetModels(), log.GetLogger())

	conn, err := getGRPCClient("auth")
	if err != nil {
		log.WithError(err).Fatal("connecting with auth service")
	}

	db := database.MustConnectCustom(&database.DBConfig{
		Host:           config.PGHost,
		User:           config.PGUser,
		Pass:           config.PGPass,
		Port:           config.PGPort,
		RetryOnFailure: true,
		TestMode:       database.TestModeOff,
		TimeOut:        int(time.Minute.Seconds()),
		DBName:         "luarunner",
		LogMode:        false,
		AutoMigrate:    true,
		Extensions:     []string{"uuid-ossp", "pg_trgm", "postgis"},
	}, models.GetModels())

	models.AddForeignKeys(db)

	luaRunnerDB := NewLuaRunnerDB(db)
	if err != nil {
		log.WithError(err).Fatal("setting up luarunner database")
	}

	authClient := authProto.NewAuthClient(conn)
	luaRunnerService := NewLuaRunnerService(authDB, authClient, luaRunnerDB)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		luaRunnerService.HTTPCallback(w, r)
	})

	http.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			log.Errorf("unable to return v1/health value, error: %s", err.Error())
		}
	})

	log.Infof("running on http://localhost:%d", config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}

func getGRPCClient(service string) (*grpc.ClientConn, error) {
	var connection *grpc.ClientConn
	connection, err := dial.Service(fmt.Sprintf("%s:80", service), grpc.WithInsecure())

	return connection, err
}
