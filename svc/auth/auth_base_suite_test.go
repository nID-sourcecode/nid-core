package main

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
)

type AuthServiceBaseTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	db     *gorm.DB
	tx     *gorm.DB
	authDB *AuthDB

	passwordManager password.IManager
	TestModels
}

type TestModels struct {
	client              *models.Client
	client2             *models.Client
	redirectTarget      *models.RedirectTarget
	audience            *models.Audience
	audience2           *models.Audience
	accessModelGql1     *models.AccessModel
	accessModelGql2     *models.AccessModel
	accessModelGql3     *models.AccessModel
	accessModelGql4     *models.AccessModel
	accessModelRestPOST *models.AccessModel
	accessModelRestGET  *models.AccessModel
}

func (s *AuthServiceBaseTestSuite) setupDB() *AuthDB {
	s.tx = s.db.BeginTx(s.Ctx, nil)
	// s.Require().Len(models.AddForeignKeys(s.tx), 0) FIXME reintroduce foreign keys when generator bug is fixed
	s.Require().NoError(s.tx.AutoMigrate(models.GetModels()...).Error)
	return &AuthDB{
		db:               s.tx,
		AccessModelDB:    models.NewAccessModelDB(s.tx),
		AudienceDB:       models.NewAudienceDB(s.tx),
		ClientDB:         models.NewClientDB(s.tx),
		RedirectTargetDB: models.NewRedirectTargetDB(s.tx),
		SessionDB:        models.NewSessionDB(s.tx),
	}
}

func (s *AuthServiceBaseTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
	s.authDB = s.setupDB()
	s.Require().NoError(s.initModels())
}

func (s *AuthServiceBaseTestSuite) SetupSuite() {
	s.passwordManager = password.NewDefaultManager()
}

func (s *AuthServiceBaseTestSuite) initModels() error {
	client1Password, err := s.passwordManager.GenerateHash("test^123")
	s.Require().NoError(err)
	client2Password, err := s.passwordManager.GenerateHash("456#%test")
	s.Require().NoError(err)

	s.TestModels.client = &models.Client{
		Color:    "blue",
		Name:     "testclient",
		Password: client1Password,
		Metadata: postgres.Jsonb{RawMessage: json.RawMessage(`{"oin":"000012345"}`)},
	}
	if err := s.tx.Create(s.client).Error; err != nil {
		return errors.Wrap(err, "initialising client 1")
	}

	s.TestModels.client2 = &models.Client{
		Color:    "red",
		Name:     "testclient2",
		Password: client2Password,
	}
	if err := s.tx.Create(s.client2).Error; err != nil {
		return errors.Wrap(err, "initialising client 2")
	}
	s.TestModels.redirectTarget = &models.RedirectTarget{
		ClientID:       s.client.ID,
		RedirectTarget: "https://weave.nl/code",
	}

	if err := s.tx.Create(s.redirectTarget).Error; err != nil {
		return errors.Wrap(err, "initialising redirect target 1")
	}
	if err := s.tx.Create(&models.RedirectTarget{
		ClientID:       s.client2.ID,
		RedirectTarget: "https://weave2.nl/code",
		UpdatedAt:      time.Time{},
	}).Error; err != nil {
		return errors.Wrap(err, "initialising redirect target 2")
	}

	s.TestModels.audience = &models.Audience{
		Audience:  "https://test.com/gql",
		Namespace: "alice",
	}
	if err := s.tx.Create(s.audience).Error; err != nil {
		return errors.Wrap(err, "initialising audience 1")
	}
	s.TestModels.audience2 = &models.Audience{
		Audience:  "https://test2.com/gql",
		Namespace: "bob",
	}
	if err := s.tx.Create(s.audience2).Error; err != nil {
		return errors.Wrap(err, "initialising audience 2")
	}

	s.TestModels.accessModelGql1 = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "abc",
		Name:       "test:stuff",
		Type:       models.AccessModelTypeGQL,
		GqlAccessModel: &models.GqlAccessModel{
			Path:      "/gql",
			JSONModel: `{"r":"somestuff"}`,
		},
	}
	if err := s.tx.Create(s.accessModelGql1).Error; err != nil {
		return errors.Wrap(err, "initialising access model 1")
	}
	s.TestModels.accessModelGql2 = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "ghi",
		Name:       "test:stuff2",
		Type:       models.AccessModelTypeGQL,
		GqlAccessModel: &models.GqlAccessModel{
			Path:      "/graphql",
			JSONModel: `{"r":"somemorestuff"}`,
		},
	}
	if err := s.tx.Create(s.accessModelGql2).Error; err != nil {
		return errors.Wrap(err, "initialising access model 2")
	}
	s.TestModels.accessModelGql3 = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "abcsubject",
		Name:       "test:stuff3",
		Type:       models.AccessModelTypeGQL,
		GqlAccessModel: &models.GqlAccessModel{
			Path:      "/gql",
			JSONModel: `{"subject":"$$nid:subject$$"}`,
		},
	}
	if err := s.tx.Create(s.accessModelGql3).Error; err != nil {
		return errors.Wrap(err, "initialising access model 3")
	}
	s.TestModels.accessModelGql4 = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "abcbsn",
		Name:       "test:stuff4",
		Type:       models.AccessModelTypeGQL,
		GqlAccessModel: &models.GqlAccessModel{
			Path:      "/gql",
			JSONModel: `{"bsn":"$$nid:bsn$$"}`,
		},
	}
	if err := s.tx.Create(s.accessModelGql4).Error; err != nil {
		return errors.Wrap(err, "initialising access model 4")
	}
	s.TestModels.accessModelRestPOST = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "jkl",
		Name:       "test:stuff3",
		Type:       models.AccessModelTypeREST,
		RestAccessModel: &models.RestAccessModel{
			Body:   "something",
			Method: "POST",
			Path:   "/some/rest/endpoint",
			Query:  "{}",
		},
	}
	if err := s.tx.Create(s.accessModelRestPOST).Error; err != nil {
		return errors.Wrap(err, "initialising access model 3")
	}

	s.TestModels.accessModelRestGET = &models.AccessModel{
		AudienceID: s.audience2.ID,
		Hash:       "def",
		Name:       "test2:stuff2",
		Type:       models.AccessModelTypeREST,
		RestAccessModel: &models.RestAccessModel{
			Body:   "",
			Method: "GET",
			Path:   "/some/rest/endpoint",
			Query:  `{"something":"somethingelse"}`,
		},
	}
	if err := s.tx.Create(s.accessModelRestGET).Error; err != nil {
		return errors.Wrap(err, "initialising access model 4")
	}

	return nil
}

// 599fc504-4075-411f-b4e0-eaba0d4b6d59
