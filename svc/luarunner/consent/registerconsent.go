// Package consent provides functionality for retrieving and registering consents
package consent

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

// ErrServiceNotRecognized error definitions
var errServiceNotRecognized = fmt.Errorf("given service name was not recognised")

const consentMutation = `mutation addConsent($input: CreateRequest!){
  createRequest(input: $input) {
    id
  }
}`

type createRequestInput struct {
	ServiceID string                 `json:"serviceID"`
	Token     string                 `json:"token"`
	Metadata  map[string]interface{} `json:"metadata"`
}

const serviceIDQuery = `query getService($filter: ServiceFilterInput!) {
  services(filter: $filter) {
    id
  }
}`

type getServicesResponse struct {
	Services []struct {
		ID string `json:"id"`
	} `json:"services"`
}

// RegisterConsent registers a consent
func RegisterConsent(ctx context.Context, clientName, token string, metadata map[string]interface{}) error {
	gqlClient := graphql.NewClient("http://consentregister.nid/gql")
	idReq := graphql.NewRequest(serviceIDQuery)
	idReq.Var("filter", map[string]interface{}{
		"name": map[string]string{
			"eq": clientName,
		},
	})
	idRes := getServicesResponse{}
	if err := gqlClient.Run(ctx, idReq, &idRes); err != nil {
		return errors.Wrap(err, "unable to get service response from consent register")
	}

	if len(idRes.Services) == 0 {
		return errServiceNotRecognized
	}
	serviceID := idRes.Services[0].ID

	registerReq := graphql.NewRequest(consentMutation)
	registerReq.Var("input", createRequestInput{
		serviceID,
		token,
		metadata,
	})
	registerRes := make(map[string]interface{})
	if err := gqlClient.Run(ctx, registerReq, &registerRes); err != nil {
		return errors.Wrap(err, "unable to execute consent registry request")
	}

	return nil
}
