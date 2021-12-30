package main

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers"
	headersmock "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers/mock"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
)

var ErrAllesIsLek error = fmt.Errorf("alles is lek")

type AuthorizationServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	db *gorm.DB
	tx *gorm.DB

	authServer *AuthorizationServiceServer
}

func (s *AuthorizationServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
	priv, pub, err := jwt.GenerateTestKeys()
	s.Require().NoError(err)

	mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
	s.authServer = &AuthorizationServiceServer{
		metadataHelper: mockedMetadataHelper,
		db:             s.setupDB(),
		jwtClient:      jwt.NewJWTClient(priv, pub),
		pwManager:      password.NewDefaultManager(),
	}
	mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("wim@weave.nl", "alsdkfj#@4!", nil)
}

func (s *AuthorizationServiceTestSuite) setupDB() *DashboardDB {
	s.tx = s.db.BeginTx(s.Ctx, nil)
	s.Require().NoError(s.tx.AutoMigrate(models.GetModels()...).Error)
	dashboardDB := &DashboardDB{
		db:     s.tx,
		UserDB: models.NewUserDB(s.tx),
	}

	dashboardDB.migrate(&DashBoardConfig{
		DefaultUser: "test@weave.nl",
		PilotUser:   "pilot@weave.nl",
	})
	return dashboardDB
}

func (s *AuthorizationServiceTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *AuthorizationServiceTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
	fmt.Println("Closing db connection")
}

func (s *AuthorizationServiceTestSuite) TestSignInSuccess() {
	s.createDummyUser("wim@weave.nl", "alsdkfj#@4!")
	res, err := s.authServer.Signin(s.Ctx, &empty.Empty{})
	s.Require().NoError(err)

	claims, err := s.authServer.jwtClient.GetClaims(res.Bearer)
	s.Require().NoError(err)
	s.Require().NoError(claims.Valid())
	s.Equal("weave.nl", claims["iss"].(string))
}

func (s *AuthorizationServiceTestSuite) TestCantSignin() {
	s.createDummyUser("wim@weave.nl", "alsdkfj#@4!")

	tests := []struct {
		Name        string
		MetaHelper  func() headers.MetadataHelper
		ExpectedErr codes.Code
	}{
		{
			Name: "Email not found",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("wim2@weave.nl", "alsdkfj#@4!", nil)

				return mockedMetadataHelper
			},
			ExpectedErr: codes.InvalidArgument,
		},
		{
			Name: "Incorrect password",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("wim@weave.nl", "aldsdkfj#@4!", nil)

				return mockedMetadataHelper
			},
			ExpectedErr: codes.InvalidArgument,
		},
		{
			Name: "No basic auth",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("", "", ErrAllesIsLek)

				return mockedMetadataHelper
			},
			ExpectedErr: codes.InvalidArgument,
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			s.authServer.metadataHelper = test.MetaHelper()
			_, err := s.authServer.Signin(s.Ctx, &empty.Empty{})
			s.Require().Error(err)
			s.VerifyStatusError(err, test.ExpectedErr)
		})
	}
}

func TestAuthorizationServiceTestSuite(t *testing.T) {
	suite.Run(t, &AuthorizationServiceTestSuite{
		// Intentionally do not supply models to automigrate, this should be done inside the transaction
		db: database.MustConnectTest("auth", nil),
	})
}

func (s *AuthorizationServiceTestSuite) createDummyUser(email, pass string) {
	pwManager := password.NewDefaultManager()
	password, err := pwManager.GenerateHash(pass)
	s.Require().NoError(err)
	err = s.tx.Create(&models.User{
		Email:    email,
		Password: password,
	}).Error
	s.Require().NoError(err)
}
