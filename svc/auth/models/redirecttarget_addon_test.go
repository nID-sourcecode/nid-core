package models

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
)

type RedirectTargetAddOnTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	db               *gorm.DB
	tx               *gorm.DB
	RedirectTargetDB *RedirectTargetDB
}

func (s *RedirectTargetAddOnTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	s.tx = s.db.BeginTx(s.Ctx, nil)
	s.Require().NoError(s.tx.AutoMigrate(GetModels()...).Error)

	type Column struct {
		ColumnName string
		// more fields if needed...
	}
	var columns []*Column
	if err := s.tx.Table("information_schema.columns").Select("column_name").Where("table_schema = ? AND table_name = ?", "public", "accepted_access_models_sessions").Find(&columns).Error; err != nil {
		panic(err)
	}

	// s.Require().Len(AddForeignKeys(s.tx), 0) FIXME reintroduce foreign keys when generator bug is fixed

	s.RedirectTargetDB = NewRedirectTargetDB(s.tx)
}

func (s *RedirectTargetAddOnTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *RedirectTargetAddOnTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
}

func (s *RedirectTargetAddOnTestSuite) TestGetRedirectTarget() {
	redirectURI := "https://weave.nl/code"
	clientID := uuid.Must(uuid.NewV4())
	_, err := s.RedirectTargetDB.GetRedirectTarget(redirectURI, clientID.String())
	s.Error(err)
	s.EqualError(err, gorm.ErrRecordNotFound.Error())

	client := &Client{
		ID: clientID,
	}
	err = s.tx.Create(client).Error
	s.Require().NoError(err)

	redirectTarget := &RedirectTarget{
		ClientID:       clientID,
		RedirectTarget: redirectURI,
	}
	err = s.tx.Create(redirectTarget).Error
	s.Require().NoError(err)

	insertedRedirectTarget, err := s.RedirectTargetDB.GetRedirectTarget(redirectURI, clientID.String())
	s.NoError(err)
	s.Equal(insertedRedirectTarget.ID, redirectTarget.ID)
}

func TestRedirectTargetAddOnTestSuite(t *testing.T) {
	suite.Run(t, &RedirectTargetAddOnTestSuite{
		// Intentionally do not supply models to automigrate, this should be done inside the transaction
		db: database.MustConnectTest(databaseName, nil),
	})
}
