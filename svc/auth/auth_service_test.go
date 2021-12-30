package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"lab.weave.nl/nid/nid-core/pkg/authtoken"
	"lab.weave.nl/nid/nid-core/pkg/gqlutil"
	"lab.weave.nl/nid/nid-core/pkg/pseudonym"
	"lab.weave.nl/nid/nid-core/pkg/utilities/database/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	errgrpc "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	headersmock "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers/mock"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/metrics"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v3"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	callbackHandlerMock "lab.weave.nl/nid/nid-core/svc/auth/internal/callbackhandler/mock"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
	pb "lab.weave.nl/nid/nid-core/svc/auth/proto"
	walletPB "lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
	walletMock "lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto/mock"
)

type AuthServiceTestSuite struct {
	*AuthServiceBaseTestSuite

	srv                 *AuthServiceServer
	conf                *AuthConfig
	metadataHelperMock  *headersmock.GRPCMetadataHelperMock
	jwtClient           *jwt.Client
	mockPseudonymizer   *pseudonym.MockPseudonymizer
	mockWalletClient    *walletMock.WalletClient
	callbackHandlerMock *callbackHandlerMock.CallbackHandler
}

type schemaFetcher struct {
	mock.Mock

	testSchema *gqlutil.Schema
}

func (f *schemaFetcher) FetchSchema(ctx context.Context, url string) (*gqlutil.Schema, error) {
	args := f.Called(ctx, url)

	return f.testSchema, args.Error(1)
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.AuthServiceBaseTestSuite.SetupTest()
	s.srv.db = s.authDB
	s.SetupTokenTest()
	s.SetupRegisterTest()
}

func (s *AuthServiceTestSuite) SetupSuite() {
	s.AuthServiceBaseTestSuite.SetupSuite()
	c := &AuthConfig{}
	if err := envconfig.InitWithOptions(c, envconfig.Options{AllOptional: true}); err != nil {
		s.Failf("init conf failed", "%+v", err)
	}
	s.conf = c
	s.conf.AuthorizationCodeExpirationTime = time.Minute
	s.conf.AuthorizationCodeLength = 32
	s.conf.AuthRequestURI = "https://authrequest.com"

	s.callbackHandlerMock = &callbackHandlerMock.CallbackHandler{}

	s.srv = &AuthServiceServer{
		stats:           CreateStats(metrics.NewNopeScope()),
		wk:              nil,
		pseudonymizer:   nil,
		jwtClient:       s.jwtClient,
		schemaFetcher:   nil,
		walletClient:    nil,
		conf:            s.conf,
		metadataHelper:  nil,
		passwordManager: s.passwordManager,
		callbackhandler: s.callbackHandlerMock,
	}

	s.srv.metadataHelper = s.metadataHelperMock
}

func (s *AuthServiceTestSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *AuthServiceTestSuite) TearDownSuite() {
	s.NoError(s.db.Close())
	fmt.Println("Closing db connection")
}

func (s *AuthServiceTestSuite) mockSubject(subject string) {
	md := metadata.MD{}
	s.metadataHelperMock.On("MetadataFromCtx", s.Ctx).Return(md, nil)

	claimsJSON := fmt.Sprintf(`{"sub":"%s"}`, subject)
	claimsb64 := base64.RawURLEncoding.EncodeToString([]byte(claimsJSON))
	s.metadataHelperMock.On("GetMetadataValue", md, "claims").Return(claimsb64, nil)
}

// Accept
const (
	dummySubject = "sadasdasjkdhaiouysdg867ig672315471r23t7"
	issuer       = "issuer.nl"
)

func (s *AuthServiceTestSuite) createAcceptDummySession(claimed bool) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           s.audience.ID,
		ClientID:             s.client.ID,
		RedirectTargetID:     s.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		Subject:              dummySubject,
	}
	if claimed {
		session.State = models.SessionStateClaimed
	} else {
		session.State = models.SessionStateUnclaimed
	}
	err := s.tx.Create(&session).Error
	s.Require().NoError(err, "Did not expect an error on create session")
	err = s.tx.Model(&session).Association("RequiredAccessModels").Append(s.accessModelGql1, s.accessModelGql2).Error
	s.Require().NoError(err, "Did not expect an error on create association")
	err = s.tx.Model(&session).Association("OptionalAccessModels").Append(s.accessModelRestPOST).Error
	s.Require().NoError(err, "Did not expect an error on create association")

	return session
}

func (s *AuthServiceTestSuite) TestAcceptedClaimsUpdated() {
	s.mockSubject(dummySubject)
	testSession := s.createAcceptDummySession(true)
	resp, err := s.srv.Accept(s.Ctx, &pb.AcceptRequest{
		SessionId:      testSession.ID.String(),
		AccessModelIds: []string{s.accessModelRestPOST.ID.String()},
	})
	s.Require().NoError(err)

	var session models.Session
	err = s.tx.Where("id = ?", testSession.ID).Preload("AcceptedAccessModels").Find(&session).Error
	s.NoError(err)
	s.Len(session.AcceptedAccessModels, 3, "sessions should have 3 accepted access models (1) optional and (2) required")

	s.Equal(resp.Id, testSession.ID.String())
	s.Equal(resp.State, pb.SessionState_ACCEPTED)
	s.Require().NotNil(resp.Audience)
	s.Equal(resp.Audience.Id, s.audience.ID.String())
	s.Equal(resp.Audience.Audience, s.audience.Audience)
	s.Equal(resp.Audience.Namespace, s.audience.Namespace)

	if s.NotNil(resp.Client) {
		s.Equal(resp.Client.Id, s.client.ID.String())
		s.Equal(resp.Client.Name, s.client.Name)
		s.Equal(resp.Client.Icon, s.client.Icon)
		s.Equal(resp.Client.Logo, s.client.Logo)
		s.Equal(resp.Client.Color, s.client.Color)
	}
	if s.NotNil(resp.RequiredAccessModels) && s.NotNil(resp.RequiredAccessModels[0]) {
		s.Equal(resp.RequiredAccessModels[0].Id, s.accessModelGql1.ID.String())
		s.Equal(resp.RequiredAccessModels[0].Name, s.accessModelGql1.Name)
		s.Equal(resp.RequiredAccessModels[0].Hash, s.accessModelGql1.Hash)
		s.Equal(resp.RequiredAccessModels[0].Description, s.accessModelGql1.Description)
	}
	if s.NotNil(resp.OptionalAccessModels) && s.NotNil(resp.OptionalAccessModels[0]) {
		s.Equal(resp.OptionalAccessModels[0].Id, s.accessModelRestPOST.ID.String())
		s.Equal(resp.OptionalAccessModels[0].Name, s.accessModelRestPOST.Name)
		s.Equal(resp.OptionalAccessModels[0].Hash, s.accessModelRestPOST.Hash)
		s.Equal(resp.OptionalAccessModels[0].Description, s.accessModelRestPOST.Description)
	}
	if s.NotNil(resp.AcceptedAccessModels) && s.NotNil(resp.AcceptedAccessModels[0]) {
		s.Equal(resp.AcceptedAccessModels[0].Id, s.accessModelRestPOST.ID.String())
		s.Equal(resp.AcceptedAccessModels[0].Name, s.accessModelRestPOST.Name)
		s.Equal(resp.AcceptedAccessModels[0].Hash, s.accessModelRestPOST.Hash)
		s.Equal(resp.AcceptedAccessModels[0].Description, s.accessModelRestPOST.Description)
	}
}

