//go:generate env GO111MODULE=on GOBIN=$PWD/bin go install lab.weave.nl/weave/generator/cmd/gen
//go:generate env GO111MODULE=on GOBIN=$PWD/bin bin/gen lab.weave.nl/nid/nid-core/services/wallet
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	messagebird "github.com/messagebird/go-rest-api"
	"github.com/vrischmann/envconfig"

	"lab.weave.nl/nid/nid-core/pkg/cors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/httpserver/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	postmarkUtils "lab.weave.nl/nid/nid-core/pkg/utilities/postmark/v2"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/auth"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/graphql"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/models"
)

func initialise() (graphql.Resolver, *WalletConfig) {
	// Init conf
	conf := &WalletConfig{}
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to load environment config")
	}

	err := log.SetFormat(log.Format(conf.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	log.Info("Wallet-gql initialising")

	graphql.MessageBirdClient = messagebird.New(conf.Messagebird)
	graphql.PostmarkClient = postmarkUtils.NewClient(conf.Postmark.API, conf.Postmark.Account)
	graphql.AuthorizationURI = conf.AuthorizationURI

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
			DBName:         "wallet",
		},
		models.GetModels(),
		log.GetLogger(),
	)
	return graphql.Resolver{DB: db}, conf
}

func main() {
	resolver, conf := initialise()

	serverOpts := httpserver.DefaultServerOptions()
	serverOpts.UseLogMiddleware = false

	r := httpserver.NewGinServerWithOpts(serverOpts)

	r.POST("/gql", gqlHandler(resolver))
	r.GET("/gql", gqlHandler(resolver))

	log.Infof("Wallet-gql running on port %d", conf.Port)
	log.Fatal(r.Run(":" + strconv.Itoa(conf.Port)))

	// FakeDigiD() -- TODO: Depricated -- remove if app has switched to gql (mutation { createDigid(...) }

	// This is probably not needed anymore since consents are saved in the db now.
	// But it's useful for testing purposes.
	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse HTTP body: "+err.Error(), http.StatusBadRequest)

			return
		}

		since, err := time.Parse(r.PostForm["since"][0], time.RFC3339)
		if err != nil {
			http.Error(w, "Could not retrieve history: "+err.Error(), 500)

			return
		}

		history, err := getConsentHistory(r.Context(),
			r.PostForm["wid"][0],
			since,
		)
		if err != nil {
			http.Error(w, "Could not retrieve history: "+err.Error(), 500)

			return
		}

		reqBodyBytes := new(bytes.Buffer)
		if err := json.NewEncoder(reqBodyBytes).Encode(history); err != nil {
			http.Error(w, "Could not encode history: "+err.Error(), 500)

			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(reqBodyBytes.Bytes()); err != nil {
			log.Errorf("unable to return history, error: %s", err.Error())
		}
	})

	http.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			log.Errorf("unable to return v1/health value, error: %s", err.Error())
		}
	})
}

// Defining the Graphql handler
func gqlHandler(resolver graphql.Resolver) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := cors.Cors(cors.ReflectOrigin())(auth.NewCustomIstioAuthMiddleware(resolver.DB)(handler.NewDefaultServer(graphql.NewExecutableSchema(resolver.DefaultConfig()))))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
