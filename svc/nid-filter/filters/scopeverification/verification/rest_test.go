package verification

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/accessmodel"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

type RESTVerifierTestSuite struct {
	suite.Suite
}

func TestRESTVerifierTestSuite(t *testing.T) {
	suite.Run(t, &RESTVerifierTestSuite{})
}

func (s *RESTVerifierTestSuite) TestRestVerifier() {
	tests := []*struct {
		name          string
		request       *Request
		errorExpected bool
	}{
		{
			name: "NoErrorGET",
			request: &Request{
				Scopes: map[string]*accessmodel.AccessModel{
					"foods": {
						Type: accessmodel.RESTType,
						RESTAccessModel: accessmodel.RESTAccessModel{
							Path: "/users/recipes",
							Query: map[string]string{
								"fruit":  "lemon",
								"greens": "broccoli",
							},
							Body:   "",
							Method: "GET",
						},
					},
				},
				Method: "GET",
				Path:   "/users/recipes",
				Query: url.Values{
					"fruit":  []string{"lemon"},
					"greens": []string{"broccoli"},
				},
				Body: "",
			},
			errorExpected: false,
		},
		{
			name: "NoErrorPOST",
			request: &Request{
				Scopes: map[string]*accessmodel.AccessModel{
					"foods": {
						Type: accessmodel.RESTType,
						RESTAccessModel: accessmodel.RESTAccessModel{
							Path:   "/users/recipes",
							Query:  map[string]string{},
							Body:   `{"fruit": "lemon", "greens": ["broccoli"]}`,
							Method: "POST",
						},
					},
				},
				Method: "POST",
				Path:   "/users/recipes",
				Query:  url.Values{},
				Body:   `{"fruit": "lemon", "greens": ["broccoli"]}`,
			},
			errorExpected: false,
		},
		{
			name: "ErrorOnDifferentPath",
			request: &Request{
				Scopes: map[string]*accessmodel.AccessModel{
					"foods": {
						Type: accessmodel.RESTType,
						RESTAccessModel: accessmodel.RESTAccessModel{
							Path:   "/users/recipes",
							Query:  map[string]string{},
							Body:   `{"fruit": "lemon", "greens": ["broccoli"]}`,
							Method: "POST",
						},
					},
				},
				Method: "POST",
				Path:   "/users/cars",
				Query:  url.Values{},
				Body:   `{"fruit": "lemon", "greens": ["broccoli"]}`,
			},
			errorExpected: true,
		},
		{
			name: "ErrorOnDifferentBody",
			request: &Request{
				Scopes: map[string]*accessmodel.AccessModel{
					"foods": {
						Type: accessmodel.RESTType,
						RESTAccessModel: accessmodel.RESTAccessModel{
							Path:   "/users/recipes",
							Query:  map[string]string{},
							Body:   `{"fruit": "lemon", " greens": ["broccoli"]}`,
							Method: "POST",
						},
					},
				},
				Method: "POST",
				Path:   "/users/recipes",
				Query:  url.Values{},
				Body:   `{"fruit": "lemon", "greens": ["broccoli"]}`,
			},
			errorExpected: true,
		},
		{
			name: "ErrorOnDifferentMethod",
			request: &Request{
				Scopes: map[string]*accessmodel.AccessModel{
					"foods": {
						Type: accessmodel.RESTType,
						RESTAccessModel: accessmodel.RESTAccessModel{
							Path:   "/users/recipes",
							Query:  map[string]string{},
							Body:   `{"fruit": "lemon", "greens": ["broccoli"]}`,
							Method: "PUT",
						},
					},
				},
				Method: "POST",
				Path:   "/users/recipes",
				Query:  url.Values{},
				Body:   `{"fruit": "lemon", "greens": ["broccoli"]}`,
			},
			errorExpected: true,
		},
		{
			name: "ErrorOnDifferentQuery",
			request: &Request{
				Scopes: map[string]*accessmodel.AccessModel{
					"foods": {
						Type: accessmodel.RESTType,
						RESTAccessModel: accessmodel.RESTAccessModel{
							Path: "/users/recipes",
							Query: map[string]string{
								"fruit":  "lemon",
								"greens": "broccoli",
							},
							Body:   "",
							Method: "GET",
						},
					},
				},
				Method: "GET",
				Path:   "/users/recipes",
				Query: url.Values{
					"fruit":  []string{"apple"},
					"greens": []string{"broccoli"},
				},
				Body: "",
			},
			errorExpected: true,
		},
		{
			name: "NoErrorOnOneMatchingModel",
			request: &Request{
				Scopes: map[string]*accessmodel.AccessModel{
					"foods": {
						Type: accessmodel.RESTType,
						RESTAccessModel: accessmodel.RESTAccessModel{
							Path: "/users/recipes",
							Query: map[string]string{
								"fruit":  "lemon",
								"greens": "broccoli",
							},
							Body:   "",
							Method: "GET",
						},
					},
					"foods:apple": {
						Type: accessmodel.RESTType,
						RESTAccessModel: accessmodel.RESTAccessModel{
							Path: "/users/recipes",
							Query: map[string]string{
								"fruit":  "apple",
								"greens": "broccoli",
							},
							Body:   "",
							Method: "GET",
						},
					},
				},
				Method: "GET",
				Path:   "/users/recipes",
				Query: url.Values{
					"fruit":  []string{"apple"},
					"greens": []string{"broccoli"},
				},
				Body: "",
			},
			errorExpected: false,
		},
		{
			name: "ErrorOnOnlyGQLScopes",
			request: &Request{
				Scopes: map[string]*accessmodel.AccessModel{
					"scope": {
						Type: accessmodel.GQLType,
					},
				},
				Method: "GET",
				Path:   "/users/recipes",
				Query: url.Values{
					"fruit":  []string{"lemon"},
					"greens": []string{"broccoli"},
				},
				Body: "",
			},
			errorExpected: true,
		},
	}

	verifier := RESTVerifier{}

	for _, test := range tests {
		s.Run(test.name, func() {
			err := verifier.Verify(test.request)
			if test.errorExpected {
				s.Require().Error(err)
				s.Require().NotNil(err)
				s.True(errors.Is(err, ErrNotValid), err)
			} else {
				s.NoError(err)
			}
		})
	}
}