func (s *AuthServiceTestSuite) TestAcceptAccessModels() {
	s.mockSubject(dummySubject)
	testSession := s.createAcceptDummySession(true)
	resp, err := s.srv.Accept(s.Ctx, &pb.AcceptRequest{
		SessionId:      testSession.ID.String(),
		AccessModelIds: []string{s.accessModelRestPOST.ID.String()},
	})
	s.Require().NoError(err)
	s.NotNil(resp)
	s.NotNil(resp.AcceptedAccessModels)
	s.Len(resp.AcceptedAccessModels, 3)
	s.NotEqualValues(resp.AcceptedAccessModels[0].Id, resp.AcceptedAccessModels[1].Id, resp.AcceptedAccessModels[2].Id)
	for _, a := range resp.AcceptedAccessModels {
		if s.True(a.Id == s.accessModelGql1.ID.String() || a.Id == s.accessModelGql2.ID.String() || a.Id == s.accessModelRestPOST.ID.String()) {
			continue
		}
		s.Failf("access_model id %s not found in response", a.Id)
	}
}

func (s *AuthServiceTestSuite) TestNoOptionalAccessModelsNotProvided() {
	s.mockSubject(dummySubject)
	testSession := s.createAcceptDummySession(true)

	resp, err := s.srv.Accept(s.Ctx, &pb.AcceptRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().NoError(err)
	s.NotNil(resp)
	s.NotNil(resp.AcceptedAccessModels)
	s.Len(resp.AcceptedAccessModels, 2)
	s.NotEqualValues(resp.AcceptedAccessModels[0].Id, resp.AcceptedAccessModels[1].Id)
	for _, a := range resp.AcceptedAccessModels {
		if s.True(a.Id == s.accessModelGql1.ID.String() || a.Id == s.accessModelGql2.ID.String()) {
			continue
		}
		s.Failf("access_model id %s not found in response", a.Id)
	}
}

func (s *AuthServiceTestSuite) TestErrorNonExistingAccessModel() {
	s.mockSubject(dummySubject)
	testSession := s.createAcceptDummySession(true)

	_, err := s.srv.Accept(s.Ctx, &pb.AcceptRequest{
		SessionId:      testSession.ID.String(),
		AccessModelIds: []string{uuid.Must(uuid.NewV4()).String()},
	})
	s.Require().Error(err)
}

func (s *AuthServiceTestSuite) TestErrorRequiredAccessModelIdInPayload() {
	s.mockSubject(dummySubject)
	testSession := s.createAcceptDummySession(false)
	_, err := s.srv.Accept(s.Ctx, &pb.AcceptRequest{
		SessionId:      testSession.ID.String(),
		AccessModelIds: []string{s.accessModelGql1.ID.String()},
	})
	s.Error(err)
	s.VerifyStatusError(err, codes.FailedPrecondition)
}

// Authorize
func (s *AuthServiceTestSuite) TestSessionCreated() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	s.Require().NoError(err)
	s.VerifyCreatedHeader("grpc-statuscode", "302")

	sessionID := s.RetrieveHeader("location")[len("https://authrequest.com#"):]
	session, err := s.srv.db.SessionDB.GetSessionByID(s.Ctx, preloadAll, sessionID, s.conf.AuthorizationCodeExpirationTime)
	s.Require().NoError(err)
	s.Require().NotZero(session.ID)
	s.VerifyCreatedHeader("location", fmt.Sprintf("https://authrequest.com#%s", session.ID.String()))
	s.Equal(s.redirectTarget.ID, session.RedirectTargetID)
	s.Equal(models.SessionStateUnclaimed, session.State)
	s.Equal(s.audience.ID, session.AudienceID)
	s.Equal("", session.Subject)
	s.Len(session.AcceptedAccessModels, 0)
	s.Len(session.OptionalAccessModels, 0)
	s.Require().Len(session.RequiredAccessModels, 1)
	s.Equal(s.accessModelGql1.ID, session.RequiredAccessModels[0].ID)
}

func (s *AuthServiceTestSuite) TestSessionCreatedWithOptionalScopes() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:          "openid test:stuff@abc",
		ResponseType:   "code",
		ClientId:       s.client.ID.String(),
		RedirectUri:    "https://weave.nl/code",
		Audience:       "https://test.com/gql",
		OptionalScopes: "test:stuff2@ghi test:stuff3@jkl",
	})
	s.Require().NoError(err)
	s.VerifyCreatedHeader("grpc-statuscode", "302")

	sessionID := s.RetrieveHeader("location")[len("https://authrequest.com#"):]
	session, err := s.srv.db.SessionDB.GetSessionByID(s.Ctx, preloadAll, sessionID, s.conf.AuthorizationCodeExpirationTime)
	s.Require().NoError(err)
	s.Require().NotZero(session.ID)

	s.VerifyCreatedHeader("location", fmt.Sprintf("https://authrequest.com#%s", session.ID.String()))
	s.Equal(s.redirectTarget.ID, session.RedirectTargetID)
	s.Equal(models.SessionStateUnclaimed, session.State)
	s.Equal(s.audience.ID, session.AudienceID)
	s.Equal("", session.Subject)
	s.Len(session.OptionalAccessModels, 2)
	foundIDs := make([]uuid.UUID, 0)
	for _, optAccessModel := range session.OptionalAccessModels {
		foundIDs = append(foundIDs, optAccessModel.ID)
	}
	s.Contains(foundIDs, s.accessModelGql2.ID)
	s.Contains(foundIDs, s.accessModelRestPOST.ID)

	s.Len(session.AcceptedAccessModels, 0)
	s.Require().Len(session.RequiredAccessModels, 1)
	s.Equal(s.accessModelGql1.ID, session.RequiredAccessModels[0].ID)
}

func (s *AuthServiceTestSuite) TestErrorOnNoScopes() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), "at least one access model-scope must be specified")
}

func (s *AuthServiceTestSuite) TestErrorOnNonExistentRedirectURI() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://wee.nl/code",
		Audience:     "https://test.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), fmt.Sprintf("client with id \"%s\" does not have redirect URI \"https://wee.nl/code\"", s.client.ID.String()))
}

