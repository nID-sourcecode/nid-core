package postmark

import (
	"github.com/stretchr/testify/mock"
)

// MockedClient mock for the postmark verifier interface
type MockedClient struct {
	mock.Mock
}

// NewEmailVerification mock for the new verification functionality
func (c *MockedClient) NewEmailVerification(email string) (string, error) {
	args := c.Called(email)
	return args.String(0), args.Error(1)
}

// CheckEmailVerification mock for the check email verification functionality
func (c *MockedClient) CheckEmailVerification(token string, code string) error {
	args := c.Called(token, code)
	return args.Error(0)
}
