// Package pseudonymization provides functionality to encode and decode pseudonyms
package pseudonymization

import (
	"encoding/base64"
	"fmt"

	"github.com/miscreant/miscreant.go"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/svc/pseudonymization/keyutil"
)

const nonceSize = 96

// ErrInternalIDLargerThanServiceID Error definitions
var (
	ErrInternalIDLargerThanServiceID = fmt.Errorf("internal ID should be larger than service ID")
)

// Encode encodes a
func Encode(internalID, serviceID []byte) ([]byte, error) {
	if len(serviceID) > len(internalID) {
		return nil, ErrInternalIDLargerThanServiceID
	}

	key, err := keyutil.LoadKey()
	if err != nil {
		return nil, errors.Wrap(err, "unable to load key for encoding pseudonym")
	}

	c, err := miscreant.NewAEAD("AES-SIV", key, nonceSize)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create AES key for encoding pseudonym")
	}

	nonce, err := keyutil.LoadNonce(c)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load nonce for encoding pseudonym")
	}

	plaintext := make([]byte, len(internalID))
	copy(plaintext, internalID)
	for i, b := range serviceID {
		plaintext[i] ^= b
	}

	pseudonym := make([]byte, len(plaintext)+c.Overhead())

	c.Seal(pseudonym[:0], nonce, plaintext, nil)

	return pseudonym, nil
}

// Decode decodes pseudonym for given service
func Decode(pseudonym string, serviceID []byte) ([]byte, error) {
	key, err := keyutil.LoadKey()
	if err != nil {
		return nil, errors.Wrap(err, "unable to load key")
	}

	c, err := miscreant.NewAEAD("AES-SIV", key, nonceSize)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create AES cipher")
	}

	nonce, err := keyutil.LoadNonce(c)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load nonce")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(pseudonym)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode pseudonym string")
	}

	l := len(ciphertext) - c.Overhead()
	if l < 0 {
		return nil, PseudonymTooShortError{pseudonym}
	}
	plaintext := make([]byte, l)

	if _, err := c.Open(plaintext[:0], nonce, ciphertext, nil); err != nil {
		return nil, errors.Wrap(err, "unable to decrypt cipher text")
	}

	internalID := make([]byte, len(plaintext))
	copy(internalID, plaintext)
	for i, b := range serviceID {
		internalID[i] ^= b
	}

	return internalID, nil
}

// PseudonymTooShortError error definition for incorrect pseudonym length
type PseudonymTooShortError struct {
	Pseudonym string
}

func (err PseudonymTooShortError) Error() string {
	return fmt.Sprintf("Pseudonym \"%s\" too short", err.Pseudonym)
}
