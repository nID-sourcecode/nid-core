package authtoken

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthTokenTestSuite struct {
	suite.Suite
}

func TestAuthTokenTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTokenTestSuite))
}

func (s *AuthTokenTestSuite) TestGenerateBytes() {
	// generateBytes function only returns an error when the systems generator fails to function correctly
	// i.e. we do not make a test case for that
	n := 32
	bytes, err := generateBytes(n)
	s.Require().NoError(err)
	s.Equal(n, len(bytes))
}

func (s *AuthTokenTestSuite) TestNewTokenHashing() {
	n := 32
	token, err := NewToken(n)
	s.Require().NoError(err)

	// DecodeString will return error when token encoding is invalid
	_, err = base64.RawURLEncoding.DecodeString(token)
	s.Require().NoError(err)

	_, err = Hash(token)
	s.Require().NoError(err)
}
