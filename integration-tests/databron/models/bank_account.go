// Code generated by lab.weave.nl/weave/generator, DO NOT EDIT.

package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// BankAccount Relational Model
type BankAccount struct {
	ID              uuid.UUID         `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	AccountNumber   string            `json:"account_number"`
	Amount          int               `json:"amount"`
	SavingsAccounts []*SavingsAccount `json:"savings_accounts"`
	User            *User             `json:"user"`
	UserID          uuid.UUID         `gorm:"index:idx_bank_account_user_id" json:"user_id"`
	CreatedAt       time.Time         `json:"created_at"`
	DeletedAt       *time.Time        `json:"deleted_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m BankAccount) TableName() string {
	return "bank_accounts"
}

// BankAccountDB is the implementation of the storage interface for
// BankAccount.
type BankAccountDB struct {
	Db *gorm.DB // Deprecated: Use BankAccountDB.DB() instead.
}

// NewBankAccountDB creates a new storage type.
func NewBankAccountDB(db *gorm.DB) *BankAccountDB {
	return &BankAccountDB{Db: db}
}

// DB returns the underlying database.
func (m *BankAccountDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(BankAccount{}) instead.
func (m *BankAccountDB) TableName() string {
	return "bank_accounts"
}

// CRUD Functions

// Get returns a single BankAccount as a Database Model
func (m *BankAccountDB) Get(ctx context.Context, id uuid.UUID) (*BankAccount, error) {
	var native BankAccount
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return &native, err
}

// List returns an array of BankAccount
func (m *BankAccountDB) List(ctx context.Context) ([]*BankAccount, error) {
	var objs []*BankAccount
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *BankAccountDB) Add(ctx context.Context, model *BankAccount) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *BankAccountDB) Update(ctx context.Context, model *BankAccount) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *BankAccountDB) Delete(ctx context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&BankAccount{}).Error
	if err != nil {
		return err
	}

	return nil
}
