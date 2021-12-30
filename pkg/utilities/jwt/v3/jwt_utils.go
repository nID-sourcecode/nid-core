package jwt

import (
	"crypto/rand"
	"crypto/rsa"
)

// GenerateTestKeys generate JWT key pair
func GenerateTestKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	return key, &key.PublicKey, nil
}

// HasScope checks if given scope is present in list of scopes
func HasScope(scope string, scopes []string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}
