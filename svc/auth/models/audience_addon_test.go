package models

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
)

type AudienceAddOnTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	db         *gorm.DB
	tx         *gorm.DB
	AudienceDB *AudienceDB
}

func (s *AudienceAddOnTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
	s.tx = s.db.BeginTx(s.Ctx, nil)
	s.Require().NoError(s.tx.AutoMigrate(GetModels()...).Error)
	// s.Require().Len(AddForeignKeys(s.tx), 0) FIXME reintroduce foreign keys when generator bug is fixed
	s.AudienceDB = NewAudienceDB(s.tx)
}

func (s *AudienceAddOnTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *AudienceAddOnTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
}

func (s *AudienceAddOnTestSuite) TestGetAudienceByURI() {
	audienceURI := "https://test.com/gql"
	audience, err := s.AudienceDB.GetAudienceByURI(audienceURI)
	s.Require().Nil(audience)
	s.Require().Error(err)
	s.EqualError(err, gorm.ErrRecordNotFound.Error())

	audience = &Audience{
		Audience: audienceURI,
	}
	err = s.tx.Create(audience).Error

	s.Require().NoError(err)

	insertedAudience, err := s.AudienceDB.GetAudienceByURI(audienceURI)
	s.NoError(err)
	s.Equal(audience.ID, insertedAudience.ID)
}

func TestAudienceAddOnTestSuite(t *testing.T) {
	suite.Run(t, &AudienceAddOnTestSuite{
		// Intentionally do not supply models to automigrate, this should be done inside the transaction
		db: database.MustConnectTest(databaseName, nil),
	})
}