func (s *AuthServiceTestSuite) TestErrorOnWrongRedirectURI() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave2.nl/code",
		Audience:     "https://test.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), fmt.Sprintf("client with id \"%s\" does not have redirect URI \"https://weave2.nl/code\"", s.client.ID.String()))
}

func (s *AuthServiceTestSuite) TestErrorOnNonExistentAudience() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://te.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), "audience \"https://te.com/gql\" does not exist")
}

func (s *AuthServiceTestSuite) TestErrorOnWrongAudience() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test2.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), "audience \"https://test2.com/gql\" does not support access model \"test:stuff@abc\"")
}

func (s *AuthServiceTestSuite) TestErrorOnResponseTypeOtherThanCode() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "token",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), "no response type other than \"code\" is supported")
}

func (s *AuthServiceTestSuite) TestErrorOnNonExistentAccessModel() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test:stf@ac",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), "audience \"https://test.com/gql\" does not support access model \"test:stf@ac\"")
}

func (s *AuthServiceTestSuite) TestErrorOnWrongAccessModelHash() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test:stuff@ac",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), "audience \"https://test.com/gql\" does not support access model \"test:stuff@ac\"")
}

func (s *AuthServiceTestSuite) TestErrorOnWrongAccessModel() {
	_, err := s.srv.Authorize(s.Ctx, &pb.AuthorizeRequest{
		Scope:        "openid test2:stuff2@def",
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), "audience \"https://test.com/gql\" does not support access model \"test2:stuff2@def\"")
}

// AuthorizeHeadeless
func (s *AuthServiceTestSuite) TestErrorOnWrongAudienceHeadless() {
	_, err := s.srv.AuthorizeHeadless(s.Ctx, &pb.AuthorizeHeadlessRequest{
		ResponseType:   "code",
		ClientId:       s.client.ID.String(),
		RedirectUri:    "https://weave.nl/code",
		Audience:       "audience",
		QueryModelJson: "{}",
		QueryModelPath: "/gql",
	})
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Contains(err.Error(), "audience \"audience\" does not exist")
}

// FIXME: Schema checker does not work with Vecozo's schema  (TODO: Create issue for this on lab.weave.nl) Related with auth_service.go AuthorizeHeadless function
//func (s *AuthServiceTestSuite) TestErrorOnInvalidAccessModelHeadless() {
//	_, err := s.srv.AuthorizeHeadless(s.Ctx, &pb.AuthorizeHeadlessRequest{
//		ResponseType:   "code",
//		ClientId:       s.client.ID.String(),
//		RedirectUri:    "https://weave.nl/code",
//		Audience:       "https://test.com/gql",
//		QueryModelJson: "{}",
//	})
//	s.VerifyStatusError(err, codes.InvalidArgument)
//	s.Contains(err.Error(), "Root Model \"R\" Not Specified")
//}

func (s *AuthServiceTestSuite) TestErrorOnHandleCallbackHeadless() {
	s.callbackHandlerMock.On("HandleCallback", mock.Anything, "https://weave.nl/code", mock.Anything).Return(errors.New("error"))

	_, err := s.srv.AuthorizeHeadless(s.Ctx, &pb.AuthorizeHeadlessRequest{
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
		QueryModelJson: `{
			"r": {
				"m": {
					"users": "#U"
				}
			},
			"U": {
				"f": ["firstName", "lastName"],
				"p": {
					"filter": {
						"bsn": {
							"eq": "$$nid:bsn$$"
						}
					}
				}
			}
		}`,
		QueryModelPath: "/gql",
	})
	s.VerifyStatusError(err, codes.Internal)
	s.Contains(err.Error(), "internal server error")
}

func (s *AuthServiceTestSuite) TestAuthHeadless() {
	s.callbackHandlerMock.On("HandleCallback", mock.Anything, "https://weave.nl/code", mock.Anything).Return(nil)

	_, err := s.srv.AuthorizeHeadless(s.Ctx, &pb.AuthorizeHeadlessRequest{
		ResponseType: "code",
		ClientId:     s.client.ID.String(),
		RedirectUri:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
		QueryModelJson: `{
			"r": {
				"m": {
					"users": "#U"
				}
			},
			"U": {
				"f": ["firstName", "lastName"],
				"p": {
					"filter": {
						"bsn": {
							"eq": "$$nid:bsn$$"
						}
					}
				}
			}
		}
	`,
		QueryModelPath: "/gql",
	})
	s.Nil(err)
}

// Claim
func (s *AuthServiceTestSuite) createClaimDummySession(claimed bool) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           s.audience.ID,
		ClientID:             s.client.ID,
		RedirectTargetID:     s.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		Subject:              dummySubject,
	}
	if claimed {
		session.State = models.SessionStateClaimed
	} else {
		session.State = models.SessionStateUnclaimed
	}
	err := s.tx.Create(&session).Error
	s.Require().NoError(err, "Did not expect an error on create session")
	err = s.tx.Model(&session).Association("RequiredAccessModels").Append(s.accessModelGql1, s.accessModelGql2).Error
	s.Require().NoError(err, "Did not expect an error on create association")
	err = s.tx.Model(&session).Association("OptionalAccessModels").Append(s.accessModelRestPOST).Error
	s.Require().NoError(err, "Did not expect an error on create association")

	return session
}

