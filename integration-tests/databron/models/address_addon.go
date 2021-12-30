package models

// DefaultModel returns a new default address struct
func (m *AddressDB) DefaultModel() *Address {
	houseNumberAddon := "xyz"
	return &Address{
		HouseNumber:      1001, //nolint
		HouseNumberAddon: &houseNumberAddon,
		PostalCode:       "4321 YX",
	}
}
