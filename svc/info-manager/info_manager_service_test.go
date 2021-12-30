package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vrischmann/envconfig"

	"lab.weave.nl/nid/nid-core/pkg/authtoken"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	errgrpc "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	mocklog "lab.weave.nl/nid/nid-core/pkg/utilities/log/v2/mock"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
	mockobjectstorage "lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage/mock"
	mockinforestarter "lab.weave.nl/nid/nid-core/svc/info-manager/inforestarter/mock"
	"lab.weave.nl/nid/nid-core/svc/info-manager/models"
	"lab.weave.nl/nid/nid-core/svc/info-manager/proto"
)

// InfoManagerServiceTestSuite suite struct
type InfoManagerServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	tx                *gorm.DB
	srv               *InfoManagerServiceServer
	conf              *InfoManagerConfig
	objectStorageMock *mockobjectstorage.Client
	models            InfoManagerServiceTestModels
	infoRestarterMock *mockinforestarter.InfoRestarter
	loggerMock        *mocklog.LoggerUtility
}

// InfoManagerServiceTestModels struct storing the dummy script and script sources
type InfoManagerServiceTestModels struct {
	script              *models.Script
	scriptSources       []models.ScriptSource
	scriptSourcesByHash map[string]string
}

// SetupSuite runs first every time the suite gets tested
func (s *InfoManagerServiceTestSuite) SetupSuite() {
	// setup config
	c := &InfoManagerConfig{}
	err := envconfig.InitWithOptions(c, envconfig.Options{AllOptional: true})
	if err != nil {
		s.Failf("init conf failed", "%+v", err)
	}
	s.conf = c
	// setup server
	s.objectStorageMock = &mockobjectstorage.Client{}
	s.srv = &InfoManagerServiceServer{
		db:            s.setupDB(c),
		storageClient: s.objectStorageMock,
	}

	// setup models
	if err := s.setupModels(); err != nil {
		s.Failf("error inserting test records", "%+v", err)
	}
}

// TearDownSuite runs lastly every time the suite gets tested
func (s *InfoManagerServiceTestSuite) TearDownSuite() {
	s.tx.Rollback()
	s.T().Log("Tearing down test suite")
}

func (s *InfoManagerServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
	s.infoRestarterMock = &mockinforestarter.InfoRestarter{}
	s.srv.infoRestarter = s.infoRestarterMock

	s.loggerMock = &mocklog.LoggerUtility{}
	// s.srv.logger = s.loggerMock
}

// setupDB returns the infomanager database and returns the started transaction
func (s *InfoManagerServiceTestSuite) setupDB(conf *InfoManagerConfig) *InfoManagerDB {
	db := initDB(conf, true).db
	models.AddForeignKeys(db)
	s.tx = db.Begin()
	return &InfoManagerDB{
		db:             s.tx,
		ScriptDB:       models.NewScriptDB(s.tx),
		ScriptSourceDB: models.NewScriptSourceDB(s.tx),
	}
}

