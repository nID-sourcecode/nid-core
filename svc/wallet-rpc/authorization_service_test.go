package main

import (
	"errors"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"

	pw "github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/headers"
	headersmock "github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/headers/mock"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/jwt/v2"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
)

var errAllesIsLek = errors.New("alles is lek")

type WalletAuthorizationServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	db *gorm.DB
	tx *gorm.DB

	authServer           *AuthorizationServer
	mockedMetadataHelper *headersmock.GRPCMetadataHelperMock
	pwManager            pw.IManager
}

func (s *WalletAuthorizationServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
	priv, pub, err := jwt.GenerateTestKeys()
	s.Require().NoError(err)

	s.mockedMetadataHelper = &headersmock.GRPCMetadataHelperMock{}

	opts := jwt.DefaultOpts()
	opts.ClaimsOpts.Issuer = "wallet.ns"

	s.authServer = &AuthorizationServer{
		metadataHelper: s.mockedMetadataHelper,
		db:             s.setupDB(),
		jwtClient:      jwt.NewJWTClientWithOpts(priv, pub, opts),
		pwManager:      s.pwManager,
	}
}

func (s *WalletAuthorizationServiceTestSuite) setupDB() *WalletDB {
	s.tx = s.db.BeginTx(s.Ctx, nil)
	s.Require().NoError(s.tx.AutoMigrate(models.GetModels()...).Error)

	return &WalletDB{
		db:       s.tx,
		UserDB:   models.NewUserDB(s.tx),
		DeviceDB: models.NewDeviceDB(s.tx),
	}
}

func (s *WalletAuthorizationServiceTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *WalletAuthorizationServiceTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
}

func (s *WalletAuthorizationServiceTestSuite) TestRegisterApp() {
	user := s.createDummyUser()
	s.mockedMetadataHelper.On("GetBasicAuth", s.Ctx).Return("123456789", "darthvader123#", nil)
	res, err := s.authServer.RegisterDevice(s.Ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	device := &models.Device{}
	s.Require().NoError(s.tx.Find(&device).Where("code = ?", res.Code).Error)

	secretMatches, err := s.pwManager.ComparePassword(res.Secret, device.Secret)
	s.Require().NoError(err)
	s.True(secretMatches)
	s.Equal(user.ID, device.UserID)
}

func (s *WalletAuthorizationServiceTestSuite) TestCantRegisterApp() {
	s.createDummyUser()

	tests := []struct {
		Name               string
		MetaHelper         func() headers.MetadataHelper
		ExpectedErrCode    codes.Code
		ExpectedErrMessage string
	}{
		{
			Name: "Bsn not found",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("987654321", "alsdkfj#@4!", nil)

				return mockedMetadataHelper
			},
			ExpectedErrCode:    codes.InvalidArgument,
			ExpectedErrMessage: "incorrect username or password",
		},
		{
			Name: "Incorrect password",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("123456789", "lukeskywalker456!", nil)

				return mockedMetadataHelper
			},
			ExpectedErrCode:    codes.InvalidArgument,
			ExpectedErrMessage: "incorrect username or password",
		},
		{
			Name: "No basic auth",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("", "", errAllesIsLek)

				return mockedMetadataHelper
			},
			ExpectedErrCode:    codes.InvalidArgument,
			ExpectedErrMessage: "retrieving basic auth: alles is lek",
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			s.authServer.metadataHelper = test.MetaHelper()
			_, err := s.authServer.RegisterDevice(s.Ctx, &emptypb.Empty{})
			s.Require().Error(err)
			s.VerifyStatusError(err, test.ExpectedErrCode)
			s.Contains(err.Error(), test.ExpectedErrMessage)
		})
	}
}

func (s *WalletAuthorizationServiceTestSuite) TestSignInSuccess() {
	s.mockedMetadataHelper.On("GetBasicAuth", s.Ctx).Return("device1", "alsdkfj#@4!", nil)

	user := s.createDummyUser()
	s.createDummyDevice(user.ID, "device1", "alsdkfj#@4!")

	res, err := s.authServer.SignIn(s.Ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	claims, err := s.authServer.jwtClient.GetClaims(res.Bearer)
	s.Require().NoError(err)

	issuerClaim, ok := claims["iss"]
	s.Require().True(ok, "no issuer in token")
	issuer, ok := issuerClaim.(string)
	s.Require().True(ok, "issuer is not a string: %+v", issuerClaim)

	s.Equal("wallet.ns", issuer)

	subjectClaim, ok := claims["sub"]
	s.Require().True(ok, "no subject in token")
	subject, ok := subjectClaim.(string)
	s.Require().True(ok, "subject is not a string: %+v", subjectClaim)

	s.Equal("abcdefghijk", subject)
}

func (s *WalletAuthorizationServiceTestSuite) TestCantSignin() {
	user := s.createDummyUser()
	s.createDummyDevice(user.ID, "device1", "alsdkfj#@4!")

	tests := []struct {
		Name               string
		MetaHelper         func() headers.MetadataHelper
		ExpectedErrCode    codes.Code
		ExpectedErrMessage string
	}{
		{
			Name: "Email not found",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("device2", "alsdkfj#@4!", nil)

				return mockedMetadataHelper
			},
			ExpectedErrCode:    codes.InvalidArgument,
			ExpectedErrMessage: "incorrect username or password",
		},
		{
			Name: "Incorrect password",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("device1", "aldsdkfj#@4!", nil)

				return mockedMetadataHelper
			},
			ExpectedErrCode:    codes.InvalidArgument,
			ExpectedErrMessage: "incorrect username or password",
		},
		{
			Name: "No basic auth",
			MetaHelper: func() headers.MetadataHelper {
				mockedMetadataHelper := &headersmock.GRPCMetadataHelperMock{}
				mockedMetadataHelper.On("GetBasicAuth", mock.Anything).Return("", "", errAllesIsLek)

				return mockedMetadataHelper
			},
			ExpectedErrCode:    codes.InvalidArgument,
			ExpectedErrMessage: "retrieving basic auth: alles is lek",
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			s.authServer.metadataHelper = test.MetaHelper()
			_, err := s.authServer.SignIn(s.Ctx, &emptypb.Empty{})
			s.Require().Error(err)
			s.VerifyStatusError(err, test.ExpectedErrCode)
			s.Contains(err.Error(), test.ExpectedErrMessage)
		})
	}
}

func TestWalletAuthorizationServiceTestSuite(t *testing.T) {
	suite.Run(t, &WalletAuthorizationServiceTestSuite{
		db:        database.MustConnectTest("auth", nil),
		pwManager: pw.NewDefaultManager(),
	})
}

func (s *WalletAuthorizationServiceTestSuite) createDummyUser() *models.User {
	hashedPassword, err := s.pwManager.GenerateHash("darthvader123#")
	s.Require().NoError(err)
	user := &models.User{
		Pseudonym: "abcdefghijk",
		Bsn:       "123456789",
		Password:  hashedPassword,
	}
	err = s.tx.Create(user).Error
	s.Require().NoError(err)
	return user
}

func (s *WalletAuthorizationServiceTestSuite) createDummyDevice(userID uuid.UUID, deviceCode, deviceSecret string) {
	hashedSecret, err := s.pwManager.GenerateHash(deviceSecret)
	s.Require().NoError(err)
	err = s.tx.Create(&models.Device{
		UserID: userID,
		Code:   deviceCode,
		Secret: hashedSecret,
	}).Error
	s.Require().NoError(err)
}
