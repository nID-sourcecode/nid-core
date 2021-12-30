//go:build integration || to || files
// +build integration to files

package integration

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	gql "lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
)

/**
 * Setup suite
 */
type InfoIntegrationTestSuite struct {
	BaseTestSuite
}

// func TestInfoIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(InfoIntegrationTestSuite))
// }

func (s *InfoIntegrationTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
}

/**
* Tests
 */
func (s *InfoIntegrationTestSuite) TestStress() {
	tests := []func(){
		s.TestEndpointContactable,
		s.TestEndpointPerson,
		s.TestEndpointPersonHasBackAccounts,
	}
	rounds := 100
	if testing.Short() {
		rounds = 5
	}
	for i := 0; i < rounds; i++ {
		for _, test := range tests {
			tokens := strings.Split(runtime.FuncForPC(reflect.ValueOf(test).Pointer()).Name(), ".")
			funcName := tokens[len(tokens)-1]
			name := fmt.Sprintf("%s_%d", funcName, i)
			s.Run(name, test)
		}
	}
}

func (s *InfoIntegrationTestSuite) TestEndpointContactable() {
	res := struct {
		PersonContactable struct {
			Answer bool `json:"answer"`
		} `json:"personContactable"`
	}{}
	req := gql.NewRequest(`query($id: String!) {
  personContactable(id: $id) {
    answer
  }
}`)
	req.Variables["id"] = "1d3efd44-e72a-49c7-888d-3fd1146a4c93"
	s.Require().NoError(s.infoClient.Post(context.Background(), req, &res))

	s.Require().NotNil(res)
	s.Require().NotNil(res.PersonContactable)
	s.Equal(true, res.PersonContactable.Answer)
}

func (s *InfoIntegrationTestSuite) TestEndpointPerson() {
	res := struct {
		Person struct {
			ID      string `json:"id"`
			Name    string
			Address struct {
				Line string
			}
			Pseudonymised bool
		}
	}{}

	req := gql.NewRequest(`query {
	person {
		id
		name
		address {
			line
		}
		pseudonymised
	}
}`)

	s.Require().NoError(s.infoClient.Post(context.Background(), req, &res))

	// Validate answer
	s.Equal("1d3efd44-e72a-49c7-888d-3fd1146a4c93", res.Person.ID)
	s.Equal("John Dummy Doe Ho", res.Person.Name)
	s.Equal("1001 xyz 4321 YX", res.Person.Address.Line)
	s.Equal(true, res.Person.Pseudonymised)
}

func (s *InfoIntegrationTestSuite) TestEndpointPersonHasBackAccounts() {
	res := struct {
		PersonHasBankAccounts struct {
			Answer bool
		}
	}{}

	req := gql.NewRequest(`query {
	personHasBankAccounts {
		answer
	}
}`)

	s.Require().NoError(s.infoClient.Post(context.Background(), req, &res))

	// Validate answer
	s.Equal(true, res.PersonHasBankAccounts.Answer)
}
