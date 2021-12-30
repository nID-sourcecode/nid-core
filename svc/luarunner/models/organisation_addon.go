// Package models generated for the database
package models

import (
	"github.com/jinzhu/gorm"
)

// GetOrganisationWithUzoviName retrieves a organisation with the given uzovi name.
func (m *OrganisationDB) GetOrganisationWithUzoviName(uzoviNami string) (Organisation, error) {
	var organisation Organisation
	err := m.DB().(*gorm.DB).Where("uzovi = ?", uzoviNami).Preload("Scripts").First(&organisation).Error
	return organisation, err
}
