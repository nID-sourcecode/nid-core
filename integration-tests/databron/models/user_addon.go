// Package models provides all database models for the databron service and also includes some addons
package models

import (
	"github.com/gofrs/uuid"
)

// DefaultModel returns a new default user struct
func (m *UserDB) DefaultModel() *User {
	return &User{
		ID:        uuid.FromStringOrNil("1d3efd44-e72a-49c7-888d-3fd1146a4c93"),
		FirstName: "John Dummy",
		LastName:  "Doe Ho",
		Pseudonym: "Ajmjkuq6JKiCWEevkB1V7SPNCRd8uHAKh0nABF7BSTXZV9a7k2eZ2iE2BxBdVI60",
		Bsn:       "123456789",
	}
}
