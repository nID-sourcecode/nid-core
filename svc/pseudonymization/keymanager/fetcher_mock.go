package keymanager

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/mock"
)

// JWKSFetcherMock mocks a client for fetching JWKS endpoints
type JWKSFetcherMock struct {
	mock.Mock
}

// Fetch mocks the Fetch method
func (a *JWKSFetcherMock) Fetch(url string) (*jwk.Set, error) {
	args := a.Called(url)

	return args.Get(0).(*jwk.Set), args.Error(1)
}