func (s *AuthServiceTestSuite) TestClaimSessionUpdated() {
	s.mockSubject(dummySubject)
	testSession := s.createClaimDummySession(false)
	resp, err := s.srv.Claim(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().NoError(err)

	var session models.Session
	s.tx.First(&session, "id = ?", testSession.ID)

	s.NotZero(session.ID)
	s.Equal(session.ID, testSession.ID)
	s.Equal(session.State, models.SessionStateClaimed)
	s.Equal(resp.State, pb.SessionState_CLAIMED)
}

func (s *AuthServiceTestSuite) TestClaimSessionResponse() {
	s.mockSubject(dummySubject)
	testSession := s.createClaimDummySession(false)
	resp, err := s.srv.Claim(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().NoError(err)

	s.Equal(resp.Id, testSession.ID.String())
	s.Equal(resp.State, pb.SessionState_CLAIMED)
	if s.NotNil(resp.Audience) {
		s.Equal(resp.Audience.Id, s.audience.ID.String())
		s.Equal(resp.Audience.Audience, s.audience.Audience)
		s.Equal(resp.Audience.Namespace, s.audience.Namespace)
	}
	if s.NotNil(resp.Client) {
		s.Equal(resp.Client.Id, s.client.ID.String())
		s.Equal(resp.Client.Name, s.client.Name)
		s.Equal(resp.Client.Icon, s.client.Icon)
		s.Equal(resp.Client.Logo, s.client.Logo)
		s.Equal(resp.Client.Color, s.client.Color)
	}
	if s.NotNil(resp.RequiredAccessModels) && s.NotNil(resp.RequiredAccessModels[0]) {
		s.Equal(resp.RequiredAccessModels[0].Id, s.accessModelGql1.ID.String())
		s.Equal(resp.RequiredAccessModels[0].Name, s.accessModelGql1.Name)
		s.Equal(resp.RequiredAccessModels[0].Hash, s.accessModelGql1.Hash)
		s.Equal(resp.RequiredAccessModels[0].Description, s.accessModelGql1.Description)
	}
	if s.NotNil(resp.OptionalAccessModels) && s.NotNil(resp.OptionalAccessModels[0]) {
		s.Equal(resp.OptionalAccessModels[0].Id, s.accessModelRestPOST.ID.String())
		s.Equal(resp.OptionalAccessModels[0].Name, s.accessModelRestPOST.Name)
		s.Equal(resp.OptionalAccessModels[0].Hash, s.accessModelRestPOST.Hash)
		s.Equal(resp.OptionalAccessModels[0].Description, s.accessModelRestPOST.Description)
	}
	s.Nil(resp.AcceptedAccessModels)
}

func (s *AuthServiceTestSuite) TestErrorClaimedStateSession() {
	s.mockSubject(dummySubject)
	testSession := s.createClaimDummySession(true)
	_, err := s.srv.Claim(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().Error(err)
}

func (s *AuthServiceTestSuite) TestErrorEmptySubject() {
	s.mockSubject("")
	testSession := s.createClaimDummySession(true)
	_, err := s.srv.Claim(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().Error(err)
}

// createFinaliseDummySession returns the new session with a token.
func (s *AuthServiceTestSuite) createFinaliseDummySession(accepted bool) (*models.Session, string) {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           s.audience.ID,
		ClientID:             s.client.ID,
		RedirectTargetID:     s.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
	}
	if accepted {
		session.State = models.SessionStateAccepted
	} else {
		session.State = models.SessionStateRejected
	}
	err := s.tx.Create(&session).Error
	s.Require().NoError(err, "Did not expect an error")

	token, err := s.setupSessionFinaliseToken(session)
	s.Require().NoError(err, "Did not expect an error")

	return session, token
}

func (s *AuthServiceTestSuite) TestAcceptedSessionStateUpdated() {
	testSession, token := s.createFinaliseDummySession(true)

	response, err := s.srv.Finalise(s.Ctx, &pb.FinaliseRequest{SessionId: testSession.ID.String(), SessionFinaliseToken: token})
	s.Require().NoError(err)

	s.NotNil(response)
	s.NotEmpty(response.RedirectLocation)
	u, err := url.Parse(response.GetRedirectLocation())
	s.Require().NoError(err)
	authCode := u.Query().Get("authorization_code")
	s.NotEmpty(authCode)

	var session models.Session
	s.tx.First(&session, "id = ?", testSession.ID)
	s.NotZero(session.ID)
	s.NotNil(session.AuthorizationCode)
	s.Equal(session.ID, testSession.ID)

	hash, err := authtoken.Hash(authCode)
	s.Require().NoError(err)
	s.Equal(*session.AuthorizationCode, hash)

	s.Equal(s.redirectTarget.RedirectTarget+"/?authorization_code="+authCode, response.GetRedirectLocation())
	s.Equal(session.State, models.SessionStateCodeGranted)
}

func (s *AuthServiceTestSuite) TestNonAcceptedSessionError() {
	testSession, token := s.createFinaliseDummySession(false)

	_, err := s.srv.Finalise(s.Ctx, &pb.FinaliseRequest{SessionId: testSession.ID.String(), SessionFinaliseToken: token})
	s.Error(err)
}

// Get_Session_Details
func (s *AuthServiceTestSuite) createGetSessionDetailsDummySession(state models.SessionState) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           s.audience.ID,
		ClientID:             s.client.ID,
		RedirectTargetID:     s.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		Subject:              "sadasdasjkdhaiouysdg867ig672315471r23t7",
		State:                state,
	}

	err := s.tx.Create(&session).Error
	s.Require().NoError(err, "Did not expect an error on create session")
	err = s.tx.Model(&session).Association("RequiredAccessModels").Append(s.accessModelGql1, s.accessModelGql2).Error
	s.Require().NoError(err, "Did not expect an error on create association")
	err = s.tx.Model(&session).Association("OptionalAccessModels").Append(s.accessModelRestPOST).Error
	s.Require().NoError(err, "Did not expect an error on create association")

	return session
}

func (s *AuthServiceTestSuite) TestDetailsSessionResponse() {
	testSession := s.createGetSessionDetailsDummySession(models.SessionStateUnclaimed)
	resp, err := s.srv.GetSessionDetails(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
}

func (s *AuthServiceTestSuite) TestWrongState() {
	testSession := s.createGetSessionDetailsDummySession(models.SessionStateAccepted)
	_, err := s.srv.GetSessionDetails(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().Error(err)
}

// Register
func (s *AuthServiceTestSuite) SetupRegisterTest() {
	testSchema := &gqlutil.Schema{
		Types: map[string]*gqlutil.Type{
			"AddressFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"Geography": {
				Fields: map[string]*gqlutil.Field{},
			},
			"SavingsAccountFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"UUIDFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"GeographyFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"FloatFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"UUID": {
				Fields: map[string]*gqlutil.Field{},
			},
			"ContactDetailFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"Float": {
				Fields: map[string]*gqlutil.Field{},
			},
			"Time": {
				Fields: map[string]*gqlutil.Field{},
			},
			"BankAccount": {
				Fields: map[string]*gqlutil.Field{
					"id": {
						IsModel:  false,
						TypeName: "UUID",
					},
					"savingsAccounts": {
						IsModel:  true,
						TypeName: "SavingsAccount",
					},
					"userId": {
						IsModel:  false,
						TypeName: "UUID",
					},
					"accountNumber": {
						IsModel:  false,
						TypeName: "String",
					},
					"amount": {
						IsModel:  false,
						TypeName: "Int",
					},
					"createdAt": {
						IsModel:  false,
						TypeName: "Time",
					},
					"deletedAt": {
						IsModel:  false,
						TypeName: "Time",
					},
					"updatedAt": {
						IsModel:  false,
						TypeName: "Time",
					},
					"user": {
						IsModel:  true,
						TypeName: "User",
					},
				},
			},
			"Decimal": {
				Fields: map[string]*gqlutil.Field{},
			},
			"Map": {
				Fields: map[string]*gqlutil.Field{},
			},
			"Query": {
				Fields: map[string]*gqlutil.Field{
					"address": {
						IsModel:  true,
						TypeName: "Address",
					},
					"bankAccount": {
						IsModel:  true,
						TypeName: "BankAccount",
					},
					"contactDetails": {
						IsModel:  true,
						TypeName: "ContactDetail",
					},
					"user": {
						IsModel:  true,
						TypeName: "User",
					},
					"users": {
						IsModel:  true,
						TypeName: "User",
					},
					"addresses": {
						IsModel:  true,
						TypeName: "Address",
					},
					"bankAccounts": {
						IsModel:  true,
						TypeName: "BankAccount",
					},
					"contactDetail": {
						IsModel:  true,
						TypeName: "ContactDetail",
					},
					"savingsAccount": {
						IsModel:  true,
						TypeName: "SavingsAccount",
					},
					"savingsAccounts": {
						IsModel:  true,
						TypeName: "SavingsAccount",
					},
				},
			},
			"JSONFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"ID": {
				Fields: map[string]*gqlutil.Field{},
			},
			"ContactDetail": {
				Fields: map[string]*gqlutil.Field{},
			},
			"BankAccountFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"UserFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"Boolean": {
				Fields: map[string]*gqlutil.Field{},
			},
			"StringFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"TimeFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"JSON": {
				Fields: map[string]*gqlutil.Field{},
			},
			"String": {
				Fields: map[string]*gqlutil.Field{},
			},
			"BooleanFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"DecimalFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
			"Int": {
				Fields: map[string]*gqlutil.Field{},
			},
			"Address": {
				Fields: map[string]*gqlutil.Field{},
			},
			"SavingsAccount": {
				Fields: map[string]*gqlutil.Field{},
			},
			"User": {
				Fields: map[string]*gqlutil.Field{
					"updatedAt": {
						IsModel:  false,
						TypeName: "Time",
					},
					"contactDetails": {
						IsModel:  false,
						TypeName: "ContactDetail",
					},
					"createdAt": {
						IsModel:  false,
						TypeName: "Time",
					},
					"deletedAt": {
						IsModel:  false,
						TypeName: "Time",
					},
					"firstName": {
						IsModel:  false,
						TypeName: "String",
					},
					"pseudonym": {
						IsModel:  false,
						TypeName: "String",
					},
					"id": {
						IsModel:  false,
						TypeName: "UUID",
					},
					"bankAccounts": {
						IsModel:  true,
						TypeName: "BankAccount",
					},
					"lastName": {
						IsModel:  false,
						TypeName: "String",
					},
				},
			},
			"IntFilterInput": {
				Fields: map[string]*gqlutil.Field{},
			},
		},
	}

	fetcher := &schemaFetcher{}
	fetcher.testSchema = testSchema
	fetcher.On("FetchSchema", mock.Anything, mock.Anything).Return(testSchema, nil)
	s.srv.schemaFetcher = fetcher
}

