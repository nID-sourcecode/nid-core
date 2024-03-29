package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Audience Relational Model
type Audience struct {
	ID           uuid.UUID      `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	AccessModels []*AccessModel `json:"access_models"`
	Audience     string         `json:"audience"`
	Namespace    string         `json:"namespace"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    *time.Time     `json:"deleted_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Scopes       []*Scope       `gorm:"many2many:scopes_audiences;" json:"scopes"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m Audience) TableName() string {
	return "audiences"
}

// AudienceDB is the implementation of the storage interface for
// Audience.
type AudienceDB struct {
	Db *gorm.DB // Deprecated: Use AudienceDB.DB() instead.
}

// NewAudienceDB creates a new storage type.
func NewAudienceDB(db *gorm.DB) *AudienceDB {
	return &AudienceDB{Db: db}
}

// DB returns the underlying database.
func (m *AudienceDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(Audience{}) instead.
func (m *AudienceDB) TableName() string {
	return "audiences"
}

// CRUD Functions

// Get returns a single Audience as a Database Model
func (m *AudienceDB) Get(ctx context.Context, id uuid.UUID) (*Audience, error) {
	var native Audience
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of Audience
func (m *AudienceDB) List(ctx context.Context) ([]*Audience, error) {
	var objs []*Audience
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *AudienceDB) Add(ctx context.Context, model *Audience) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *AudienceDB) Update(ctx context.Context, model *Audience) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *AudienceDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&Audience{}).Error
	if err != nil {
		return err
	}

	return nil
}
