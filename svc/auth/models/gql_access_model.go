package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// GqlAccessModel Relational Model
type GqlAccessModel struct {
	ID            uuid.UUID    `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	AccessModel   *AccessModel `json:"access_model"`
	AccessModelID uuid.UUID    `gorm:"index:idx_gql_access_model_access_model_id" json:"access_model_id"`
	JSONModel     string       `json:"json_model"`
	Path          string       `json:"path"`
	CreatedAt     time.Time    `json:"created_at"`
	DeletedAt     *time.Time   `json:"deleted_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m GqlAccessModel) TableName() string {
	return "gql_access_models"
}

// GqlAccessModelDB is the implementation of the storage interface for
// GqlAccessModel.
type GqlAccessModelDB struct {
	Db *gorm.DB // Deprecated: Use GqlAccessModelDB.DB() instead.
}

// NewGqlAccessModelDB creates a new storage type.
func NewGqlAccessModelDB(db *gorm.DB) *GqlAccessModelDB {
	return &GqlAccessModelDB{Db: db}
}

// DB returns the underlying database.
func (m *GqlAccessModelDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(GqlAccessModel{}) instead.
func (m *GqlAccessModelDB) TableName() string {
	return "gql_access_models"
}

// CRUD Functions

// Get returns a single GqlAccessModel as a Database Model
func (m *GqlAccessModelDB) Get(ctx context.Context, id uuid.UUID) (*GqlAccessModel, error) {
	var native GqlAccessModel
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of GqlAccessModel
func (m *GqlAccessModelDB) List(ctx context.Context) ([]*GqlAccessModel, error) {
	var objs []*GqlAccessModel
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *GqlAccessModelDB) Add(ctx context.Context, model *GqlAccessModel) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *GqlAccessModelDB) Update(ctx context.Context, model *GqlAccessModel) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *GqlAccessModelDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&GqlAccessModel{}).Error
	if err != nil {
		return err
	}

	return nil
}
