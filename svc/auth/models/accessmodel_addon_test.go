package models

import (
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
)

type AccessModelAddOnTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	db            *gorm.DB
	tx            *gorm.DB
	AccessModelDB *AccessModelDB
	AudienceDB    *AudienceDB
}

func (s *AccessModelAddOnTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	s.tx = s.db.BeginTx(s.Ctx, nil)
	s.Require().NoError(s.tx.AutoMigrate(GetModels()...).Error)
	// errs := AddForeignKeys(s.tx) FIXME reintroduce foreign keys when generator bug is fixed
	// s.Require().Len(errs, 0)

	s.AccessModelDB = NewAccessModelDB(s.tx)
}

func (s *AccessModelAddOnTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *AccessModelAddOnTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
}

func (s *AccessModelAddOnTestSuite) TestGetAccessModelsByIDs() {
	_, err := s.AccessModelDB.GetAccessModelsByIDs([]string{uuid.Must(uuid.NewV4()).String()})
	s.Error(err)
	s.EqualError(err, gorm.ErrRecordNotFound.Error())

	audience := &Audience{}
	err = s.tx.Create(audience).Error
	s.Require().NoError(err)

	model := &AccessModel{
		AudienceID: audience.ID,
		Hash:       "ghi",
		Name:       "test:stuff2",
		Type:       AccessModelTypeGQL,
		GqlAccessModel: &GqlAccessModel{
			Path:      "/graphql",
			JSONModel: `{"r":"somemorestuff"}`,
		},
	}
	err = s.AccessModelDB.CreateAccessModel(model)
	s.Require().NoError(err)

	accessModels, err := s.AccessModelDB.GetAccessModelsByIDs([]string{model.ID.String()})
	s.NoError(err)
	s.Equal(accessModels[0].ID, model.ID)
}

func (s *AccessModelAddOnTestSuite) TestGetAccessModelByAudience() {
	audience := &Audience{
		Audience: "https://test.com/gql",
	}
	scope := "openid test:stuff@abc"
	_, err := s.AccessModelDB.GetAccessModelByAudienceWithScope("openid test:stuff", "abc", audience)
	s.Error(err)
	s.EqualError(err, "record not found")

	err = s.tx.Create(audience).Error
	s.Require().NoError(err)

	parts := strings.Split(scope, "@")

	accessModel := &AccessModel{
		AudienceID: audience.ID,
		Hash:       parts[1],
		Name:       parts[0],
	}
	err = s.tx.Create(accessModel).Error
	s.Require().NoError(err)
	insertedAccessModel, err := s.AccessModelDB.GetAccessModelByAudienceWithScope(accessModel.Name, accessModel.Hash, audience)
	s.NoError(err)
	s.Equal(accessModel.ID, insertedAccessModel.ID)
}

func (s *AccessModelAddOnTestSuite) TestGetAccessModelsByAudiencePreloadModels() {
	audience := &Audience{
		Audience: "https://test.com/gql",
	}
	_, err := s.AccessModelDB.GetAccessModelsByAudience(true, audience)
	s.Error(err)
	s.EqualError(err, gorm.ErrRecordNotFound.Error())

	err = s.tx.Create(audience).Error
	s.Require().NoError(err)

	gqlAccessModel := &GqlAccessModel{}
	restAccessModel := &RestAccessModel{}
	accessModel := &AccessModel{
		ID:              uuid.UUID{},
		AudienceID:      audience.ID,
		GqlAccessModel:  gqlAccessModel,
		RestAccessModel: restAccessModel,
	}
	err = s.tx.Create(accessModel).Error
	s.Require().NoError(err)

	preloadedAccessModels, err := s.AccessModelDB.GetAccessModelsByAudience(true, audience)
	s.NoError(err)
	s.Equal(preloadedAccessModels[0].ID, accessModel.ID)
	s.Equal(preloadedAccessModels[0].GqlAccessModel.ID, gqlAccessModel.ID)
	s.Equal(preloadedAccessModels[0].RestAccessModel.ID, restAccessModel.ID)

	noPreloadAccessModels, err := s.AccessModelDB.GetAccessModelsByAudience(false, audience)
	s.NoError(err)
	s.Equal(noPreloadAccessModels[0].ID, accessModel.ID)
	s.Empty(noPreloadAccessModels[0].GqlAccessModel)
	s.Empty(noPreloadAccessModels[0].RestAccessModel)
}

func TestAccessModelAddOnTestSuite(t *testing.T) {
	db := database.MustConnectTest(databaseName, nil)
	db.DB().SetMaxOpenConns(50)
	db.DB().SetMaxIdleConns(50)
	suite.Run(t, &AccessModelAddOnTestSuite{
		db: db,
	})
}
