//go:build integration || to || files
// +build integration to files

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	databron "lab.weave.nl/nid/nid-core/integration-tests/databron/models"
	auth "lab.weave.nl/nid/nid-core/svc/auth/models"
)

/**
 * Setup suite
 */
type InformationIntegrationTestSuite struct {
	BaseTestSuite
	accessTokenContactable     string
	accessTokenHasAddress      string
	accessTokenPositiveBalance string
	client                     *http.Client
}

func TestInformationIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(InformationIntegrationTestSuite))
}

func (s *InformationIntegrationTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	authTestConfig, err := GetAuthTestConfig()
	s.Require().NoError(err)

	// Run auth setup test to make sure pseudonym is set
	audienceDB := &auth.AudienceDB{}
	accessModelDB := &auth.AccessModelDB{}
	s.client = &http.Client{}

	// Get access token for primary audience
	s.accessTokenContactable, _ = AuthorizeTokens(
		&s.BaseTestSuite,
		&s.authClient,
		s.authHTTPClient,
		NewAuthTest(
			audienceDB.DefaultModelInformationService(s.envConfig.Namespace),
			accessModelDB.DefaultModelContactable(s.envConfig.Namespace),
			accessModelDB.DefaultModelHasAddress(s.envConfig.Namespace),
			false,
			authTestConfig))

	// Get access token for secondary audience
	s.accessTokenHasAddress, _ = AuthorizeTokens(
		&s.BaseTestSuite,
		&s.authClient,
		s.authHTTPClient,
		NewAuthTest(
			audienceDB.DefaultModelInformationService(s.envConfig.Namespace),
			accessModelDB.DefaultModelHasAddress(s.envConfig.Namespace),
			accessModelDB.DefaultModelContactable(s.envConfig.Namespace),
			false,
			authTestConfig))

	s.accessTokenPositiveBalance, _ = AuthorizeTokens(
		&s.BaseTestSuite,
		&s.authClient,
		s.authHTTPClient,
		NewAuthTest(
			audienceDB.DefaultModelInformationService(s.envConfig.Namespace),
			accessModelDB.DefaultModelHasPositiveBankAccountBalance(s.envConfig.Namespace),
			accessModelDB.DefaultModelHasAddress(s.envConfig.Namespace),
			false,
			authTestConfig))

	// Make sure access token and pseudonym are set after auth flow tests are run
	s.Require().NotEmpty(s.accessTokenContactable)
	s.Require().NotEmpty(s.accessTokenHasAddress)
	s.Require().NotEmpty(s.accessTokenPositiveBalance)
	s.Require().NotEmpty(authTestConfig.UserPseudo)
}

/**
* Tests
 */
func (s *InformationIntegrationTestSuite) TestStress() {
	tests := []func(){
		s.TestEndpointContactable,
		s.TestEndpointHasAddress,
		s.TestEndpointBankAccountPositiveBalance,
		s.TestErrorOnWrongAccessToken,
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

func (s *InformationIntegrationTestSuite) TestEndpointContactable() {
	url := s.BaseTestSuite.httpURL(s.envConfig.Service.Information, "/v1/info/contact-details/contactable")
	request, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessTokenContactable))
	s.Require().NoError(err)
	response, err := s.client.Do(request)
	s.Require().NoError(err)
	defer func() {
		s.NoError(response.Body.Close())
	}()
	mediaType := struct {
		Contactable bool `json:"contactable"`
	}{}
	s.GetResponse(response, &mediaType)

	s.NotNil(response.Body)
	s.Equal(http.StatusOK, response.StatusCode)

	// Validate answers
	s.Equal(true, mediaType.Contactable)
}

func (s *InformationIntegrationTestSuite) TestErrorOnWrongAccessToken() {
	url := s.BaseTestSuite.httpURL(s.envConfig.Service.Information, "/v1/info/contact-details/contactable")
	request, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessTokenPositiveBalance))
	s.Require().NoError(err)
	response, err := s.client.Do(request)
	s.Require().NoError(err)
	defer func() {
		s.NoError(response.Body.Close())
	}()

	mediaType := struct {
		Errors []struct {
			Message string
		}
	}{}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	s.Require().NoError(err)
	err = json.Unmarshal(bodyBytes, &mediaType)
	s.Require().NoError(err)
	s.Require().Len(mediaType.Errors, 1)

	s.Equal(http.StatusForbidden, response.StatusCode)

	// Validate answers
	s.Equal("request does not match scopes", mediaType.Errors[0].Message)
}

func (s *InformationIntegrationTestSuite) TestEndpointHasAddress() {
	url := s.BaseTestSuite.httpURL(s.envConfig.Service.Information, "/v1/info/contact-details/has-address")
	request, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessTokenHasAddress))
	s.Require().NoError(err)
	response, err := s.client.Do(request)
	s.Require().NoError(err)
	defer func() {
		s.NoError(response.Body.Close())
	}()
	mediaType := struct {
		AddressFound bool `json:"address_found"`
	}{}
	s.GetResponse(response, &mediaType)

	s.NotNil(response)
	s.NotNil(response.Body)
	s.Equal(http.StatusOK, response.StatusCode)

	// Validate answer
	s.Equal(true, mediaType.AddressFound)
}

func (s *InformationIntegrationTestSuite) TestEndpointBankAccountPositiveBalance() {
	url := s.BaseTestSuite.httpURL(s.envConfig.Service.Information, "/v1/info/bank-account/positive")
	request, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessTokenPositiveBalance))
	s.Require().NoError(err)
	response, err := s.client.Do(request)
	s.Require().NoError(err)
	defer func() {
		s.NoError(response.Body.Close())
	}()
	mediaType := struct {
		PositiveBalance bool `json:"balance_positive"`
	}{}
	s.GetResponse(response, &mediaType)

	// Validate response
	s.NotNil(response)
	s.NotNil(response.Body)
	s.Equal(http.StatusOK, response.StatusCode)

	// Validate answer
	bankAccountDB := databron.BankAccountDB{}
	savingsAccountDB := databron.SavingsAccountDB{}
	savingsAmount := savingsAccountDB.DefaultModel().Amount
	bankAmount := bankAccountDB.DefaultModel().Amount
	s.Equal(bankAmount > 0 || savingsAmount > 0, mediaType.PositiveBalance)
}
