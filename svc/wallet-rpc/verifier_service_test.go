package main

import (
	"testing"
	"time"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	suite "github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/database/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/messagebird"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/postmark"
	"github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
)

type EmailVerifierServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	db *gorm.DB
	tx *gorm.DB

	verifierServer *VerifierServer
}

func (s *EmailVerifierServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	s.verifierServer = &VerifierServer{
		db:            s.setupDB(),
		emailVerifier: &postmark.MockedClient{},
		phoneVerifier: &messagebird.MockedClient{},
	}
}

func (s *EmailVerifierServiceTestSuite) setupDB() *WalletDB {
	s.tx = s.db.BeginTx(s.Ctx, nil)
	s.Require().NoError(s.tx.AutoMigrate(models.GetModels()...).Error)
	return &WalletDB{
		db:             s.tx,
		UserDB:         models.NewUserDB(s.tx),
		DeviceDB:       models.NewDeviceDB(s.tx),
		ConsentDB:      models.NewConsentDB(s.tx),
		ClientDB:       models.NewClientDB(s.tx),
		EmailAddressDB: models.NewEmailAddressDB(s.tx),
		PhoneNumberDB:  models.NewPhoneNumberDB(s.tx),
	}
}

func (s *EmailVerifierServiceTestSuite) TestVerifyEmailIdNotFound() {
	in := &proto.VerifyRequest{
		Id:   uuid.Must(uuid.NewV4()).String(),
		Code: "code",
	}

	_, err := s.verifierServer.VerifyEmail(s.Ctx, in)
	s.True(errors.Is(err, ErrInternal))
}

func (s *EmailVerifierServiceTestSuite) TestEmailAlreadyValid() {
	id, _ := uuid.NewV4()
	s.createDummyEmailRecord(id, true, time.Now())

	in := &proto.VerifyRequest{
		Id:   id.String(),
		Code: "code",
	}

	r, err := s.verifierServer.VerifyEmail(s.Ctx, in)
	s.NoError(err)
	s.Equal(id.String(), r.GetId())
}

func (s *EmailVerifierServiceTestSuite) TestEmailTokenExpired() {
	id, _ := uuid.NewV4()
	s.createDummyEmailRecord(id, false, time.Now().Add(-time.Minute*50000))

	in := &proto.VerifyRequest{
		Id:   id.String(),
		Code: "code",
	}

	_, err := s.verifierServer.VerifyEmail(s.Ctx, in)
	s.True(errors.Is(err, errTokenExpired))
}

func (s *EmailVerifierServiceTestSuite) TestEmailTokenDidNotMatch() {
	id, _ := uuid.NewV4()
	s.createDummyEmailRecord(id, false, time.Now())

	emailVerifierMock := &postmark.MockedClient{}

	emailVerifierMock.On("CheckEmailVerification", "token", "code").Times(1).Return(postmark.ErrTokenDidNotMatch)
	s.verifierServer.emailVerifier = emailVerifierMock

	in := &proto.VerifyRequest{
		Id:   id.String(),
		Code: "code",
	}

	_, err := s.verifierServer.VerifyEmail(s.Ctx, in)
	s.True(errors.Is(err, ErrInternal))
}

func (s *EmailVerifierServiceTestSuite) TestVerifyEmailServiceUnknownError() {
	id, _ := uuid.NewV4()
	s.createDummyEmailRecord(id, false, time.Now())

	emailVerifierMock := &postmark.MockedClient{}

	emailVerifierMock.On("CheckEmailVerification", "token", "code").Times(1).Return(errors.New("error"))
	s.verifierServer.emailVerifier = emailVerifierMock

	in := &proto.VerifyRequest{
		Id:   id.String(),
		Code: "code",
	}

	_, err := s.verifierServer.VerifyEmail(s.Ctx, in)
	s.True(errors.Is(err, ErrInternal))
}

func (s *EmailVerifierServiceTestSuite) TestVerifyEmailSucceed() {
	id, _ := uuid.NewV4()
	s.createDummyEmailRecord(id, false, time.Now())

	emailVerifierMock := &postmark.MockedClient{}

	emailVerifierMock.On("CheckEmailVerification", "token", "token").Times(1).Return(nil)
	s.verifierServer.emailVerifier = emailVerifierMock

	in := &proto.VerifyRequest{
		Id:   id.String(),
		Code: "token",
	}

	r, err := s.verifierServer.VerifyEmail(s.Ctx, in)
	s.NoError(err)
	s.Equal(id.String(), r.GetId())
}

func (s *EmailVerifierServiceTestSuite) TestRetryEmailIdNotFound() {
	in := &proto.RetryVerifyRequest{
		Id: uuid.Must(uuid.NewV4()).String(),
	}

	_, err := s.verifierServer.RetryVerifyEmail(s.Ctx, in)
	s.True(errors.Is(err, ErrInternal))
}