func (s *AuthServiceTestSuite) TestRegisterAccessModel() {
	tests := []struct {
		testName            string
		inputAudience       string
		inputQueryModelJSON string
		inputScopeName      string
		inputDescription    string
		err                 bool
	}{
		{
			testName:            "successfullyInsert",
			inputAudience:       "https://test.com/gql",
			inputQueryModelJSON: "{ \"r\": { \"M\": { \"users\": \"#U\" } }, \"U\": { \"M\": { \"bankAccounts\": \"#B\" }, \"F\": [ \"firstName\", \"lastName\" ] }, \"B\": { \"F\": [ \"accountNumber\", \"amount\" ] } }",
			inputScopeName:      "databron:user",
			inputDescription:    "lorem ipsum",
			err:                 false,
		},
		{
			testName:            "missingField",
			inputAudience:       "https://test2.com/gql",
			inputQueryModelJSON: "{ \"r\": { \"M\": { \"users\": \"#U\" } }, \"U\": { \"M\": { \"bankAccounts\": \"#B\" }, \"F\": [ \"firstName\", \"lastName\", \"licensePlate\" ] }, \"B\": { \"F\": [ \"accountNumber\", \"amount\" ] } }",
			inputScopeName:      "databron:user2",
			inputDescription:    "lorem ipsum",
			err:                 true,
		},
	}
	for _, test := range tests {
		_, err := s.srv.RegisterAccessModel(s.Ctx, &pb.AccessModelRequest{
			Audience:       test.inputAudience,
			QueryModelJson: test.inputQueryModelJSON,
			ScopeName:      test.inputScopeName,
			Description:    test.inputDescription,
		})
		if test.err {
			s.Require().Error(err)
		} else {
			s.Require().NoError(err)
		}
	}
}

// Reject
func (s *AuthServiceTestSuite) createRejectDummySession(rejected bool) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           s.audience.ID,
		ClientID:             s.client.ID,
		RedirectTargetID:     s.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		Subject:              dummySubject,
	}
	if rejected {
		session.State = models.SessionStateRejected
	} else {
		session.State = models.SessionStateClaimed
	}
	err := s.tx.Create(&session).Error
	s.Require().NoError(err, "Did not expect an error")

	return session
}

