package verification

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/accessmodel"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

type GQLVerificationTestSuite struct {
	suite.Suite
	accessModelNames *accessmodel.AccessModel
}

func (s *GQLVerificationTestSuite) SetupSuite() {
	s.accessModelNames = s.parseAccessModel(
		`{
			"t": "GQL",
			"p": "/gql",
			"m": {
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"f": ["firstName", "lastName"],
					"p": {
						"filter": {
							"pseudonym": {
								"eq": "$$nid:subject$$"
							}
						}
					}
				}
			}
		}`)
}

func TestGQLVerificationTestSuite(t *testing.T) {
	suite.Run(t, &GQLVerificationTestSuite{})
}

func (s *GQLVerificationTestSuite) marshal(body gqlRequestBody) string {
	bytes, err := json.Marshal(body)
	s.Require().NoError(err, "error marshalling test body")

	return string(bytes)
}

func (s *GQLVerificationTestSuite) toURIQuery(body gqlRequestBody) url.Values {
	variablesJSONBytes, err := json.Marshal(body.Variables)
	s.Require().NoError(err, "error marshalling test variables")

	return url.Values{
		"query":     []string{body.Query},
		"variables": []string{string(variablesJSONBytes)},
	}
}

func (s *GQLVerificationTestSuite) parseAccessModel(modelJSON string) *accessmodel.AccessModel {
	model := &accessmodel.AccessModel{}
	err := json.Unmarshal([]byte(modelJSON), model)
	s.Require().NoError(err, "error unmarshalling test access model")

	return model
}

func (s *GQLVerificationTestSuite) TestShouldReturnErrBadRequestOnInvalidJSON() {
	req := Request{
		Scopes: map[string]*accessmodel.AccessModel{
			"somemodel": s.accessModelNames,
		},
		Method: "POST",
		Path:   "/gql",
		Query:  url.Values{},
		Body:   "{\"some invalid json\":\"indeed}",
	}

	verifier := GQLVerifier{}
	err := verifier.Verify(&req)

	s.Error(err)
	s.True(errors.Is(err, ErrBadRequest), err)
	s.Equal("bad request: parsing body: unexpected end of JSON input", err.Error())
}

func (s *GQLVerificationTestSuite) TestShouldReturnErrBadRequestOnInvalidVariablesJSON() {
	req := Request{
		Scopes: map[string]*accessmodel.AccessModel{
			"somemodel": s.accessModelNames,
		},
		Method: "GET",
		Path:   "/gql",
		Query: url.Values{
			"query":     []string{"something"},
			"variables": []string{"{\"some\":\"bad json}"},
		},
		Body: "",
	}

	verifier := GQLVerifier{}
	err := verifier.Verify(&req)

	s.Error(err)
	s.True(errors.Is(err, ErrBadRequest), err)
	s.Equal("bad request: parsing variables: unexpected end of JSON input", err.Error())
}

func (s *GQLVerificationTestSuite) TestShouldReturnErrBadRequestOnUnsupportedMethod() {
	req := Request{
		Scopes: map[string]*accessmodel.AccessModel{
			"somemodel": s.accessModelNames,
		},
		Method: "PUT",
		Path:   "/gql",
		Query:  url.Values{},
		Body:   "",
	}

	verifier := GQLVerifier{}
	err := verifier.Verify(&req)

	s.Error(err)
	s.True(errors.Is(err, ErrBadRequest), err)
	s.Equal("bad request: method PUT is not supported, use POST or GET", err.Error())
}

