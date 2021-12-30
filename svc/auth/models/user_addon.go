// Package models provides database functionality
package models

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// GetOnEmail get user on email
func (m *UserDB) GetOnEmail(ctx context.Context, email string) (*User, error) {
	var native User
	err := m.Db.Table(m.TableName()).Where("email = ?", strings.ToLower(email)).Find(&native).Error
	if err != nil {
		return nil, err
	}

	return &native, nil
}

// WalletUserModel returns the wallet user model
func (m *UserDB) WalletUserModel() *User {
	return &User{
		ID:     uuid.Must(uuid.FromString("7e880dfe-d77a-4477-ade7-66a96f76c0b2")),
		Scopes: postgres.Jsonb{RawMessage: []byte("[\"api:access\",\"api:clients:read\"]")},
	}
}
