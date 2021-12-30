// Package authswap contains the autopseudo filter logic
package authswap

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	authpb "lab.weave.nl/nid/nid-core/svc/auth/proto"
)

// Error definitions
var (
	ErrHeaderNotSpecified = errors.New("header not specified")
)

// FilterInitializer creates new authswap filters
type FilterInitializer struct {
	authClient authpb.AuthClient
}

// Name returns the filter name
func (s *FilterInitializer) Name() string {
	return "authswap"
}

// NewFilter creates a new filter
func (s *FilterInitializer) NewFilter() (filter.Filter, error) {
	return &Filter{authClient: s.authClient}, nil
}

// NewFilterInitializer creates a new authswap filter initializer
func NewFilterInitializer(authClient authpb.AuthClient) *FilterInitializer {
	return &FilterInitializer{authClient: authClient}
}

// Filter contains the auditlog filter logic
type Filter struct {
	filter.DefaultFilter
	authClient authpb.AuthClient
}

// OnHTTPRequest handles an http request
func (s *Filter) OnHTTPRequest(ctx context.Context, body []byte, headers map[string]string) (*filter.ProcessingResponse, error) {
	authHeader, ok := headers["authorization"]
	if !ok {
		return nil, nil
	}

	requestProtocol, ok := headers["x-forwarded-proto"]
	if !ok {
		return nil, errors.Errorf("%w: x-forwarded-proto", ErrHeaderNotSpecified)
	}
	requestAuthority, ok := headers[":authority"]
	if !ok {
		return nil, errors.Errorf("%w: :authority", ErrHeaderNotSpecified)
	}
	requestPath, ok := headers[":path"]
	if !ok {
		return nil, errors.Errorf("%w: :path", ErrHeaderNotSpecified)
	}
	audience := fmt.Sprintf("%s://%s%s", requestProtocol, requestAuthority, requestPath)
	audienceURI, err := url.Parse(audience)
	if err != nil {
		return nil, errors.Wrap(err, "parsing audience to URI")
	}

	token := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(authHeader, "Bearer"), "bearer"))

	res, err := s.authClient.SwapToken(ctx, &authpb.SwapTokenRequest{
		CurrentToken: token,
		Query:        "stub",
		Audience:     fmt.Sprintf("%s://%s%s", audienceURI.Scheme, audienceURI.Host, audienceURI.Path),
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting swap token from auth")
	}

	headers["authorization"] = res.GetTokenType() + " " + res.GetAccessToken()

	return &filter.ProcessingResponse{
		NewHeaders: headers,
	}, nil
}

// Name returns the filter name
func (s *Filter) Name() string {
	return "authswap"
}
