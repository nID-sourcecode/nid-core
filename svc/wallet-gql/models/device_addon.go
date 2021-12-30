package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// GetByCode gets device by code
func (m *DeviceDB) GetByCode(code string, preloadUser bool) (*Device, error) {
	query := m.Db
	if preloadUser {
		query = query.Preload("User")
	}

	var device Device
	err := query.Find(&device, "code = ?", code).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &device, err
}
