package scopeverification

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3alpha"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/accessmodel"
	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	gql "lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
)

type ScopeVerificationServerTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	filterInitializer *FilterInitializer
	filter            filter.Filter

	queryNamesAndBankAccount gql.Request
	queryNames               gql.Request
	queryInvalid             gql.Request
	accessModelNames         map[string]interface{}
	accessModelBankAccounts  map[string]interface{}
}

func (s *ScopeVerificationServerTestSuite) SetupTest() {
	filter, err := s.filterInitializer.NewFilter()
	s.Require().NoError(err)
	s.filter = filter
}

func (s *ScopeVerificationServerTestSuite) SetupSuite() {
	s.filterInitializer = NewScopeVerificationFilterInitializer()

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

	body := []byte("{\"some invalid json\":\"indeed}")

	ctx := context.TODO()

	res, err := s.filter.OnHTTPRequest(ctx, body, headers)

	expectedRes := &filter.ProcessingResponse{
		NewHeaders: nil,
		NewBody:    nil,
		ImmediateResponse: &ext_proc_pb.ImmediateResponse{
			Status: &envoy_type_v3.HttpStatus{Code: envoy_type_v3.StatusCode_BadRequest},
			Headers: &ext_proc_pb.HeaderMutation{
				SetHeaders: []*envoy_core_v3.HeaderValueOption{
					{
						Header: &envoy_core_v3.HeaderValue{
							Key:   "Content-Type",
							Value: "application/json",
						},
					},
				},
				RemoveHeaders: nil,
			},
			Body: `{
    "errors": [
        {"message": "verifying request: bad request: parsing body: unexpected end of JSON input"}
    ]
}`,
			GrpcStatus: nil,
			Details:    "",
		},
	}

	s.Require().NoError(err)
	s.Require().Equal(expectedRes, res)
}

func (s *ScopeVerificationServerTestSuite) TestShouldReturnErrBadRequestOnInvalidAuthHeader() {
	headers := map[string]string{
		":path":         "/gql",
		":method":       "POST",
		"authorization": "really bad auth header",
	}

	body := []byte(s.marshal(gql.Request{
		Query: `{
					usersfilter: {pseudonym: {eq: "$$nid:subject$$"}}) {
						firstName
						lastName
					}
				}`,
	}))

	ctx := context.TODO()
	res, err := s.filter.OnHTTPRequest(ctx, body, headers)

	expectedRes := &filter.ProcessingResponse{
		NewHeaders: nil,
		NewBody:    nil,
		ImmediateResponse: &ext_proc_pb.ImmediateResponse{
			Status: &envoy_type_v3.HttpStatus{Code: envoy_type_v3.StatusCode_BadRequest},
			Headers: &ext_proc_pb.HeaderMutation{
				SetHeaders: []*envoy_core_v3.HeaderValueOption{
					{
						Header: &envoy_core_v3.HeaderValue{
							Key:   "Content-Type",
							Value: "application/json",
						},
					},
				},
				RemoveHeaders: nil,
			},
			Body: `{
    "errors": [
        {"message": "invalid authorization header: "re..." does not adhere to the Bearer scheme"}
    ]
}`,
			GrpcStatus: nil,
			Details:    "",
		},
	}

	s.Require().NoError(err)
	s.Require().Equal(expectedRes, res)
}

