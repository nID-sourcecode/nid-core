package models

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// Session Relational Model
type Session struct {
	ID                         uuid.UUID       `sql:"default:uuid_generate_v4()" gorm:"primary_key" json:"id"` // primary key
	AcceptedAccessModels       []*AccessModel  `gorm:"many2many:accepted_access_models_sessions;" json:"accepted_access_models"`
	Audience                   *Audience       `json:"audience"`
	AudienceID                 uuid.UUID       `gorm:"index:idx_session_audience_id" json:"audience_id"`
	AuthorizationCode          *string         `json:"authorization_code"`
	AuthorizationCodeGrantedAt *time.Time      `json:"authorization_code_granted_at"`
	Client                     *Client         `json:"client"`
	ClientID                   uuid.UUID       `gorm:"index:idx_session_client_id" json:"client_id"`
	FinaliseToken              string          `json:"finalise_token"`
	OptionalAccessModels       []*AccessModel  `gorm:"many2many:optional_access_models_sessions;" json:"optional_access_models"`
	RedirectTarget             *RedirectTarget `json:"redirect_target"`
	RedirectTargetID           uuid.UUID       `gorm:"index:idx_session_redirect_target_id" json:"redirect_target_id"`
	RequiredAccessModels       []*AccessModel  `gorm:"many2many:required_access_models_sessions;" json:"required_access_models"`
	State                      SessionState    `json:"state"`
	Subject                    string          `json:"subject"`
	CreatedAt                  time.Time       `json:"created_at"`
	DeletedAt                  *time.Time      `json:"deleted_at"`
	UpdatedAt                  time.Time       `json:"updated_at"`
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (m *Session) TableName() string {
	return "sessions"
}

// SessionDB is the implementation of the storage interface for
// Session.
type SessionDB struct {
	// nolint:stylecheck
	Db *gorm.DB // Deprecated: Use SessionDB.DB() instead.
}

// NewSessionDB creates a new storage type.
func NewSessionDB(db *gorm.DB) *SessionDB {
	return &SessionDB{Db: db}
}

// DB returns the underlying database.
func (m *SessionDB) DB() interface{} {
	return m.Db
}

// TableName returns the table name of the associated model
//
// Deprecated: Use db.Model(Session{}) instead.
func (m *SessionDB) TableName() string {
	return "sessions"
}

// CRUD Functions

// Get returns a single Session as a Database Model
func (m *SessionDB) Get(_ context.Context, id uuid.UUID) (*Session, error) {
	var native Session
	err := m.Db.Table(m.TableName()).Where("id = ?", id).Find(&native).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &native, err
}

// List returns an array of Session
func (m *SessionDB) List(_ context.Context) ([]*Session, error) {
	var objs []*Session
	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return objs, nil
}

// Add creates a new record.
func (m *SessionDB) Add(_ context.Context, model *Session) error {
	err := m.Db.Create(model).Error
	if err != nil {
		return err
	}

	return nil
}

// Update modifies a single record.
func (m *SessionDB) Update(ctx context.Context, model *Session) error {
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// Delete removes a single record.
func (m *SessionDB) Delete(_ context.Context, id uuid.UUID) error {
	err := m.Db.Where("id = ?", id).Delete(&Session{}).Error
	if err != nil {
		return err
	}

	return nil
}

// IsSessionExpired validates if the session is still valid with the given expiration time.
func (m *Session) IsSessionExpired(authorizationCodeExpirationTime time.Duration) bool {
	if m.AuthorizationCodeGrantedAt != nil {
		grantedAt := *m.AuthorizationCodeGrantedAt
		deadline := grantedAt.Add(authorizationCodeExpirationTime)
		now := time.Now()
		if !(now.After(grantedAt) && now.Before(deadline)) {
			return true
		}
	}

	return false
}
