package main

import (
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/httpserver"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/info-manager-gql/graphql"
	"lab.weave.nl/nid/nid-core/svc/info-manager/models"
	"lab.weave.nl/nid/nid-core/svc/info-manager/proto"
)

func initialise() (resolver graphql.Resolver, conf InfoManagerGQLConfig, client proto.InfoManagerClient) {
	log.Info("Info-manager-gql initialising")

	// Init conf
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to load environment config")
	}

	// Connect to db
	db := database.MustConnectCustom(
		&database.DBConfig{
			TestMode:       database.TestModeOff,
			AutoMigrate:    false,
			RetryOnFailure: true,
			Extensions:     nil,
			LogMode:        conf.LogMode,
			User:           conf.PGUser,
			Host:           conf.PGHost,
			Port:           conf.PGPort,
			Pass:           conf.PGPass,
			DBName:         "infomanager",
		},
		models.GetModels(),
	)

	connection, err := grpc.Dial(conf.InfoManagerURI, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).WithField("url", conf.InfoManagerURI).Fatal("unable to dial info-manager service")
	}

	return graphql.Resolver{DB: db}, conf, proto.NewInfoManagerClient(connection)
}

func main() {
	resolver, conf, client := initialise()

	r := httpserver.NewGinServer()

	graphql.Init(client)

	r.POST("/gql", gqlHandler(resolver))
	r.GET("/gql", gqlHandler(resolver))
	if conf.GqlPlaygroundEnabled {
		r.GET("/", playgroundHandler())
	}

	log.Infof("Auth-gql running on port %d", conf.Port)
	log.Fatal(r.Run(":" + strconv.Itoa(conf.Port)))
}

// Defining the Graphql handler
func gqlHandler(resolver graphql.Resolver) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(resolver.DefaultConfig()))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/gql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
