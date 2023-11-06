package models

import (
	"errors"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/jwt/v3"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

// GetWithClaims returns the refresh token from database that corresponds to the given claims
func (m *RefreshTokenDB) GetWithClaims(claims *jwt.DefaultClaims) (*RefreshToken, error) {
	var token RefreshToken
	err := m.Db.Table(m.TableName()).
		Where("session_id = ? AND id = ?", claims.Subject, claims.ID).
		Preload("Session").
		First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &token, err
}

// DeleteClientsRefreshTokens deletes all refresh tokens related to the ID of the client.
func (m *RefreshTokenDB) DeleteClientsRefreshTokens(clientID uuid.UUID) error {
	return m.Db.Where("client_id = ?", clientID.String()).Delete(RefreshToken{}).Error
}

// DeleteBySessionID deletes all refresh tokens related to the ID of the session.
func (m *RefreshTokenDB) DeleteBySessionID(sessionID uuid.UUID) error {
	return m.Db.Where("session_id = ?", sessionID.String()).Delete(RefreshToken{}).Error
}
