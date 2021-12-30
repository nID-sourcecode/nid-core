package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"

	"lab.weave.nl/nid/nid-core/pkg/accessmodel"
	gql "lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	"lab.weave.nl/nid/nid-core/svc/scopeverification/proto"
)

type ScopeVerificationServerTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	server *ScopeVerificationServer

	queryNamesAndBankAccount gql.Request
	queryNames               gql.Request
	queryInvalid             gql.Request
	accessModelNames         map[string]interface{}
	accessModelBankAccounts  map[string]interface{}
}

func (s *ScopeVerificationServerTestSuite) SetupSuite() {
	s.server = NewScopeVerificationServer()

	s.queryNamesAndBankAccount = gql.Request{
		Query: `{
			users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
				firstName
				lastName
				bankAccounts {
					savingsAccounts {
						name
						amount
					}
					accountNumber
					amount
				}
			}
		}
					`,
	}

	s.queryNames = gql.Request{
		Query: `{
			users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
				firstName
				lastName
			}
		}`,
	}

	s.queryInvalid = gql.Request{
		Query: `{
			usersfilter: {pseudonym: {eq: "$$nid:subject$$"}}) {
				firstName
				lastName
			}
		}`,
	}

	s.accessModelNames =
		map[string]interface{}{
			"t": "GQL", "p": "/gql",
			"m": map[string]interface{}{
				"r": map[string]interface{}{
					"m": map[string]interface{}{
						"users": "#U",
					},
				},
				"U": map[string]interface{}{
					"f": []string{"firstName", "lastName"},
					"p": map[string]interface{}{
						"filter": map[string]interface{}{
							"pseudonym": map[string]interface{}{
								"eq": "$$nid:subject$$",
							},
						},
					},
				},
			},
		}

	s.accessModelBankAccounts =
		map[string]interface{}{
			"t": "GQL", "p": "/gql",
			"m": map[string]interface{}{
				"r": map[string]interface{}{
					"m": map[string]interface{}{
						"users": "#U",
					},
				},
				"U": map[string]interface{}{
					"m": map[string]interface{}{
						"bankAccounts": "#B",
					},
					"f": []string{},
					"p": map[string]interface{}{
						"filter": map[string]interface{}{
							"pseudonym": map[string]interface{}{
								"eq": "$$nid:subject$$",
							},
						},
					},
				},
				"B": map[string]interface{}{
					"m": map[string]interface{}{
						"savingsAccounts": "#S",
					},
					"f": []string{
						"accountNumber",
						"amount",
					},
				},
				"S": map[string]interface{}{
					"f": []string{
						"amount",
						"name",
					},
				},
			},
		}
}

func TestScopeVerificationServerTestSuite(t *testing.T) {
	suite.Run(t, &ScopeVerificationServerTestSuite{})
}

func (s *ScopeVerificationServerTestSuite) marshal(body gql.Request) string {
	bytes, err := json.Marshal(body)
	s.Require().NoError(err, "error marshalling test body")

	return string(bytes)
}

func (s *ScopeVerificationServerTestSuite) toQuery(body gql.Request) string {
	variablesJSONBytes, err := json.Marshal(body.Variables)
	s.Require().NoError(err, "error marshalling test variables")

	return url.Values{
		"query":     []string{body.Query},
		"variables": []string{string(variablesJSONBytes)},
	}.Encode()
}

const authHeaderFormat = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.%s.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func (s *ScopeVerificationServerTestSuite) toAuthHeader(scopes map[string]interface{}) string {
	claims := map[string]interface{}{
		"scopes": scopes,
	}
	claimsJSON, err := json.Marshal(claims)
	s.Require().NoError(err, "error marshalling test claims")

	return fmt.Sprintf(authHeaderFormat, base64.URLEncoding.EncodeToString(claimsJSON))
}

func (s *ScopeVerificationServerTestSuite) TestShouldReturnErrBadRequestOnInvalidJSON() {
	req := proto.VerifyRequest{
		AuthHeader: s.toAuthHeader(map[string]interface{}{
			"somemodel": map[string]interface{}{
				"t": accessmodel.GQLType,
				"p": "/gql",
				"m": map[string]interface{}{
					"r": map[string]interface{}{
						"f": []string{"id"},
					},
				},
			},
		}),
		Method: "POST",
		Path:   "/gql",
		Body:   "{\"some invalid json\":\"indeed}",
	}

	_, err := s.server.Verify(context.Background(), &req)

	s.Error(err)
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Equal(`rpc error: code = InvalidArgument desc = verifying request: bad request: parsing body: unexpected end of JSON input`, err.Error())
}

