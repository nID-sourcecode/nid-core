package gqlclient

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

// AuthClientMock mocks the IAuthClient interface
type AuthClientMock struct {
	mock.Mock
}

// FetchClient mocks the FetchClient request
func (m *AuthClientMock) FetchClient(ctx context.Context, clientID uuid.UUID) (*Client, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).(*Client), args.Error(1)
}
