// //go:generate env GO111MODULE=on GOBIN=$PWD/bin go install github.com/goadesign/goa/goagen
// //go:generate env GO111MODULE=on GOBIN=$PWD/bin bin/goagen -d github.com/nID-sourcecode/nid-core/services/gqlciz/design gen --pkg-path=lab.weave.nl/weave/generator
//
//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install lab.weave.nl/weave/generator/cmd/gen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin bin/gen github.com/nID-sourcecode/nid-core/services/luarunner
package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/nID-sourcecode/nid-core/svc/luarunner/internal"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"github.com/nID-sourcecode/nid-core/pkg/environment"
	"github.com/nID-sourcecode/nid-core/pkg/httpserver"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/dial"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	authProto "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
	"github.com/nID-sourcecode/nid-core/svc/luarunner/models"
)

type contextKeyType string

const ctxB3Key = contextKeyType("B3Trace")

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

	luaRunnerDB := internal.NewLuaRunnerDB(db)
	if err != nil {
		log.WithError(err).Fatal("setting up luarunner database")
	}

	authClient := authProto.NewAuthClient(conn)
	luaRunnerService := internal.NewLuaRunnerService(authDB, &authClient, luaRunnerDB)

	server := httpserver.NewGinServer()
	server.Use(func(c *gin.Context) {
		headers := c.Request.Header

		ctx := context.WithValue(c.Request.Context(), ctxB3Key, map[string]string{
			"x-b3-traceid":      headers["X-B3-Traceid"][0],
			"x-b3-spanid":       headers["X-B3-Spanid"][0],
			"x-b3-parentspanid": headers["X-B3-Parentspanid"][0],
			"x-request-id":      headers["X-Request-Id"][0],
		})

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	server.POST("/callback", func(c *gin.Context) {
		luaRunnerService.HTTPCallback(c.Writer, c.Request)
	})

	server.POST("/:organisation/:script/run", luaRunnerService.Run)

	err = server.Run(":" + strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}
}

func getGRPCClient(service string) (*grpc.ClientConn, error) {
	return dial.Service(fmt.Sprintf("%s:8080", service), grpc.WithTransportCredentials(
		insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				b3Trace := ctx.Value(ctxB3Key).(map[string]string)

				for k := range b3Trace {
					ctx = metadata.AppendToOutgoingContext(ctx, k, b3Trace[k])
				}

				return invoker(ctx, method, req, reply, cc, opts...)
			},
		),
	)
}
