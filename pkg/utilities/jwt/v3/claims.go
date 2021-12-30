package jwt

import (
	"time"

	"github.com/gofrs/uuid"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

var (
	// ErrIssuerIsEmpty means the issuer was empty when validating
	ErrIssuerIsEmpty = errors.New("issuer cannot be empty")
	// ErrAudienceIsEmpty means the audience was empty when validating
	ErrAudienceIsEmpty = errors.New("audience cannot be empty")
	// ErrTokenExpired means the expired at has passed
	ErrTokenExpired = errors.New("token expired")
	// ErrTokenUsedBeforeIssued means the token was used before the issued at
	ErrTokenUsedBeforeIssued = errors.New("token used before issued")
	// ErrTokenNotValidYet means the token was used before the not before time
	ErrTokenNotValidYet = errors.New("token is not valid yet")
)

// Claims specifies the methods for a JWT claim, it must just have a Valid method that determines if the token is invalid for any supported reason
type Claims interface {
	Valid() error
}

// ClaimsWithDefaultTimeFields specifies the methods for validation. These can be overwritten
type ClaimsWithDefaultTimeFields interface {
	GetExpiresAt() int64
	GetIssuedAt() int64
	GetNotBefore() int64
}

// DefaultClaims based on the https://tools.ietf.org/html/rfc7519#section-4.1 specification
type DefaultClaims struct {
	Audience  []string `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	JWTID     string   `json:"jti,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	Subject   string   `json:"sub,omitempty"`
}

// NewDefaultClaims will initiate the DefaultClaims struct
func NewDefaultClaims() *DefaultClaims {
	return &DefaultClaims{
		Issuer:    "weave.nl",
		Audience:  []string{"weave"},
		ExpiresAt: time.Now().Add(DurationDay).Unix(),
		JWTID:     uuid.Must(uuid.NewV4()).String(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Add(-2 * time.Minute).Unix(),
	}
}

// Valid will verify if the DefaultClaims is valid
func (d *DefaultClaims) Valid() error {
	now := time.Now().Unix()

	if d.Issuer == "" {
		return ErrIssuerIsEmpty
	}

	if len(d.Audience) == 0 {
		return ErrAudienceIsEmpty
	}

	if !VerifyExpiresAt(d, now) {
		delta := time.Unix(now, 0).Sub(time.Unix(d.ExpiresAt, 0))
		return errors.Wrapf(ErrTokenExpired, "by %v", delta)
	}

	if !VerifyIssuedAt(d, now) {
		return ErrTokenUsedBeforeIssued
	}

	if !VerifyNotBefore(d, now) {
		return ErrTokenNotValidYet
	}

	return nil
}

// GetExpiresAt will retrieve the exp value
func (d *DefaultClaims) GetExpiresAt() int64 {
	return d.ExpiresAt
}

// GetIssuedAt will retrieve the iat value
func (d *DefaultClaims) GetIssuedAt() int64 {
	return d.IssuedAt
}

// GetNotBefore will retrieve the nbf value
func (d *DefaultClaims) GetNotBefore() int64 {
	return d.NotBefore
}

// VerifyExpiresAt will verify the exp value
func VerifyExpiresAt(claims ClaimsWithDefaultTimeFields, now int64) bool {
	exp := claims.GetExpiresAt()
	if exp == 0 {
		return false
	}
	return now <= exp
}

// VerifyIssuedAt will verify the iat value
func VerifyIssuedAt(claims ClaimsWithDefaultTimeFields, now int64) bool {
	iat := claims.GetIssuedAt()
	if iat == 0 {
		return false
	}
	return now >= iat
}

// VerifyNotBefore will verify the nbf value
func VerifyNotBefore(claims ClaimsWithDefaultTimeFields, now int64) bool {
	nbf := claims.GetNotBefore()
	if nbf == 0 {
		return false
	}
	return now >= nbf
}
