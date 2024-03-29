package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// AccessModel Relational Model
type AccessModel struct {
	ID              uuid.UUID        `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	Audience        *Audience        `json:"audience"`
	AudienceID      uuid.UUID        `gorm:"index:idx_access_model_audience_id" json:"audience_id"`
	Description     string           `json:"description"`
	GqlAccessModel  *GqlAccessModel  `json:"gql_access_model"`
	Hash            string           `json:"hash"`
	JSONModel       string           `json:"json_model"`
	Name            string           `json:"name"`
	RestAccessModel *RestAccessModel `json:"rest_access_model"`
	Type            AccessModelType  `json:"type"`
	CreatedAt       time.Time        `json:"created_at"`
	DeletedAt       *time.Time       `json:"deleted_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m AccessModel) TableName() string {
	return "access_models"
}

// AccessModelDB is the implementation of the storage interface for
// AccessModel.
type AccessModelDB struct {
	Db *gorm.DB // Deprecated: Use AccessModelDB.DB() instead.
}

// NewAccessModelDB creates a new storage type.
func NewAccessModelDB(db *gorm.DB) *AccessModelDB {
	return &AccessModelDB{Db: db}
}

// DB returns the underlying database.
func (m *AccessModelDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(AccessModel{}) instead.
func (m *AccessModelDB) TableName() string {
	return "access_models"
}

// CRUD Functions

// Get returns a single AccessModel as a Database Model
func (m *AccessModelDB) Get(ctx context.Context, id uuid.UUID) (*AccessModel, error) {
	var native AccessModel
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of AccessModel
func (m *AccessModelDB) List(ctx context.Context) ([]*AccessModel, error) {
	var objs []*AccessModel
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *AccessModelDB) Add(ctx context.Context, model *AccessModel) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *AccessModelDB) Update(ctx context.Context, model *AccessModel) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *AccessModelDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&AccessModel{}).Error
	if err != nil {
		return err
	}

	return nil
}
