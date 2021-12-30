package graphql

import (
	"context"

	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
)

// PasswordManager is the password manager responsible for hashing
// FIXME PasswordManager should not be a global variable but we cannot change the signature of the MutateAppPassword function https://lab.weave.nl/nid/nid-core/-/issues/41
// nolint: gochecknoglobals
var PasswordManager password.IManager

// MutateAppPassword hashes the password before it is put in the database
func (cwh *CustomUserHooks) MutateAppPassword(ctx context.Context, password string) (string, error) {
	return PasswordManager.GenerateHash(password)
}
