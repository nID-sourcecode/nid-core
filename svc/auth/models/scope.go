package models

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Scope struct {
	ID        uuid.UUID   `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	Resource  string      `json:"resource"`
	Scope     string      `json:"scope"`
	Audiences []*Audience `gorm:"many2many:scopes_audiences;" json:"audiences"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m Scope) TableName() string {
	return "scopes"
}

// ScopeDB is the implementation of the storage interface for
// Scope.
type ScopeDB struct {
	Db *gorm.DB // Deprecated: Use ScopeDB.DB() instead.
}

// NewScopeDB creates a new storage type.
func NewScopeDB(db *gorm.DB) *ScopeDB {
	return &ScopeDB{Db: db}
}

// DB returns the underlying database.
func (m *ScopeDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(Scope{}) instead.
func (m *ScopeDB) TableName() string {
	return "scopes"
}

// CRUD Functions

// Get returns a single Scope as a Database Model
func (m *ScopeDB) Get(ctx context.Context, id uuid.UUID) (*Scope, error) {
	var native Scope
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// ListAllMatching returns all of the Scope records in the database that match input list of scopes and include audiences.
func (m *ScopeDB) ListAllMatching(ctx context.Context, scopes []string) ([]*Scope, error) {
	var objs []*Scope
	err := m.Db.Table(m.TableName()).Where("scope IN (?)", scopes).Preload("Audiences").Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// List returns an array of Scope
func (m *ScopeDB) List(ctx context.Context) ([]*Scope, error) {
	var objs []*Scope
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *ScopeDB) Add(ctx context.Context, model *Scope) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *ScopeDB) Update(ctx context.Context, model *Scope) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *ScopeDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&Scope{}).Error
	if err != nil {
		return err
	}

	return nil
}
