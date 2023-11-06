// Package retryhttp contains the http implementation of the callback handler.
package retryhttp

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	httpClient "github.com/nID-sourcecode/nid-core/pkg/utilities/http"
)

const (
	authCodeParamKey  = "authorization_code"
	successStatusCode = http.StatusAccepted
	retryMaxWait      = 4 * time.Hour
	requestTimeout    = 15 * time.Second
)

// CallbackHandler implements te callback handler with a http client.
type CallbackHandler struct {
	Client httpClient.Client
}

// NewCallbackHandler creates a new callback handler with a http retry client.
func NewCallbackHandler(retryMax int) *CallbackHandler {
	retryClient := retryablehttp.NewClient()
	retryClient.CheckRetry = checkRetry
	retryClient.RetryMax = retryMax
	retryClient.RetryWaitMax = retryMaxWait

	httpClient := retryClient.StandardClient()
	httpClient.Timeout = requestTimeout

	return &CallbackHandler{
		Client: httpClient,
	}
}

// HandleCallback handles a callback request with a GET request.
func (c *CallbackHandler) HandleCallback(ctx context.Context, location string, authorizationCode string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", location, nil)
	if err != nil {
		return errors.Wrapf(err, "creating request", location)
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		req.Header.Set("x-b3-traceid", md.Get("x-b3-traceid")[0])
		req.Header.Set("x-b3-spanid", md.Get("x-b3-spanid")[0])
		req.Header.Set("x-b3-parentspandid", md.Get("x-b3-parentspanid")[0])
	}

	q := req.URL.Query()
	q.Add(authCodeParamKey, authorizationCode)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "executing request for url %s", location)
	}

	if resp.StatusCode != successStatusCode {
		return errors.Errorf("request returned error status code: %d", resp.StatusCode)
	}

	err = resp.Body.Close()
	if err != nil {
		return errors.Wrap(err, "closing response body")
	}

	return nil
}

func checkRetry(ctx context.Context, resp *http.Response, _ error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	return resp.StatusCode != http.StatusAccepted, nil
}
