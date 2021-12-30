package main

import (
	"context"
	"net/url"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"

	"lab.weave.nl/nid/nid-core/pkg/accessmodel"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/svc/nid-filter/filters/scopeverification/verification"
	"lab.weave.nl/nid/nid-core/svc/scopeverification/proto"
)

// ScopeVerificationServer verifies that a request matches the scopes in its auth header
type ScopeVerificationServer struct {
	verifiers map[string]verification.Verifier
}

// NewScopeVerificationServer creates a new scope verification server with default verifiers
func NewScopeVerificationServer() *ScopeVerificationServer {
	return &ScopeVerificationServer{
		verifiers: map[string]verification.Verifier{
			"gql":  &verification.GQLVerifier{},
			"rest": &verification.RESTVerifier{},
		},
	}
}

const bearerScheme = "bearer "

var (
	// ErrRequestDoesNotMatchScopes is returned if all verifiers deny the request
	ErrRequestDoesNotMatchScopes = errors.New("request does not match scopes")
	// ErrInvalidAuthorizationHeader is returned on a malformed authorization header
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
)

// Verify verifies that request matches the scopes in its JWT
func (s ScopeVerificationServer) Verify(ctx context.Context, req *proto.VerifyRequest) (*empty.Empty, error) {
	logger := log.Extract(ctx)

	authHeader := req.GetAuthHeader()

	// Note: this check could be done by a proto validator, but we need clear error messages since this service is called from a filter
	if !strings.HasPrefix(strings.ToLower(authHeader), bearerScheme) {
		return nil, grpcerrors.ErrInvalidArgument(
			errors.Errorf(`%w: "%s..." does not adhere to the Bearer scheme`, ErrInvalidAuthorizationHeader, authHeader[:2]))
	}
	token := authHeader[len(bearerScheme):]

	scopes, err := accessmodel.ExtractScopesFromJWT(token)
	if err != nil {
		return nil, grpcerrors.ErrInvalidArgument(errors.Wrap(err, "extracting scopes from jwt"))
	}

	uri, err := url.Parse(req.GetPath())
	if err != nil {
		return nil, grpcerrors.ErrInvalidArgument(errors.Wrap(err, "parsing path"))
	}

	request := verification.Request{
		Scopes: scopes,
		Method: req.GetMethod(),
		Path:   uri.Path,
		Query:  uri.Query(),
		Body:   req.GetBody(),
	}

	for name, verifier := range s.verifiers {
		err := verifier.Verify(&request)
		if err == nil {
			// Access granted
			return &empty.Empty{}, nil
		}
		// We need if statements since wrapper errors cant be switched
		// nolint: gocritic
		if errors.Is(err, verification.ErrBadRequest) {
			return nil, grpcerrors.ErrInvalidArgument(errors.Wrap(err, "verifying request"))
		} else if errors.Is(err, verification.ErrNotValid) {
			continue
		} else {
			logger.WithError(err).Errorf("verifying using %s verifier", name)

			return nil, grpcerrors.ErrInternalServer()
		}
	}

	return nil, grpcerrors.ErrPermissionDenied(ErrRequestDoesNotMatchScopes)
}
