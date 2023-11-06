package integration

import (
	"context"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type AuthorizationRuleImpl struct{}

var (
	errCheckError        = errors.New("failure on check")
	errNoReturnTypeMatch = errors.New("failure on check")
)

const headerReturnCode = "Header-Return-Code"

type returnType int

const (
	OK returnType = iota
	Fail
)

func NewAuthzImpl() *AuthorizationRuleImpl {
	return &AuthorizationRuleImpl{}
}

func (a AuthorizationRuleImpl) Name() string {
	return "FakeAuthorizationRuleImpl"
}

func (a AuthorizationRuleImpl) Check(_ context.Context, request *authv3.CheckRequest) error {
	switch request.GetAttributes().GetRequest().GetHttp().GetHeaders()[headerReturnCode] {
	case string(rune(OK)):
		return nil
	case string(rune(Fail)):
		return errCheckError
	}

	return errNoReturnTypeMatch
}
