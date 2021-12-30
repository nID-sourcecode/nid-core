package models

// DefaultModel returns a new default savings account struct
func (m *SavingsAccountDB) DefaultModel() *SavingsAccount {
	return &SavingsAccount{
		Name: "Piggy",
	}
}
