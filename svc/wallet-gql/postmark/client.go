// Package postmark contains helper functions for sending emails using the weave utilities postmark client
package postmark

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/keighl/postmark"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	postmarkUtils "lab.weave.nl/nid/nid-core/pkg/utilities/postmark/v2"
)

const verificationEmailTemplateID = 19258979

// Token did not match error
var (
	ErrTokenDidNotMatch error = fmt.Errorf("token did not match")
)

// EmailVerifier is a definition for a service that can create and validate email verification requests
type EmailVerifier interface {
	NewEmailVerification(email string) (string, error)
	CheckEmailVerification(token, code string) error
}

// Postmark implementation of an EmailVerifier
type Postmark struct {
	Client postmarkUtils.EmailClient
}

// NewEmailVerification sends a new verification email to the provided email address
func (p *Postmark) NewEmailVerification(email string) (string, error) {
	randomUUID, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(err, "issue while creating a UUID")
	}

	token := randomUUID.String()[0:8] // Dirty way to get 8 random readable bytes

	if _, err := p.Client.SendTemplatedEmail(postmark.TemplatedEmail{
		TemplateId: verificationEmailTemplateID,
		TemplateModel: map[string]interface{}{
			"code": token,
		},
		From: "no-reply@n-id.network",
		To:   email,
		Tag:  "verify",
	}); err != nil {
		return "", errors.Wrap(err, "unable to send verification email")
	}

	pwManager := password.NewDefaultManager()
	hash, err := pwManager.GenerateHash(token)
	if err != nil {
		return "", errors.Wrap(err, "issue while hashing token")
	}

	return hash, nil
}

// CheckEmailVerification checks if the code is correct for the given token
func (p *Postmark) CheckEmailVerification(token string, code string) error {
	pwManager := password.NewDefaultManager()
	match, err := pwManager.ComparePassword(code, token)
	if err != nil {
		return errors.Wrap(err, "unable to verify email code")
	}

	if !match {
		return ErrTokenDidNotMatch
	}

	return nil
}
