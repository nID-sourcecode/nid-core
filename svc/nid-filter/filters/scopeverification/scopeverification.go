// Package scopeverification contains the scopeverification filter logic
package scopeverification

import (
	"context"
	"net/url"
	"strings"

	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"

	"lab.weave.nl/nid/nid-core/pkg/accessmodel"
	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/svc/nid-filter/filters/scopeverification/verification"
	"lab.weave.nl/nid/nid-core/svc/nid-filter/filters/utils"
)

const bearerScheme = "bearer "

// Error definitions
var (
	ErrRequestDoesNotMatchScopes  = errors.New("request does not match scopes")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrHeaderNotSpecified         = errors.New("header not specified")
)

// FilterInitializer creates new scope verification filters
type FilterInitializer struct {
	verifiers map[string]verification.Verifier
}

// Name returns the filter name
func (s *FilterInitializer) Name() string {
	return "scopeverification"
}

// NewFilter creates a new filter
func (s *FilterInitializer) NewFilter() (filter.Filter, error) {
	return &Filter{verifiers: s.verifiers}, nil
}

// NewScopeVerificationFilterInitializer creates a new scope verification filter initializer with default verifiers
func NewScopeVerificationFilterInitializer() *FilterInitializer {
	return &FilterInitializer{
		verifiers: map[string]verification.Verifier{
			"gql":  &verification.GQLVerifier{},
			"rest": &verification.RESTVerifier{},
		},
	}
}

// Filter is responsible for checking whether HTTP requests match the scopes provided in the JWT
type Filter struct {
	filter.DefaultFilter
	verifiers  map[string]verification.Verifier
	authHeader string
	path       string
	method     string
}

// OnHTTPRequest handles an http request
func (s *Filter) OnHTTPRequest(ctx context.Context, body []byte, headers map[string]string) (*filter.ProcessingResponse, error) {
	authHeader, ok := headers["authorization"]
	if !ok {
		return utils.GraphqlError("authorization header not found or empty", envoy_type_v3.StatusCode_BadRequest), nil
	}

	s.authHeader = authHeader

	path, ok := headers[":path"]
	if !ok {
		return nil, errors.Errorf("%w: :path", ErrHeaderNotSpecified)
	}

	s.path = path

	method, ok := headers[":method"]
	if !ok {
		return nil, errors.Errorf("%w: :method", ErrHeaderNotSpecified)
	}

	s.method = method

	stringBody := ""
	if body != nil {
		stringBody = string(body)
	}

	return s.verify(stringBody)
}

// Name returns the filter name
func (s *Filter) Name() string {
	return "scopeverification"
}

const amountOfCharactersForShortening = 3

func (s *Filter) verify(body string) (*filter.ProcessingResponse, error) {
	if !strings.HasPrefix(strings.ToLower(s.authHeader), bearerScheme) { // FIXME abstract away token reading logic since it happens in many filters
		shortAuthHeader := s.authHeader
		if len(s.authHeader) > amountOfCharactersForShortening {
			shortAuthHeader = s.authHeader[:2] + "..."
		}
		return utils.GraphqlError(
			errors.Errorf(`%w: "%s" does not adhere to the Bearer scheme`, ErrInvalidAuthorizationHeader, shortAuthHeader).Error(),
			envoy_type_v3.StatusCode_BadRequest), nil
	}
	token := s.authHeader[len(bearerScheme):]

	scopes, err := accessmodel.ExtractScopesFromJWT(token)
	if err != nil {
		return utils.GraphqlError(errors.Wrap(err, "extracting scopes from jwt").Error(), envoy_type_v3.StatusCode_BadRequest), nil
	}

	uri, err := url.Parse(s.path)
	if err != nil {
		return nil, grpcerrors.ErrInvalidArgument(errors.Wrap(err, "parsing path"))
	}

	request := verification.Request{
		Scopes: scopes,
		Method: s.method,
		Path:   uri.Path,
		Query:  uri.Query(),
		Body:   body,
	}

	for name, verifier := range s.verifiers {
		err := verifier.Verify(&request)
		if err == nil {
			// Access granted
			return nil, nil
		}
		// We need if statements since wrapper errors cant be switched
		// nolint: gocritic
		if errors.Is(err, verification.ErrBadRequest) {
			return utils.GraphqlError(errors.Wrap(err, "verifying request").Error(), envoy_type_v3.StatusCode_BadRequest), nil
		} else if errors.Is(err, verification.ErrNotValid) {
			continue
		} else {
			return nil, errors.Wrapf(err, "verifying using %s verifier", name)
		}
	}

	return utils.GraphqlError(ErrRequestDoesNotMatchScopes.Error(), envoy_type_v3.StatusCode_Forbidden), nil
}
