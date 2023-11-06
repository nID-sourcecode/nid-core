package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nID-sourcecode/nid-core/pkg/authtoken"
	"github.com/nID-sourcecode/nid-core/pkg/gqlutil"
	"github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/pseudonym"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/metrics"
	utilJWT "github.com/nID-sourcecode/nid-core/pkg/utilities/jwt/v3"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	contractMock "github.com/nID-sourcecode/nid-core/svc/auth/contract/mocks"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/audienceprovider"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/config"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/identityprovider"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/repository"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/stats"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	pb "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
	walletPB "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
	walletMock "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto/mock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/vrischmann/envconfig"
)

type TestAppSuite struct {
	repository.AuthDBSuite

	app *App

	jwtClient *utilJWT.Client

	stats *stats.Stats
	conf  *config.AuthConfig

	passwordManager     password.IManager
	callbackHandlerMock *contractMock.CallbackHandler
	mockPseudonymizer   *pseudonym.MockPseudonymizer
	schemaFetcher       *schemaFetcher

	TestModels
}

type TestModels struct {
	client              *models.Client
	client2             *models.Client
	redirectTarget      *models.RedirectTarget
	audience            *models.Audience
	audience2           *models.Audience
	scope               *models.Scope
	scope2              *models.Scope
	accessModelGql1     *models.AccessModel
	accessModelGql2     *models.AccessModel
	accessModelGql3     *models.AccessModel
	accessModelGql4     *models.AccessModel
	accessModelRestPOST *models.AccessModel
	accessModelRestGET  *models.AccessModel
}

func (a *TestAppSuite) SetupSuite() {
	a.AuthDBSuite.SetupSuite()

	a.conf = &config.AuthConfig{}
	if err := envconfig.InitWithOptions(a.conf, envconfig.Options{AllOptional: true}); err != nil {
		a.Failf("init conf failed", "%+v", err)
	}

	scope := metrics.NewPromScope(prometheus.DefaultRegisterer, "auth")
	a.stats = stats.CreateStats(scope)

	a.conf.AuthorizationCodeExpirationTime = time.Minute
	a.conf.AuthorizationCodeLength = 32
	a.conf.AuthRequestURI = "https://authrequest.com"

	priv, pub, err := utilJWT.GenerateTestKeys()
	a.Require().NoError(err, "error generating test keys")
	a.conf.Namespace = "nid"
	a.conf.Issuer = issuer
	opts := utilJWT.DefaultOpts()
	a.jwtClient = utilJWT.NewJWTClientWithOpts(priv, pub, opts)

	a.passwordManager = password.NewDefaultManager()
	a.callbackHandlerMock = &contractMock.CallbackHandler{}

	a.setupRegisterTest()
}

func (a *TestAppSuite) SetupTest() {
	a.AuthDBSuite.SetupTest()

	a.mockPseudonymizer = setupMockPseudonymizer(context.Background())
	audienceProvider := audienceprovider.NewDatabaseAudienceProvider(a.conf)
	identityProvider := identityprovider.NewDatabaseIdentityProvider(a.AuthDBSuite.Repo.ClientDB, a.passwordManager)

	a.app = a.createApp(audienceProvider, identityProvider)
	a.seedModels(a.Require(), a.AuthDBSuite.Tx, a.passwordManager)
}

func (a *TestAppSuite) createApp(audienceProvider contract.AudienceProvider, identityProvider contract.IdentityProvider) *App {
	mockWalletClient := &walletMock.WalletClient{}
	mockWalletClient.On("CreateConsent", context.Background(), mock.Anything).Return(&walletPB.ConsentResponse{}, nil)

	return New(a.conf, a.AuthDBSuite.Repo, a.schemaFetcher, a.stats, a.callbackHandlerMock, a.passwordManager, a.jwtClient, a.mockPseudonymizer, mockWalletClient, audienceProvider, identityProvider)
}

func (s *TestModels) seedModels(suite *require.Assertions, tx *gorm.DB, passwordManager password.IManager) {
	client1Password, err := passwordManager.GenerateHash("test^123")
	suite.NoError(err)
	client2Password, err := passwordManager.GenerateHash("456#%test")
	suite.NoError(err)

	s.client = &models.Client{
		Color:    "blue",
		Name:     "testclient",
		Password: client1Password,
		Metadata: postgres.Jsonb{RawMessage: json.RawMessage(`{"oin":"000012345"}`)},
	}
	err = tx.Create(s.client).Error
	suite.NoError(err)

	s.client2 = &models.Client{
		Color:    "red",
		Name:     "testclient2",
		Password: client2Password,
	}
	err = tx.Create(s.client2).Error
	suite.NoError(err)

	s.redirectTarget = &models.RedirectTarget{
		ClientID:       s.client.ID,
		RedirectTarget: "https://weave.nl/code",
	}

	err = tx.Create(s.redirectTarget).Error
	suite.NoError(err)

	err = tx.Create(&models.RedirectTarget{
		ClientID:       s.client2.ID,
		RedirectTarget: "https://weave2.nl/code",
		UpdatedAt:      time.Time{},
	}).Error
	suite.NoError(err)

	s.audience = &models.Audience{
		Audience:  "https://test.com/gql",
		Namespace: "alice",
	}
	err = tx.Create(s.audience).Error
	suite.NoError(err)

	s.audience2 = &models.Audience{
		Audience:  "https://test2.com/gql",
		Namespace: "bob",
	}
	err = tx.Create(s.audience2).Error
	suite.NoError(err)

	s.accessModelGql1 = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "abc",
		Name:       "test:stuff",
		Type:       models.AccessModelTypeGQL,
		GqlAccessModel: &models.GqlAccessModel{
			Path:      "/gql",
			JSONModel: `{"r":"somestuff"}`,
		},
	}
	err = tx.Create(s.accessModelGql1).Error
	suite.NoError(err)

	s.accessModelGql2 = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "ghi",
		Name:       "test:stuff2",
		Type:       models.AccessModelTypeGQL,
		GqlAccessModel: &models.GqlAccessModel{
			Path:      "/graphql",
			JSONModel: `{"r":"somemorestuff"}`,
		},
	}
	err = tx.Create(s.accessModelGql2).Error
	suite.NoError(err)

	s.accessModelGql3 = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "abcsubject",
		Name:       "test:stuff3",
		Type:       models.AccessModelTypeGQL,
		GqlAccessModel: &models.GqlAccessModel{
			Path:      "/gql",
			JSONModel: `{"subject":"$$nid:subject$$"}`,
		},
	}

	err = tx.Create(s.accessModelGql3).Error
	s.accessModelGql4 = &models.AccessModel{
		AudienceID: s.audience.ID,
		Hash:       "abcbsn",
		Name:       "test:stuff4",
		Type:       models.AccessModelTypeGQL,
		GqlAccessModel: &models.GqlAccessModel{
			Path:      "/gql",
			JSONModel: `{"bsn":"$$nid:bsn$$"}`,
		},
	}
	suite.NoError(err)

	err = tx.Create(s.accessModelGql4).Error
	suite.NoError(err)

	s.accessModelRestPOST = &models.AccessModel{
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
	err = tx.Create(s.accessModelRestPOST).Error
	suite.NoError(err)

	s.accessModelRestGET = &models.AccessModel{
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

	err = tx.Create(s.accessModelRestGET).Error
	suite.NoError(err)

	s.scope = &models.Scope{Resource: "test", Scope: "test:stuff", Audiences: []*models.Audience{s.audience}}
	err = tx.Create(s.scope).Error
	suite.NoError(err)

	s.scope2 = &models.Scope{Resource: "test2", Scope: "test2:stuff", Audiences: []*models.Audience{s.audience, s.audience2}}
	err = tx.Create(s.scope2).Error
	suite.NoError(err)
}

func setupMockPseudonymizer(ctx context.Context) *pseudonym.MockPseudonymizer {
	mockPseudonymizer := &pseudonym.MockPseudonymizer{}
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym1", "alice").Return("translatedPseudo1ForAlice", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym1", "nid").Return("translatedPseudo1ForNid", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym1", "bob").Return("translatedPseudo1ForBob", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym2", "alice").Return("translatedPseudo2ForAlice", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym2", "nid").Return("translatedPseudo2ForNid", nil)
	mockPseudonymizer.On("GetPseudonym", ctx, "pseudonym2", "bob").Return("translatedPseudo2ForBob", nil)
	return mockPseudonymizer
}

func (a *TestAppSuite) TearDownTest() {
	a.AuthDBSuite.TearDownTest()
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(TestAppSuite))
}

