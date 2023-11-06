package identityprovider

import (
	"context"
	"testing"

	"github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/password/mock"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/repository"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DatabaseProviderTestSuite struct {
	repository.AuthDBSuite

	provider        *DatabaseIdentityProvider
	passwordManager password.IManager
}

func (s *DatabaseProviderTestSuite) SetupSuite() {
	s.AuthDBSuite.SetupSuite()

	s.passwordManager = password.NewDefaultManager()
	s.provider = NewDatabaseIdentityProvider(s.Repo.ClientDB, s.passwordManager)
}

func (s *DatabaseProviderTestSuite) SetupTest()     { s.AuthDBSuite.SetupTest() }
func (s *DatabaseProviderTestSuite) TearDownTest()  { s.AuthDBSuite.TearDownTest() }
func (s *DatabaseProviderTestSuite) TearDownSuite() { s.AuthDBSuite.TearDownSuite() }

func TestDatabaseProviderTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseProviderTestSuite))
}

func (s *DatabaseProviderTestSuite) DatabaseProviderGetIdentity(t *testing.T) {
	testClient := &models.Client{
		Name:     "testclient",
		Password: "testpassword",
	}
	err := s.AuthDBSuite.Tx.Create(testClient).Error
	s.Require().NoError(err)

	tests := []struct {
		Scenario               string
		Metadata               *models.TokenRequestMetadata
		ComparePasswordMatches bool
		ComparePasswordError   error
		Expected               string
		Error                  error
	}{
		{
			Scenario: "Client not found",
			Metadata: &models.TokenRequestMetadata{
				Username: "notfound",
			},
			Expected: "",
			Error:    contract.ErrUnauthenticated,
		},
		{
			Scenario: "Password does not match",
			Metadata: &models.TokenRequestMetadata{
				Username: testClient.Name,
				Password: "wrong",
			},
			ComparePasswordMatches: false,
			ComparePasswordError:   nil,
			Expected:               "",
			Error:                  contract.ErrUnauthenticated,
		},
		{
			Scenario: "Password manager error",
			Metadata: &models.TokenRequestMetadata{
				Username: testClient.Name,
				Password: "wrong",
			},
			ComparePasswordMatches: false,
			ComparePasswordError:   contract.ErrInternalError,
			Expected:               "",
			Error:                  contract.ErrUnauthenticated,
		},
		{
			Scenario: "Password matches",
			Metadata: &models.TokenRequestMetadata{
				Username: testClient.Name,
				Password: testClient.Password,
			},
			ComparePasswordMatches: true,
			ComparePasswordError:   nil,
			Expected:               testClient.Name,
			Error:                  nil,
		},
	}

	for _, test := range tests {
		passwordManagerMock := new(mock.PasswordManagerMock)
		passwordManagerMock.
			On("ComparePassword", test.Metadata.Password, testClient.Password).
			Return(test.ComparePasswordMatches, test.ComparePasswordError)

		t.Run(test.Scenario, func(t *testing.T) {
			provider := NewDatabaseIdentityProvider(s.Repo.ClientDB, nil)
			result, err := provider.GetIdentity(context.Background(), test.Metadata)

			assert.Equal(t, test.Expected, result)
			assert.ErrorIs(t, err, test.Error)
		})
	}
}
