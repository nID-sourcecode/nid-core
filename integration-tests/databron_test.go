//go:build integration || to || files
// +build integration to files

package integration

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	databron "lab.weave.nl/nid/nid-core/integration-tests/databron/models"
	gql "lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
	auth "lab.weave.nl/nid/nid-core/svc/auth/models"
)

/**
 * Setup suite
 */
type DatabronIntegrationTestSuite struct {
	BaseTestSuite
	accessTokenPrimaryAudience    string
	accessTokenPrimaryAudienceBSN string
	accessTokenSecondaryAudience  string
	userDB                        databron.UserDB
	contactDetailDB               databron.ContactDetailDB
	addresslDB                    databron.AddressDB
	bankAccountDB                 databron.BankAccountDB
	savingsAccountDB              databron.SavingsAccountDB
	Method                        gql.Method
}

func TestDatabronIntegrationTestSuitePost(t *testing.T) {
	suite.Run(t, &DatabronIntegrationTestSuite{Method: gql.MethodPost})
}

func TestDatabronIntegrationTestSuiteGet(t *testing.T) {
	suite.Run(t, &DatabronIntegrationTestSuite{Method: gql.MethodGet})
}

func (s *DatabronIntegrationTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	// Run auth setup test to make sure pseudonym is set
	audienceDB := auth.AudienceDB{}
	accessModeDB := auth.AccessModelDB{}

	authTestConfig, err := GetAuthTestConfig()
	s.Require().NoError(err)

	// Get access token for primary audience
	s.accessTokenPrimaryAudience, _ = AuthorizeTokens(&s.BaseTestSuite, &s.authClient, s.authHTTPClient,
		NewAuthTest(audienceDB.DefaultModelPrimary(s.envConfig.Namespace),
			accessModeDB.DefaultModelUserFirstAndLastName(s.envConfig.Namespace),
			accessModeDB.DefaultModelOptionalBankAccounts(s.envConfig.Namespace),
			true,
			authTestConfig))

	s.accessTokenPrimaryAudienceBSN, _ = AuthorizeTokens(&s.BaseTestSuite, &s.authClient, s.authHTTPClient,
		NewAuthTest(audienceDB.DefaultModelPrimary(s.envConfig.Namespace),
			accessModeDB.DefaultModelUserFirstAndLastNameByBSN(s.envConfig.Namespace),
			accessModeDB.DefaultModelOptionalBankAccountsByBSN(s.envConfig.Namespace),
			true,
			authTestConfig))

	// Get access token for secondary audience
	s.accessTokenSecondaryAudience, _ = AuthorizeTokens(&s.BaseTestSuite, &s.authClient, s.authHTTPClient, NewAuthTest(audienceDB.DefaultModelSecondary(s.envConfig.Namespace), accessModeDB.DefaultModelUserAddresses(s.envConfig.Namespace), accessModeDB.DefaultModelOptionalUserAddressContactDetails(s.envConfig.Namespace), true, authTestConfig))

	// Make sure access token and pseudonym are set after auth flow tests are run
	s.Require().NotEmpty(s.accessTokenPrimaryAudience)
	s.Require().NotEmpty(s.accessTokenSecondaryAudience)
	s.Require().NotEmpty(authTestConfig.UserPseudo)

	// Init db's
	s.userDB = databron.UserDB{}
	s.bankAccountDB = databron.BankAccountDB{}
	s.savingsAccountDB = databron.SavingsAccountDB{}
}

/**
 * Tests
 */
