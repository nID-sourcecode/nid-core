package scopeverification

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	gql "github.com/nID-sourcecode/nid-core/pkg/gqlclient"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/nID-sourcecode/nid-core/svc/nid-filter/contract"

	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/accessmodel"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
)

type ScopeVerificationServerTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	filter contract.AuthorizationRule

	queryNamesAndBankAccount gql.Request
	queryNames               gql.Request
	queryInvalid             gql.Request
	accessModelNames         map[string]interface{}
	accessModelBankAccounts  map[string]interface{}
}

func (s *ScopeVerificationServerTestSuite) SetupTest() {
	s.filter = New()
}

func (s *ScopeVerificationServerTestSuite) SetupSuite() {
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

	s.accessModelNames = map[string]interface{}{
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

	s.accessModelBankAccounts = map[string]interface{}{
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

	return fmt.Sprintf(authHeaderFormat, base64.RawURLEncoding.EncodeToString(claimsJSON))
}

func (s *ScopeVerificationServerTestSuite) TestShouldReturnErrBadRequestOnInvalidJSON() {
	headers := map[string]string{
		":path":   "/gql",
		":method": "POST",
		"authorization": s.toAuthHeader(map[string]interface{}{
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
	}

	body := "{\"some invalid json\":\"indeed}"

	ctx := context.TODO()

	err := s.filter.Check(ctx, returnAuthV3CheckRequest(body, headers))

	s.Error(err)
	s.ErrorContains(err, "verifying request: bad request: parsing body: unexpected end of JSON input")
}

func (s *ScopeVerificationServerTestSuite) TestShouldReturnErrBadRequestOnInvalidAuthHeader() {
	headers := map[string]string{
		":path":         "/gql",
		":method":       "POST",
		"authorization": "really bad auth header",
	}

	body := s.marshal(gql.Request{
		Query: `{
					usersfilter: {pseudonym: {eq: "$$nid:subject$$"}}) {
						firstName
						lastName
					}
				}`,
	})

	ctx := context.TODO()
	err := s.filter.Check(ctx, returnAuthV3CheckRequest(body, headers))

	s.Error(err)
	s.ErrorContains(err, `invalid authorization header: "re..." does not adhere to the Bearer scheme`)
}

func (s *ScopeVerificationServerTestSuite) TestVerifyCombinations() {
	tests := []*struct {
		name                 string
		scopes               map[string]interface{}
		body                 string
		method               string
		path                 string
		errorExpected        bool
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
			expectedErrorMessage: "verifying request: bad request: parsing body: unexpected end of JSON input",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()

			headers := map[string]string{
				"authorization": s.toAuthHeader(test.scopes),
				":method":       test.method,
				":path":         test.path,
			}

			ctx := context.TODO()
			err := s.filter.Check(ctx, returnAuthV3CheckRequest(test.body, headers))

			if test.errorExpected {
				s.Error(err)
				s.ErrorContains(err, test.expectedErrorMessage)
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
			expectedErrorMessage: "verifying request: bad request: parsing GQL query: input:2: Expected Name, found {",
		},
		{
			name: "ShouldReturnErrNotValidOnFieldOutsideScope",
			scopes: map[string]interface{}{
				"names": s.accessModelBankAccounts,
			},
			req:                  s.queryNamesAndBankAccount,
			errorExpected:        true,
			expectedErrorMessage: "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidOnModelOutsideScope",
			scopes: map[string]interface{}{
				"names": s.accessModelNames,
			},
			req:                  s.queryNamesAndBankAccount,
			errorExpected:        true,
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
			expectedErrorMessage: "request does not match scopes",
		},
	}

	for _, test := range tests {
		s.Run(fmt.Sprintf("%s_POST", test.name), func() {
			s.SetupTest()

			headers := map[string]string{
				"authorization": s.toAuthHeader(test.scopes),
				":method":       "POST",
				":path":         "/gql",
			}
			body := s.marshal(test.req)

			ctx := context.TODO()
			err := s.filter.Check(ctx, returnAuthV3CheckRequest(body, headers))

			if test.errorExpected {
				s.Error(err)
				s.ErrorContains(err, test.expectedErrorMessage)
			} else {
				s.NoError(err)
			}
		})

		s.Run(fmt.Sprintf("%s_GET", test.name), func() {
			s.SetupTest()

			headers := map[string]string{
				"authorization": s.toAuthHeader(test.scopes),
				":method":       "GET",
				":path":         "/gql?" + s.toQuery(test.req),
			}

			ctx := context.TODO()
			err := s.filter.Check(ctx, returnAuthV3CheckRequest("", headers))

			if test.errorExpected {
				s.Error(err)
				s.ErrorContains(err, test.expectedErrorMessage)
			} else {
				s.NoError(err)
			}
		})
	}
}

func returnAuthV3CheckRequest(body string, headers map[string]string) *authv3.CheckRequest {
	return &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: headers,
					Body:    body,
				},
			},
		},
	}
}
