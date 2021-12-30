package main

import (
	"context"
	"time"

	"github.com/machinebox/graphql"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

const authURL = "http://authorization.nid"

const consentQuery = `
query ($sub: SessionFilterInput!) {
  sessions(filter: $sub) {
	updatedAt
    request {
      requestedAt
      client {
        name
        color
      }
      requestQueryModels {
        granted
        queryModel {
          name
          description
        }
      }
      requestScopes {
        granted
        scope {
          scope
        }
      }
    }
  }
}
`

type consentResponse struct {
	Sessions []struct {
		UpdatedAt string `json:"updatedAt"`
		Request   struct {
			RequestedAt string `json:"requestedAt"`
			Client      struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"client"`
			RequestQueryModels []struct {
				Granted    bool `json:"granted"`
				QueryModel struct {
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"queryModel"`
			} `json:"requestQueryModels"`
			RequestScopes []struct {
				Granted bool `json:"granted"`
				Scope   struct {
					Scope string `json:"scope"`
				} `json:"scope"`
			} `json:"requestScopes"`
		} `json:"request"`
	} `json:"sessions"`
}

func getConsentHistory(ctx context.Context, pseudo string, since time.Time) (consentResponse, error) {
	walletGqlClient := graphql.NewClient(authURL + "/gql")
	consentReq := graphql.NewRequest(consentQuery)
	consentReq.Var("sub", map[string]interface{}{
		"subject": map[string]string{
			"eq": pseudo,
		},
		"type": map[string]string{
			"eq": "CODE",
		},
		"updatedAt": map[string]string{
			"ge": since.Format(time.RFC3339),
		},
	})

	res := consentResponse{} // maybe just use a map for when format ever changes?

	if err := walletGqlClient.Run(ctx, consentReq, &res); err != nil {
		return consentResponse{}, errors.Wrap(err, "unable to retrieve consent from wallet")
	}

	return res, nil
}
