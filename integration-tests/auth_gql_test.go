// +build integration to files

package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"

	gql "lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
)

type AuthGQLIntegrationTestSuite struct {
	BaseTestSuite
}

type introspectResponse struct {
	Schema schema `json:"__schema"`
}

type schema struct {
	Types []gqlType
}

type gqlType struct {
	Name   string
	Fields []string
}

func TestAuthGQLIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthGQLIntegrationTestSuite))
}

func (s *AuthGQLIntegrationTestSuite) TestAuthGQLListNamespaces() {
	introspectQuery := `{
		 __schema {
		   types {
			 name
		   }
		 }
		}`
	res := &introspectResponse{}
	s.Require().NoError(s.clients.authGQLClient.Run(s.ctx, gql.NewRequest(introspectQuery), res, gql.MethodPost))
	s.Require().LessOrEqual(1, len(res.Schema.Types))
}
