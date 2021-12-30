package models

// DefaultModel returns a new default contact detail struct
func (m *ContactDetailDB) DefaultModel() *ContactDetail {
	return &ContactDetail{
		Phone: "+31688776655",
	}
}
