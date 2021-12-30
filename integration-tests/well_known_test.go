// +build integration to files

package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vrischmann/envconfig"
)

type WellKnownIntegrationTestSuiteConfig struct {
	Namespace string `envconfig:"NAMESPACE"`
}

type WellKnownIntegrationTestSuite struct {
	BaseTestSuite
	conf       *WellKnownIntegrationTestSuiteConfig
	httpClient *http.Client
}

func (s *WellKnownIntegrationTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.httpClient = http.DefaultClient

	s.conf = &WellKnownIntegrationTestSuiteConfig{}
	err := envconfig.Init(&s.conf)
	s.Require().NoError(err, "reading env")
}

func TestWellKnownIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WellKnownIntegrationTestSuite))
}

func (s *WellKnownIntegrationTestSuite) TestWellKnownOpenIDConfiguration() {
	scheme := schemeHTTPS
	if !s.envConfig.IsTLS {
		scheme = schemeHTTP
	}

	url := fmt.Sprintf("%s://auth.%s/.well-known/openid-configuration", scheme, s.envConfig.BackendURL)
	req, err := http.NewRequestWithContext(s.ctx, http.MethodGet, url, nil)
	s.Require().NoError(err)

	resp, err := s.httpClient.Do(req)
	s.Require().NoError(err, "expected no error (this is probably due the http transcoder)")
	defer func(test *WellKnownIntegrationTestSuite) {
		test.Require().NoError(resp.Body.Close())
	}(s)

	// nolint: govet
	respBody := struct {
		Issuer                           string   `json:"issuer"`
		AuthorizationEndpoint            string   `json:"authorization_endpoint`
		TokenEndpoint                    string   `json:"token_endpoint`
		JWKSURI                          string   `json:"jwks_uri`
		ScopesSupported                  []string `json:"scopes_supported"`
		ResponseTypesSupported           []string `json:"response_types_supported"`
		GrantTypesSupported              []string `json:"grant_types_supported"`
		IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
		ClaimTypesSupported              []string `json:"claim_types_supported"`
	}{}

	s.Require().Equal(http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)

	s.Require().NoError(json.Unmarshal(body, &respBody))
	s.Equal("auth."+s.conf.Namespace, respBody.Issuer)
	s.Require().Len(respBody.ScopesSupported, 1)
	s.Equal("openid", respBody.ScopesSupported[0])
	s.Require().Len(respBody.ResponseTypesSupported, 1)
	s.Equal("code", respBody.ResponseTypesSupported[0])
	s.Require().Len(respBody.GrantTypesSupported, 1)
	s.Equal("authorization_code", respBody.GrantTypesSupported[0])
	s.Require().Len(respBody.IDTokenSigningAlgValuesSupported, 1)
	s.Equal("RS256", respBody.IDTokenSigningAlgValuesSupported[0])
}
