package postmark

import (
	"github.com/keighl/postmark"
	"github.com/stretchr/testify/mock"
)

// MockedClient mock for the postmark client interface
type MockedClient struct {
	mock.Mock
}

// SendTemplatedEmail mock for the send templated email functionality
func (c *MockedClient) SendTemplatedEmail(mail postmark.TemplatedEmail) (postmark.EmailResponse, error) {
	args := c.Called(mail)
	return args.Get(0).(postmark.EmailResponse), args.Error(1)
}

// NewMockedClient returns mocked client for postmark
func NewMockedClient() EmailClient {
	return &MockedClient{}
}