func (s *AuthServiceTestSuite) TestRejectSession() {
	s.mockSubject(dummySubject)
	testSession := s.createRejectDummySession(false)
	_, err := s.srv.Reject(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().NoError(err)

	var session models.Session
	s.tx.First(&session, "id = ?", testSession.ID)

	s.NotZero(session.ID)
	s.Equal(session.ID, testSession.ID)
	s.Equal(session.State, models.SessionStateRejected)
}

func (s *AuthServiceTestSuite) TestErrorIncorrectSessionState() {
	s.mockSubject(dummySubject)
	testSession := s.createRejectDummySession(true)
	_, err := s.srv.Reject(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Error(err)
}

// Session Password
func (s *AuthServiceTestSuite) TestCreateSessionFinaliseToken() {
	session := s.createAcceptDummySession(false)
	_, err := s.setupSessionFinaliseToken(session)
	s.NoError(err, "Could not setup password for session")
}

func (s *AuthServiceTestSuite) TestCorrectSessionFinaliseToken() {
	testSession := s.createAcceptDummySession(false)
	token, err := s.setupSessionFinaliseToken(testSession)
	s.Require().NoError(err)

	session, err := s.authDB.SessionDB.GetSessionByID(s.Ctx, noPreload, testSession.ID.String(), s.conf.AuthorizationCodeExpirationTime)
	s.Require().NoError(err)
	ok, err := s.passwordManager.ComparePassword(token, session.FinaliseToken)
	s.True(ok)
	s.NoError(err)
}

// createPasswordDummySession adds token as a password to session
// returns the password
func (s *AuthServiceTestSuite) setupSessionFinaliseToken(session *models.Session) (string, error) {
	token, err := authtoken.NewToken(s.conf.AuthorizationCodeLength)
	if err != nil {
		return "", err
	}

	hash, err := s.passwordManager.GenerateHash(token)
	if err != nil {
		return "", err
	}

	err = s.authDB.SessionDB.SetSessionFinaliseToken(session, hash)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Status
func (s *AuthServiceTestSuite) createStatusDummySession(state models.SessionState) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           s.audience.ID,
		ClientID:             s.client.ID,
		RedirectTargetID:     s.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		State:                state,
	}
	err := s.tx.Create(&session).Error
	s.Require().NoError(err, "Did not expect an error")

	return session
}

// Swap
func (s *AuthServiceTestSuite) TestSwapTokenValidate() {
	missingToken := pb.SwapTokenRequest{
		Query:    "query",
		Audience: "test",
	}
	s.Error(missingToken.Validate())
	requestMissingQuery := pb.SwapTokenRequest{
		CurrentToken: "Token",
		Audience:     "someone",
	}
	s.Error(requestMissingQuery.Validate())
	missingAudience := pb.SwapTokenRequest{
		CurrentToken: "Token",
		Query:        "{}",
	}
	s.Error(missingAudience.Validate())
}

func (s *AuthServiceTestSuite) TestInvalidToken() {
	invalidToken := pb.SwapTokenRequest{
		CurrentToken: "someinvalidtoken",
		Query:        "{}",
	}
	resp, err := s.srv.SwapToken(context.Background(), &invalidToken)
	s.VerifyStatusError(err, codes.InvalidArgument)
	s.Nil(resp)
}

func (s *AuthServiceTestSuite) TestValidToken() {
	clientID := uuid.Must(uuid.NewV4())
	subID := uuid.Must(uuid.NewV4())

	defaultClaims := jwt.NewDefaultClaims()
	defaultClaims.Subject = subID.String()
	defaultClaims.Audience = []string{"informationservice-A"}
	defaultClaims.Issuer = issuer
	claims := &TokenClaims{
		DefaultClaims: defaultClaims,
		ClientID:      clientID.String(),
	}

	token, err := s.srv.jwtClient.SignToken(claims)
	s.Require().NoError(err)

	req := &pb.SwapTokenRequest{
		CurrentToken: token,
		Query:        "{age}",
		Audience:     s.audience.Audience,
	}

	resp, err := s.srv.SwapToken(context.Background(), req)
	s.Require().NoError(err)

	s.Equal("Bearer", resp.TokenType)
	newClaims := &TokenClaims{}
	err = s.srv.jwtClient.ValidateAndParseClaims(resp.AccessToken, newClaims)
	s.Require().NoError(err)
	s.Require().Len(newClaims.Audience, 1)
	s.Equal(s.audience.Audience, newClaims.Audience[0])
	// s.Equal("[{age}]", newToken["queries"])
	s.Equal(clientID.String(), newClaims.ClientID)
	s.Equal(subID.String(), newClaims.Subject)
}

func (s *AuthServiceTestSuite) TestGetStatusNewlyCreatedSession() {
	testSession := s.createStatusDummySession(models.SessionStateRejected)
	resp, err := s.srv.Status(s.Ctx, &pb.SessionRequest{
		SessionId: testSession.ID.String(),
	})
	s.Require().NoError(err)

	var session models.Session
	s.tx.First(&session, "id = ?", testSession.ID)
	s.NotZero(session.ID)
	s.Equal(session.ID, testSession.ID)
	s.Equal(session.State, models.SessionStateRejected)
	s.Equal(resp.State, pb.SessionState_REJECTED)
}

func setupMockPseudonimizer(ctx context.Context) *pseudonym.MockPseudonymizer {
	mockPseudonymizer := &pseudonym.MockPseudonymizer{}
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym1", "alice").Return("translatedPseudo1ForAlice", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym1", "nid").Return("translatedPseudo1ForNid", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym1", "bob").Return("translatedPseudo1ForBob", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym2", "alice").Return("translatedPseudo2ForAlice", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym2", "nid").Return("translatedPseudo2ForNid", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym2", "bob").Return("translatedPseudo2ForBob", nil)
	return mockPseudonymizer
}

// Token
func (s *AuthServiceTestSuite) SetupTokenTest() {
	priv, pub, err := jwt.GenerateTestKeys()
	s.Require().NoError(err, "error generating test keys")
	s.conf.Namespace = "nid"
	s.conf.Issuer = issuer
	opts := jwt.DefaultOpts()
	s.jwtClient = jwt.NewJWTClientWithOpts(priv, pub, opts)

	s.mockPseudonymizer = setupMockPseudonimizer(s.Ctx)

	s.mockWalletClient = &walletMock.WalletClient{}
	// FIXME test this mock properly (https://lab.weave.nl/twi/core/-/issues/109)
	s.mockWalletClient.On("CreateConsent", s.Ctx, mock.Anything).Return(&walletPB.ConsentResponse{}, nil)

	s.metadataHelperMock = &headersmock.GRPCMetadataHelperMock{}

	s.callbackHandlerMock = &callbackHandlerMock.CallbackHandler{}

	s.srv = &AuthServiceServer{
		db:              s.authDB,
		stats:           CreateStats(metrics.NewNopeScope()),
		wk:              nil,
		pseudonymizer:   s.mockPseudonymizer,
		jwtClient:       s.jwtClient,
		schemaFetcher:   nil,
		walletClient:    s.mockWalletClient,
		conf:            s.conf,
		metadataHelper:  s.metadataHelperMock,
		passwordManager: s.passwordManager,
		callbackhandler: s.callbackHandlerMock,
	}
}

