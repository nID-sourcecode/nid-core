package graphql

import (
	"context"
	"strings"

	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/models"
)

// AfterReadSetToken after read create the token field
// FIXME: Temporary hook for frontend due to the redundant token and access token in consent https://lab.weave.nl/ibnext/core/-/issues/19
func (cwh *CustomConsentHooks) AfterReadSetToken(ctx context.Context, tx *gorm.DB, model *models.Consent) error {
	jwtPayloadEncoded := strings.Split(model.AccessToken, ".")[1]

	jwtPayloadBytes, err := base64url.Decode(jwtPayloadEncoded)
	if err != nil {
		return errors.Wrap(err, "unable to decode jwt payload")
	}
	model.Token = postgres.Jsonb{
		RawMessage: jwtPayloadBytes,
	}
	return nil
}
