package messagebird

import (
	"github.com/stretchr/testify/mock"
)

// MockedClient mock for the postmark verifier interface
type MockedClient struct {
	mock.Mock
}

// NewPhoneLookup mock for the new phone lookup functionality
func (c *MockedClient) NewPhoneLookup(number string) (string, error) {
	args := c.Called(number)
	return args.String(0), args.Error(1)
}

// NewPhoneVerification mock for the new verification functionality
func (c *MockedClient) NewPhoneVerification(number, verificationType string) (string, error) {
	args := c.Called(number, verificationType)
	return args.String(0), args.Error(1)
}

// CheckPhoneVerification mock for the check verification functionality
func (c *MockedClient) CheckPhoneVerification(token, code string) error {
	args := c.Called(token, code)
	return args.Error(0)
}
