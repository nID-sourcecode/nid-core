package graphql

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	messagebirdRest "github.com/messagebird/go-rest-api"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/messagebird"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
)

// MessageBirdClient is a client for the postmark utils package
// FIXME MessageBirdClient should not be a global variable https://lab.weave.nl/nid/nid-core/-/issues/41
// nolint: gochecknoglobals
var MessageBirdClient *messagebirdRest.Client

// BeforeCreateHook hook called before created phone number
func (h *CustomPhoneNumberHooks) BeforeCreateHook(_ context.Context, tx *gorm.DB, input *CreatePhoneNumber) error {
	// Normalise phone number
	messagebird := &messagebird.Messagebird{Client: MessageBirdClient}
	normalised, err := messagebird.NewPhoneLookup(input.PhoneNumber)
	if err != nil {
		return errors.Wrapf(err, "failed to normalise phone number \"%s\"", input.PhoneNumber)
	}

	// Write normalised phone number to model for database insertion
	input.PhoneNumber = normalised

	// Try to find a match for phone_number-user_id combination
	var found int64
	if err := tx.Model(&models.PhoneNumber{}).Where("user_id = ? AND phone_number = ?", input.UserID, input.PhoneNumber).Count(&found).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return errors.Wrapf(err, "failed to run query to find phone number %s for user %s", input.PhoneNumber, input.UserID.String())
	}

	// If a phone number is found break
	if found > 0 {
		return errors.New("phone number already exists for this user")
	}

	return nil
}

// AfterCreateHook hook executed after creating a phone number
func (h *CustomPhoneNumberHooks) AfterCreateHook(_ context.Context, tx *gorm.DB, model *models.PhoneNumber) error {
	// Send verification SMS or call
	messagebird := &messagebird.Messagebird{Client: MessageBirdClient}
	token, err := messagebird.NewPhoneVerification(model.PhoneNumber, model.VerificationType.String())
	if err != nil {
		log.Errorf("error creating verification_token for phone_number %s, error %s", model.ID.String(), err)

		return errors.Wrap(err, "unable to create phone verification")
	}

	// Write verification session token
	if err := tx.Table(model.TableName()).Where("id = ?", model.ID).Updates(map[string]interface{}{"verification_token": token, "updated_at": time.Now()}).Error; err != nil {
		log.Errorf("error setting verification_token for phone_number %s, error %s", model.ID.String(), err)

		return err
	}

	return nil
}
