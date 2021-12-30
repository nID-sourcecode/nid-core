package grpcserver

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt"
)

//nolint: gochecknoglobals
var (
	publicKey *rsa.PublicKey
)

// BearerPrefix prefix of jwt bearer
const BearerPrefix string = "Bearer "

// Error messages and variables used in parsing jwt's
var (
	ErrInvalidToken          error = fmt.Errorf("invalid token")
	ErrIncorrectBearerLength error = fmt.Errorf("incorrect bearer length")
)

// GetClaimsWithoutValidation only returns the claims for a jwt token and does not check token validity.
func GetClaimsWithoutValidation(bearer string) (jwt.MapClaims, error) {
	token, err := parseToken(bearer)
	if err != nil {
		return nil, err
	}
	return getClaims(token), nil
}

// parseToken parse a jwt string token
func parseToken(bearer string) (*jwt.Token, error) {
	if len(bearer) == 0 {
		return nil, ErrIncorrectBearerLength
	}
	return jwt.ParseWithClaims(
		bearer,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, ErrInvalidToken
			}
			// Unpack key from PEM encoded PKCS8
			return getPublicKey(), nil
		},
	)
}

// getClaims retrieve claims from token
func getClaims(token *jwt.Token) jwt.MapClaims {
	return token.Claims.(jwt.MapClaims)
}

// publicKey returns the public key
func getPublicKey() *rsa.PublicKey {
	return publicKey
}

// SetPubKey set pub key used for verifying the JWT token.
func SetPubKey(key *rsa.PublicKey) {
	publicKey = key
}
