// Package keyutil provides utility functions for parsing the RSA key and its corresponding JWK
package keyutil

import (
	"crypto/rsa"
	goErr "errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
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
func CreateJWKSet(key *rsa.PublicKey) (jwk.Set, error) {
	jwkKey, err := jwk.FromRaw(key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create jwk from rsa pub key")
	}

	jwkKey.Set(jwk.AlgorithmKey, "RSA1_5")
	jwkKey.Set(jwk.KeyUsageKey, string(jwk.ForEncryption))

	jwkSet := jwk.NewSet()
	err = jwkSet.AddKey(jwkKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to add jwk to jwk set")
	}

	return jwkSet, nil
}
