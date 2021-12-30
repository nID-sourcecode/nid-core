// Package auditlog contains the auditlog filter logic. The audit log filter logs requests with JWT claims and response code.
package auditlog

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/dvsekhvalnov/jose2go/base64url"

	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
)

// Error definitions
var (
	ErrHeaderNotSpecified = errors.New("header not specified")
)

// FilterInitializer is responseible for creating new filters
type FilterInitializer struct {
	logger log.LoggerUtility
}

// Name returns the filter name
func (s *FilterInitializer) Name() string {
	return "auditlog"
}

// NewFilter creates a new filter
func (s *FilterInitializer) NewFilter() (filter.Filter, error) {
	return &Filter{logger: s.logger}, nil
}

// NewFilterInitializer creates a new filter initializer
func NewFilterInitializer(logger log.LoggerUtility) *FilterInitializer {
	return &FilterInitializer{logger: logger}
}

// Filter is responsible for processing a single HTTP request and response
type Filter struct {
	filter.DefaultFilter
	logger    log.LoggerUtility
	requestID string
}

// OnHTTPRequest processes a single HTTP request
func (f *Filter) OnHTTPRequest(ctx context.Context, body []byte, headers map[string]string) (*filter.ProcessingResponse, error) {
	var token string
	authHeader, ok := headers["authorization"]
	if ok {
		tokenSplit := strings.Split(authHeader, ".")

		// Do not send actual token to auditlog service, this creates a risk if the auditlogs are compromised
		if len(tokenSplit) > 1 {
			token = tokenSplit[1]
		}
	} // else token remains empty

	claims := decodeToken(token)

	host, ok := headers[":authority"]
	if !ok {
		return nil, errors.Errorf("%w: :authority", ErrHeaderNotSpecified)
	}
	path, ok := headers[":path"]
	if !ok {
		return nil, errors.Errorf("%w: :path", ErrHeaderNotSpecified)
	}
	method, ok := headers[":method"]
	if !ok {
		return nil, errors.Errorf("%w: :method", ErrHeaderNotSpecified)
	}
	f.requestID, ok = headers["x-request-id"]
	if !ok {
		return nil, errors.Errorf("%w: x-request-id", ErrHeaderNotSpecified)
	}

	url, err := url.PathUnescape(host + path)
	if err != nil {
		log.WithField("url", url).WithError(err).Error("unable to unescape path")
		return nil, grpcerrors.ErrInvalidArgument("url not pathunescapable")
	}

	f.logger.WithFields(log.Fields{
		// Just print the claims as an object, let the log formatter handle the formatting
		"token":       claims,
		"url":         url,
		"body":        string(body),
		"http_method": method,
		"request_id":  f.requestID,
	}).Info("received request")

	return nil, nil
}

// OnHTTPResponse processes a single HTTP response
func (f *Filter) OnHTTPResponse(ctx context.Context, body []byte, headers map[string]string) (*filter.ProcessingResponse, error) {
	log.Infof("response headers: %+v", headers)
	status, ok := headers[":status"]
	if !ok {
		return nil, errors.Errorf("%w: :status", ErrHeaderNotSpecified)
	}

	f.logger.WithFields(log.Fields{
		// Just print the claims as an object, let the log formatter handle the formatting
		"request_id":  f.requestID,
		"status_code": status,
	}).Info("received response")

	return nil, nil
}

// Name returns the filter name
func (f *Filter) Name() string {
	return "auditlog"
}

func decodeToken(token string) interface{} {
	if token == "" {
		return "no valid token"
	}
	var claims interface{}
	rawTokenBody, err := base64url.Decode(token)
	if err != nil {
		claims = "token not parsable"
	} else {
		decodedToken := map[string]interface{}{}

		err = json.Unmarshal(rawTokenBody, &decodedToken)
		if err != nil {
			log.WithField("raw_token_body", rawTokenBody).Info("unable to json unmarshal token body token")
			// Do not return an error, just log that it is invalid
			claims = "token not parsable"
		} else {
			// Only print claims, we do not want to print the actual signature since this would be a security risk
			claims = decodedToken
		}
	}
	return claims
}