const (
	dummySubject = "sadasdasjkdhaiouysdg867ig672315471r23t7"
	issuer       = "issuer.nl"
)

func (a *TestAppSuite) createAcceptDummySession(claimed bool) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           a.audience.ID,
		ClientID:             a.client.ID,
		RedirectTargetID:     a.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		Subject:              dummySubject,
	}
	if claimed {
		session.State = models.SessionStateClaimed
	} else {
		session.State = models.SessionStateUnclaimed
	}
	err := a.AuthDBSuite.Tx.Create(&session).Error
	a.Require().NoError(err, "Did not expect an error on create session")
	err = a.AuthDBSuite.Tx.Model(&session).Association("RequiredAccessModels").Append(a.accessModelGql1, a.accessModelGql2).Error
	a.Require().NoError(err, "Did not expect an error on create association")
	err = a.AuthDBSuite.Tx.Model(&session).Association("OptionalAccessModels").Append(a.accessModelRestPOST).Error
	a.Require().NoError(err, "Did not expect an error on create association")

	return session
}

func (a *TestAppSuite) TestAcceptedClaimsUpdated() {
	testSession := a.createAcceptDummySession(true)
	token := a.newToken(dummySubject)
	claim := strings.Split(token, ".")[1]

	resp, err := a.app.Accept(context.Background(), claim, &models.AcceptRequest{
		SessionID:      testSession.ID.String(),
		AccessModelIds: []string{a.accessModelRestPOST.ID.String()},
	})
	a.Require().NoError(err)

	var session models.Session
	err = a.AuthDBSuite.Tx.Where("id = ?", testSession.ID).Preload("AcceptedAccessModels").Find(&session).Error
	a.NoError(err)
	a.Len(session.AcceptedAccessModels, 3, "sessions should have 3 accepted access models (1) optional and (2) required")

	a.Equal(resp.ID, testSession.ID.String())
	a.Equal(resp.State, models.SessionStateAccepted)
	a.Require().NotNil(resp.Audience)
	a.Equal(resp.Audience.ID, a.audience.ID)
	a.Equal(resp.Audience.Audience, a.audience.Audience)
	a.Equal(resp.Audience.Namespace, a.audience.Namespace)

	if a.NotNil(resp.Client) {
		a.Equal(resp.Client.ID, a.client.ID)
		a.Equal(resp.Client.Name, a.client.Name)
		a.Equal(resp.Client.Icon, a.client.Icon)
		a.Equal(resp.Client.Logo, a.client.Logo)
		a.Equal(resp.Client.Color, a.client.Color)
	}
	if a.NotNil(resp.RequiredAccessModels) && a.NotNil(resp.RequiredAccessModels[0]) {
		a.Equal(resp.RequiredAccessModels[0].ID, a.accessModelGql1.ID)
		a.Equal(resp.RequiredAccessModels[0].Name, a.accessModelGql1.Name)
		a.Equal(resp.RequiredAccessModels[0].Hash, a.accessModelGql1.Hash)
		a.Equal(resp.RequiredAccessModels[0].Description, a.accessModelGql1.Description)
	}
	if a.NotNil(resp.OptionalAccessModels) && a.NotNil(resp.OptionalAccessModels[0]) {
		a.Equal(resp.OptionalAccessModels[0].ID, a.accessModelRestPOST.ID)
		a.Equal(resp.OptionalAccessModels[0].Name, a.accessModelRestPOST.Name)
		a.Equal(resp.OptionalAccessModels[0].Hash, a.accessModelRestPOST.Hash)
		a.Equal(resp.OptionalAccessModels[0].Description, a.accessModelRestPOST.Description)
	}
	if a.NotNil(resp.AcceptedAccessModels) && a.NotNil(resp.AcceptedAccessModels[0]) {
		a.Equal(resp.AcceptedAccessModels[0].ID, a.accessModelRestPOST.ID)
		a.Equal(resp.AcceptedAccessModels[0].Name, a.accessModelRestPOST.Name)
		a.Equal(resp.AcceptedAccessModels[0].Hash, a.accessModelRestPOST.Hash)
		a.Equal(resp.AcceptedAccessModels[0].Description, a.accessModelRestPOST.Description)
	}
}

func (a *TestAppSuite) TestAcceptAccessModels() {
	testSession := a.createAcceptDummySession(true)
	token := a.newToken(dummySubject)
	claim := strings.Split(token, ".")[1]
	resp, err := a.app.Accept(context.Background(), claim, &models.AcceptRequest{
		SessionID:      testSession.ID.String(),
		AccessModelIds: []string{a.accessModelRestPOST.ID.String()},
	})

	a.Require().NoError(err)
	a.NotNil(resp)
	a.NotNil(resp.AcceptedAccessModels)
	a.Len(resp.AcceptedAccessModels, 3)
	a.NotEqualValues(resp.AcceptedAccessModels[0].ID, resp.AcceptedAccessModels[1].ID, resp.AcceptedAccessModels[2].ID)
	for _, model := range resp.AcceptedAccessModels {
		if a.True(model.ID == a.accessModelGql1.ID ||
			model.ID == a.accessModelGql2.ID ||
			model.ID == a.accessModelRestPOST.ID) {
			continue
		}
		a.Failf("access_model id %s not found in response", model.ID.String())
	}
}

