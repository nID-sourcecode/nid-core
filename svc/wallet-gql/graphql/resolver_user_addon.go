package graphql

import (
	"context"

	"github.com/nID-sourcecode/nid-core/pkg/password"
)

// PasswordManager is the password manager responsible for hashing
// FIXME PasswordManager should not be a global variable but we cannot change the signature of the MutateAppPassword function https://lab.weave.nl/nid/nid-core/-/issues/41
// nolint: gochecknoglobals
var PasswordManager password.IManager

// MutateAppPassword hashes the password before it is put in the database
func (cwh *CustomUserHooks) MutateAppPassword(_ context.Context, password string) (string, error) {
	return PasswordManager.GenerateHash(password)
}
