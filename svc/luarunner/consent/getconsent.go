package consent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/machinebox/graphql"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

// error definitions
var (
	errUnableToRetrieveConsent = fmt.Errorf("unable to retrieve consent")
	errClientNameNotRecognized = fmt.Errorf("given client name was not recognised")
)

type consentResponse struct {
	AccessToken string `json:"access_token"`
}

const clientIDQuery = `query getClient($filter: ClientFilterInput!) {
  clients(filter: $filter) {
    id
  }
}`

type getClientsResponse struct {
	Clients []struct {
		ID string `json:"id"`
	} `json:"clients"`
}

func getClientIDFromAuthorizationService(ctx context.Context, clientName string) (string, error) {
	gqlClient := graphql.NewClient("http://authorization.nid/gql")
	req := graphql.NewRequest(clientIDQuery)
	req.Var("filter", map[string]interface{}{
		"name": map[string]string{
			"eq": clientName,
		},
	})
	res := getClientsResponse{}
	if err := gqlClient.Run(ctx, req, &res); err != nil {
		return "", errors.Wrap(err, "unable to retrieve client id from authorization service")
	}
	if len(res.Clients) == 0 {
		return "", errClientNameNotRecognized
	}

	return res.Clients[0].ID, nil
}

// GetConsentToken retrieves a consent
func GetConsentToken(ctx context.Context, clientName, scopeString, audience, objectID, scriptID string) (string, error) {
	clientID, err := getClientIDFromAuthorizationService(ctx, clientName)
	if err != nil {
		return "", errors.Wrap(err, "unable to get client id from authorization service")
	}

	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}, // FIXME: https://lab.weave.nl/nid/nid-core/-/issues/62
	}

	authReq, err := http.NewRequest("GET", "http://authorization.nid/oidc/consent", nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to create consent request")
	}
	authReq.Header.Add("Accept", "application/json")
	params := authReq.URL.Query()
	params.Set("client_id", clientID)
	params.Set("scope", scopeString)
	params.Set("audience", audience)
	params.Set("nonce", "jjjjjjjjjjjjj")
	params.Set("state", "magic_state")
	params.Set("response_type", "code")
	params.Set("object_id", objectID)
	params.Set("script_id", scriptID)
	authReq.URL.RawQuery = params.Encode()

	authReq = authReq.WithContext(ctx)
	resp, err := httpClient.Do(authReq)
	if err != nil {
		return "", errors.Wrap(err, "unable to perform consent request")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != 302 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Wrap(err, "unable to read failed consent requests response body")
		}

		return "", fmt.Errorf("%s, %w", b, errUnableToRetrieveConsent)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.WithError(err).Error("unable to close response body")
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read consent response body")
	}
	tokenWrapper := consentResponse{}
	if err := json.Unmarshal(body, &tokenWrapper); err != nil {
		return "", errors.Wrap(err, "unable to unmarshal consent response")
	}

	return tokenWrapper.AccessToken, nil
}