func (a *TestAppSuite) TestNoOptionalAccessModelsNotProvided() {
	testSession := a.createAcceptDummySession(true)

	token := a.newToken(dummySubject)
	claim := strings.Split(token, ".")[1]

	resp, err := a.app.Accept(context.Background(), claim, &models.AcceptRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().NoError(err)
	a.NotNil(resp)
	a.NotNil(resp.AcceptedAccessModels)
	a.Len(resp.AcceptedAccessModels, 2)
	a.NotEqualValues(resp.AcceptedAccessModels[0].ID, resp.AcceptedAccessModels[1].ID)
	for _, model := range resp.AcceptedAccessModels {
		if a.True(model.ID == a.accessModelGql1.ID || model.ID == a.accessModelGql2.ID) {
			continue
		}
		a.Failf("access_model id %s not found in response", model.ID.String())
	}
}

func (a *TestAppSuite) TestErrorNonExistingAccessModel() {
	testSession := a.createAcceptDummySession(true)

	_, err := a.app.Accept(context.Background(), a.newToken(dummySubject), &models.AcceptRequest{
		SessionID:      testSession.ID.String(),
		AccessModelIds: []string{uuid.Must(uuid.NewV4()).String()},
	})
	a.Require().Error(err)
}

func (a *TestAppSuite) TestErrorRequiredAccessModelIdInPayload() {
	testSession := a.createAcceptDummySession(false)
	token := a.newToken(dummySubject)
	claim := strings.Split(token, ".")[1]
	_, err := a.app.Accept(context.Background(), claim, &models.AcceptRequest{
		SessionID:      testSession.ID.String(),
		AccessModelIds: []string{a.accessModelGql1.ID.String()},
	})
	a.Error(err)
	a.ErrorIs(err, contract.ErrInvalidArguments)
}

func (a *TestAppSuite) newToken(subject string) string {
	mySigningKey := []byte("AllYourBase")

	type ClaimWithSubject struct {
		Sub string `json:"sub"`
		jwt.StandardClaims
	}

	// Create the Claims
	claims := ClaimWithSubject{
		subject,
		jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	a.Require().NoError(err)

	return ss
}

func (a *TestAppSuite) TestAuthHeadless() {
	a.callbackHandlerMock.On("HandleCallback", mock.Anything, "https://weave.nl/code", mock.Anything).
		Return(nil).Once()

	err := a.app.AuthorizeHeadless(context.Background(), &models.AuthorizeHeadlessRequest{
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
		QueryModelJSON: `{
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
	a.Nil(err)
}

func (a *TestAppSuite) TestErrorOnHandleCallbackHeadless() {
	a.callbackHandlerMock.On("HandleCallback", mock.Anything, "https://weave.nl/code", mock.Anything).
		Return(errors.New("error")).Once()

	err := a.app.AuthorizeHeadless(context.Background(), &models.AuthorizeHeadlessRequest{
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
		QueryModelJSON: `{
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
	a.Equal(err, contract.ErrInternalError)
}

func (a *TestAppSuite) TestAuthorize_Session() {
	redirectURL, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	a.Require().NoError(err)

	sessionID := redirectURL[len("https://authrequest.com#"):]
	session, err := a.AuthDBSuite.Repo.SessionDB.GetSessionByID(repository.PreloadAll, sessionID)
	a.Require().NoError(err)
	a.Require().NotZero(session.ID)

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	a.Require().False(isExpired)

	a.Equal(fmt.Sprintf("https://authrequest.com#%s", session.ID.String()), redirectURL)
	a.Equal(a.redirectTarget.ID, session.RedirectTargetID)
	a.Equal(models.SessionStateUnclaimed, session.State)
	a.Equal(a.audience.ID, session.AudienceID)
	a.Equal("", session.Subject)
	a.Len(session.AcceptedAccessModels, 0)
	a.Len(session.OptionalAccessModels, 0)
	a.Require().Len(session.RequiredAccessModels, 1)
	a.Equal(a.accessModelGql1.ID, session.RequiredAccessModels[0].ID)
}

func (a *TestAppSuite) TestAuthorize_SessionCreatedWithOptionalScopes() {
	redirectURL, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:          "openid test:stuff@abc",
		ResponseType:   "code",
		ClientID:       a.client.ID.String(),
		RedirectURI:    "https://weave.nl/code",
		Audience:       "https://test.com/gql",
		OptionalScopes: "test:stuff2@ghi test:stuff3@jkl",
	})
	a.Require().NoError(err)

	sessionID := redirectURL[len("https://authrequest.com#"):]
	session, err := a.AuthDBSuite.Repo.SessionDB.GetSessionByID(repository.PreloadAll, sessionID)
	a.Require().NoError(err)
	a.Require().NotZero(session.ID)

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	a.Require().False(isExpired)

	a.Equal(redirectURL, fmt.Sprintf("https://authrequest.com#%s", session.ID.String()))
	a.Equal(a.redirectTarget.ID, session.RedirectTargetID)
	a.Equal(models.SessionStateUnclaimed, session.State)
	a.Equal(a.audience.ID, session.AudienceID)
	a.Equal("", session.Subject)
	a.Len(session.OptionalAccessModels, 2)
	foundIDs := make([]uuid.UUID, 0)
	for _, optAccessModel := range session.OptionalAccessModels {
		foundIDs = append(foundIDs, optAccessModel.ID)
	}
	a.Contains(foundIDs, a.accessModelGql2.ID)
	a.Contains(foundIDs, a.accessModelRestPOST.ID)

	a.Len(session.AcceptedAccessModels, 0)
	a.Require().Len(session.RequiredAccessModels, 1)
	a.Equal(a.accessModelGql1.ID, session.RequiredAccessModels[0].ID)
}

func (a *TestAppSuite) TestAuthorize_ErrorOnNoScopes() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), "at least one access model-scope must be specified")
}

func (a *TestAppSuite) TestErrorOnNonExistentRedirectURI() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://wee.nl/code",
		Audience:     "https://test.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), fmt.Sprintf("client with id \"%s\" does not have redirect URI \"https://wee.nl/code\"", a.client.ID.String()))
}

func (a *TestAppSuite) TestAuthorize_ErrorOnWrongRedirectURI() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave2.nl/code",
		Audience:     "https://test.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), fmt.Sprintf("client with id \"%s\" does not have redirect URI \"https://weave2.nl/code\"", a.client.ID.String()))
}

func (a *TestAppSuite) TestAuthorize_ErrorOnNonExistentAudience() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://te.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), "audience \"https://te.com/gql\" does not exist")
}

func (a *TestAppSuite) TestAuthorize_ErrorOnWrongAudience() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test2.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), "audience \"https://test2.com/gql\" does not support access model \"test:stuff@abc\"")
}

func (a *TestAppSuite) TestAuthorize_ErrorOnResponseTypeOtherThanCode() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test:stuff@abc",
		ResponseType: "token",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), "no response type other than \"code\" is supported")
}

func (a *TestAppSuite) TestAuthorize_ErrorOnNonExistentAccessModel() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test:stf@ac",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), "audience \"https://test.com/gql\" does not support access model \"test:stf@ac\"")
}

func (a *TestAppSuite) TestAuthorize_ErrorOnWrongAccessModelHash() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test:stuff@ac",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), "audience \"https://test.com/gql\" does not support access model \"test:stuff@ac\"")
}

