package jwt

import (
	"crypto/rsa"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

// error declarations
var (
	ErrJWTPrivateNotFound error = fmt.Errorf("jwt private key not provided")
	ErrJWTPublicNotFound  error = fmt.Errorf("jwt public key not provided")
)

// ParseKeys parse priv pub key pair as bytes
// Deprecated: use ParseKey
func ParseKeys(jwtKey, jwtPub []byte) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if len(jwtKey) == 0 {
		return nil, nil, ErrJWTPrivateNotFound
	}
	if len(jwtPub) == 0 {
		return nil, nil, ErrJWTPublicNotFound
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
		return nil, ErrJWTPrivateNotFound
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