func (s *GQLVerificationTestSuite) TestVerifyGQLRequest() {
	queryNamesAndBankAccount := gqlRequestBody{
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

	queryNames := gqlRequestBody{
		Query: `{
					users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
						firstName
						lastName
					}
				}`,
	}

	queryInvalid := gqlRequestBody{
		Query: `{
					usersfilter: {pseudonym: {eq: "$$nid:subject$$"}}) {
						firstName
						lastName
					}
				}`,
	}

	queryIntrospect := gqlRequestBody{
		Query: `{
					__schema {
						types {
							name
						}
					}
				}`,
	}

	accessModelNamesOfAllUsers := s.parseAccessModel(
		`{
			"t": "GQL",
			"p": "/gql",
			"m": {
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"f": ["firstName", "lastName"]
				}
			}
		}`)

	accessModelWithMapParameter := s.parseAccessModel(
		`{
			"t": "GQL",
			"p": "/gql",
			"m": {
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"f": ["firstName", "lastName"],
					"p": {
						"something": {
							"that": true,
							"need": 2,
							"not": "five",
							"beOrdered": null
						}
					}
				}
			}
		}`)

	accessModelBankAccounts := s.parseAccessModel(
		`{
			"t": "GQL",
			"p": "/gql",
			"m": {
				"r": {
					"m": {
						"users": "#U"
					}
				},
				"U": {
					"m": {
						"bankAccounts": "#B"
					},
					"f": [],
					"p": {
						"filter": {
							"pseudonym": {
								"eq": "$$nid:subject$$"
							}
						}
					}
				},
				"B": {
					"m": {
						"savingsAccounts": "#S"
					},
					"f": [
						"accountNumber",
						"amount"
					] 
				},
				"S": {
					"f": [
						"amount",
						"name"
					]
				}
			}
		}`)

	tests := []*struct {
		name                  string
		scopes                map[string]*accessmodel.AccessModel
		req                   gqlRequestBody
		errorExpected         bool
		expectedErrorIdentity error
		expectedErrorMessage  string
	}{
		{
			name: "ShouldReturnErrNotValidOnNoGQLScopes",
			scopes: map[string]*accessmodel.AccessModel{
				"somemodel": {
					Type:            accessmodel.RESTType,
					RESTAccessModel: accessmodel.RESTAccessModel{},
				},
			},
			req:                   queryNamesAndBankAccount,
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes: no gql scopes found",
		},
		{
			name: "ShouldReturnErrNotValidOnIntrospectQueryWithNoGQLScopes",
			scopes: map[string]*accessmodel.AccessModel{
				"somemodel": {
					Type:            accessmodel.RESTType,
					RESTAccessModel: accessmodel.RESTAccessModel{},
				},
			},
			req:                   queryNamesAndBankAccount,
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes: no gql scopes found",
		},
		{
			name: "ShouldReturnNoErrorOnQueryInSingleScope",
			scopes: map[string]*accessmodel.AccessModel{
				"somemodel": s.accessModelNames,
			},
			req:           queryNames,
			errorExpected: false,
		},
		{
			name: "ShouldReturnNoErrorOnQueryInOneOfMultipleScopes",
			scopes: map[string]*accessmodel.AccessModel{
				"names":       s.accessModelNames,
				"bankaccount": accessModelBankAccounts,
			},
			req:           queryNames,
			errorExpected: false,
		},
		{
			name: "ShouldReturnNoErrorOnQueryInScopesCombined",
			scopes: map[string]*accessmodel.AccessModel{
				"names":       s.accessModelNames,
				"bankaccount": accessModelBankAccounts,
			},
			req:           queryNamesAndBankAccount,
			errorExpected: false,
		},
		{
			name: "ShouldReturnErrBadRequestOnInvalidGQLQuery",
			scopes: map[string]*accessmodel.AccessModel{
				"names":       s.accessModelNames,
				"bankaccount": accessModelBankAccounts,
			},
			req:                   queryInvalid,
			errorExpected:         true,
			expectedErrorIdentity: ErrBadRequest,
			expectedErrorMessage:  "bad request: parsing GQL query: input:2: Expected Name, found {",
		},
		{
			name: "ShouldReturnNoErrorOnIntrospectionWithOneScope",
			scopes: map[string]*accessmodel.AccessModel{
				"names": s.accessModelNames,
			},
			req:           queryIntrospect,
			errorExpected: false,
		},
		{
			name: "ShouldReturnErrNotValidOnFieldOutsideScope",
			scopes: map[string]*accessmodel.AccessModel{
				"names": accessModelBankAccounts,
			},
			req:                   queryNamesAndBankAccount,
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidOnModelOutsideScope",
			scopes: map[string]*accessmodel.AccessModel{
				"names": s.accessModelNames,
			},
			req:                   queryNamesAndBankAccount,
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidOnNotMatchingParameter",
			scopes: map[string]*accessmodel.AccessModel{
				"names":        s.accessModelNames,
				"bankaccounts": accessModelBankAccounts,
			},
			req: gqlRequestBody{
				Query: `
					{
						users(filter: {pseudonym: {eq: "$$nid:other$$"}}) {
							firstName
							lastName
						}
					}
				`,
			},
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidOnNotMatchingParameter",
			scopes: map[string]*accessmodel.AccessModel{
				"names":        s.accessModelNames,
				"bankaccounts": accessModelBankAccounts,
			},
			req: gqlRequestBody{
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
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes",
		},
		{
			name: "ShouldReturnNoErrorOnParameterWithVariableInsideScope",
			scopes: map[string]*accessmodel.AccessModel{
				"names":        s.accessModelNames,
				"bankaccounts": accessModelBankAccounts,
			},
			req: gqlRequestBody{
				Query: `
					query getNames($pseudo: String){
						users(filter: {pseudonym: {eq: $pseudo}}) {
							firstName
							lastName
						}
					}
				`,
				Variables: map[string]interface{}{
					"pseudo": "$$nid:subject$$",
				},
			},
			errorExpected: false,
		},
		{
			name: "ShouldReturnErrNotValidIfVariableIsNotSpecifiedButShouldBe",
			scopes: map[string]*accessmodel.AccessModel{
				"names":        s.accessModelNames,
				"bankaccounts": accessModelBankAccounts,
			},
			req: gqlRequestBody{
				Query: `
					query getNames($pseudo: String){
						users(filter: {pseudonym: {eq: $pseudo}}) {
							firstName
							lastName
						}
					}
				`,
			},
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidIfParameterIsNotSpecifiedButShouldBe",
			scopes: map[string]*accessmodel.AccessModel{
				"names":        s.accessModelNames,
				"bankaccounts": accessModelBankAccounts,
			},
			req: gqlRequestBody{
				Query: `
					query getNames($pseudo: String){
						users {
							firstName
							lastName
						}
					}
				`,
			},
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidIfSpecifiedParameterIsNotInScope",
			scopes: map[string]*accessmodel.AccessModel{
				"names": accessModelNamesOfAllUsers,
			},
			req: gqlRequestBody{
				Query: `
					{
						users(filter: {pseudonym: {eq: "$$nid:subject"}}) {
							firstName
							lastName
						}
					}
				`,
			},
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes",
		},
		{
			name: "ShouldReturnErrNotValidIfSpecifiedParameterIsNotInScopeTwoParameters",
			scopes: map[string]*accessmodel.AccessModel{
				"names": s.accessModelNames,
			},
			req: gqlRequestBody{
				Query: `
					{
						users(someOther: "something", filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
							firstName
							lastName
						}
					}
				`,
			},
			errorExpected:         true,
			expectedErrorIdentity: ErrNotValid,
			expectedErrorMessage:  "request does not match scopes",
		},
		{
			name: "ShouldReturnNoErrorOnMapParameterWithDifferentOrder",
			scopes: map[string]*accessmodel.AccessModel{
				"somemodel": accessModelWithMapParameter,
			},
			req: gqlRequestBody{
				Query: `{
					users(something: $map) {
						firstName
						lastName
					}
				}`,
				Variables: map[string]interface{}{
					"map": map[string]interface{}{
						"not":       "five",
						"need":      2,
						"beOrdered": nil,
						"that":      true,
					},
				},
			},
			errorExpected: false,
		},
	}

	verifier := GQLVerifier{}
	for _, test := range tests {
		s.Run(fmt.Sprintf("%s_POST", test.name), func() {
			req := Request{
				Scopes: test.scopes,
				Method: "POST",
				Path:   "/gql",
				Query:  url.Values{},
				Body:   s.marshal(test.req),
			}
			err := verifier.Verify(&req)
			if test.errorExpected {
				s.Require().Error(err)
				s.Require().NotNil(err)
				s.True(errors.Is(err, test.expectedErrorIdentity), err)
				s.Equal(test.expectedErrorMessage, err.Error())
			} else {
				s.NoError(err)
			}
		})

		s.Run(fmt.Sprintf("%s_GET", test.name), func() {
			req := Request{
				Scopes: test.scopes,
				Method: "GET",
				Path:   "/gql",
				Query:  s.toURIQuery(test.req),
			}
			err := verifier.Verify(&req)
			if test.errorExpected {
				s.Require().Error(err)
				s.Require().NotNil(err)
				s.True(errors.Is(err, test.expectedErrorIdentity), err)
				s.Equal(test.expectedErrorMessage, err.Error())
			} else {
				s.NoError(err)
			}
		})
	}
}