func (s *AuthServiceTestSuite) TestTokenCreated() {
	tests := []struct {
		Name                   string
		Session                *models.Session
		ClientPassword         string
		ExpectedSubjects       map[string]interface{}
		ExpectedScopes         map[string]interface{}
		ExpectedAudience       string
		ExpectedSubject        string
		ExpectedClientID       string
		ExpectedClientMetadata map[string]interface{}
		PseudoCalls            int
	}{
		{
			Name:           "SuccessCase1",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{s.accessModelGql1},
				Audience:             s.audience,
				Client:               s.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       s.redirectTarget,
				Subject:              "pseudonym1",
			},
			ExpectedSubjects: map[string]interface{}{
				"alice": "translatedPseudo1ForAlice",
			},
			ExpectedScopes: map[string]interface{}{
				"test:stuff": map[string]interface{}{
					"t": "GQL",
					"p": "/gql",
					"m": map[string]interface{}{
						"r": "somestuff",
					},
				},
			},
			ExpectedAudience:       s.audience.Audience,
			ExpectedClientID:       s.client.ID.String(),
			ExpectedSubject:        "translatedPseudo1ForNid",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            2,
		},
		{
			Name:           "SuccessCase2",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{s.accessModelGql1},
				RequiredAccessModels: []*models.AccessModel{s.accessModelGql2},
				Audience:             s.audience,
				Client:               s.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       s.redirectTarget,
				Subject:              "pseudonym2",
			},
			ExpectedSubjects: map[string]interface{}{
				"alice": "translatedPseudo2ForAlice",
			},
			ExpectedScopes: map[string]interface{}{
				"test:stuff": map[string]interface{}{
					"t": "GQL",
					"p": "/gql",
					"m": map[string]interface{}{
						"r": "somestuff",
					},
				},
				"test:stuff2": map[string]interface{}{
					"t": "GQL",
					"p": "/graphql",
					"m": map[string]interface{}{
						"r": "somemorestuff",
					},
				},
			},
			ExpectedAudience:       s.audience.Audience,
			ExpectedClientID:       s.client.ID.String(),
			ExpectedSubject:        "translatedPseudo2ForNid",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            2,
		},
		{
			Name:           "SuccessCase3",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{s.accessModelGql1},
				RequiredAccessModels: []*models.AccessModel{s.accessModelRestPOST},
				Audience:             s.audience,
				Client:               s.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       s.redirectTarget,
				Subject:              "pseudonym2",
			},
			ExpectedSubjects: map[string]interface{}{
				"alice": "translatedPseudo2ForAlice",
			},
			ExpectedScopes: map[string]interface{}{
				"test:stuff": map[string]interface{}{
					"t": "GQL",
					"p": "/gql",
					"m": map[string]interface{}{
						"r": "somestuff",
					},
				},
				"test:stuff3": map[string]interface{}{
					"t": "REST",
					"p": "/some/rest/endpoint",
					"m": "POST",
					"b": "something",
					"q": map[string]interface{}{},
				},
			},
			ExpectedAudience:       s.audience.Audience,
			ExpectedClientID:       s.client.ID.String(),
			ExpectedSubject:        "translatedPseudo2ForNid",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            2,
		},
		{
			Name:           "SuccessCase4",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{s.accessModelRestPOST},
				RequiredAccessModels: []*models.AccessModel{s.accessModelRestGET},
				Audience:             s.audience,
				Client:               s.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       s.redirectTarget,
				Subject:              "pseudonym2",
			},
			ExpectedSubjects: map[string]interface{}{
				"alice": "translatedPseudo2ForAlice",
			},
			ExpectedScopes: map[string]interface{}{
				"test:stuff3": map[string]interface{}{
					"t": "REST",
					"p": "/some/rest/endpoint",
					"m": "POST",
					"b": "something",
					"q": map[string]interface{}{},
				},
				"test2:stuff2": map[string]interface{}{
					"t": "REST",
					"p": "/some/rest/endpoint",
					"m": "GET",
					"b": "",
					"q": map[string]interface{}{
						"something": "somethingelse",
					},
				},
			},
			ExpectedAudience:       s.audience.Audience,
			ExpectedClientID:       s.client.ID.String(),
			ExpectedSubject:        "translatedPseudo2ForNid",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            2,
		},
		{
			Name:           "SuccessCase5EmptySubject",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{s.accessModelRestPOST},
				RequiredAccessModels: []*models.AccessModel{s.accessModelRestGET},
				Audience:             s.audience,
				Client:               s.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       s.redirectTarget,
				Subject:              "",
			},
			ExpectedSubjects: map[string]interface{}{
				"alice": "",
			},
			ExpectedScopes: map[string]interface{}{
				"test:stuff3": map[string]interface{}{
					"t": "REST",
					"p": "/some/rest/endpoint",
					"m": "POST",
					"b": "something",
					"q": map[string]interface{}{},
				},
				"test2:stuff2": map[string]interface{}{
					"t": "REST",
					"p": "/some/rest/endpoint",
					"m": "GET",
					"b": "",
					"q": map[string]interface{}{
						"something": "somethingelse",
					},
				},
			},
			ExpectedAudience:       s.audience.Audience,
			ExpectedClientID:       s.client.ID.String(),
			ExpectedSubject:        "",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            0,
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			s.mockPseudonymizer = setupMockPseudonimizer(s.Ctx)
			s.srv.pseudonymizer = s.mockPseudonymizer

			authorizationCode, err := authtoken.NewToken(s.conf.AuthorizationCodeLength)
			s.Require().NoError(err, "error creating test authorization code")
			hash, err := authtoken.Hash(authorizationCode)
			s.Require().NoError(err)
			test.Session.AuthorizationCode = &hash
			err = s.tx.Create(test.Session).Error
			s.Require().NoError(err, "error creating test session")

			s.metadataHelperMock.On("GetBasicAuth", s.Ctx).Return(test.Session.Client.ID.String(), test.ClientPassword, nil)

			res, err := s.srv.Token(s.Ctx, &pb.TokenRequest{
				AuthorizationCode: authorizationCode,
				GrantType:         "authorization_code",
			})
			s.Require().NoError(err, "unexpected error calling token endpoint")

			s.metadataHelperMock.AssertCalled(s.T(), "GetBasicAuth", s.Ctx)
			s.mockPseudonymizer.AssertNumberOfCalls(s.T(), "GetPseudonym", test.PseudoCalls)

			s.Equal("Bearer", res.GetTokenType())

			claims := &TokenClaims{}
			err = s.jwtClient.ValidateAndParseClaims(res.GetAccessToken(), claims)
			s.Require().NoError(err, "error parsing claims from token")

			s.Require().Len(claims.Audience, 1)
			s.Equal(test.ExpectedAudience, claims.Audience[0])
			s.Equal(test.ExpectedClientID, claims.ClientID)

			scopes := claims.Scopes

			s.Equal(test.ExpectedScopes, scopes)

			s.Equal(test.ExpectedSubjects, claims.Subjects)

			s.Equal(test.ExpectedSubject, claims.Subject)

			s.Equal(issuer, claims.Issuer)

			s.Equal(test.ExpectedClientMetadata, claims.ClientMetadata)
		})
	}
}

func (s *AuthServiceTestSuite) TestErrorOnEmptySubjectToken() {
	tests := []struct {
		Name           string
		Session        *models.Session
		ClientPassword string
	}{
		{
			Name:           "ErrorSubject",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{s.accessModelGql3},
				Audience:             s.audience,
				Client:               s.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       s.redirectTarget,
				Subject:              "",
			},
		},
		{
			Name:           "ErrorBsn",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{s.accessModelGql4},
				Audience:             s.audience,
				Client:               s.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       s.redirectTarget,
				Subject:              "",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			authorizationCode, err := authtoken.NewToken(s.conf.AuthorizationCodeLength)
			s.Require().NoError(err, "error creating test authorization code")
			hash, err := authtoken.Hash(authorizationCode)
			s.Require().NoError(err)
			test.Session.AuthorizationCode = &hash
			err = s.tx.Create(test.Session).Error
			s.Require().NoError(err, "error creating test session")

			s.metadataHelperMock.On("GetBasicAuth", s.Ctx).Return(test.Session.Client.ID.String(), test.ClientPassword, nil)

			_, err = s.srv.Token(s.Ctx, &pb.TokenRequest{
				AuthorizationCode: authorizationCode,
				GrantType:         "authorization_code",
			})
			s.Error(err)
		})
	}
}