func (a *TestAppSuite) TestAuthorize_ErrorOnWrongAccessModel() {
	_, err := a.app.Authorize(context.Background(), &models.AuthorizeRequest{
		Scope:        "openid test2:stuff2@def",
		ResponseType: "code",
		ClientID:     a.client.ID.String(),
		RedirectURI:  "https://weave.nl/code",
		Audience:     "https://test.com/gql",
	})
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Contains(err.Error(), "audience \"https://test.com/gql\" does not support access model \"test2:stuff2@def\"")
}

func (a *TestAppSuite) createClaimDummySession(claimed bool) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           a.audience.ID,
		ClientID:             a.client.ID,
		RedirectTargetID:     a.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		Subject:              dummySubject,
	}
	if claimed {
		session.State = models.SessionStateClaimed
	} else {
		session.State = models.SessionStateUnclaimed
	}
	err := a.AuthDBSuite.Tx.Create(&session).Error
	a.Require().NoError(err, "Did not expect an error on create session")
	err = a.AuthDBSuite.Tx.Model(&session).Association("RequiredAccessModels").Append(a.accessModelGql1, a.accessModelGql2).Error
	a.Require().NoError(err, "Did not expect an error on create association")
	err = a.AuthDBSuite.Tx.Model(&session).Association("OptionalAccessModels").Append(a.accessModelRestPOST).Error
	a.Require().NoError(err, "Did not expect an error on create association")

	return session
}

