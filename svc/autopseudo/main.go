// Package autopseudo
package main

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/vrischmann/envconfig"

	"github.com/nID-sourcecode/nid-core/pkg/httpserver"
	"github.com/nID-sourcecode/nid-core/pkg/keyutil"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

type accessClaims struct {
	Subjects map[string]string
}

func (claims accessClaims) Valid() error {
	return nil
}

const (
	bearerScheme      = "Bearer "
	subjectIdentifier = "$$nid:subject$$"
)

type decryptAndApplyResponse struct {
	Body  string `json:"body"`
	Query string `json:"query"`
}

func main() {
	var conf AutoPseudoConfig
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to load config from environment")
	}

	err := log.SetFormat(log.Format(conf.LogFormat))
	if err != nil {
		log.WithError(err).Fatal("unable to set log format")
	}

	key, err := keyutil.ParseKeypair(conf.RSAPriv)
	if err != nil {
		log.Fatal(err)
	}
	jwkSet, err := keyutil.CreateJWKSet(&key.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	router := initRouter(jwkSet, key, &conf)

	log.Fatal(router.Run(fmt.Sprintf(":%d", conf.Port)))
}

func initRouter(jwkSet jwk.Set, key *rsa.PrivateKey, conf *AutoPseudoConfig) *gin.Engine {
	server := NewAutoPseudoServer(key, conf)
	router := httpserver.NewGinServer()
	router.GET("/jwks", func(c *gin.Context) {
		c.JSON(http.StatusOK, jwkSet)
	})

	router.Any("/decryptAndApply", server.DecryptAndApply)
	router.GET("/decrypt", server.Decrypt)

	return router
}
