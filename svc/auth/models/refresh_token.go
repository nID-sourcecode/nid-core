package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// RefreshToken Relational Model
type RefreshToken struct {
	ID        uuid.UUID  `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	Session   *Session   `json:"session"`
	SessionID uuid.UUID  `gorm:"index:idx_refresh_token_session_id" json:"session_id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m RefreshToken) TableName() string {
	return "refresh_tokens"
}

// RefreshTokenDB is the implementation of the storage interface for
// RefreshToken.
type RefreshTokenDB struct {
	Db *gorm.DB // Deprecated: Use RefreshTokenDB.DB() instead.
}

// NewRefreshTokenDB creates a new storage type.
func NewRefreshTokenDB(db *gorm.DB) *RefreshTokenDB {
	return &RefreshTokenDB{Db: db}
}

// DB returns the underlying database.
func (m *RefreshTokenDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(RefreshToken{}) instead.
func (m *RefreshTokenDB) TableName() string {
	return "refresh_tokens"
}

// CRUD Functions

// Get returns a single RefreshToken as a Database Model
func (m *RefreshTokenDB) Get(ctx context.Context, id uuid.UUID) (*RefreshToken, error) {
	var native RefreshToken
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of RefreshToken
func (m *RefreshTokenDB) List(ctx context.Context) ([]*RefreshToken, error) {
	var objs []*RefreshToken
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *RefreshTokenDB) Add(ctx context.Context, model *RefreshToken) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *RefreshTokenDB) Update(ctx context.Context, model *RefreshToken) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *RefreshTokenDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&RefreshToken{}).Error
	if err != nil {
		return err
	}

	return nil
}
