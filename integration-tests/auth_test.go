//go:build integration || to || files
// +build integration to files

package integration

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc/metadata"

	databron "lab.weave.nl/nid/nid-core/integration-tests/databron/models"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	auth "lab.weave.nl/nid/nid-core/svc/auth/models"
	authpb "lab.weave.nl/nid/nid-core/svc/auth/proto"
)

type AuthFlowTest struct {
	*AuthTestConfig
	inScope                 string
	inResponseType          string
	inClientID              string
	inRedirectURI           string
	inAudience              string
	inOptionalScopes        string
	optOptionalAccessModels bool
}

type AuthTestConfig struct {
	ClientPassword string `envconfig:"CLIENT_PASSWORD"`
	DeviceCode     string `envconfig:"DEVICE_CODE"`
	DeviceSecret   string `envconfig:"DEVICE_SECRET"`
	UserPseudo     string `envconfig:"optional,USER_PSEUDO"`
}

func NewAuthTest(audience auth.Audience, required, optional auth.AccessModel, acceptOptional bool, testConfig *AuthTestConfig) AuthFlowTest {
	clientDB := auth.ClientDB{}
	defaultClient := clientDB.DefaultModel()
	redirectTargetDB := auth.RedirectTargetDB{}
	defaultRedirectTarget := redirectTargetDB.DefaultModel(defaultClient.ID)
	defaultAudience := audience
	defaultAccessModel := required
	defaultAccessModelopt := optional

	scope := fmt.Sprintf("openid %s@%s", defaultAccessModel.Name, defaultAccessModel.Hash)
	optscope := fmt.Sprintf("%s@%s", defaultAccessModelopt.Name, defaultAccessModelopt.Hash)

	defaultTest := AuthFlowTest{
		inScope:                 scope,
		inResponseType:          "code",
		inClientID:              defaultClient.ID.String(),
		inRedirectURI:           defaultRedirectTarget.RedirectTarget,
		inAudience:              defaultAudience.Audience,
		inOptionalScopes:        optscope,
		optOptionalAccessModels: acceptOptional,
		AuthTestConfig:          testConfig,
	}

	return defaultTest
}