func (a *TestAppSuite) TestClaimSessionUpdated() {
	testSession := a.createClaimDummySession(false)
	token := a.newToken(dummySubject)
	claim := strings.Split(token, ".")[1]
	resp, err := a.app.Claim(context.Background(), claim, &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().NoError(err)

	var session models.Session
	a.AuthDBSuite.Tx.First(&session, "id = ?", testSession.ID)

	a.NotZero(session.ID)
	a.Equal(session.ID, testSession.ID)
	a.Equal(session.State, models.SessionStateClaimed)
	a.Equal(resp.State, models.SessionStateClaimed)
}

func (a *TestAppSuite) TestClaimSessionResponse() {
	testSession := a.createClaimDummySession(false)
	token := a.newToken(dummySubject)
	claim := strings.Split(token, ".")[1]
	resp, err := a.app.Claim(context.Background(), claim, &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().NoError(err)

	a.Equal(resp.ID, testSession.ID.String())
	a.Equal(resp.State, models.SessionStateClaimed)
	if a.NotNil(resp.Audience) {
		a.Equal(resp.Audience.ID, a.audience.ID)
		a.Equal(resp.Audience.Audience, a.audience.Audience)
		a.Equal(resp.Audience.Namespace, a.audience.Namespace)
	}
	if a.NotNil(resp.Client) {
		a.Equal(resp.Client.ID, a.client.ID)
		a.Equal(resp.Client.Name, a.client.Name)
		a.Equal(resp.Client.Icon, a.client.Icon)
		a.Equal(resp.Client.Logo, a.client.Logo)
		a.Equal(resp.Client.Color, a.client.Color)
	}
	if a.NotNil(resp.RequiredAccessModels) && a.NotNil(resp.RequiredAccessModels[0]) {
		a.Equal(resp.RequiredAccessModels[0].ID, a.accessModelGql1.ID)
		a.Equal(resp.RequiredAccessModels[0].Name, a.accessModelGql1.Name)
		a.Equal(resp.RequiredAccessModels[0].Hash, a.accessModelGql1.Hash)
		a.Equal(resp.RequiredAccessModels[0].Description, a.accessModelGql1.Description)
	}
	if a.NotNil(resp.OptionalAccessModels) && a.NotNil(resp.OptionalAccessModels[0]) {
		a.Equal(resp.OptionalAccessModels[0].ID, a.accessModelRestPOST.ID)
		a.Equal(resp.OptionalAccessModels[0].Name, a.accessModelRestPOST.Name)
		a.Equal(resp.OptionalAccessModels[0].Hash, a.accessModelRestPOST.Hash)
		a.Equal(resp.OptionalAccessModels[0].Description, a.accessModelRestPOST.Description)
	}
	a.Nil(resp.AcceptedAccessModels)
}

func (a *TestAppSuite) TestErrorClaimedStateSession() {
	testSession := a.createClaimDummySession(true)
	_, err := a.app.Claim(context.Background(), a.newToken(dummySubject), &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().Error(err)
}

func (a *TestAppSuite) TestErrorEmptySubject() {
	testSession := a.createClaimDummySession(true)
	_, err := a.app.Claim(context.Background(), a.newToken(dummySubject), &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().Error(err)
}

// Get_Session_Details
func (a *TestAppSuite) createGetSessionDetailsDummySession(state models.SessionState) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           a.audience.ID,
		ClientID:             a.client.ID,
		RedirectTargetID:     a.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		Subject:              "sadasdasjkdhaiouysdg867ig672315471r23t7",
		State:                state,
	}

	err := a.AuthDBSuite.Tx.Create(&session).Error
	a.Require().NoError(err, "Did not expect an error on create session")
	err = a.AuthDBSuite.Tx.Model(&session).Association("RequiredAccessModels").Append(a.accessModelGql1, a.accessModelGql2).Error
	a.Require().NoError(err, "Did not expect an error on create association")
	err = a.AuthDBSuite.Tx.Model(&session).Association("OptionalAccessModels").Append(a.accessModelRestPOST).Error
	a.Require().NoError(err, "Did not expect an error on create association")

	return session
}

func (a *TestAppSuite) TestDetailsSessionResponse() {
	testSession := a.createGetSessionDetailsDummySession(models.SessionStateUnclaimed)
	resp, err := a.app.GetSessionDetails(context.Background(), &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().NoError(err)
	a.Require().NotNil(resp)
}

func (a *TestAppSuite) TestWrongState() {
	testSession := a.createGetSessionDetailsDummySession(models.SessionStateAccepted)
	_, err := a.app.GetSessionDetails(context.Background(), &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().Error(err)
}

func (a *TestAppSuite) createFinaliseDummySession(accepted bool) (*models.Session, string) {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           a.audience.ID,
		ClientID:             a.client.ID,
		RedirectTargetID:     a.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
	}
	if accepted {
		session.State = models.SessionStateAccepted
	} else {
		session.State = models.SessionStateRejected
	}
	err := a.AuthDBSuite.Tx.Create(&session).Error
	a.Require().NoError(err, "Did not expect an error")

	token, err := a.setupSessionFinaliseToken(session)
	a.Require().NoError(err, "Did not expect an error")

	return session, token
}

func (a *TestAppSuite) TestAcceptedSessionStateUpdated() {
	testSession, token := a.createFinaliseDummySession(true)

	response, err := a.app.Finalise(context.Background(), &models.FinaliseRequest{SessionID: testSession.ID.String(), SessionFinaliseToken: token})
	a.Require().NoError(err)

	a.NotNil(response)
	a.NotEmpty(response.RedirectLocation)
	u, err := url.Parse(response.RedirectLocation)
	a.Require().NoError(err)
	authCode := u.Query().Get("authorization_code")
	a.NotEmpty(authCode)

	var session models.Session
	a.AuthDBSuite.Tx.First(&session, "id = ?", testSession.ID)
	a.NotZero(session.ID)
	a.NotNil(session.AuthorizationCode)
	a.Equal(session.ID, testSession.ID)

	hash, err := authtoken.Hash(authCode)
	a.Require().NoError(err)
	a.Equal(*session.AuthorizationCode, hash)

	a.Equal(a.redirectTarget.RedirectTarget+"/?authorization_code="+authCode, response.RedirectLocation)
	a.Equal(session.State, models.SessionStateCodeGranted)
}

func (a *TestAppSuite) TestNonAcceptedSessionError() {
	testSession, token := a.createFinaliseDummySession(false)

	_, err := a.app.Finalise(context.Background(), &models.FinaliseRequest{SessionID: testSession.ID.String(), SessionFinaliseToken: token})
	a.Error(err)
}

func (a *TestAppSuite) TestCreateSessionFinaliseToken() {
	session := a.createAcceptDummySession(false)
	_, err := a.setupSessionFinaliseToken(session)
	a.NoError(err, "Could not setup password for session")
}

func (a *TestAppSuite) TestCorrectSessionFinaliseToken() {
	testSession := a.createAcceptDummySession(false)
	token, err := a.setupSessionFinaliseToken(testSession)
	a.Require().NoError(err)

	session, err := a.AuthDBSuite.Repo.SessionDB.GetSessionByID(repository.NoPreload, testSession.ID.String())
	a.Require().NoError(err)

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	a.Require().False(isExpired)

	ok, err := a.passwordManager.ComparePassword(token, session.FinaliseToken)
	a.True(ok)
	a.NoError(err)
}

// createPasswordDummySession adds token as a password to session
// returns the password
func (a *TestAppSuite) setupSessionFinaliseToken(session *models.Session) (string, error) {
	token, err := authtoken.NewToken(a.conf.AuthorizationCodeLength)
	if err != nil {
		return "", err
	}

	hash, err := a.passwordManager.GenerateHash(token)
	if err != nil {
		return "", err
	}

	err = a.AuthDBSuite.Repo.SessionDB.SetSessionFinaliseToken(session, hash)
	if err != nil {
		return "", err
	}

	return token, nil
}

type schemaFetcher struct {
	mock.Mock

	testSchema *gqlutil.Schema
}

func (f *schemaFetcher) FetchSchema(ctx context.Context, url string) (*gqlutil.Schema, error) {
	args := f.Called(ctx, url)

	return f.testSchema, args.Error(1)
}

func (a *TestAppSuite) setupRegisterTest() {
	testSchema := &gqlutil.Schema{
		Types: map[string]*gqlutil.Type{
			"AddressFilterInput":        {Fields: map[string]*gqlutil.Field{}},
			"Geography":                 {Fields: map[string]*gqlutil.Field{}},
			"SavingsAccountFilterInput": {Fields: map[string]*gqlutil.Field{}},
			"UUIDFilterInput":           {Fields: map[string]*gqlutil.Field{}},
			"GeographyFilterInput":      {Fields: map[string]*gqlutil.Field{}},
			"FloatFilterInput":          {Fields: map[string]*gqlutil.Field{}},
			"UUID":                      {Fields: map[string]*gqlutil.Field{}},
			"ContactDetailFilterInput":  {Fields: map[string]*gqlutil.Field{}},
			"Float":                     {Fields: map[string]*gqlutil.Field{}},
			"Time":                      {Fields: map[string]*gqlutil.Field{}},
			"BankAccount": {
				Fields: map[string]*gqlutil.Field{
					"id":              {IsModel: false, TypeName: "UUID"},
					"savingsAccounts": {IsModel: true, TypeName: "SavingsAccount"},
					"userId":          {IsModel: false, TypeName: "UUID"},
					"accountNumber":   {IsModel: false, TypeName: "String"},
					"amount":          {IsModel: false, TypeName: "Int"},
					"createdAt":       {IsModel: false, TypeName: "Time"},
					"deletedAt":       {IsModel: false, TypeName: "Time"},
					"updatedAt":       {IsModel: false, TypeName: "Time"},
					"user":            {IsModel: true, TypeName: "User"},
				},
			},
			"Decimal": {Fields: map[string]*gqlutil.Field{}},
			"Map":     {Fields: map[string]*gqlutil.Field{}},
			"Query": {
				Fields: map[string]*gqlutil.Field{
					"address":         {IsModel: true, TypeName: "Address"},
					"bankAccount":     {IsModel: true, TypeName: "BankAccount"},
					"contactDetails":  {IsModel: true, TypeName: "ContactDetail"},
					"user":            {IsModel: true, TypeName: "User"},
					"users":           {IsModel: true, TypeName: "User"},
					"addresses":       {IsModel: true, TypeName: "Address"},
					"bankAccounts":    {IsModel: true, TypeName: "BankAccount"},
					"contactDetail":   {IsModel: true, TypeName: "ContactDetail"},
					"savingsAccount":  {IsModel: true, TypeName: "SavingsAccount"},
					"savingsAccounts": {IsModel: true, TypeName: "SavingsAccount"},
				},
			},
			"JSONFilterInput":        {Fields: map[string]*gqlutil.Field{}},
			"ID":                     {Fields: map[string]*gqlutil.Field{}},
			"ContactDetail":          {Fields: map[string]*gqlutil.Field{}},
			"BankAccountFilterInput": {Fields: map[string]*gqlutil.Field{}},
			"UserFilterInput":        {Fields: map[string]*gqlutil.Field{}},
			"Boolean":                {Fields: map[string]*gqlutil.Field{}},
			"StringFilterInput":      {Fields: map[string]*gqlutil.Field{}},
			"TimeFilterInput":        {Fields: map[string]*gqlutil.Field{}},
			"JSON":                   {Fields: map[string]*gqlutil.Field{}},
			"String":                 {Fields: map[string]*gqlutil.Field{}},
			"BooleanFilterInput":     {Fields: map[string]*gqlutil.Field{}},
			"DecimalFilterInput":     {Fields: map[string]*gqlutil.Field{}},
			"Int":                    {Fields: map[string]*gqlutil.Field{}},
			"Address":                {Fields: map[string]*gqlutil.Field{}},
			"SavingsAccount":         {Fields: map[string]*gqlutil.Field{}},
			"User": {
				Fields: map[string]*gqlutil.Field{
					"updatedAt":      {IsModel: false, TypeName: "Time"},
					"contactDetails": {IsModel: false, TypeName: "ContactDetail"},
					"createdAt":      {IsModel: false, TypeName: "Time"},
					"deletedAt":      {IsModel: false, TypeName: "Time"},
					"firstName":      {IsModel: false, TypeName: "String"},
					"pseudonym":      {IsModel: false, TypeName: "String"},
					"id":             {IsModel: false, TypeName: "UUID"},
					"bankAccounts": {
						IsModel:  true,
						TypeName: "BankAccount",
					},
					"lastName": {IsModel: false, TypeName: "String"},
				},
			},
			"IntFilterInput": {Fields: map[string]*gqlutil.Field{}},
		},
	}

	fetcher := &schemaFetcher{}
	fetcher.testSchema = testSchema
	fetcher.On("FetchSchema", mock.Anything, mock.Anything).Return(testSchema, nil)
	a.schemaFetcher = fetcher
}

func (a *TestAppSuite) TestRegisterAccessModel() {
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
		err := a.app.RegisterAccessModel(context.Background(), &models.AccessModelRequest{
			Audience:       test.inputAudience,
			QueryModelJSON: test.inputQueryModelJSON,
			ScopeName:      test.inputScopeName,
			Description:    test.inputDescription,
		})
		if test.err {
			a.Require().Error(err)
		} else {
			a.Require().NoError(err)
		}
	}
}

func (a *TestAppSuite) createRejectDummySession(rejected bool) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           a.audience.ID,
		ClientID:             a.client.ID,
		RedirectTargetID:     a.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		Subject:              dummySubject,
	}
	if rejected {
		session.State = models.SessionStateRejected
	} else {
		session.State = models.SessionStateClaimed
	}
	err := a.AuthDBSuite.Tx.Create(&session).Error
	a.Require().NoError(err, "Did not expect an error")

	return session
}

