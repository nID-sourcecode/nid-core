// Package gqlclient contains all the gql client wrappers
package gqlclient

import (
	"context"

	"github.com/gofrs/uuid"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
)

// IAuthClient specifies the AuthClient methods
type IAuthClient interface {
	FetchClient(ctx context.Context, clientID uuid.UUID) (*Client, error)
}

// AuthClient is base client for GQL
type AuthClient struct {
	client gqlclient.Client
}

// NewAuthClient will initialise a new AuthClient
func NewAuthClient(authURL string) *AuthClient {
	return &AuthClient{
		client: gqlclient.NewClient(authURL),
	}
}

const clientQuery = `
query($id: UUID!) {
  client (id: $id) {
    id
    color
    icon
    logo
    name
  }
}
`

// Client minimal version of a client
type Client struct {
	Color string `json:"color"`
	Icon  string `json:"icon"`
	Logo  string `json:"logo"`
	Name  string `json:"name"`
}

// ErrEmptyResult is returned if there were no clients in the response
var ErrEmptyResult = errors.New("empty result, something went wrong with gql call")

// FetchClient will retrieve client by ID
func (a *AuthClient) FetchClient(ctx context.Context, clientID uuid.UUID) (*Client, error) {
	clientReq := gqlclient.NewRequest(clientQuery)
	clientReq.Variables["id"] = clientID.String()

	res := make(map[string]Client)

	err := a.client.Get(ctx, clientReq, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unable to execute authorization client request")
	} else if len(res) == 0 {
		return nil, ErrEmptyResult
	}
	out := res["client"]
	return &out, nil
}