func (s *AuthServiceTestSuite) TestErrorOnWrongSessionStateAndDeadline() {
	tests := []struct {
		name            string
		state           models.SessionState
		errorExpected   error
		codeGrantedTime time.Time
	}{
		{
			name:          "Unclaimed",
			state:         models.SessionStateUnclaimed,
			errorExpected: errgrpc.ErrFailedPrecondition(models.ErrUnableToRetrieveTokenInvalidState),
		},
		{
			name:          "Claimed",
			state:         models.SessionStateClaimed,
			errorExpected: errgrpc.ErrFailedPrecondition(models.ErrUnableToRetrieveTokenInvalidState),
		},
		{
			name:          "Accepted",
			state:         models.SessionStateAccepted,
			errorExpected: errgrpc.ErrFailedPrecondition(models.ErrUnableToRetrieveTokenInvalidState),
		},
		{
			name:          "Rejected",
			state:         models.SessionStateRejected,
			errorExpected: errgrpc.ErrFailedPrecondition(models.ErrUnableToRetrieveTokenInvalidState),
		},
		{
			name:          "TokenGranted",
			state:         models.SessionStateTokenGranted,
			errorExpected: errgrpc.ErrFailedPrecondition(models.ErrUnableToRetrieveTokenInvalidState),
		},
		{
			name:  "CodeGranted",
			state: models.SessionStateCodeGranted,
		},
		{
			name:            "CodeGranted",
			state:           models.SessionStateCodeGranted,
			errorExpected:   errgrpc.ErrDeadlineExceeded(models.ErrUnableToRetrieveTokenExpiration),
			codeGrantedTime: time.Now().Add(time.Minute - 20),
		},
		{
			name:            "CodeGranted",
			state:           models.SessionStateCodeGranted,
			errorExpected:   errgrpc.ErrDeadlineExceeded(models.ErrUnableToRetrieveTokenExpiration),
			codeGrantedTime: time.Now().Add(time.Minute - 10),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			authorizationCode, err := authtoken.NewToken(s.conf.AuthorizationCodeLength)
			s.Require().NoError(err, "error creating test authorization code")
			hash, err := authtoken.Hash(authorizationCode)
			s.Require().NoError(err)
			dummySession := &models.Session{
				AcceptedAccessModels: []*models.AccessModel{s.accessModelGql1},
				RequiredAccessModels: []*models.AccessModel{s.accessModelGql2},
				Audience:             s.audience,
				Client:               s.client,
				State:                test.state,
				RedirectTarget:       s.redirectTarget,
				Subject:              "pseudonym2",
				AuthorizationCode:    &hash,
			}
			err = s.tx.Create(dummySession).Error
			s.Require().NoError(err, "error creating test session")

			s.metadataHelperMock.On("GetBasicAuth", s.Ctx).Return(s.client.ID.String(), "test^123", nil)

			if !test.codeGrantedTime.IsZero() {
				log.Infof("%s %v", test.name, test.codeGrantedTime)
				err = s.tx.Model(models.Session{}).Where("id = ?", dummySession.ID).Update("authorization_code_granted_at", test.codeGrantedTime).Error
				s.Require().NoError(err, "error updating authorization_code_granted_at for test session")
			}

			_, err = s.srv.Token(s.Ctx, &pb.TokenRequest{
				AuthorizationCode: authorizationCode,
				GrantType:         "authorization_code",
			})

			if test.errorExpected != nil {
				s.Require().Error(err)
				s.EqualError(err, test.errorExpected.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *AuthServiceTestSuite) TestErrorOnNonExistingSession() {
	code, err := uuid.NewV4()
	s.metadataHelperMock.On("GetBasicAuth", s.Ctx).Return(s.client.ID.String(), "test^123", nil)
	s.Require().NoError(err)
	_, err = s.srv.Token(s.Ctx, &pb.TokenRequest{
		AuthorizationCode: code.String(),
		GrantType:         "authorization_code",
	})
	s.Require().Error(err)
	s.Contains(err.Error(), "session not found")
}

func (s *AuthServiceTestSuite) TestErrorNotFoundOnOtherClientsSession() {
	authorizationCode, err := authtoken.NewToken(s.conf.AuthorizationCodeLength)
	s.Require().NoError(err, "error creating test authorization code")
	hash, err := authtoken.Hash(authorizationCode)
	s.Require().NoError(err)
	dummySession := &models.Session{
		AcceptedAccessModels: []*models.AccessModel{s.accessModelGql1},
		RequiredAccessModels: []*models.AccessModel{s.accessModelGql2},
		Audience:             s.audience,
		Client:               s.client,
		State:                models.SessionStateCodeGranted,
		RedirectTarget:       s.redirectTarget,
		Subject:              "pseudonym2",
		AuthorizationCode:    &hash,
	}
	err = s.tx.Create(dummySession).Error
	s.Require().NoError(err, "error creating test session")

	s.metadataHelperMock.On("GetBasicAuth", s.Ctx).Return(s.client2.ID.String(), "456#%test", nil)
	s.Require().NoError(err)
	_, err = s.srv.Token(s.Ctx, &pb.TokenRequest{
		AuthorizationCode: authorizationCode,
		GrantType:         "authorization_code",
	})
	s.Require().Error(err)
	s.Contains(err.Error(), "session not found")
}

func (s *AuthServiceTestSuite) TestErrorOnWrongPassword() {
	authorizationCode, err := authtoken.NewToken(s.conf.AuthorizationCodeLength)
	s.Require().NoError(err, "error creating test authorization code")
	hash, err := authtoken.Hash(authorizationCode)
	s.Require().NoError(err)
	dummySession := &models.Session{
		AcceptedAccessModels: []*models.AccessModel{s.accessModelGql1},
		RequiredAccessModels: []*models.AccessModel{s.accessModelGql2},
		Audience:             s.audience,
		Client:               s.client,
		State:                models.SessionStateCodeGranted,
		RedirectTarget:       s.redirectTarget,
		Subject:              "pseudonym2",
		AuthorizationCode:    &hash,
	}
	err = s.tx.Create(dummySession).Error
	s.Require().NoError(err, "error creating test session")

	s.metadataHelperMock.On("GetBasicAuth", s.Ctx).Return(s.client.ID.String(), "456#%test", nil)
	s.Require().NoError(err)
	_, err = s.srv.Token(s.Ctx, &pb.TokenRequest{
		AuthorizationCode: authorizationCode,
		GrantType:         "authorization_code",
	})
	s.Require().Error(err)
	s.Contains(err.Error(), "incorrect password")
}

func (s *AuthServiceTestSuite) TestErrorOnNotAuthorizationCodeGrantType() {
	code, err := uuid.NewV4()
	s.Require().NoError(err)
	err = (&pb.TokenRequest{
		AuthorizationCode: code.String(),
		GrantType:         "implicit",
	}).Validate()
	s.Require().Error(err)
	fmt.Println("Finsihed case")
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, &AuthServiceTestSuite{
		AuthServiceBaseTestSuite: &AuthServiceBaseTestSuite{
			db: database.MustConnectTest("auth", nil),
		},
	})
}