func (a *TestAppSuite) TestRejectSession() {
	testSession := a.createRejectDummySession(false)
	err := a.app.Reject(context.Background(), &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().NoError(err)

	var session models.Session
	a.AuthDBSuite.Tx.First(&session, "id = ?", testSession.ID)

	a.NotZero(session.ID)
	a.Equal(session.ID, testSession.ID)
	a.Equal(session.State, models.SessionStateRejected)
}

func (a *TestAppSuite) TestErrorIncorrectSessionState() {
	testSession := a.createRejectDummySession(true)
	err := a.app.Reject(context.Background(), &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Error(err)
}

func (a *TestAppSuite) createStatusDummySession(state models.SessionState) *models.Session {
	session := &models.Session{
		AcceptedAccessModels: nil,
		AudienceID:           a.audience.ID,
		ClientID:             a.client.ID,
		RedirectTargetID:     a.redirectTarget.ID,
		OptionalAccessModels: nil,
		RequiredAccessModels: nil,
		State:                state,
	}
	err := a.AuthDBSuite.Tx.Create(&session).Error
	a.Require().NoError(err, "Did not expect an error")

	return session
}

func (a *TestAppSuite) TestGetStatusNewlyCreatedSession() {
	testSession := a.createStatusDummySession(models.SessionStateRejected)
	resp, err := a.app.Status(context.Background(), &models.SessionRequest{
		SessionID: testSession.ID.String(),
	})
	a.Require().NoError(err)

	var session models.Session
	a.AuthDBSuite.Tx.First(&session, "id = ?", testSession.ID)
	a.NotZero(session.ID)
	a.Equal(session.ID, testSession.ID)
	a.Equal(session.State, models.SessionStateRejected)
	a.Equal(resp.State, models.SessionStateRejected)
}

func (a *TestAppSuite) TestSwapTokenValidate() {
	missingToken := pb.SwapTokenRequest{
		Query:    "query",
		Audience: "test",
	}
	a.Error(missingToken.Validate())
	requestMissingQuery := pb.SwapTokenRequest{
		CurrentToken: "Token",
		Audience:     "someone",
	}
	a.Error(requestMissingQuery.Validate())
	missingAudience := pb.SwapTokenRequest{
		CurrentToken: "Token",
		Query:        "{}",
	}
	a.Error(missingAudience.Validate())
}

func (a *TestAppSuite) TestInvalidToken() {
	invalidToken := models.SwapTokenRequest{
		CurrentToken: "someinvalidtoken",
		Query:        "{}",
	}
	resp, err := a.app.SwapToken(context.Background(), &invalidToken)
	a.ErrorIs(err, contract.ErrInvalidArguments)
	a.Nil(resp)
}

func (a *TestAppSuite) TestValidToken() {
	clientID := uuid.Must(uuid.NewV4())
	subID := uuid.Must(uuid.NewV4())

	defaultClaims := utilJWT.NewDefaultClaims(time.Duration(a.conf.JWTExpirationHours))
	defaultClaims.Subject = subID.String()
	defaultClaims.Audience = []string{"informationservice-A"}
	defaultClaims.Issuer = issuer
	claims := &models.TokenClaims{
		DefaultClaims: defaultClaims,
		ClientID:      clientID.String(),
	}

	token, err := a.app.jwtClient.SignToken(claims)
	a.Require().NoError(err)

	req := &models.SwapTokenRequest{
		CurrentToken: token,
		Query:        "{age}",
		Audience:     a.audience.Audience,
	}

	resp, err := a.app.SwapToken(context.Background(), req)
	a.Require().NoError(err)

	a.Equal("Bearer", resp.TokenType)
	newClaims := &models.TokenClaims{}
	err = a.app.jwtClient.ValidateAndParseClaims(resp.AccessToken, newClaims)
	a.Require().NoError(err)
	a.Require().Len(newClaims.Audience, 1)
	a.Equal(a.audience.Audience, newClaims.Audience[0])
	a.Equal(clientID.String(), newClaims.ClientID)
	a.Equal(subID.String(), newClaims.Subject)
}

type mockAudienceProvider struct {
	Res []string
	Err error
}

func (m *mockAudienceProvider) GetAudience(ctx context.Context, req *models.TokenClientFlowRequest, scopes []*models.Scope) ([]string, error) {
	return m.Res, m.Err
}

func (a *TestAppSuite) TestTokenClientFlowTokenCreated() {
	tests := []struct {
		Scenario          string
		GetAudienceResult []string
		GetAudienceError  error
	}{
		{
			Scenario:          "Nil audience",
			GetAudienceResult: nil,
			GetAudienceError:  nil,
		},
		{
			Scenario:          "Single Audience",
			GetAudienceResult: []string{"single_audience"},
		},
		{
			Scenario: "Multiple Audiences",
		},
		{
			Scenario: "Audience error",
		},
	}

	for _, test := range tests {
		audienceProvider := &mockAudienceProvider{
			Res: test.GetAudienceResult,
			Err: test.GetAudienceError,
		}
		app := a.createApp(audienceProvider, a.app.identityProvider)

		token, err := app.TokenClientFlow(context.Background(),
			&models.TokenClientFlowRequest{
				GrantType: "client_credentials",
				Metadata: models.TokenRequestMetadata{
					Username: a.client.ID.String(),
					Password: "test^123",
				},
			})

		a.Require().NoError(err)

		a.Equal(token.RefreshToken, "", "RefreshToken should not be given for client flow")
		a.NotEmpty(token.AccessToken)

		claims := &models.TokenClaims{}
		err = a.jwtClient.ValidateAndParseClaims(token.AccessToken, claims)
		a.Require().NoError(err, "error parsing claims from token")

		a.Equal(a.client.ID.String(), claims.ClientID)
		a.Equal(a.client.ID.String(), claims.Subject)
	}
}

func (a *TestAppSuite) TestParseScopeString() {
	tests := []struct {
		Scenario string
		Scope    string
		Expected []string
	}{
		{
			Scenario: "empty scope",
			Scope:    "",
			Expected: []string{},
		},
		{
			Scenario: "single scope",
			Scope:    "test:stuff",
			Expected: []string{"test:stuff"},
		},
		{
			Scenario: "multiple scopes",
			Scope:    "test:stuff test:stuff2 test:stuff3",
			Expected: []string{"test:stuff", "test:stuff2", "test:stuff3"},
		},
		{
			Scenario: "duplicate scopes",
			Scope:    "test:stuff test:stuff test:stuff2",
			Expected: []string{"test:stuff", "test:stuff2"},
		},
	}

	for _, test := range tests {
		a.Run(test.Scenario, func() {
			scopes := a.app.parseScopesString(test.Scope)
			a.Equal(test.Expected, scopes)
		})
	}
}

func (a *TestAppSuite) TestClientTokenErrorOnWrongPassword() {
	_, err := a.app.TokenClientFlow(context.Background(), &models.TokenClientFlowRequest{
		GrantType: "client_credentials",
		Scope:     "",
		Metadata: models.TokenRequestMetadata{
			Username: a.client.ID.String(),
			Password: "wrong_password",
		},
	})

	a.Require().Error(err)
	a.Contains(err.Error(), "incorrect password")
}

func (a *TestAppSuite) TestClientTokenErrorWongGrantType() {
	err := (&pb.TokenClientFlowRequest{
		GrantType: "authorization_code",
		Scope:     "",
	}).Validate()
	a.Require().Error(err)
}

func (a *TestAppSuite) TestTokenCreated() {
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
				AcceptedAccessModels: []*models.AccessModel{a.accessModelGql1},
				Audience:             a.audience,
				Client:               a.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       a.redirectTarget,
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
			ExpectedAudience:       a.audience.Audience,
			ExpectedClientID:       a.client.ID.String(),
			ExpectedSubject:        "translatedPseudo1ForNid",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            2,
		},
		{
			Name:           "SuccessCase2",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{a.accessModelGql1},
				RequiredAccessModels: []*models.AccessModel{a.accessModelGql2},
				Audience:             a.audience,
				Client:               a.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       a.redirectTarget,
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
			ExpectedAudience:       a.audience.Audience,
			ExpectedClientID:       a.client.ID.String(),
			ExpectedSubject:        "translatedPseudo2ForNid",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            2,
		},
		{
			Name:           "SuccessCase3",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{a.accessModelGql1},
				RequiredAccessModels: []*models.AccessModel{a.accessModelRestPOST},
				Audience:             a.audience,
				Client:               a.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       a.redirectTarget,
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
			ExpectedAudience:       a.audience.Audience,
			ExpectedClientID:       a.client.ID.String(),
			ExpectedSubject:        "translatedPseudo2ForNid",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            2,
		},
		{
			Name:           "SuccessCase4",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{a.accessModelRestPOST},
				RequiredAccessModels: []*models.AccessModel{a.accessModelRestGET},
				Audience:             a.audience,
				Client:               a.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       a.redirectTarget,
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
			ExpectedAudience:       a.audience.Audience,
			ExpectedClientID:       a.client.ID.String(),
			ExpectedSubject:        "translatedPseudo2ForNid",
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            2,
		},
		{
			Name:           "SuccessCase5EmptySubject",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{a.accessModelRestPOST},
				RequiredAccessModels: []*models.AccessModel{a.accessModelRestGET},
				Audience:             a.audience,
				Client:               a.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       a.redirectTarget,
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
			ExpectedAudience:       a.audience.Audience,
			ExpectedClientID:       a.client.ID.String(),
			ExpectedClientMetadata: map[string]interface{}{"oin": "000012345"},
			PseudoCalls:            0,
		},
	}

	for _, test := range tests {
		a.Run(test.Name, func() {
			a.mockPseudonymizer = setupMockPseudonymizer(context.Background())
			a.app.pseudonymizer = a.mockPseudonymizer

			authorizationCode, err := authtoken.NewToken(a.conf.AuthorizationCodeLength)
			a.Require().NoError(err, "error creating test authorization code")
			hash, err := authtoken.Hash(authorizationCode)
			a.Require().NoError(err)
			test.Session.AuthorizationCode = &hash
			err = a.AuthDBSuite.Tx.Create(test.Session).Error
			a.Require().NoError(err, "error creating test session")

			res, err := a.app.Token(context.Background(),
				&models.TokenRequest{
					AuthorizationCode: authorizationCode,
					GrantType:         "authorization_code",
					Metadata: models.TokenRequestMetadata{
						Username: test.Session.Client.ID.String(),
						Password: test.ClientPassword,
					},
				})
			a.Require().NoError(err, "unexpected error calling token endpoint")

			a.mockPseudonymizer.AssertNumberOfCalls(a.T(), "GetPseudonym", test.PseudoCalls)

			a.Equal("Bearer", res.TokenType)

			claims := &models.TokenClaims{}
			err = a.jwtClient.ValidateAndParseClaims(res.AccessToken, claims)
			a.Require().NoError(err, "error parsing claims from token")

			a.Require().Len(claims.Audience, 1)
			a.Equal(test.ExpectedAudience, claims.Audience[0])
			a.Equal(test.ExpectedClientID, claims.ClientID)

			scopes := claims.Scopes

			a.Equal(test.ExpectedScopes, scopes)

			a.Equal(test.ExpectedSubjects, claims.Subjects)

			a.NotEqual("", claims.Subject)

			a.Equal(issuer, claims.Issuer)

			a.Equal(test.ExpectedClientMetadata, claims.ClientMetadata)
		})
	}
}

