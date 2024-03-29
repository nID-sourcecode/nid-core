package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// User Relational Model
type User struct {
	ID        uuid.UUID      `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	Email     string         `sql:"unique" json:"email"`
	Password  string         `json:"password"`
	Scopes    postgres.Jsonb `sql:"default:'[\"api:access\"]'::jsonb" json:"scopes"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m User) TableName() string {
	return "users"
}

// UserDB is the implementation of the storage interface for
// User.
type UserDB struct {
	Db *gorm.DB // Deprecated: Use UserDB.DB() instead.
}

// NewUserDB creates a new storage type.
func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{Db: db}
}

// DB returns the underlying database.
func (m *UserDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(User{}) instead.
func (m *UserDB) TableName() string {
	return "users"
}

// CRUD Functions

// Get returns a single User as a Database Model
func (m *UserDB) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	var native User
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of User
func (m *UserDB) List(ctx context.Context) ([]*User, error) {
	var objs []*User
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *UserDB) Add(ctx context.Context, model *User) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *UserDB) Update(ctx context.Context, model *User) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *UserDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		return err
	}

	return nil
}
