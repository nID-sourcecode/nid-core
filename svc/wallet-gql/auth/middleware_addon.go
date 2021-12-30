// Package auth provides auth functionalities
package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/vektah/gqlparser/v2/ast"

	"lab.weave.nl/nid/nid-core/svc/wallet-gql/models"
	generr "lab.weave.nl/weave/generator/pkg/errors"
)

// NewCustomIstioAuthMiddleware creates middleware for JWT's already checked by Istio
func NewCustomIstioAuthMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := ast.Path{ast.PathName("_request"), ast.PathName("headers"), ast.PathName("claims")}
			// We expect istio to return the claims in header JWT
			jwtBase64Claims := r.Header.Get("claims")
			// We are only interested in the subject
			claims := struct {
				Subject string `json:"sub"`
			}{}
			ctx := r.Context()
			claimsJSON, err := base64.RawURLEncoding.DecodeString(jwtBase64Claims)
			if err != nil {
				sendGraphQLError(ctx, w, path, fmt.Errorf("%v: can't parse claims from base 64: %w", ErrUnauthorized, err))
				return
			}

			err = json.Unmarshal(claimsJSON, &claims)
			if err != nil {
				sendGraphQLError(ctx, w, path, fmt.Errorf("%v: can't unmarshal claims: %w", ErrUnauthorized, err))
				return
			}

			// Find subject
			var user models.User
			if err := db.First(&user, "pseudonym = ?", claims.Subject).Error; err != nil {
				sendGraphQLError(ctx, w, path, generr.WrapAsInternal(err, "getting claims user"))
				return
			}
			// Add user to request context
			r = r.WithContext(context.WithValue(ctx, userCtxKey, &user))

			next.ServeHTTP(w, r)
		})
	}
}
