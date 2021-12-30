// Package models provides database functionality
package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// GetByBsn gets a user by bsn
func (m *UserDB) GetByBsn(bsn string) (*User, error) {
	var user User
	err := m.Db.Find(&user, "bsn = ?", bsn).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, err
}

// GetByPseudo gets a user by pseudo
func (m *UserDB) GetByPseudo(pseudo string) (*User, error) {
	var user User
	err := m.Db.Find(&user, "pseudonym = ?", pseudo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, err
}
