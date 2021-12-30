// Package models provides the model design
package models

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// GetWithScriptSourcesByVersion returns a single script with related script sources matched on version as a Database Model
func (m *ScriptDB) GetWithScriptSourcesByVersion(ctx context.Context, id uuid.UUID, version string) (*Script, error) {
	var native Script
	err := m.Db.Table(m.TableName()).
		Where("id = ?", id).
		Preload("ScriptSources", func(db *gorm.DB) *gorm.DB {
			return db.Where("version = ?", version)
		}).
		Find(&native).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &native, err
}