func AuthorizeTokens(s *BaseTestSuite, protoClient *authpb.AuthClient, httpClient *resty.Client, test AuthFlowTest) (string, string) {
	s.Require().NotNil(s, "suite can't be nil")
	s.Require().NotNil(protoClient, "proto client can't be nil")
	s.Require().NotNil(httpClient, "http client can't be nil")
	s.Require().NotEmpty(test, "test can't be empty")

	c := *protoClient

	var accessToken, tokenType, sessionID, authorizationCode, sessionFinaliseToken string
	var accessModelIds []string

	s.Run("authorize", func() {
		req := httpClient.R()
		query := url.Values{
			"scope":         []string{test.inScope},
			"response_type": []string{test.inResponseType},
			"client_id":     []string{test.inClientID},
			"redirect_uri":  []string{test.inRedirectURI},
			"audience":      []string{test.inAudience},
		}

		scopeURL := url.URL{Path: test.inScope}
		scopeEncoded := scopeURL.String()[2:]

		optionalScopesURL := url.URL{Path: test.inOptionalScopes}
		optionalScopesEncoded := optionalScopesURL.String()[2:]

		res, err := req.Get("/authorize?" + query.Encode() + "&scope=" + scopeEncoded + "&optional_scopes=" + optionalScopesEncoded)
		s.Require().NoError(err)

		headers := res.Header()
		s.Require().NotNil(headers)

		s.Require().Equal(http.StatusFound, res.StatusCode(), "body: %s, headers: %v", string(res.Body()), headers)

		// get session id
		locationList, ok := headers["Location"]
		s.Require().True(ok, "%v", headers)
		s.Require().NotEmpty(locationList)
		s.Require().Len(locationList, 1)
		location := locationList[0]
		splitLocation := strings.Split(location, "#")
		s.Require().Len(splitLocation, 2)
		sessionID = splitLocation[len(splitLocation)-1]
		requestSessionState(s, c, sessionID, authpb.SessionState_UNCLAIMED)
	})

	s.Require().NotEmpty(sessionID)

	basicAuthHeader := "basic " + base64.StdEncoding.EncodeToString([]byte(test.DeviceCode+":"+test.DeviceSecret))
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", basicAuthHeader)
	signInResponse, err := s.walletAuthClient.SignIn(ctx, &empty.Empty{})
	s.Require().NoError(err)

	walletAuthHeader := "Bearer " + signInResponse.Bearer

	s.Run("session-token", func() {
		ctx = metadata.AppendToOutgoingContext(s.ctx, "authorization", walletAuthHeader)
		tokenReq := &authpb.SessionRequest{
			SessionId: sessionID,
		}
		resp, err := c.GenerateSessionFinaliseToken(ctx, tokenReq)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().NotNil(resp.GetFinaliseToken())

		sessionFinaliseToken = resp.GetFinaliseToken()
	})

	s.Run("claim", func() {
		ctx = metadata.AppendToOutgoingContext(s.ctx, "authorization", walletAuthHeader)
		claimReq := &authpb.SessionRequest{
			SessionId: sessionID,
		}
		resp, err := c.Claim(ctx, claimReq)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().NotNil(resp.GetClient())
		s.Require().Equal(test.inClientID, resp.GetClient().Id)
		s.Require().NotNil(resp.GetAudience())
		s.Require().Equal(test.inAudience, resp.GetAudience().Audience)
		s.Require().NotNil(resp.GetRequiredAccessModels())
		s.Require().Len(resp.GetOptionalAccessModels(), len(strings.Split(test.inOptionalScopes, " ")))
		requestSessionState(s, c, sessionID, authpb.SessionState_CLAIMED)

		if test.optOptionalAccessModels {
			for _, model := range resp.GetOptionalAccessModels() {
				accessModelIds = append(accessModelIds, model.Id)
			}
		}
	})

	s.Run("accept", func() {
		ctx = metadata.AppendToOutgoingContext(s.ctx, "authorization", walletAuthHeader)
		acceptReq := &authpb.AcceptRequest{
			SessionId:      sessionID,
			AccessModelIds: accessModelIds,
		}
		resp, err := c.Accept(ctx, acceptReq)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().NotNil(resp.GetClient())
		s.Require().Equal(test.inClientID, resp.GetClient().Id)
		s.Require().NotNil(resp.GetAudience())
		s.Require().Equal(test.inAudience, resp.GetAudience().Audience)
		s.Require().NotNil(resp.GetRequiredAccessModels())
		requestSessionState(s, c, sessionID, authpb.SessionState_ACCEPTED)
	})

	s.Run("finalise", func() {
		sessionRequest := &authpb.FinaliseRequest{SessionId: sessionID, SessionFinaliseToken: sessionFinaliseToken}

		response, err := c.Finalise(s.ctx, sessionRequest)
		s.Require().NoError(err)
		s.Require().NotNil(response)
		s.NotEmpty(response.GetRedirectLocation())

		location := response.GetRedirectLocation()
		u, err := url.Parse(location)
		s.Require().NoError(err)
		m, _ := url.ParseQuery(u.RawQuery)
		codeList, ok := m["authorization_code"]
		s.Require().True(ok)
		s.Require().Len(codeList, 1)
		authorizationCode = codeList[0]
		requestSessionState(s, c, sessionID, authpb.SessionState_CODE_GRANTED)
	})

	s.Run("token", func() {
		tokenReq := &authpb.TokenRequest{
			GrantType:         "authorization_code",
			AuthorizationCode: authorizationCode,
		}

		basicAuthHeader := "basic " + base64.StdEncoding.EncodeToString([]byte(test.inClientID+":"+test.ClientPassword))
		ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", basicAuthHeader)

		tokenResp, err := c.Token(ctx, tokenReq)
		s.Require().NoError(err)
		s.Require().NotNil(tokenResp)
		s.Require().NotEmpty(tokenResp.GetAccessToken())
		s.Require().NotEmpty(tokenResp.GetTokenType())
		accessToken = tokenResp.GetAccessToken()
		tokenType = tokenResp.GetTokenType()
		requestSessionState(s, c, sessionID, authpb.SessionState_TOKEN_GRANTED)
	})

	return accessToken, tokenType
}

func requestSessionState(s *BaseTestSuite, c authpb.AuthClient, sessionID string, state authpb.SessionState) {
	statusReq := &authpb.SessionRequest{SessionId: sessionID}
	statusResp, err := c.Status(context.Background(), statusReq)
	s.Require().NoError(err)
	s.Require().Equal(state, statusResp.GetState())
}

func GetAuthTestConfig() (*AuthTestConfig, error) {
	config := &AuthTestConfig{}
	err := envconfig.Init(config)
	if err != nil {
		return nil, errors.Wrap(err, "initialising env")
	}

	if config.UserPseudo == "" {
		config.UserPseudo = DefaultUserPseudonym()
	}

	return config, nil
}

// UserPseudonym retrieves user pseudonym
func DefaultUserPseudonym() string {
	userDB := databron.UserDB{}
	defaultUser := userDB.DefaultModel()
	return defaultUser.Pseudonym
}
