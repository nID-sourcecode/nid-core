package gqlclient

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type gqlResponse struct {
	Data   interface{}
	Errors []gqlErrMsg `json:"errors"`
}

type gqlErrMsg struct {
	Message string `json:"message"`
}

// RestyGQLClientFactory creates resty gql clients
type RestyGQLClientFactory struct{}

// NewClient creates a new resty GQL client
func (*RestyGQLClientFactory) NewClient(url string) Client {
	return &restyGQLClient{
		client: resty.New().SetHostURL(url),
	}
}

// NewRestyGQLClientFactory creates a new RestyGQLClientFactory
func NewRestyGQLClientFactory() *RestyGQLClientFactory {
	return &RestyGQLClientFactory{}
}

// restyGQLClient is a graphql client implementation using resty
type restyGQLClient struct {
	client *resty.Client
}

// Run executes a graphql request with the supplied method
func (c *restyGQLClient) Run(ctx context.Context, req Request, resp interface{}, method Method) error {
	switch method {
	case MethodPost:
		return c.Post(ctx, req, resp)
	case MethodGet:
		return c.Get(ctx, req, resp)
	}

	return ErrMethodNotSupported
}

// Post executes a POST graphql request
func (c *restyGQLClient) Post(ctx context.Context, req Request, resp interface{}) error {
	body := gqlRequest{
		Query:     req.Query,
		Variables: req.Variables,
	}

	res, err := c.client.R().SetContext(ctx).SetHeader("Content-Type", "application/json").SetHeaders(req.Headers).SetBody(body).Post("")
	if err != nil {
		return errors.Wrap(err, "error calling graphql client")
	}

	return c.handleResponse(res, resp)
}

// Get executes a GET graphql request
func (c *restyGQLClient) Get(ctx context.Context, req Request, resp interface{}) error {
	jsonVars, err := json.Marshal(req.Variables)
	if err != nil {
		return errors.Wrap(err, "failed to marshal variables")
	}
	params := map[string]string{
		"query":     req.Query,
		"variables": string(jsonVars),
	}

	res, err := c.client.R().SetContext(ctx).SetQueryParams(params).SetHeaders(req.Headers).Get("")
	if err != nil {
		return errors.Wrap(err, "error calling graphql client")
	}

	return c.handleResponse(res, resp)
}

func (c *restyGQLClient) handleResponse(res *resty.Response, resp interface{}) error {
	gqlRes := gqlResponse{
		Data: resp,
	}
	err := json.Unmarshal(res.Body(), &gqlRes)
	if err != nil {
		if res.IsError() {
			return errors.Errorf("%w (%d): %s", ErrRemoteErrorResponse, res.StatusCode(), res.Body())
		}
		return errors.Wrap(err, "error unmarshaling graphql response")
	}

	if len(gqlRes.Errors) > 0 {
		return errors.Errorf("%w (%d): %s", ErrRemoteErrorResponse, res.StatusCode(), gqlRes.Errors[0].Message)
	}

	return nil
}
