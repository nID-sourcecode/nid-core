package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/nID-sourcecode/nid-core/svc/auth/models"
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

// UserHasScope checks if the scope contains the user model
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
