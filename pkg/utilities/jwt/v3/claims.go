package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofrs/uuid"
)

// DefaultClaims based on the https://tools.ietf.org/html/rfc7519#section-4.1 specification
type DefaultClaims struct {
	jwt.RegisteredClaims
}

// NewDefaultClaims will initiate the DefaultClaims struct
func NewDefaultClaims(expDurationHours time.Duration) *DefaultClaims {
	return &DefaultClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "weave.nl",
			Audience:  []string{"weave"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * expDurationHours)),
			ID:        uuid.Must(uuid.NewV4()).String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
		},
	}
}

// NewDefaultRefreshTokenClaims  returns claims for the refresh token
func NewDefaultRefreshTokenClaims(subjectID, tokenID string, expDurationHours time.Duration) DefaultClaims {
	return DefaultClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * expDurationHours)),
			ID:        tokenID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   subjectID,
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
}
