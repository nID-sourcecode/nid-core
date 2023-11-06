package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	"github.com/stretchr/testify/suite"
)

// AuthDBSuite is the test-suite for the auth database.
type AuthDBSuite struct {
	suite.Suite

	Db   *gorm.DB
	Repo *AuthDB
	Tx   *gorm.DB
}

// SetupSuite is ran before the first in the test-suite.
func (s *AuthDBSuite) SetupSuite() {
	err := database.CreateTestDatabase("auth")
	if err != nil {
		s.T().Fatal(err)
	}

	s.Db = database.MustConnectTest("auth", nil)
}

// SetupTest is ran before each test in the test-suite.
func (s *AuthDBSuite) SetupTest() {
	s.Tx = s.Db.BeginTx(context.Background(), nil)
	s.Require().NoError(s.Tx.AutoMigrate(models.GetModels()...).Error)

	s.Repo = NewAuthDB(s.Tx)
}

// TearDownTest is ran after each test in the test-suite.
func (s *AuthDBSuite) TearDownTest() {
	s.Tx.Rollback()
}

// TearDownSuite is ran after the last test in the test-suite.
func (s *AuthDBSuite) TearDownSuite() {
	err := s.Db.Close()
	s.Require().NoError(err)
}