func (s *ScopeVerificationServerTestSuite) TestVerifyCombinations() {
	tests := []*struct {
		name                 string
		scopes               map[string]interface{}
		body                 []byte
		method               string
		path                 string
		errorExpected        bool
		expectedErrorCode    envoy_type_v3.StatusCode
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
			body:          nil,
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
			body:          nil,
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
			body:                 nil,
			method:               "GET",
			path:                 "/v1/recipes?fruit=blueberry",
			errorExpected:        true,
			expectedErrorCode:    envoy_type_v3.StatusCode_Forbidden,
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
			body:                 []byte(""),
			method:               "POST",
			path:                 "/gql",
			errorExpected:        true,
			expectedErrorCode:    envoy_type_v3.StatusCode_BadRequest,
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
			res, err := s.filter.OnHTTPRequest(ctx, test.body, headers)

			s.Require().NoError(err)

			if test.errorExpected {
				expectedRes := &filter.ProcessingResponse{
					NewHeaders: nil,
					NewBody:    nil,
					ImmediateResponse: &ext_proc_pb.ImmediateResponse{
						Status: &envoy_type_v3.HttpStatus{Code: test.expectedErrorCode},
						Headers: &ext_proc_pb.HeaderMutation{
							SetHeaders: []*envoy_core_v3.HeaderValueOption{
								{
									Header: &envoy_core_v3.HeaderValue{
										Key:   "Content-Type",
										Value: "application/json",
									},
								},
							},
							RemoveHeaders: nil,
						},
						Body: fmt.Sprintf(`{
    "errors": [
        {"message": "%s"}
    ]
}`, test.expectedErrorMessage),
						GrpcStatus: nil,
						Details:    "",
					},
				}
				s.Require().Equal(expectedRes, res)
			} else {
				s.Nil(res)
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
		expectedErrorCode    envoy_type_v3.StatusCode
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
			expectedErrorCode:    envoy_type_v3.StatusCode_BadRequest,
			expectedErrorMessage: "verifying request: bad request: parsing GQL query: input:2: Expected Name, found {",
		},
		{
			name: "ShouldReturnErrNotValidOnFieldOutsideScope",
			scopes: map[string]interface{}{
				"names": s.accessModelBankAccounts,
			},
			req:                  s.queryNamesAndBankAccount,
			errorExpected:        true,
			expectedErrorCode:    envoy_type_v3.StatusCode_Forbidden,
			expectedErrorMessage: "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidOnModelOutsideScope",
			scopes: map[string]interface{}{
				"names": s.accessModelNames,
			},
			req:                  s.queryNamesAndBankAccount,
			errorExpected:        true,
			expectedErrorCode:    envoy_type_v3.StatusCode_Forbidden,
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
			expectedErrorCode:    envoy_type_v3.StatusCode_Forbidden,
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
			expectedErrorCode:    envoy_type_v3.StatusCode_Forbidden,
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
			body := []byte(s.marshal(test.req))

			ctx := context.TODO()
			res, err := s.filter.OnHTTPRequest(ctx, body, headers)

			s.Require().NoError(err)

			if test.errorExpected {
				expectedRes := &filter.ProcessingResponse{
					NewHeaders: nil,
					NewBody:    nil,
					ImmediateResponse: &ext_proc_pb.ImmediateResponse{
						Status: &envoy_type_v3.HttpStatus{Code: test.expectedErrorCode},
						Headers: &ext_proc_pb.HeaderMutation{
							SetHeaders: []*envoy_core_v3.HeaderValueOption{
								{
									Header: &envoy_core_v3.HeaderValue{
										Key:   "Content-Type",
										Value: "application/json",
									},
								},
							},
							RemoveHeaders: nil,
						},
						Body: fmt.Sprintf(`{
    "errors": [
        {"message": "%s"}
    ]
}`, test.expectedErrorMessage),
						GrpcStatus: nil,
						Details:    "",
					},
				}
				s.Require().Equal(expectedRes, res)
			} else {
				s.Nil(res)
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
			res, err := s.filter.OnHTTPRequest(ctx, nil, headers)
			s.Require().NoError(err)

			if test.errorExpected {
				expectedRes := &filter.ProcessingResponse{
					NewHeaders: nil,
					NewBody:    nil,
					ImmediateResponse: &ext_proc_pb.ImmediateResponse{
						Status: &envoy_type_v3.HttpStatus{Code: test.expectedErrorCode},
						Headers: &ext_proc_pb.HeaderMutation{
							SetHeaders: []*envoy_core_v3.HeaderValueOption{
								{
									Header: &envoy_core_v3.HeaderValue{
										Key:   "Content-Type",
										Value: "application/json",
									},
								},
							},
							RemoveHeaders: nil,
						},
						Body: fmt.Sprintf(`{
    "errors": [
        {"message": "%s"}
    ]
}`, test.expectedErrorMessage),
						GrpcStatus: nil,
						Details:    "",
					},
				}
				s.Require().Equal(expectedRes, res)
			} else {
				s.Nil(res)
			}
		})
	}
}
