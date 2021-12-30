package main

import (
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vrischmann/envconfig"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/httpserver/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/auth-gql/auth"
	"lab.weave.nl/nid/nid-core/svc/auth-gql/graphql"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
)

func initialise() (*graphql.Resolver, *AuthGQLConfig, *gorm.DB) {
	// Init conf
	conf := &AuthGQLConfig{}
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to load environment config")
	}

	err := log.SetFormat(log.Format(conf.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	log.Info("Auth-gql initialising")

	// Connect to db
	db := database.MustConnectCustomWithCustomLogger(
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
			DBName:         "auth",
		},
		models.GetModels(),
		log.GetLogger(),
	)
	return &graphql.Resolver{DB: db}, conf, db
}

func main() {
	resolver, conf, db := initialise()

	serverOpts := httpserver.DefaultServerOptions()
	serverOpts.UseLogMiddleware = false

	r := httpserver.NewGinServerWithOpts(serverOpts)
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(resolver.DefaultConfig()))
	ginHandler := auth.CheckServiceAccount(resolver.DB, h, models.NewUserDB(db).WalletUserModel(), "wallet", conf.Namespace)

	r.POST("/gql", ginHandler)
	r.GET("/gql", ginHandler)
	if conf.GqlPlaygroundEnabled {
		r.GET("/", playgroundHandler())
	}

	log.Infof("Auth-gql running on port %d", conf.Port)
	log.Fatal(r.Run(":" + strconv.Itoa(conf.Port)))
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/gql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