func (s *EmailVerifierServiceTestSuite) TestRetryEmailAlreadyValid() {
	id, _ := uuid.NewV4()

	in := &proto.RetryVerifyRequest{
		Id: id.String(),
	}

	s.createDummyEmailRecord(id, true, time.Now())
	_, err := s.verifierServer.RetryVerifyEmail(s.Ctx, in)
	s.True(errors.Is(err, errRetryAlreadyValid))
}

func (s *EmailVerifierServiceTestSuite) TestRetryEmailDebounce() {
	id, _ := uuid.NewV4()

	in := &proto.RetryVerifyRequest{
		Id: id.String(),
	}

	s.createDummyEmailRecord(id, false, time.Now().Add(-time.Second*179))
	_, err := s.verifierServer.RetryVerifyEmail(s.Ctx, in)
	s.NotNil(err)
}

func (s *EmailVerifierServiceTestSuite) TestRetryEmailNewVerificationFailed() {
	id, _ := uuid.NewV4()

	in := &proto.RetryVerifyRequest{
		Id: id.String(),
	}

	s.createDummyEmailRecord(id, false, time.Now().Add(-time.Second*181))

	emailVerifierMock := &postmark.MockedClient{}

	emailVerifierMock.On("NewEmailVerification", "email@email.com").Times(1).Return("", errors.New("error"))

	s.verifierServer.emailVerifier = emailVerifierMock

	_, err := s.verifierServer.RetryVerifyEmail(s.Ctx, in)

	emailVerifierMock.AssertNumberOfCalls(s.T(), "NewEmailVerification", 1)

	s.NotNil(err)
}

func (s *EmailVerifierServiceTestSuite) TestRetryEmailSucceed() {
	id, _ := uuid.NewV4()

	in := &proto.RetryVerifyRequest{
		Id: id.String(),
	}

	s.createDummyEmailRecord(id, false, time.Now().Add(-time.Second*181))

	emailVerifierMock := &postmark.MockedClient{}

	emailVerifierMock.On("NewEmailVerification", "email@email.com").Times(1).Return("token", nil)

	s.verifierServer.emailVerifier = emailVerifierMock

	r, err := s.verifierServer.RetryVerifyEmail(s.Ctx, in)

	emailVerifierMock.AssertNumberOfCalls(s.T(), "NewEmailVerification", 1)

	s.NoError(err)
	s.Equal(id.String(), r.GetId())
}

func (s *EmailVerifierServiceTestSuite) TestVerifyPhoneIdNotFound() {
	in := &proto.VerifyRequest{
		Id:   uuid.Must(uuid.NewV4()).String(),
		Code: "code",
	}

	_, err := s.verifierServer.VerifyPhoneNumber(s.Ctx, in)
	s.True(errors.Is(err, ErrInternal))
}

func (s *EmailVerifierServiceTestSuite) TestPhoneAlreadyValid() {
	id, _ := uuid.NewV4()
	s.createDummyPhoneRecord(id, true, time.Now())

	in := &proto.VerifyRequest{
		Id:   id.String(),
		Code: "code",
	}

	r, err := s.verifierServer.VerifyPhoneNumber(s.Ctx, in)
	s.NoError(err)
	s.Equal(id.String(), r.GetId())
}

func (s *EmailVerifierServiceTestSuite) TestVerifyPhoneServiceUnknownError() {
	id, _ := uuid.NewV4()
	s.createDummyPhoneRecord(id, false, time.Now())

	phoneVerifierMock := &messagebird.MockedClient{}

	phoneVerifierMock.On("CheckPhoneVerification", "token", "code").Times(1).Return(errors.New("error"))
	s.verifierServer.phoneVerifier = phoneVerifierMock

	in := &proto.VerifyRequest{
		Id:   id.String(),
		Code: "code",
	}

	_, err := s.verifierServer.VerifyPhoneNumber(s.Ctx, in)
	s.True(errors.Is(err, ErrInternal))
}

func (s *EmailVerifierServiceTestSuite) TestVerifyPhoneSucceed() {
	id, _ := uuid.NewV4()
	s.createDummyPhoneRecord(id, false, time.Now())

	phoneVerifierMock := &messagebird.MockedClient{}

	phoneVerifierMock.On("CheckPhoneVerification", "token", "token").Times(1).Return(nil)
	s.verifierServer.phoneVerifier = phoneVerifierMock

	in := &proto.VerifyRequest{
		Id:   id.String(),
		Code: "token",
	}

	r, err := s.verifierServer.VerifyPhoneNumber(s.Ctx, in)
	s.NoError(err)
	s.Equal(id.String(), r.GetId())
}

func (s *EmailVerifierServiceTestSuite) TestRetryPhoneIdNotFound() {
	in := &proto.RetryPhoneRequest{
		Id:               uuid.Must(uuid.NewV4()).String(),
		VerificationType: 0,
	}

	_, err := s.verifierServer.RetryVerifyPhoneNumber(s.Ctx, in)
	s.True(errors.Is(err, ErrInternal))
}

