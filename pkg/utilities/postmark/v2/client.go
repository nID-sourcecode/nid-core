// Package postmark provides an interface on top of the postmark client
package postmark

import (
	"github.com/keighl/postmark"
)

// EmailClient interface for postmark client
type EmailClient interface {
	SendTemplatedEmail(postmark.TemplatedEmail) (postmark.EmailResponse, error)
}

// DefaultClient holds postmark client information in it
type DefaultClient struct {
	client *postmark.Client
}

// SendTemplatedEmail via postmark client
// Returns the messageID for the email as tracked in postmark, on error, the messageID is also returned.
// When the messageID is set, there was an error on postmark's side
// nolint: gocritic
func (c *DefaultClient) SendTemplatedEmail(mail postmark.TemplatedEmail) (postmark.EmailResponse, error) {
	return c.client.SendTemplatedEmail(mail)
}

// NewClient returns client for postmark
func NewClient(apiToken, accountToken string) EmailClient {
	client := postmark.NewClient(apiToken, accountToken)
	return &DefaultClient{client: client}
}
