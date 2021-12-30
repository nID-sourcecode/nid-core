package models

// DefaultModel returns a new default bank account struct
func (m *BankAccountDB) DefaultModel() *BankAccount {
	return &BankAccount{
		AccountNumber: "000000010101",
		Amount:        7395000, // nolint
	}
}
