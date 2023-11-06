// Package authtoken provides functionality for generating random tokens
package authtoken

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

// GenerateBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, errors.Wrap(err, "unable to read random bytes")
	}

	return b, nil
}

// NewToken returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func NewToken(n int) (string, error) {
	b, err := generateBytes(n)

	return base64.RawURLEncoding.EncodeToString(b), err
}

// Hash returns a sha256 hash based on given token
func Hash(token string) (string, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(strings.ToLower(token))); err != nil {
		return "", errors.Wrap(err, "unable to create sh256 hash from token")
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
