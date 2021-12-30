// Package auth provides auth functionalities
package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"lab.weave.nl/nid/nid-core/pkg/gql"
	"lab.weave.nl/nid/nid-core/pkg/spiffeparser"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
)

// CheckServiceAccount is a middleware for verifying if the client is wallet
func CheckServiceAccount(db *gorm.DB, handler http.Handler, user *models.User, serviceAccount, namespace string) gin.HandlerFunc {
	spiffeParser := spiffeparser.NewDefaultSpiffeParser()

	return func(c *gin.Context) {
		// checks if client is wallet
		clientCert := c.Request.Header.Get("x-forwarded-client-cert")
		cert, err := spiffeParser.Parse(clientCert)
		if err != nil {
			log.WithError(err).Error("parsing client cert")
			c.JSON(http.StatusInternalServerError, gql.ErrorResponse{Errors: []gql.Error{{Message: "internal server error"}}})
			return
		}
		if cert.URI.GetNamespace() == namespace && cert.URI.GetServiceAccount() == serviceAccount {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), userCtxKey, user))
		}
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
