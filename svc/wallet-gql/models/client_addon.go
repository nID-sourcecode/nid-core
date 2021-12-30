package models

import (
	"github.com/gofrs/uuid"
)

// GetByExtClientID will retrieve the client by ext_client_id property
func (m *ClientDB) GetByExtClientID(id uuid.UUID) (*Client, error) {
	var native Client
	err := m.Db.Model(&Client{}).Where("ext_client_id = ?", id).First(&native).Error
	if err != nil {
		return nil, err
	}
	return &native, nil
}