func (s *DatabronIntegrationTestSuite) TestStress() {
	tests := []func(){
		s.TestIntrospectQuery,
		s.TestPrimaryFieldOutsideScope,
		s.TestPrimaryFieldsWithinScope,
		s.TestPrimaryFieldsWithinScope_BSN,
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

func (s *DatabronIntegrationTestSuite) TestIntrospectQuery() {
	introspectQuery := `{
	  __schema {
		types {
		  name
		}
	  }
	}`
	request := gql.NewRequest(introspectQuery)
	request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", s.accessTokenPrimaryAudience)
	res := &introspectResponse{}
	s.Require().NoError(s.clients.databronClient.Run(s.ctx, request, res, s.Method))
	s.Require().LessOrEqual(1, len(res.Schema.Types))
}

func (s *DatabronIntegrationTestSuite) TestPrimaryFieldsWithinScope() {
	// Init query
	request := gql.NewRequest(`{
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
	`)
	request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", s.accessTokenPrimaryAudience)

	// Do request and check data
	response := &PrimaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Require().NoError(err)
	s.NotEmpty(response)
	s.NotEmpty(response.Users)
	s.Require().Len(response.Users, 1)
	defaultUser := s.userDB.DefaultModel()
	s.Equal(defaultUser.FirstName, response.Users[0].FirstName)
	s.Equal(defaultUser.LastName, response.Users[0].LastName)
	s.NotEmpty(response.Users[0].BankAccounts)
	s.Len(response.Users[0].BankAccounts, 1)
	defaultBankAccount := s.bankAccountDB.DefaultModel()
	s.Equal(defaultBankAccount.AccountNumber, response.Users[0].BankAccounts[0].AccountNumber)
	s.Equal(defaultBankAccount.Amount, response.Users[0].BankAccounts[0].Amount)
	s.NotEmpty(response.Users[0].BankAccounts[0].SavingsAccounts)
	s.Len(response.Users[0].BankAccounts[0].SavingsAccounts, 1)
	defaultSavingsAccount := s.savingsAccountDB.DefaultModel()
	s.Equal(defaultSavingsAccount.Amount, response.Users[0].BankAccounts[0].SavingsAccounts[0].Amount)
	s.Equal(defaultSavingsAccount.Name, response.Users[0].BankAccounts[0].SavingsAccounts[0].Name)
}

func (s *DatabronIntegrationTestSuite) TestPrimaryFieldsWithinScope_BSN() {
	// Init query
	request := gql.NewRequest(`{
	  users(filter: {bsn: {eq: "$$nid:bsn$$"}}) {
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
	`)
	request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", s.accessTokenPrimaryAudienceBSN)

	// Do request and check data
	response := &PrimaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Require().NoError(err)
	s.NotEmpty(response)
	s.NotEmpty(response.Users)
	s.Require().Len(response.Users, 1)
	defaultUser := s.userDB.DefaultModel()
	s.Equal(defaultUser.FirstName, response.Users[0].FirstName)
	s.Equal(defaultUser.LastName, response.Users[0].LastName)
	s.NotEmpty(response.Users[0].BankAccounts)
	s.Len(response.Users[0].BankAccounts, 1)
	defaultBankAccount := s.bankAccountDB.DefaultModel()
	s.Equal(defaultBankAccount.AccountNumber, response.Users[0].BankAccounts[0].AccountNumber)
	s.Equal(defaultBankAccount.Amount, response.Users[0].BankAccounts[0].Amount)
	s.NotEmpty(response.Users[0].BankAccounts[0].SavingsAccounts)
	s.Len(response.Users[0].BankAccounts[0].SavingsAccounts, 1)
	defaultSavingsAccount := s.savingsAccountDB.DefaultModel()
	s.Equal(defaultSavingsAccount.Amount, response.Users[0].BankAccounts[0].SavingsAccounts[0].Amount)
	s.Equal(defaultSavingsAccount.Name, response.Users[0].BankAccounts[0].SavingsAccounts[0].Name)
}

func (s *DatabronIntegrationTestSuite) TestPrimaryFieldOutsideScope() {
	// Init query
	request := gql.NewRequest(`{
	  users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
		id
	  }
	}
	`)
	request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", s.accessTokenPrimaryAudience)

	// Do request and check data
	response := &PrimaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Error(err)
	s.EqualError(err, "remote graphql error response (403): request does not match scopes")
}

func (s *DatabronIntegrationTestSuite) TestPrimaryWithoutAuthorizationHeader() {
	// Init query
	request := gql.NewRequest(`{
	  users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
		id
	  }
	}
	`)
	// Do request and check data
	response := &PrimaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Error(err)
	s.EqualError(err, `remote graphql error response (400): authorization header not found or empty`)
}

func (s *DatabronIntegrationTestSuite) TestPrimaryWithBadAuthorizationHeader() {
	// Init query
	request := gql.NewRequest(`{
	  users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
		id
	  }
	}
	`)
	request.Headers["Authorization"] = "Beare abcde"
	// Do request and check data
	response := &PrimaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Error(err)
	s.ErrorIs(err, gql.ErrRemoteErrorResponse)
	s.True(strings.Contains(err.Error(), `invalid authorization header`), err.Error())
}

func (s *DatabronIntegrationTestSuite) TestSecondaryFieldsWithinScope() {
	// Init query
	request := gql.NewRequest(`
		{
		  users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
				contactDetails {
			  phone
			  address {
				houseNumber
			  }
			}
		  }
		}
	`)
	request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", s.accessTokenSecondaryAudience)

	// Do request and check data
	response := &SecondaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Require().NoError(err)
	s.Require().NotEmpty(response)
	s.Require().NotEmpty(response.Users)
	s.Require().Len(response.Users, 1)
	s.Require().NotEmpty(response.Users[0].ContactDetails)
	s.Require().Len(response.Users[0].ContactDetails, 1)
	defaultContactDetail := s.contactDetailDB.DefaultModel()
	s.Equal(defaultContactDetail.Phone, response.Users[0].ContactDetails[0].Phone)
	s.NotEmpty(response.Users[0].ContactDetails[0].Address)
	defaultAddress := s.addresslDB.DefaultModel()
	s.Equal(defaultAddress.HouseNumber, response.Users[0].ContactDetails[0].Address.HouseNumber)
}

func (s *DatabronIntegrationTestSuite) TestSecondaryFieldOutsideScope() {
	// Init query
	request := gql.NewRequest(`
		{
		  users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
			contactDetails {
			  phone
			  address {
				houseNumber
				houseNumberAddon
			  }
			}
		  }
		}
	`)
	request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", s.accessTokenSecondaryAudience)

	// Do request and check data
	response := &PrimaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Error(err)
	s.EqualError(err, "remote graphql error response (403): request does not match scopes")
}

func (s *DatabronIntegrationTestSuite) TestSecondaryWithPrimaryToken() {
	// Init query
	request := gql.NewRequest(`
		{
		  users(filter: {pseudonym: {eq: "$$nid:subject$$"}}) {
			contactDetails {
			  phone
			  address {
				houseNumber
				houseNumberAddon
			  }
			}
		  }
		}
	`)
	request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", s.accessTokenPrimaryAudience)

	// Do request and check data
	response := &PrimaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Error(err)
	s.EqualError(err, "remote graphql error response (403): request does not match scopes")
}

func (s *DatabronIntegrationTestSuite) TestWithoutFilter() {
	// Init query
	request := gql.NewRequest(`
		{
		  users {
			firstName
		  }
		}
	`)
	request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", s.accessTokenPrimaryAudience)

	// Do request and check data
	response := &PrimaryQueryResponse{}
	err := s.clients.databronClient.Run(s.ctx, request, response, s.Method)
	s.Error(err)
	s.EqualError(err, "remote graphql error response (403): request does not match scopes")
}

/**
 * Types and helpers
 */
type PrimaryQueryResponse struct {
	Users []UsersPrimaryQuery `json:"users"`
}

type UsersPrimaryQuery struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	BankAccounts []BankAccount
}

type BankAccount struct {
	AccountNumber   string `json:"accountNumber"`
	Amount          int    `json:"amount"`
	SavingsAccounts []SavingsAccount
}

type SavingsAccount struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type SecondaryQueryResponse struct {
	Users []UsersSecondaryQuery `json:"users"`
}

type UsersSecondaryQuery struct {
	ContactDetails []ContactDetail `json:"contactDetails"`
}

type ContactDetail struct {
	Phone   string  `json:"phone"`
	Address Address `json:"address"`
}

type Address struct {
	HouseNumber int `json:"houseNumber"`
}