func (a *TestAppSuite) TestErrorOnEmptySubjectToken() {
	tests := []struct {
		Name           string
		Session        *models.Session
		ClientPassword string
	}{
		{
			Name:           "ErrorSubject",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{a.accessModelGql3},
				Audience:             a.audience,
				Client:               a.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       a.redirectTarget,
				Subject:              "",
			},
		},
		{
			Name:           "ErrorBsn",
			ClientPassword: "test^123",
			Session: &models.Session{
				AcceptedAccessModels: []*models.AccessModel{a.accessModelGql4},
				Audience:             a.audience,
				Client:               a.client,
				State:                models.SessionStateCodeGranted,
				RedirectTarget:       a.redirectTarget,
				Subject:              "",
			},
		},
	}

	for _, test := range tests {
		a.Run(test.Name, func() {
			authorizationCode, err := authtoken.NewToken(a.conf.AuthorizationCodeLength)
			a.Require().NoError(err, "error creating test authorization code")
			hash, err := authtoken.Hash(authorizationCode)
			a.Require().NoError(err)
			test.Session.AuthorizationCode = &hash
			err = a.AuthDBSuite.Tx.Create(test.Session).Error
			a.Require().NoError(err, "error creating test session")

			_, err = a.app.Token(context.Background(), &models.TokenRequest{
				AuthorizationCode: authorizationCode,
				GrantType:         "authorization_code",
				Metadata: models.TokenRequestMetadata{
					Username: test.Session.Client.ID.String(),
					Password: test.ClientPassword,
				},
			})
			a.Error(err)
		})
	}
}

