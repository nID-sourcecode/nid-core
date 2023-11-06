package graphql

import (
	"context"
	"net/mail"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/postmark"
	postmarkUtils "github.com/nID-sourcecode/nid-core/svc/wallet-gql/postmark"
)

// PostmarkClient is a postmark email client
// FIXME PostmarkClient should not be a global variable https://lab.weave.nl/nid/nid-core/-/issues/41
// nolint: gochecknoglobals
var PostmarkClient postmarkUtils.EmailClient

// BeforeCreateHook hook called before creating an email address
func (h *CustomEmailAddressHooks) BeforeCreateHook(_ context.Context, tx *gorm.DB, input *CreateEmailAddress) error {
	// check for email validity
	_, err := mail.ParseAddress(input.EmailAddress)
	if err != nil {
		return errors.Wrapf(err, "invalid email address %s", input.EmailAddress)
	}

	// Try to find a match for email_address-user_id combination
	var found int64
	if err := tx.Model(&models.EmailAddress{}).Where("user_id = ? AND email_address = ?", input.UserID, input.EmailAddress).Count(&found).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return errors.Wrapf(err, "failed to run query to find email address %s for user %s", input.EmailAddress, input.UserID.String())
	}

	// If an email address is found break
	if found > 0 {
		return errors.New("email address already exists for this user")
	}

	return nil
}

// AfterCreateHook hook called after creating an email address
func (h *CustomEmailAddressHooks) AfterCreateHook(_ context.Context, tx *gorm.DB, model *models.EmailAddress) error {
	// Send verification email
	postmark := &postmark.Postmark{Client: PostmarkClient}
	token, err := postmark.NewEmailVerification(model.EmailAddress)
	if err != nil {
		log.Errorf("error creating verification_token for email_address %s, error %s", model.ID.String(), err)

		return errors.Wrap(err, "unable to create new email verification")
	}

	// Write verification session token
	if err := tx.Table(model.TableName()).Where("id = ?", model.ID).Updates(map[string]interface{}{"verification_token": token, "updated_at": time.Now()}).Error; err != nil {
		log.Errorf("error setting verification_token for email_address %s, error %s", model.ID.String(), err)

		return err
	}

	return nil
}