func (s *EmailVerifierServiceTestSuite) TestRetryPhoneAlreadyValid() {
	id, _ := uuid.NewV4()

	in := &proto.RetryPhoneRequest{
		Id:               id.String(),
		VerificationType: 0,
	}

	s.createDummyPhoneRecord(id, true, time.Now())
	_, err := s.verifierServer.RetryVerifyPhoneNumber(s.Ctx, in)
	s.True(errors.Is(err, errRetryAlreadyValid))
}

func (s *EmailVerifierServiceTestSuite) TestRetryPhoneDebounce() {
	id, _ := uuid.NewV4()

	in := &proto.RetryPhoneRequest{
		Id:               id.String(),
		VerificationType: 0,
	}

	s.createDummyPhoneRecord(id, false, time.Now().Add(-time.Second*79))
	_, err := s.verifierServer.RetryVerifyPhoneNumber(s.Ctx, in)
	s.NotNil(err)
}

func (s *EmailVerifierServiceTestSuite) TestRetryPhoneNewVerificationFailed() {
	id, _ := uuid.NewV4()

	in := &proto.RetryPhoneRequest{
		Id:               id.String(),
		VerificationType: 1,
	}

	phoneVerifierMock := &messagebird.MockedClient{}

	phoneVerifierMock.On("NewPhoneVerification", "0612345678", "SMS").Times(1).Return("", errors.New("error"))
	s.verifierServer.phoneVerifier = phoneVerifierMock

	s.createDummyPhoneRecord(id, false, time.Now().Add(-time.Second*90))
	_, err := s.verifierServer.RetryVerifyPhoneNumber(s.Ctx, in)
	s.NotNil(err)
}

func (s *EmailVerifierServiceTestSuite) TestRetryPhoneSuccessNewMethod() {
	id, _ := uuid.NewV4()

	in := &proto.RetryPhoneRequest{
		Id:               id.String(),
		VerificationType: 1,
	}

	phoneVerifierMock := &messagebird.MockedClient{}

	phoneVerifierMock.On("NewPhoneVerification", "0612345678", "SMS").Times(1).Return("token", nil)
	s.verifierServer.phoneVerifier = phoneVerifierMock

	s.createDummyPhoneRecord(id, false, time.Now().Add(-time.Second*90))
	r, err := s.verifierServer.RetryVerifyPhoneNumber(s.Ctx, in)
	s.NoError(err)
	s.Equal(id.String(), r.GetId())
}

func (s *EmailVerifierServiceTestSuite) TestRetryPhoneSuccess() {
	id, _ := uuid.NewV4()

	in := &proto.RetryPhoneRequest{
		Id:               id.String(),
		VerificationType: 2,
	}

	phoneVerifierMock := &messagebird.MockedClient{}

	phoneVerifierMock.On("NewPhoneVerification", "0612345678", "TTS").Times(1).Return("token", nil)
	s.verifierServer.phoneVerifier = phoneVerifierMock

	s.createDummyPhoneRecord(id, false, time.Now().Add(-time.Second*90))
	r, err := s.verifierServer.RetryVerifyPhoneNumber(s.Ctx, in)
	s.NoError(err)
	s.Equal(id.String(), r.GetId())
}

func (s *EmailVerifierServiceTestSuite) createDummyEmailRecord(id uuid.UUID, verified bool, updated time.Time) {
	userID, _ := uuid.NewV4()
	emailAddress := &models.EmailAddress{
		ID:                id,
		EmailAddress:      "email@email.com",
		User:              &models.User{},
		UserID:            userID,
		VerificationToken: "token",
		Verified:          verified,
		CreatedAt:         time.Now(),
		UpdatedAt:         updated,
	}

	err := s.tx.Create(emailAddress).Error
	s.Require().NoError(err)
}

func (s *EmailVerifierServiceTestSuite) createDummyPhoneRecord(id uuid.UUID, verified bool, updated time.Time) {
	userID, _ := uuid.NewV4()
	phoneNumber := &models.PhoneNumber{
		ID:                id,
		PhoneNumber:       "0612345678",
		User:              &models.User{},
		UserID:            userID,
		VerificationToken: "token",
		Verified:          verified,
		CreatedAt:         time.Now(),
		UpdatedAt:         updated,
		VerificationType:  2,
	}

	err := s.tx.Create(phoneNumber).Error
	s.Require().NoError(err)
}

func (s *EmailVerifierServiceTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *EmailVerifierServiceTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
	log.Info("Closing db connection")
}

func TestEmailVerifierServiceTestSuite(t *testing.T) {
	suite.Run(t, &EmailVerifierServiceTestSuite{
		db: database.MustConnectTest("auth", nil),
	})
}
