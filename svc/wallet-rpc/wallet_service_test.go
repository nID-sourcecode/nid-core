package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	"lab.weave.nl/nid/nid-core/svc/wallet-gql/models"
	"lab.weave.nl/nid/nid-core/svc/wallet-rpc/gqlclient"
	"lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
)

type WalletWalletServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	db *gorm.DB
	tx *gorm.DB

	authClientMock *gqlclient.AuthClientMock

	walletServer *WalletServer
}

func (s *WalletWalletServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
	s.authClientMock = &gqlclient.AuthClientMock{}

	s.walletServer = &WalletServer{
		db:         s.setupDB(),
		authClient: s.authClientMock,
	}
}

func (s *WalletWalletServiceTestSuite) setupDB() *WalletDB {
	s.tx = s.db.BeginTx(s.Ctx, nil)
	s.Require().NoError(s.tx.AutoMigrate(models.GetModels()...).Error)
	return &WalletDB{
		db:        s.tx,
		UserDB:    models.NewUserDB(s.tx),
		DeviceDB:  models.NewDeviceDB(s.tx),
		ConsentDB: models.NewConsentDB(s.tx),
		ClientDB:  models.NewClientDB(s.tx),
	}
}

func (s *WalletWalletServiceTestSuite) TestCreateConsent_UserNotFound() {
	s.createDummyUser()

	in := &proto.CreateConsentRequest{
		UserPseudo:  "some random pseudo",
		ClientId:    uuid.Must(uuid.NewV4()).String(),
		Description: "This is the description",
		Name:        "This is the name",
		AccessToken: "asdfasdfasdf",
	}

	_, err := s.walletServer.CreateConsent(s.Ctx, in)
	s.True(errors.Is(err, ErrUserNotFound))
}

func (s *WalletWalletServiceTestSuite) TestCreateConsent_FetchClientError() {
	user := s.createDummyUser()

	// Mock the fetch client
	s.authClientMock.On("FetchClient", s.Ctx, mock.Anything).Return(&gqlclient.Client{}, errors.New("alles is lek"))

	in := &proto.CreateConsentRequest{
		UserPseudo:  user.Pseudonym,
		ClientId:    uuid.Must(uuid.NewV4()).String(),
		Description: "This is the description",
		Name:        "This is the name",
		AccessToken: "asdfasdfasdf",
	}

	_, err := s.walletServer.CreateConsent(s.Ctx, in)
	s.True(errors.Is(err, ErrInternal))
}

func (s *WalletWalletServiceTestSuite) TestCreateConsent() {
	user := s.createDummyUser()

	// Mock the fetch client
	authClientID := uuid.Must(uuid.NewV4())
	authClient := &gqlclient.Client{
		Name:  "name",
		Color: "blue",
		Icon:  "icon",
		Logo:  "logo",
	}
	s.authClientMock.On("FetchClient", s.Ctx, mock.Anything).Return(authClient, nil)

	// GrantedProtoDate
	granted := time.Now().UTC()
	pGranted := timestamppb.New(granted)

	in := &proto.CreateConsentRequest{
		UserPseudo:  user.Pseudonym,
		ClientId:    authClientID.String(),
		Description: "This is the description",
		Name:        "This is the name",
		AccessToken: "asdfasdfasdf",
		GrantedAt:   pGranted,
	}

	out, err := s.walletServer.CreateConsent(s.Ctx, in)
	s.Require().NoError(err)

	// Assert mock correctly called
	s.authClientMock.AssertCalled(s.T(), "FetchClient", s.Ctx, authClientID)

	// Validate created consent
	consent, err := s.walletServer.db.ConsentDB.Get(s.Ctx, uuid.FromStringOrNil(out.Id))
	s.Require().NoError(err)

	s.Equal(user.ID, consent.UserID)
	s.Equal(in.Description, consent.Description)
	s.Equal(in.Name, consent.Name)
	s.Equal(in.AccessToken, consent.AccessToken)
	s.NotNil(*consent.Granted)

	// Validate client is created correctly
	client, err := s.walletServer.db.ClientDB.Get(s.Ctx, uuid.FromStringOrNil(out.ClientId))
	s.Require().NoError(err)

	s.Equal(authClientID.String(), client.ExtClientID)
	s.Equal(authClient.Name, client.Name)
	s.Equal(authClient.Color, client.Color)
	s.Equal(authClient.Icon, client.Icon)
	s.Equal(authClient.Logo, client.Logo)
}

func (s *WalletWalletServiceTestSuite) TestGetBSNForPseudonym_UserNotFound() {
	s.createDummyUser()

	in := &proto.GetBSNForPseudonymRequest{
		Pseudonym: "some random pseudo",
	}

	_, err := s.walletServer.GetBSNForPseudonym(s.Ctx, in)
	s.Require().Error(err)
	s.VerifyStatusError(err, codes.NotFound)
	s.Contains(err.Error(), "user not found")
}

func (s *WalletWalletServiceTestSuite) TestGetBSNForPseudonym() {
	s.createDummyUser()

	in := &proto.GetBSNForPseudonymRequest{
		Pseudonym: "abcdefghijk",
	}

	res, err := s.walletServer.GetBSNForPseudonym(s.Ctx, in)
	s.Require().NoError(err)

	s.Equal("123456789", res.Bsn)
}

func (s *WalletWalletServiceTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *WalletWalletServiceTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
	fmt.Println("Closing db connection")
}

func TestWalletWalletServiceTestSuite(t *testing.T) {
	suite.Run(t, &WalletWalletServiceTestSuite{
		db: database.MustConnectTest("auth", nil),
	})
}

func (s *WalletWalletServiceTestSuite) createDummyUser() *models.User {
	user := &models.User{
		Pseudonym: "abcdefghijk",
		Bsn:       "123456789",
		Password:  "hashedpassword",
	}
	err := s.tx.Create(user).Error
	s.Require().NoError(err)
	return user
}
