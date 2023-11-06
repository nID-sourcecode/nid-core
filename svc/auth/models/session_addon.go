package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

// DefaultModel is a default representation of the model
func (m *SessionDB) DefaultModel(clientID, audienceID, redirectTargetID uuid.UUID, state *SessionState, subject *string) Session {
	session := Session{
		AudienceID:       audienceID,
		ClientID:         clientID,
		RedirectTargetID: redirectTargetID,
	}

	if state != nil {
		session.State = *state
	}

	if subject != nil {
		session.Subject = *subject
	}

	return session
}

// BeforeUpdate sets authorization_code_granted_at date when code is set
func (m *Session) BeforeUpdate(_ *gorm.DB) (err error) {
	if m.AuthorizationCode != nil && m.AuthorizationCodeGrantedAt == nil {
		grantedAt := time.Now()
		m.AuthorizationCodeGrantedAt = &grantedAt
	}
	return nil
}

// BeforeCreate sets authorization_code_granted_at date when code is set
func (m *Session) BeforeCreate(_ *gorm.DB) (err error) {
	if m.AuthorizationCode != nil {
		grantedAt := time.Now()
		m.AuthorizationCodeGrantedAt = &grantedAt
	}

	return nil
}

// CreateSession inserts a session in the db
func (m *SessionDB) CreateSession(session *Session) error {
	err := m.Db.Create(session).Error
	if err != nil {
		return err
	}

	return nil
}

// PreloadOption type for the different preload options
type PreloadOption uint8

const (
	preloadRequiredAndOptionalScopes PreloadOption = 1
	preloadAll                       PreloadOption = 2
)

func (m *SessionDB) getSession(query *gorm.DB, preload PreloadOption) (*Session, error) {
	var session Session

	// FIXME optimise this query (https://lab.weave.nl/twi/core/-/issues/107)
	if preload == preloadRequiredAndOptionalScopes || preload == preloadAll {
		query = query.Preload("Client").
			Preload("Audience").
			Preload("RequiredAccessModels").
			Preload("OptionalAccessModels").
			Preload("RedirectTarget")

		if preload == preloadAll {
			query = query.Preload("AcceptedAccessModels").
				Preload("RequiredAccessModels.GqlAccessModel").
				Preload("RequiredAccessModels.RestAccessModel").
				Preload("OptionalAccessModels.GqlAccessModel").
				Preload("OptionalAccessModels.RestAccessModel").
				Preload("AcceptedAccessModels.GqlAccessModel").
				Preload("AcceptedAccessModels.RestAccessModel")
		}
	}
	if err := query.Find(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSessionByID retrieves a session given it's ID
func (m *SessionDB) GetSessionByID(option PreloadOption, id string) (*Session, error) {
	return m.getSession(m.Db.Where("id = ?", id), option)
}

// GetSessionByCodeAndClientID retrieves a session given it's authorization code and client ID
func (m *SessionDB) GetSessionByCodeAndClientID(option PreloadOption, hash, clientID string) (*Session, error) {
	return m.getSession(m.Db.Where("authorization_code = ? AND client_id = ?", hash, clientID), option)
}

// GetSessionByIDAndSubject retrieves a session given it's id and subject
func (m *SessionDB) GetSessionByIDAndSubject(option PreloadOption, id, subject string) (*Session, error) {
	return m.getSession(m.Db.Where("id = ? AND subject = ?", id, subject), option)
}

// UpdateAcceptedAccessModels updates the accepted access models
func (m *SessionDB) UpdateAcceptedAccessModels(session *Session, combinedAccessModels []*AccessModel) error {
	err := m.Db.Model(&session).Association("AcceptedAccessModels").Replace(combinedAccessModels).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateSessionState updates the state of a session
func (m *SessionDB) UpdateSessionState(session *Session, state SessionState) error {
	err := m.Db.Model(&session).Select("state").Update(map[string]interface{}{"state": state}).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateSessionSubject updates the subject of the session
func (m *SessionDB) UpdateSessionSubject(session *Session, subject string) error {
	err := m.Db.Model(&session).Update("subject", subject).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateSessionAuthorizationCode updates the authorization code of a session
func (m *SessionDB) UpdateSessionAuthorizationCode(session *Session, code string) error {
	err := m.Db.Model(&session).Update("authorization_code", code).Error
	if err != nil {
		return err
	}

	return nil
}

// SetSessionFinaliseToken sets finalise token of the session.
// returns an error if a password is already set.
func (m *SessionDB) SetSessionFinaliseToken(session *Session, token string) error {
	if session.FinaliseToken != "" {
		return errors.New("session already has a finalise token")
	}

	err := m.Db.Model(&session).Update("finalise_token", token).Error
	if err != nil {
		return err
	}

	return nil
}
