package models

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
)

type ClientAddOnTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	db       *gorm.DB
	tx       *gorm.DB
	ClientDB *ClientDB
}

func (s *ClientAddOnTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	s.tx = s.db.BeginTx(s.Ctx, nil)
	s.Require().NoError(s.tx.AutoMigrate(GetModels()...).Error)
	// errs := AddForeignKeys(s.tx) FIXME reintroduce foreign keys when generator bug is fixed
	// s.Len(errs, 0)

	s.ClientDB = NewClientDB(s.tx)
}

func (s *ClientAddOnTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *ClientAddOnTestSuite) TearDownSuite() {
	s.Require().NoError(s.db.Close())
}

func (s *ClientAddOnTestSuite) TestGetClientByID() {
	clientID := uuid.Must(uuid.NewV4())
	_, err := s.ClientDB.GetClientByID(clientID.String())
	s.Error(err)
	s.EqualError(err, gorm.ErrRecordNotFound.Error())

	client := &Client{
		ID: clientID,
	}
	err = s.tx.Create(client).Error
	s.Require().NoError(err)

	insertedClient, err := s.ClientDB.GetClientByID(clientID.String())
	s.NoError(err)
	s.Equal(insertedClient.ID, client.ID)
}

func TestClientAddOnTestSuite(t *testing.T) {
	suite.Run(t, &ClientAddOnTestSuite{
		// Intentionally do not supply models to automigrate, this should be done inside the transaction
		db: database.MustConnectTest(databaseName, nil),
	})
}
