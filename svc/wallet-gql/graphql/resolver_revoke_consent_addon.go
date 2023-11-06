package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
)

// AuthorizationURI Contains the uri for the authorization service
// FIXME AuthorizationURI should not be a global variable https://lab.weave.nl/nid/nid-core/-/issues/41
// nolint: gochecknoglobals
var AuthorizationURI string

var errUnableToRevokeConsent = fmt.Errorf("unable to revoke consent")

func (r *mutationResolver) CreateRevokeConsent(ctx context.Context, input CreateRevokeConsent) (*RevokeConsent, error) {
	var consent models.Consent
	var client models.Client
	if err := r.DB.Where("id = ?", input.ID).Find(&consent).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Where("client_id = ?", consent.ClientID).Find(&client).Error; err != nil {
		return nil, err
	}
	now := time.Now()
	consent.Revoked = &now

	values := map[string]string{
		"token":           consent.AccessToken,
		"client_id":       client.ExtClientID,
		"token_type_hint": "access_token",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal consent revokation values")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, AuthorizationURI+"/oidc/userrevoke", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create revocation request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to perform recovation requests")
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.WithError(err).Error("unable to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read response body")
		}
		errorResp := RFC6749Error{}
		err = json.Unmarshal(responseData, &errorResp)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal revocation request")
		}
		log.Errorf("error while revoking consent, authorization "+
			"returned status code: %v, with error: [%v] and error description: [%v]",
			errorResp.Code, errorResp.Name, errorResp.Description)

		return nil, errUnableToRevokeConsent
	}

	if err := r.DB.Save(&consent).Error; err != nil {
		return nil, errors.Wrap(err, "unable to store consent recovation")
	}

	// Also update other querymodels access was granted for in same consent
	var batch []models.Consent
	if err := r.DB.Where("access_token = ?", consent.AccessToken).Find(&batch).Error; err != nil {
		return nil, errors.Wrap(err, "unable to list query models for revoked consent")
	}
	for i := range batch {
		batch[i].Revoked = &now
		if err := r.DB.Save(&batch[i]).Error; err != nil {
			return nil, errors.Wrap(err, "unable to update query models for revoked consent")
		}
	}

	return &RevokeConsent{ID: consent.ID, Revoked: *consent.Revoked}, nil
}

// RFC6749Error definition of RFC6749 error
type RFC6749Error struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
	Hint        string `json:"error_hint,omitempty"`
	Code        int    `json:"status_code,omitempty"`
	Debug       string `json:"error_debug,omitempty"`
}
