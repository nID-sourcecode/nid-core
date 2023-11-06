// This package contains the mock implementation of the password's IManager interface.
package mock

import "github.com/stretchr/testify/mock"

// PasswordManagerMock is a mock of the password's IManager interface.
type PasswordManagerMock struct {
	mock.Mock
}

// GenerateHash provides a mock function with given fields: password.
func (m *PasswordManagerMock) GenerateHash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

// ComparePassword provides a mock function with given fields: password, hash.
func (m *PasswordManagerMock) ComparePassword(password, hash string) (bool, error) {
	args := m.Called(password, hash)
	return args.Bool(0), args.Error(1)
}
