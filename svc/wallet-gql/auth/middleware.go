package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	generr "lab.weave.nl/weave/generator/pkg/errors"
)

// A private key for context that only this package can access. This is
// important to prevent collisions between different context uses.
var userCtxKey = &contextKey{"user"} //nolint:gochecknoglobals

type contextKey struct {
	name string
}

// GetUser finds the User from the context. REQUIRES Middleware to have run.
func GetUser(ctx context.Context) *models.User {
	raw, _ := ctx.Value(userCtxKey).(*models.User)
	return raw
}

// SetUser sets the User to the context.
func SetUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

// UserHasScope checks if the user has a scope
func UserHasScope(ctx context.Context, scope string) bool {
	user := GetUser(ctx)

	// Return false if not logged in.
	if user == nil {
		return false
	}

	// Return false if an empty scope is given.
	if len(scope) == 0 {
		return false
	}

	// Wrap scope name with " to make sure the strings.Contains function
	// doesn't find a partial match.
	scope = fmt.Sprintf("\"%s\"", scope)

	// Return true if the scopes of the user contain the given scope.
	return strings.Contains(string(user.Scopes.RawMessage), scope)
}

const defaultInternalServerError = "{\"errors\":[{\"message\":\"internal server error\"}]}"

func sendGraphQLError(ctx context.Context, w http.ResponseWriter, path ast.Path, err error) {
	// Present the error
	gqlErr := gqlerror.WrapPath(path, err)
	presenter := generr.NewPresenter(&generr.LogrusErrorLogger{}) // FIXME make this logger configurable
	presentedErr := presenter.Present(ctx, gqlErr)
	res := &graphql.Response{
		Errors: gqlerror.List{presentedErr},
	}

	// Make sure http status code is correct
	statusCode := http.StatusUnauthorized
	var isInternalErr generr.WithIsInternal
	if errors.As(err, &isInternalErr) {
		if isInternalErr.IsInternal() {
			statusCode = http.StatusInternalServerError
		}
	}

	// Marshal and return
	jsonRes, marshalErr := json.Marshal(res)
	if marshalErr != nil {
		http.Error(w, defaultInternalServerError, http.StatusInternalServerError)
		return
	}
	http.Error(w, string(jsonRes), statusCode)
}

// ErrUnauthorized graphql unauthorized error
var ErrUnauthorized = generr.NewGraphQLError("unauthorized", false, "UNAUTHORIZED")
