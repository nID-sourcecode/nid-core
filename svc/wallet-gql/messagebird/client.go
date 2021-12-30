// Package messagebird contains helper functions for the messagebird client
package messagebird

import (
	"strings"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/lookup"
	"github.com/messagebird/go-rest-api/verify"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/variables"
)

// PhoneVerifier is a definition for a service that can create and validate phone verification requests
type PhoneVerifier interface {
	NewPhoneLookup(number string) (string, error)
	NewPhoneVerification(number string, verificationType string) (string, error)
	CheckPhoneVerification(token string, code string) error
}

// Messagebird implementation of PhoneVerifier
type Messagebird struct {
	Client *messagebird.Client
}

// NewPhoneLookup looks up a phone number in messagebird
func (m *Messagebird) NewPhoneLookup(number string) (string, error) {
	v, err := lookup.Read(m.Client, number, &lookup.Params{
		CountryCode: "NL",
	})
	if err != nil {
		return "", errors.Wrap(err, "unable to read phone number")
	}
	if v == nil {
		return "", errors.New("MISSING_RESPONSE")
	}

	return v.Formats.E164, nil
}

// NewPhoneVerification verifies the provided phone number
func (m *Messagebird) NewPhoneVerification(number, verificationType string) (string, error) {
	v, err := verify.Create(m.Client, number, &verify.Params{
		// FIXME: do not hardcode originator https://lab.weave.nl/twi/core/-/issues/86
		Originator: "nID",
		Timeout:    variables.VerifyPhoneNumberTimeout,
		Type:       strings.ToLower(verificationType),
	})
	if err != nil {
		return "", errors.Wrap(err, "unable to verify phone")
	}
	if v == nil {
		return "", errors.New("MISSING_RESPONSE")
	}

	return v.ID, nil
}

// CheckPhoneVerification checks if a the given token and code combination is valid
func (m *Messagebird) CheckPhoneVerification(token, code string) error {
	v, err := verify.VerifyToken(m.Client, token, code)
	if err != nil {
		return errors.Wrap(err, "unable to verify message bird token")
	}
	if v == nil {
		return errors.New("MISSING_RESPONSE")
	}
	if v.Status != "verified" {
		return errors.New("MISSING_VERIFIED")
	}

	return nil
}
