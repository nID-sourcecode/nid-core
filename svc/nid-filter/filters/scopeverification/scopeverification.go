// Package scopeverification contains the scopeverification filter logic
package scopeverification

import (
	"context"
	"net/url"
	"strings"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/nID-sourcecode/nid-core/pkg/accessmodel"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	grpcerrors "github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/errors"
	"github.com/nID-sourcecode/nid-core/svc/nid-filter/filters/scopeverification/verification"
)

const bearerScheme = "bearer "

// Error definitions
var (
	ErrRequestDoesNotMatchScopes  = errors.New("request does not match scopes")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrHeaderNotSpecified         = errors.New("header not specified")

	errAuthorizationHeaderNotFound = errors.New("authorization header not found")
)

// New returns a new instance of Filter struct
func New() *Filter {
	return &Filter{
		verifiers: map[string]verification.Verifier{
			"gql":  &verification.GQLVerifier{},
			"rest": &verification.RESTVerifier{},
		},
	}
}

// Filter is responsible for checking whether HTTP requests match the scopes provided in the JWT
type Filter struct {
	verifiers map[string]verification.Verifier
}

type verifyParamsData struct {
	authHeader string
	path       string
	method     string
}

// Check runs scopeverification on the request.
func (s *Filter) Check(_ context.Context, request *authv3.CheckRequest) error {
	headers := request.GetAttributes().GetRequest().GetHttp().GetHeaders()

	authHeader, ok := headers["authorization"]
	if !ok {
		return errAuthorizationHeaderNotFound
	}

	path, ok := headers[":path"]

	if !ok {
		return errors.Errorf("%w: :path", ErrHeaderNotSpecified)
	}

	method, ok := headers[":method"]
	if !ok {
		return errors.Errorf("%w: :method", ErrHeaderNotSpecified)
	}

	verifyParams := &verifyParamsData{
		authHeader: authHeader,
		path:       path,
		method:     method,
	}
	body := request.GetAttributes().GetRequest().GetHttp().GetBody()

	return s.verify(body, verifyParams)
}

// Name returns the filter name
func (s *Filter) Name() string {
	return "scopeverification"
}

const amountOfCharactersForShortening = 3

func (s *Filter) verify(body string, verifyParams *verifyParamsData) error {
	authHeader := verifyParams.authHeader
	if !strings.HasPrefix(strings.ToLower(authHeader), bearerScheme) { // FIXME abstract away token reading logic since it happens in many filters
		shortAuthHeader := authHeader
		if len(authHeader) > amountOfCharactersForShortening {
			shortAuthHeader = authHeader[:2] + "..."
		}
		return errors.Errorf(`%w: "%s" does not adhere to the Bearer scheme`, ErrInvalidAuthorizationHeader, shortAuthHeader)
	}
	token := authHeader[len(bearerScheme):]

	scopes, err := accessmodel.ExtractScopesFromJWT(token)
	if err != nil {
		return errors.Wrap(err, "extracting scopes from jwt")
	}

	path := verifyParams.path
	uri, err := url.Parse(path)
	if err != nil {
		return grpcerrors.ErrInvalidArgument(errors.Wrap(err, "parsing path").Error())
	}

	request := verification.Request{
		Scopes: scopes,
		Method: verifyParams.method,
		Path:   uri.Path,
		Query:  uri.Query(),
		Body:   body,
	}

	for name, verifier := range s.verifiers {
		err := verifier.Verify(&request)
		if err == nil {
			// Access granted
			return nil
		}
		// We need if statements since wrapper errors cant be switched
		// nolint: gocritic
		if errors.Is(err, verification.ErrBadRequest) {
			return errors.Wrap(err, "verifying request")
		} else if errors.Is(err, verification.ErrNotValid) {
			log.WithError(err).Error("tried verifiying the request")
			continue
		} else {
			return errors.Wrapf(err, "verifying using %s verifier", name)
		}
	}

	return ErrRequestDoesNotMatchScopes
}