func (s *ScopeVerificationServerTestSuite) TestShouldReturnErrBadRequestOnInvalidAuthHeader() {
	req := proto.VerifyRequest{
		AuthHeader: "really bad auth header",
		Method:     "POST",
		Path:       "/gql",
		Body: s.marshal(gql.Request{
			Query: `{
					usersfilter: {pseudonym: {eq: "$$nid:subject$$"}}) {
						firstName
						lastName
					}
				}`,
		}),
	}

	_, err := s.server.Verify(context.Background(), &req)

	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Equal(`rpc error: code = InvalidArgument desc = invalid authorization header: "re..." does not adhere to the Bearer scheme`, err.Error())
}

func (s *ScopeVerificationServerTestSuite) TestVerifyCombinations() {
	tests := []*struct {
		name                 string
		scopes               map[string]interface{}
		body                 string
		method               string
		path                 string
		errorExpected        bool
		expectedErrorCode    codes.Code
		expectedErrorMessage string
	}{
		{
			name: "NoErrorOnMatchingRestScope",
			scopes: map[string]interface{}{
				"recipes:strawberry": map[string]interface{}{
					"t": "REST",
					"m": "GET",
					"q": map[string]string{
						"fruit": "strawberry",
					},
					"b": "",
					"p": "/v1/recipes",
				},
			},
			body:          "",
			method:        "GET",
			path:          "/v1/recipes?fruit=strawberry",
			errorExpected: false,
		},
		{
			name: "NoErrorOnMatchingRestScopeWithOtherPathGQLScope",
			scopes: map[string]interface{}{
				"recipes:strawberry": map[string]interface{}{
					"t": "REST",
					"m": "GET",
					"q": map[string]string{
						"fruit": "strawberry",
					},
					"b": "",
					"p": "/v1/recipes",
				},
				"data:names": s.accessModelNames,
			},
			body:          "",
			method:        "GET",
			path:          "/v1/recipes?fruit=strawberry",
			errorExpected: false,
		},
		{
			name: "ErrorOnNonMatchingRestScope",
			scopes: map[string]interface{}{
				"recipes:strawberry": map[string]interface{}{
					"t": "REST",
					"m": "GET",
					"q": map[string]string{
						"fruit": "strawberry",
					},
					"b": "",
					"p": "/v1/recipes",
				},
				"data:names": s.accessModelNames,
			},
			body:                 "",
			method:               "GET",
			path:                 "/v1/recipes?fruit=blueberry",
			errorExpected:        true,
			expectedErrorCode:    codes.PermissionDenied,
			expectedErrorMessage: "request does not match scopes",
		},
		{
			name: "ErrorOnNonGQLBodyButOneGQLScopeMatchingPath",
			scopes: map[string]interface{}{
				"recipes:strawberry": map[string]interface{}{
					"t": "REST",
					"m": "POST",
					"q": map[string]string{
						"fruit": "strawberry",
					},
					"b": "",
					"p": "/gql",
				},
				"data:names": s.accessModelNames,
			},
			body:                 "",
			method:               "POST",
			path:                 "/gql",
			errorExpected:        true,
			expectedErrorCode:    codes.InvalidArgument,
			expectedErrorMessage: "verifying request: bad request: parsing body: unexpected end of JSON input",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			req := proto.VerifyRequest{
				AuthHeader: s.toAuthHeader(test.scopes),
				Method:     test.method,
				Path:       test.path,
				Body:       test.body,
			}
			_, err := s.server.Verify(context.Background(), &req)
			if test.errorExpected {
				s.Require().Error(err)
				s.Require().NotNil(err)
				s.VerifyStatusError(err, test.expectedErrorCode)
				s.Contains(err.Error(), test.expectedErrorMessage)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *ScopeVerificationServerTestSuite) TestVerifyGQLRequest() {
	tests := []*struct {
		name                 string
		scopes               map[string]interface{}
		req                  gql.Request
		errorExpected        bool
		expectedErrorCode    codes.Code
		expectedErrorMessage string
	}{
		{
			name: "ShouldReturnNoErrorOnQueryInSingleScope",
			scopes: map[string]interface{}{
				"somemodel": s.accessModelNames,
			},
			req:           s.queryNames,
			errorExpected: false,
		},
		{
			name: "ShouldReturnNoErrorOnQueryInOneOfMultipleScopes",
			scopes: map[string]interface{}{
				"names":       s.accessModelNames,
				"bankaccount": s.accessModelBankAccounts,
			},
			req:           s.queryNames,
			errorExpected: false,
		},
		{
			name: "ShouldReturnNoErrorOnQueryInScopesCombined",
			scopes: map[string]interface{}{
				"names":       s.accessModelNames,
				"bankaccount": s.accessModelBankAccounts,
			},
			req:           s.queryNamesAndBankAccount,
			errorExpected: false,
		},
		{
			name: "ShouldReturnErrBadRequestOnInvalidGQLQuery",
			scopes: map[string]interface{}{
				"names":       s.accessModelNames,
				"bankaccount": s.accessModelBankAccounts,
			},
			req:                  s.queryInvalid,
			errorExpected:        true,
			expectedErrorCode:    codes.InvalidArgument,
			expectedErrorMessage: "bad request: parsing GQL query: input:2: Expected Name, found {",
		},
		{
			name: "ShouldReturnErrNotValidOnFieldOutsideScope",
			scopes: map[string]interface{}{
				"names": s.accessModelBankAccounts,
			},
			req:                  s.queryNamesAndBankAccount,
			errorExpected:        true,
			expectedErrorCode:    codes.PermissionDenied,
			expectedErrorMessage: "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidOnModelOutsideScope",
			scopes: map[string]interface{}{
				"names": s.accessModelNames,
			},
			req:                  s.queryNamesAndBankAccount,
			errorExpected:        true,
			expectedErrorCode:    codes.PermissionDenied,
			expectedErrorMessage: "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidOnNotMatchingParameter",
			scopes: map[string]interface{}{
				"names":        s.accessModelNames,
				"bankaccounts": s.accessModelBankAccounts,
			},
			req: gql.Request{
				Query: `
					{
						users(filter: {pseudonym: {eq: "$$nid:other$$"}}) {
							firstName
							lastName
						}
					}
				`,
			},
			errorExpected:        true,
			expectedErrorCode:    codes.PermissionDenied,
			expectedErrorMessage: "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidOnNotMatchingParameter",
			scopes: map[string]interface{}{
				"names":        s.accessModelNames,
				"bankaccounts": s.accessModelBankAccounts,
			},
			req: gql.Request{
				Query: `
					query getNames($pseudo: String){
						users(filter: {pseudonym: {eq: $pseudo}}) {
							firstName
							lastName
						}
					}
				`,
				Variables: map[string]interface{}{
					"pseudo": "$$nid:other$$",
				},
			},
			errorExpected:        true,
			expectedErrorCode:    codes.PermissionDenied,
			expectedErrorMessage: "request does not match scopes",
		},
	}

	for _, test := range tests {
		s.Run(fmt.Sprintf("%s_POST", test.name), func() {
			req := proto.VerifyRequest{
				AuthHeader: s.toAuthHeader(test.scopes),
				Method:     "POST",
				Path:       "/gql",
				Body:       s.marshal(test.req),
			}
			_, err := s.server.Verify(context.Background(), &req)
			if test.errorExpected {
				s.Require().Error(err)
				s.Require().NotNil(err)
				s.VerifyStatusError(err, test.expectedErrorCode)
				s.Contains(err.Error(), test.expectedErrorMessage)
			} else {
				s.NoError(err)
			}
		})

		s.Run(fmt.Sprintf("%s_GET", test.name), func() {
			req := proto.VerifyRequest{
				AuthHeader: s.toAuthHeader(test.scopes),
				Method:     "GET",
				Path:       "/gql?" + s.toQuery(test.req),
			}
			_, err := s.server.Verify(context.Background(), &req)
			if test.errorExpected {
				s.Require().Error(err)
				s.Require().NotNil(err)
				s.VerifyStatusError(err, test.expectedErrorCode)
				s.Contains(err.Error(), test.expectedErrorMessage)
			} else {
				s.NoError(err)
			}
		})
	}
}
