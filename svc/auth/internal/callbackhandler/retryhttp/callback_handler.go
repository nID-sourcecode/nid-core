// Package http contains the http implementation of the callback handler.
package retryhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	httpClient "lab.weave.nl/nid/nid-core/pkg/utilities/http"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
)

const (
	authCodeParamKey  = "authorization_code"
	successStatusCode = http.StatusAccepted
	retryMaxWait      = 4 * time.Hour
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

	q := req.URL.Query()
	q.Add(authCodeParamKey, authorizationCode)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Info(err)
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

func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	return resp.StatusCode != http.StatusAccepted, nil
}