// initModels creates some seeded models for the InfoManagerServiceTestModels struct
func (s *InfoManagerServiceTestSuite) setupModels() error {
	// add script sources on hash
	s.models.scriptSourcesByHash = make(map[string]string)
	hash1, err := authtoken.Hash(scriptSourceV1)
	s.Require().NoError(err)
	hash2, err := authtoken.Hash(scriptSourceV2)
	s.Require().NoError(err)
	s.models.scriptSourcesByHash[hash1] = scriptSourceV1
	s.models.scriptSourcesByHash[hash2] = scriptSourceV2

	// setup script model
	s.models.script = &models.Script{
		Name:        "Test example script",
		Description: "Test example script",
		Status:      models.ScriptStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.tx.Create(s.models.script).Error; err != nil {
		return errors.Wrap(err, "initialising script 1")
	}
	// setup script source model 1
	s.models.scriptSources = append(s.models.scriptSources, models.ScriptSource{
		ChangeDescription: "Initial script source",
		Version:           "1",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Checksum:          hash1,
		ScriptID:          s.models.script.ID,
	})
	// setup script model 2
	s.models.scriptSources = append(s.models.scriptSources, models.ScriptSource{
		ChangeDescription: "Added yolo variable",
		Version:           "2",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Checksum:          hash2,
		ScriptID:          s.models.script.ID,
	})

	for _, source := range s.models.scriptSources {
		if err := s.tx.Create(&source).Error; err != nil {
			return errors.Wrap(err, "initialising script source 1 and 2")
		}
	}

	return nil
}

// var ErrAllesIsLek = errors.New("alles is lek")

//func (s *InfoManagerServiceTestSuite) TestInfoRestartErrorIsLogged() {
//	s.infoRestarterMock.On("RestartInfoServices").Return(ErrAllesIsLek)
//	s.loggerMock.On("WithError", ErrAllesIsLek).Return(s.loggerMock)
//	s.loggerMock.On("Error", "restarting info services")
//
//	s.srv.restartInfoServices()
//
//	s.loggerMock.AssertExpectations(s.T())
//}

const scriptSourceV1 = `
	--define the new schema with a name
	name = "person"
		
	--define the schema fields
	schema = {
		id="",
		name="",
		test=schema,
		address="",
		pseudonymised=false,
	}
`

const scriptSourceV2 = `
	--define the new schema with a name
	name = "person"
		
	--define the schema fields
	schema = {
		id="",
		name="",
		test=schema,
		address="",
		pseudonymised=false,
		yolo=false,
	}
`

const scriptSourceV3 = `
	--define the new schema with a name
	name = "person"
		
	--define the schema fields
	schema = {
		id="",
		name="",
		test=schema,
		address="",
		pseudonymised=false,
		yolo=false,
		happy=true,
	}
`

func (s *InfoManagerServiceTestSuite) TestScriptsUpload() {
	// setup request
	req := &proto.ScriptsUploadRequest{
		ScriptId:          s.models.script.ID.String(),
		Script:            []byte(scriptSourceV3),
		ChangeDescription: "Added happy variable",
	}

	// create objects array for storage mock
	var objects []objectstorage.Object
	for hash := range s.models.scriptSourcesByHash {
		objects = append(objects, objectstorage.Object{
			Key: strings.Join([]string{s.models.script.ID.String(), hash}, "/"),
		})
	}
	s.Require().Len(objects, 2, "expected  objects in storage")
	var sources int
	err := s.tx.Model(&models.ScriptSource{}).Where("script_id = ?", s.models.script.ID).Count(&sources).Error
	s.Require().NoError(err)
	s.Require().NotNil(sources)
	s.Equal(sources, 2)

	// setup mocks
	s.objectStorageMock.On("List", mock.Anything, mock.Anything).Return(objects, nil)
	s.objectStorageMock.On("Write", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockCalled := make(chan struct{})
	s.infoRestarterMock.On("RestartInfoServices").Return(func() error {
		close(mockCalled)
		return nil
	})

	// upload script
	_, err = s.srv.ScriptsUpload(context.Background(), req)
	s.Require().NoError(err)

	s.objectStorageMock.AssertExpectations(s.T())

	// check if new source exists in db
	var newSources []models.ScriptSource
	err = s.tx.Model(&models.ScriptSource{}).Where("script_id = ?", s.models.script.ID).Order("version desc").Find(&newSources).Error
	s.Require().NoError(err)
	s.Require().Len(newSources, 3)

	// Check if values are set correctly
	hash, err := authtoken.Hash(scriptSourceV3)
	s.Require().NoError(err)
	s.Equal("Added happy variable", newSources[0].ChangeDescription)
	s.Equal("3", newSources[0].Version)
	s.Equal(s.models.script.ID, newSources[0].ScriptID)
	s.Equal(hash, newSources[0].Checksum)

	select {
	case <-mockCalled:
		break
	case <-time.After(time.Second):
		s.loggerMock.AssertExpectations(s.T())
		s.FailNow("RestartInfoServices not called within a second")
	}
}

func (s *InfoManagerServiceTestSuite) TestScriptsGet() {
	// setup mock
	s.objectStorageMock.On("GetPresignedObjectURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("http://myfakepresignedurl.com/object/1", nil)

	// setup requests
	testCases := []struct {
		name  string
		req   *proto.ScriptsGetRequest
		fail  bool
		error error
	}{
		{
			name: "shouldReturnScriptVersion1",
			req: &proto.ScriptsGetRequest{
				ScriptId: s.models.script.ID.String(),
				Version:  "1",
			},
			fail:  false,
			error: nil,
		},
		{
			name: "shouldReturnScriptVersion2",
			req: &proto.ScriptsGetRequest{
				ScriptId: s.models.script.ID.String(),
				Version:  "2",
			},
			fail:  false,
			error: nil,
		},
		{
			name: "shouldNotFindScriptVersion3",
			req: &proto.ScriptsGetRequest{
				ScriptId: s.models.script.ID.String(),
				Version:  "3",
			},
			fail:  true,
			error: nil,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			// get script
			resp, err := s.srv.ScriptsGet(context.Background(), testCase.req)
			if testCase.fail {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Equal("http://myfakepresignedurl.com/object/1", resp.SignedUrl)
			}
		})
	}
}

func (s *InfoManagerServiceTestSuite) TestScriptTest() {
	req := &proto.ScriptsTestRequest{
		Script: []byte("test"),
	}
	_, err := s.srv.ScriptsTest(context.Background(), req)
	s.Require().Error(err)
	s.ErrorIs(err, errgrpc.ErrUnimplemented(""))
}

func TestInfoManagerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InfoManagerServiceTestSuite))
}