func (a *TestAppSuite) TestErrorOnWrongSessionStateAndDeadline() {
	tests := []struct {
		name            string
		state           models.SessionState
		errorExpected   error
		codeGrantedTime time.Time
	}{
		{
			name:          "Unclaimed",
			state:         models.SessionStateUnclaimed,
			errorExpected: errors.Wrapf(contract.ErrInvalidArguments, contract.ErrUnableToRetrieveTokenInvalidState.Error()),
		},
		{
			name:          "Claimed",
			state:         models.SessionStateClaimed,
			errorExpected: contract.ErrUnableToRetrieveTokenInvalidState,
		},
		{
			name:          "Accepted",
			state:         models.SessionStateAccepted,
			errorExpected: contract.ErrUnableToRetrieveTokenInvalidState,
		},
		{
			name:          "Rejected",
			state:         models.SessionStateRejected,
			errorExpected: contract.ErrUnableToRetrieveTokenInvalidState,
		},
		{
			name:          "TokenGranted",
			state:         models.SessionStateTokenGranted,
			errorExpected: contract.ErrUnableToRetrieveTokenInvalidState,
		},
		{
			name:          "CodeGranted",
			state:         models.SessionStateCodeGranted,
			errorExpected: nil,
		},
		{
			name:            "CodeGranted",
			state:           models.SessionStateCodeGranted,
			errorExpected:   errors.Wrapf(contract.ErrInvalidArguments, "session is expired"),
			codeGrantedTime: time.Now().Add(time.Minute - 20),
		},
		{
			name:            "CodeGranted",
			state:           models.SessionStateCodeGranted,
			errorExpected:   errors.Wrapf(contract.ErrInvalidArguments, "session is expired"),
			codeGrantedTime: time.Now().Add(time.Minute - 10),
		},
	}

	for _, test := range tests {
		a.Run(test.name, func() {
			authorizationCode, err := authtoken.NewToken(a.conf.AuthorizationCodeLength)
			a.Require().NoError(err, "error creating test authorization code")
			hash, err := authtoken.Hash(authorizationCode)
			a.Require().NoError(err)
			dummySession := &models.Session{
				AcceptedAccessModels: []*models.AccessModel{a.accessModelGql1},
				RequiredAccessModels: []*models.AccessModel{a.accessModelGql2},
				Audience:             a.audience,
				Client:               a.client,
				State:                test.state,
				RedirectTarget:       a.redirectTarget,
				Subject:              "pseudonym2",
				AuthorizationCode:    &hash,
			}
			err = a.AuthDBSuite.Tx.Create(dummySession).Error
			a.Require().NoError(err, "error creating test session")

			if !test.codeGrantedTime.IsZero() {
				err = a.AuthDBSuite.Tx.Model(models.Session{}).Where("id = ?", dummySession.ID).Update("authorization_code_granted_at", test.codeGrantedTime).Error
				a.Require().NoError(err, "error updating authorization_code_granted_at for test session")
			}

			_, err = a.app.Token(context.Background(), &models.TokenRequest{
				AuthorizationCode: authorizationCode,
				GrantType:         "authorization_code",
				Metadata: models.TokenRequestMetadata{
					Username: a.client.ID.String(),
					Password: "test^123",
				},
			})

			if test.errorExpected != nil {
				a.Require().Error(err)
				a.ErrorContains(err, test.errorExpected.Error())
			} else {
				a.Require().NoError(err)
			}
		})
	}
}

func (a *TestAppSuite) TestErrorOnNonExistingSession() {
	code, err := uuid.NewV4()
	a.Require().NoError(err)
	_, err = a.app.Token(context.Background(), &models.TokenRequest{
		AuthorizationCode: code.String(),
		GrantType:         "authorization_code",
		Metadata: models.TokenRequestMetadata{
			Username: a.client.ID.String(),
			Password: "test^123",
		},
	})
	a.Require().Error(err)
	a.Contains(err.Error(), "session not found")
}

func (a *TestAppSuite) TestErrorNotFoundOnOtherClientsSession() {
	authorizationCode, err := authtoken.NewToken(a.conf.AuthorizationCodeLength)
	a.Require().NoError(err, "error creating test authorization code")
	hash, err := authtoken.Hash(authorizationCode)
	a.Require().NoError(err)
	dummySession := &models.Session{
		AcceptedAccessModels: []*models.AccessModel{a.accessModelGql1},
		RequiredAccessModels: []*models.AccessModel{a.accessModelGql2},
		Audience:             a.audience,
		Client:               a.client,
		State:                models.SessionStateCodeGranted,
		RedirectTarget:       a.redirectTarget,
		Subject:              "pseudonym2",
		AuthorizationCode:    &hash,
	}
	err = a.AuthDBSuite.Tx.Create(dummySession).Error
	a.Require().NoError(err, "error creating test session")

	a.Require().NoError(err)
	_, err = a.app.Token(context.Background(), &models.TokenRequest{
		AuthorizationCode: authorizationCode,
		GrantType:         "authorization_code",
		Metadata: models.TokenRequestMetadata{
			Username: a.client2.ID.String(),
			Password: "456#%test",
		},
	})
	a.Require().Error(err)
	a.Contains(err.Error(), "session not found")
}

func (a *TestAppSuite) TestErrorOnWrongPassword() {
	authorizationCode, err := authtoken.NewToken(a.conf.AuthorizationCodeLength)
	a.Require().NoError(err, "error creating test authorization code")
	hash, err := authtoken.Hash(authorizationCode)
	a.Require().NoError(err)
	dummySession := &models.Session{
		AcceptedAccessModels: []*models.AccessModel{a.accessModelGql1},
		RequiredAccessModels: []*models.AccessModel{a.accessModelGql2},
		Audience:             a.audience,
		Client:               a.client,
		State:                models.SessionStateCodeGranted,
		RedirectTarget:       a.redirectTarget,
		Subject:              "pseudonym2",
		AuthorizationCode:    &hash,
	}
	err = a.AuthDBSuite.Tx.Create(dummySession).Error
	a.Require().NoError(err, "error creating test session")

	a.Require().NoError(err)
	_, err = a.app.Token(context.Background(), &models.TokenRequest{
		AuthorizationCode: authorizationCode,
		GrantType:         "authorization_code",
		Metadata: models.TokenRequestMetadata{
			Username: a.client.ID.String(),
			Password: "456#%test",
		},
	})
	a.Require().Error(err)
	a.Contains(err.Error(), "incorrect password")
}

func (a *TestAppSuite) TestErrorOnNotAuthorizationCodeGrantType() {
	code, err := uuid.NewV4()
	a.Require().NoError(err)
	err = (&pb.TokenRequest{
		TypeValue: &pb.TokenRequest_AuthorizationCode{AuthorizationCode: code.String()},
		GrantType: "implicit",
	}).Validate()
	a.Require().Error(err)
}
