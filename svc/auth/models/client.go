package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// Client Relational Model
type Client struct {
	ID              uuid.UUID         `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	Color           string            `json:"color"`
	Icon            string            `json:"icon"`
	Logo            string            `json:"logo"`
	Metadata        postgres.Jsonb    `json:"metadata"`
	Name            string            `json:"name"`
	Password        string            `json:"password"`
	RedirectTargets []*RedirectTarget `json:"redirect_targets"`
	CreatedAt       time.Time         `json:"created_at"`
	DeletedAt       *time.Time        `json:"deleted_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m Client) TableName() string {
	return "clients"
}

// ClientDB is the implementation of the storage interface for
// Client.
type ClientDB struct {
	Db *gorm.DB // Deprecated: Use ClientDB.DB() instead.
}

// NewClientDB creates a new storage type.
func NewClientDB(db *gorm.DB) *ClientDB {
	return &ClientDB{Db: db}
}

// DB returns the underlying database.
func (m *ClientDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(Client{}) instead.
func (m *ClientDB) TableName() string {
	return "clients"
}

// CRUD Functions

// Get returns a single Client as a Database Model
func (m *ClientDB) Get(ctx context.Context, id uuid.UUID) (*Client, error) {
	var native Client
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of Client
func (m *ClientDB) List(ctx context.Context) ([]*Client, error) {
	var objs []*Client
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *ClientDB) Add(ctx context.Context, model *Client) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *ClientDB) Update(ctx context.Context, model *Client) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *ClientDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&Client{}).Error
	if err != nil {
		return err
	}

	return nil
}
