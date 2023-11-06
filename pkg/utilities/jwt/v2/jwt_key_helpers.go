package jwt

import (
	"crypto/rsa"
	"fmt"

	jwtgo "github.com/golang-jwt/jwt/v5"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

// error declarations
var (
	errJWTPrivateNotFound = fmt.Errorf("jwt private key not provided")
	errJWTPublicNotFound  = fmt.Errorf("jwt public key not provided")
)

// ParseKeys parse priv pub key pair as bytes
// Deprecated: use ParseKey
func ParseKeys(jwtKey, jwtPub []byte) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if len(jwtKey) == 0 {
		return nil, nil, errJWTPrivateNotFound
	}
	if len(jwtPub) == 0 {
		return nil, nil, errJWTPublicNotFound
	}
	keyParsed, err := jwtgo.ParseRSAPrivateKeyFromPEM(jwtKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load jwt private key error: %w", err)
	}

	pubParsed, err := ParsePubKey(jwtPub)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load jwt public key error: %w", err)
	}
	return keyParsed, pubParsed, nil
}

// ParseKey parses an RSA key encoded as pem
func ParseKey(jwtKey []byte) (*rsa.PrivateKey, error) {
	if len(jwtKey) == 0 {
		return nil, errJWTPrivateNotFound
	}
	keyParsed, err := jwtgo.ParseRSAPrivateKeyFromPEM(jwtKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load jwt private key")
	}

	return keyParsed, nil
}

// ParsePubKey parses public key PEM bytes
func ParsePubKey(jwtPub []byte) (*rsa.PublicKey, error) {
	pubParsed, err := jwtgo.ParseRSAPublicKeyFromPEM(jwtPub)
	if err != nil {
		return nil, fmt.Errorf("failed to load jwt public key error: %w", err)
	}
	return pubParsed, nil
}
