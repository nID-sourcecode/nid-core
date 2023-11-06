package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// RestAccessModel Relational Model
type RestAccessModel struct {
	ID            uuid.UUID    `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	AccessModel   *AccessModel `json:"access_model"`
	AccessModelID uuid.UUID    `gorm:"index:idx_rest_access_model_access_model_id" json:"access_model_id"`
	Body          string       `json:"body"`
	Method        string       `json:"method"`
	Path          string       `json:"path"`
	Query         string       `json:"query"`
	CreatedAt     time.Time    `json:"created_at"`
	DeletedAt     *time.Time   `json:"deleted_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m RestAccessModel) TableName() string {
	return "rest_access_models"
}

// RestAccessModelDB is the implementation of the storage interface for
// RestAccessModel.
type RestAccessModelDB struct {
	Db *gorm.DB // Deprecated: Use RestAccessModelDB.DB() instead.
}

// NewRestAccessModelDB creates a new storage type.
func NewRestAccessModelDB(db *gorm.DB) *RestAccessModelDB {
	return &RestAccessModelDB{Db: db}
}

// DB returns the underlying database.
func (m *RestAccessModelDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(RestAccessModel{}) instead.
func (m *RestAccessModelDB) TableName() string {
	return "rest_access_models"
}

// CRUD Functions

// Get returns a single RestAccessModel as a Database Model
func (m *RestAccessModelDB) Get(ctx context.Context, id uuid.UUID) (*RestAccessModel, error) {
	var native RestAccessModel
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of RestAccessModel
func (m *RestAccessModelDB) List(ctx context.Context) ([]*RestAccessModel, error) {
	var objs []*RestAccessModel
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *RestAccessModelDB) Add(ctx context.Context, model *RestAccessModel) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *RestAccessModelDB) Update(ctx context.Context, model *RestAccessModel) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *RestAccessModelDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&RestAccessModel{}).Error
	if err != nil {
		return err
	}

	return nil
}
