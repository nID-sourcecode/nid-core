//go:build integration || to || files
// +build integration to files

package integration

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/grpckeys"
	dashboard "lab.weave.nl/nid/nid-core/svc/dashboard/proto"
)

type DashboardIntegrationTestSuite struct {
	BaseTestSuite

	dashboardConfig
	httpClient *http.Client
}

type dashboardConfig struct {
	DefaultUserEmail string `envconfig:"DASHBOARD_USER_EMAIL"`
	DefaultUserPass  string `envconfig:"DASHBOARD_USER_PASS"`
}

func (s *DashboardIntegrationTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.Require().NoError(envconfig.Init(&s.dashboardConfig), "unable to initialise dashboard environment config")
	s.httpClient = http.DefaultClient
}

// func TestDashboardIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(DashboardIntegrationTestSuite))
// }

func (s *DashboardIntegrationTestSuite) TestDashboardSignin() {
	bearer, err := signinDashboard(s.ctx, s.dashboardAuthClient, s.dashboardConfig.DefaultUserEmail, s.dashboardConfig.DefaultUserPass)
	s.NoError(err)
	s.NotEmpty(bearer)
}

func (s *DashboardIntegrationTestSuite) TestDashboardHTTPSignin() {
	scheme := schemeHTTPS
	if !s.envConfig.IsTLS {
		scheme = schemeHTTP
	}
	url := fmt.Sprintf("%s://dashboard.%s/v1/signin", scheme, s.envConfig.BackendURL)
	req, err := http.NewRequestWithContext(s.ctx, http.MethodPost, url, nil)
	s.Require().NoError(err)

	authorization := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.dashboardConfig.DefaultUserEmail, s.dashboardConfig.DefaultUserPass)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authorization))

	resp, err := s.httpClient.Do(req)
	s.Require().NoError(err, "expected no error (this is probably due the http transcoder)")
	defer func(test *DashboardIntegrationTestSuite) {
		test.Require().NoError(resp.Body.Close())
	}(s)

	respBody := struct {
		Bearer string `json:"bearer"`
	}{}

	s.Require().Equal(http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err, "expected no error (this is probably due the http transcoder)")
	s.Require().NoError(json.Unmarshal(body, &respBody), "When unmarshalling body: %s", body)
	s.Require().NotEmpty(respBody.Bearer)
}

func (s *DashboardIntegrationTestSuite) TestDashboardWrongSignin() string {
	authorization := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "example@example.com", "secret")))
	ctx := metadata.AppendToOutgoingContext(s.ctx, grpckeys.AuthorizationKey.String(), fmt.Sprintf("Basic %s", authorization))

	var header metadata.MD
	resp, err := s.clients.dashboardAuthClient.Signin(ctx, &empty.Empty{}, grpc.Header(&header))
	s.Require().Error(err)
	s.Require().Nil(resp)
	s.Require().Empty(resp.GetBearer())

	return resp.GetBearer()
}

func (s *DashboardIntegrationTestSuite) TestDashboardListNamespacesAuthenticated() {
	bearer, err := signinDashboard(s.ctx, s.dashboardAuthClient, s.dashboardConfig.DefaultUserEmail, s.dashboardConfig.DefaultUserPass)
	s.Require().NoError(err)
	ctx := metadata.AppendToOutgoingContext(s.ctx, grpckeys.AuthorizationKey.String(), fmt.Sprintf("Bearer %s", bearer))
	namespaces, err := s.clients.dashboardClient.ListNamespaces(ctx, &empty.Empty{})
	s.Require().NoError(err)
	s.Require().NotNil(namespaces)
	s.Require().LessOrEqual(1, len(namespaces.GetItems()), "ListNamespaces should return namespaces")

	namespaceContainsString := false
	for _, namespace := range namespaces.GetItems() {
		if namespace == "nid" || namespace == "twi" {
			namespaceContainsString = true
		}
	}
	s.Require().True(namespaceContainsString, "namespaces dont contain either 'nid' or 'twi'")
}

func (s *DashboardIntegrationTestSuite) TestDashboardListNamespacesWrongJWT() {
	ctx := metadata.AppendToOutgoingContext(s.ctx, grpckeys.AuthorizationKey.String(), fmt.Sprintf("Bearer %s", "eyDeze.jwt.kloptniet"))
	namespaces, err := s.clients.dashboardClient.ListNamespaces(ctx, &empty.Empty{})
	s.Require().Error(err)
	s.Require().Nil(namespaces)
}

func (s *DashboardIntegrationTestSuite) TestDashboardListNamespacesUnauthenticated() {
	namespace, err := s.clients.dashboardClient.ListNamespaces(s.ctx, &empty.Empty{})
	s.Require().Error(err)
	s.Require().Nil(namespace)
	s.Require().EqualError(err, "rpc error: code = PermissionDenied desc = RBAC: access denied")
}

func signinDashboard(ctx context.Context, dashboardAuthClient dashboard.AuthorizationServiceClient, username, password string) (string, error) {
	authorization := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	ctx = metadata.AppendToOutgoingContext(ctx, grpckeys.AuthorizationKey.String(), fmt.Sprintf("Basic %s", authorization))

	var header metadata.MD

	resp, err := dashboardAuthClient.Signin(ctx, &empty.Empty{}, grpc.Header(&header))
	if err != nil {
		return "", errors.Wrap(err, "Signing in to dashboard")
	}

	return resp.GetBearer(), nil
}
