// Package gql provides an interface for executing gql queries
package gql

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

var (
	// ErrRemoteErrorResponse is returned on a remote error
	ErrRemoteErrorResponse = errors.New("remote graphql error response")
	// ErrMethodNotSupported is returned on an unsupported method
	ErrMethodNotSupported = errors.New("method not supported")
)

// Method represents the HTTP methods supported by this package
type Method string

const (
	// MethodGet represents the HTTP GET method
	MethodGet Method = "GET"
	// MethodPost represents the HTTP POST method
	MethodPost Method = "POST"
)

// Client graphql client interface
type Client interface {
	// Post executes gql query on gqlClient
	Post(ctx context.Context, req Request, resp interface{}) error // Defaults to POST
	// Get executes gql query on gqlClient
	Get(ctx context.Context, req Request, resp interface{}) error
	// Run executes gql query on gqlClient
	Run(ctx context.Context, req Request, resp interface{}, method Method) error
}

// RestyGQLClient is a graphql client implementation using resty
type RestyGQLClient struct {
	client *resty.Client
}

// RequestBody is the request body that is marshalled for GQL requests
type RequestBody struct {
	Query     string
	Variables map[string]interface{}
}

// Response is the expected response type from a GQL endpoint on success
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the expected response type from a GQL endpoint on error
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error is a single GQL error as seen in ErrorResponse
type Error struct {
	Message string `json:"message"`
}

// Request is a graphQL request
type Request struct {
	Query     string
	Variables map[string]interface{}
	Headers   map[string]string
}

// NewRequest initialises a new request
func NewRequest(query string) Request {
	return Request{
		Query:     query,
		Variables: make(map[string]interface{}),
		Headers:   make(map[string]string),
	}
}

// Run executes a graphql request with the supplied method
func (c *RestyGQLClient) Run(ctx context.Context, req Request, resp interface{}, method Method) error {
	switch method {
	case MethodPost:
		return c.Post(ctx, req, resp)
	case MethodGet:
		return c.Get(ctx, req, resp)
	}

	return errors.Wrap(ErrMethodNotSupported, "unable to execute gql request")
}

// Post executes a POST graphql request
func (c *RestyGQLClient) Post(ctx context.Context, req Request, resp interface{}) error {
	body := RequestBody{
		Query:     req.Query,
		Variables: req.Variables,
	}

	res, err := c.client.R().SetContext(ctx).SetHeaders(req.Headers).SetBody(body).Post("")
	if err != nil {
		return errors.Wrap(err, "error calling graphql client")
	}

	return c.handleResponse(res, resp)
}

// Get executes a GET graphql request
func (c *RestyGQLClient) Get(ctx context.Context, req Request, resp interface{}) error {
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

func (c *RestyGQLClient) handleResponse(res *resty.Response, resp interface{}) error {
	if res.IsError() {
		errRes := ErrorResponse{}
		err := json.Unmarshal(res.Body(), &errRes)
		if err != nil || len(errRes.Errors) == 0 {
			return errors.Errorf("%w: %s", ErrRemoteErrorResponse, string(res.Body()))
		}

		return errors.Errorf("%w: %s", ErrRemoteErrorResponse, errRes.Errors[0].Message)
	}

	gqlRes := Response{
		Data: resp,
	}
	err := json.Unmarshal(res.Body(), &gqlRes)
	if err != nil {
		return errors.Errorf("error unmarshaling graphql response: %v", ErrRemoteErrorResponse)
	}

	return nil
}

// NewClient creates GraphQL client for a new URL
func NewClient(url string) Client {
	return &RestyGQLClient{
		client: resty.New().SetHostURL(url),
	}
}
