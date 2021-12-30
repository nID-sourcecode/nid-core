package models

import (
	"github.com/gofrs/uuid"
)

// DefaultModel is a default representation of the model
func (m *RedirectTargetDB) DefaultModel(clientID uuid.UUID) RedirectTarget {
	return RedirectTarget{
		ClientID:       clientID,
		RedirectTarget: "http://localhost:8081",
	}
}

// GetRedirectTarget retrieves the redirect target given a redirect uri and client id
func (m *RedirectTargetDB) GetRedirectTarget(redirectURI, clientID string) (*RedirectTarget, error) {
	var redirectTarget RedirectTarget
	err := m.Db.Find(&redirectTarget, "redirect_target = ? AND client_id = ?", redirectURI, clientID).Error
	if err != nil {
		return nil, err
	}

	return &redirectTarget, nil
}
