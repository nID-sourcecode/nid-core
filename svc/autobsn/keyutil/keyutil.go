// Package keyutil provides utility functions for parsing the RSA key and its corresponding JWK
package keyutil

import (
	"crypto/rsa"
	goErr "errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

const minRSAKeySizeInBytes = 256 // 2048 bits, as per https://tools.ietf.org/html/rfc7518#section-4.2

// ErrInvalidRsaKeyLength is returned when the provided rsa key is too small
var ErrInvalidRsaKeyLength = goErr.New("rsa key length should be at least 2048 bits")

// ParseKeypair parses the RSA key and checks its size
func ParseKeypair(rsaPrivateKeyPEM string) (*rsa.PrivateKey, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(rsaPrivateKeyPEM))
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse rsa prive key from pem")
	}
	if privateKey.Size() < minRSAKeySizeInBytes {
		return nil, errors.Wrapf(ErrInvalidRsaKeyLength, "rsa key size %d is too small", privateKey.Size())
	}

	return privateKey, nil
}

// CreateJWKSet creates the correct jwk set from the RSA public key
func CreateJWKSet(key *rsa.PublicKey) (*jwk.Set, error) {
	jwkKey := jwk.NewRSAPublicKey()
	if err := jwkKey.FromRaw(key); err != nil {
		return nil, errors.Wrap(err, "unable to create jwk from rsa pub key")
	}
	if err := jwkKey.Set(jwk.KeyUsageKey, string(jwk.ForEncryption)); err != nil {
		return nil, errors.Wrap(err, "unable to set jwk usage key")
	}
	if err := jwkKey.Set(jwk.AlgorithmKey, "RSA1_5"); err != nil {
		return nil, errors.Wrap(err, "unable to set jwk algorithm key")
	}

	return &jwk.Set{Keys: []jwk.Key{jwkKey}}, nil
}
