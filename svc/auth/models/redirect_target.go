package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// RedirectTarget Relational Model
type RedirectTarget struct {
	ID             uuid.UUID  `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	Client         *Client    `json:"client"`
	ClientID       uuid.UUID  `gorm:"index:idx_redirect_target_client_id" json:"client_id"`
	RedirectTarget string     `json:"redirect_target"`
	CreatedAt      time.Time  `json:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m RedirectTarget) TableName() string {
	return "redirect_targets"
}

// RedirectTargetDB is the implementation of the storage interface for
// RedirectTarget.
type RedirectTargetDB struct {
	Db *gorm.DB // Deprecated: Use RedirectTargetDB.DB() instead.
}

// NewRedirectTargetDB creates a new storage type.
func NewRedirectTargetDB(db *gorm.DB) *RedirectTargetDB {
	return &RedirectTargetDB{Db: db}
}

// DB returns the underlying database.
func (m *RedirectTargetDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(RedirectTarget{}) instead.
func (m *RedirectTargetDB) TableName() string {
	return "redirect_targets"
}

// CRUD Functions

// Get returns a single RedirectTarget as a Database Model
func (m *RedirectTargetDB) Get(ctx context.Context, id uuid.UUID) (*RedirectTarget, error) {
	var native RedirectTarget
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of RedirectTarget
func (m *RedirectTargetDB) List(ctx context.Context) ([]*RedirectTarget, error) {
	var objs []*RedirectTarget
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *RedirectTargetDB) Add(ctx context.Context, model *RedirectTarget) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *RedirectTargetDB) Update(ctx context.Context, model *RedirectTarget) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *RedirectTargetDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&RedirectTarget{}).Error
	if err != nil {
		return err
	}

	return nil
}
