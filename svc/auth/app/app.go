package app

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/nID-sourcecode/nid-core/pkg/authtoken"
	"github.com/nID-sourcecode/nid-core/pkg/gqlutil"
	"github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/pseudonym"
	"github.com/nID-sourcecode/nid-core/pkg/sliceutil"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	utilJWT "github.com/nID-sourcecode/nid-core/pkg/utilities/jwt/v3"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/config"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/repository"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/stats"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	walletPB "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	accessModelIdentifierPartAmount = 2
	authCodeParamKey                = "authorization_code"
	accessModelGeneratedName        = "Generated"
	sessionNotFound                 = "session not found"
	tokenType                       = "Bearer"
)

// App deals with the business logic for auth service.
type App struct {
	repo *repository.AuthDB
	conf *config.AuthConfig

	stats               *stats.Stats
	schemaFetcher       gqlutil.SchemaFetcher
	walletClient        walletPB.WalletClient
	pseudonymizer       pseudonym.Pseudonymizer
	jwtClient           *utilJWT.Client
	jwtParserUnverified *jwt.Parser
	callbackHandler     contract.CallbackHandler
	passwordManager     password.IManager
	audienceProvider    contract.AudienceProvider
	identityProvider    contract.IdentityProvider
}

// New returns new instance of app.App.
func New(conf *config.AuthConfig, db *repository.AuthDB, schemaFetcher gqlutil.SchemaFetcher,
	stats *stats.Stats, handler contract.CallbackHandler, passwordManager password.IManager,
	jwtClient *utilJWT.Client, pseudonymizer pseudonym.Pseudonymizer, client walletPB.WalletClient,
	audienceProvider contract.AudienceProvider, identityProvider contract.IdentityProvider,
) *App {
	return &App{
		conf:                conf,
		schemaFetcher:       schemaFetcher,
		stats:               stats,
		repo:                db,
		jwtParserUnverified: jwt.NewParser(),
		callbackHandler:     handler,
		passwordManager:     passwordManager,
		jwtClient:           jwtClient,
		pseudonymizer:       pseudonymizer,
		walletClient:        client,
		audienceProvider:    audienceProvider,
		identityProvider:    identityProvider,
	}
}

func (a *App) AuthorizeHeadless(ctx context.Context, req *models.AuthorizeHeadlessRequest) error {
	session := &models.Session{
		State: models.SessionStateCodeGranted,
	}

	if err := a.fillSessionFromRequest(ctx, session, req); err != nil {
		return err
	}

	code, err := a.setSessionAuthenticationCode(session)
	if err != nil {
		return err
	}

	accessModelHash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.QueryModelJSON)))
	accessModel, err := a.repo.AccessModelDB.GetAccessModelByAudienceWithScope(accessModelGeneratedName, accessModelHash, session.Audience)
	if err != nil {
		if req.QueryModelPath == "" {
			req.QueryModelPath = "/gql"
		}
		accessModel = &models.AccessModel{
			AudienceID:  session.AudienceID,
			Description: "",
			Hash:        accessModelHash,
			GqlAccessModel: &models.GqlAccessModel{
				JSONModel: req.QueryModelJSON,
				Path:      req.QueryModelPath,
			},
			Name: accessModelGeneratedName,
			Type: models.AccessModelTypeGQL,
		}

		err = a.repo.AccessModelDB.CreateAccessModel(accessModel)
		if err != nil {
			log.Extract(ctx).WithError(err).Error("error inserting access_model")
			return contract.ErrInternalError
		}
	}

	session.RequiredAccessModels = []*models.AccessModel{accessModel}
	session.AcceptedAccessModels = []*models.AccessModel{accessModel}

	err = a.repo.SessionDB.CreateSession(session)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("inserting session")
		return contract.ErrInternalError
	}

	if err := a.callbackHandler.HandleCallback(ctx, session.RedirectTarget.RedirectTarget, code); err != nil {
		log.Extract(ctx).WithError(err).Error("error handling callback")
		return contract.ErrInternalError
	}

	return nil
}

// Authorize handles the Authorize request for authservice. Returns a redirect URL.
func (a *App) Authorize(ctx context.Context, req *models.AuthorizeRequest) (string, error) {
	if strings.TrimSpace(req.ResponseType) != "code" {
		return "", errors.Wrapf(contract.ErrInvalidArguments, "no response type other than \"code\" is supported")
	}

	session := &models.Session{
		State: models.SessionStateUnclaimed,
	}

	if err := a.fillSessionFromRequest(ctx, session, req); err != nil {
		return "", err
	}

	requiredAccessModels, err := a.getRequiredAccessModels(ctx, req.Scope, session.Audience)
	if err != nil {
		return "", err
	}
	session.RequiredAccessModels = requiredAccessModels

	optionalAccessModels, err := a.getOptionalAccessModels(ctx, req.OptionalScopes, session.Audience)
	if err != nil {
		return "", err
	}
	session.OptionalAccessModels = optionalAccessModels

	err = a.repo.SessionDB.CreateSession(session)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("inserting session")
		return "", contract.ErrInternalError
	}

	return fmt.Sprintf("%s#%s", a.conf.AuthRequestURI, session.ID.String()), nil
}

// Accept
func (a *App) Accept(ctx context.Context, jwtPayload string, req *models.AcceptRequest) (*models.SessionResponse, error) {
	subject, err := a.getSubjectFromJWTPayload(jwtPayload)
	if err != nil {
		return nil, err
	}

	session, err := a.repo.SessionDB.GetSessionByIDAndSubject(repository.PreloadRequiredAndOptionalScopes, req.SessionID, subject)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(contract.ErrNotFound, sessionNotFound)
		}
		return nil, contract.ErrInternalError
	}

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	if isExpired {
		log.Extract(ctx).Info("session is expired")

		return nil, errors.Wrapf(contract.ErrInvalidArguments, "session is expired")
	}

	if session.State != models.SessionStateClaimed {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "precondition failed session not claimed")
	}

	suppliedAccessModels := req.AccessModelIds

	// Get access_models from payload ids
	accessModels, err := a.repo.AccessModelDB.GetAccessModelsByIDs(suppliedAccessModels)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(contract.ErrNotFound, "found %d but got %d access_model_ids", len(accessModels), len(suppliedAccessModels))
		}
		log.Extract(ctx).WithField("access_model_ids", suppliedAccessModels).WithError(err).Error("querying access_models failed")
		return nil, contract.ErrInternalError
	}

	// Check if all found access_models are in sessions required or optional access_models
	err = accessModelExists(accessModels, session.RequiredAccessModels, session.OptionalAccessModels)
	if err != nil {
		return nil, err
	}

	// All supplied access models are valid save in accepted_access_models session association
	// Combine access_models with required access_models
	combinedAccessModels := append(accessModels, session.RequiredAccessModels...) // nolint:gocritic
	err = a.repo.SessionDB.UpdateAcceptedAccessModels(session, combinedAccessModels)
	if err != nil {
		log.Extract(ctx).WithField("access_model_ids", strings.Join(suppliedAccessModels, ",")).WithError(err).Error("updating accepted access models failed")
		return nil, contract.ErrInternalError
	}

	// Update accepted session to state claimed
	err = a.repo.SessionDB.UpdateSessionState(session, models.SessionStateAccepted)
	if err != nil {
		log.Extract(ctx).WithField("session_id", req.SessionID).WithError(err).Error("updating session state failed")
		return nil, contract.ErrInternalError
	}

	session.State = models.SessionStateAccepted
	session.AcceptedAccessModels = combinedAccessModels

	return sessionToResponse(session), nil
}

// Claim
func (a *App) Claim(ctx context.Context, jwtPayload string, req *models.SessionRequest) (*models.SessionResponse, error) {
	subject, err := a.getSubjectFromJWTPayload(jwtPayload)
	if err != nil {
		return nil, errors.Wrapf(err, "getting subject from jwt")
	}

	if subject == "" {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "subject cannot be empty")
	}

	session, err := a.repo.SessionDB.GetSessionByID(repository.PreloadRequiredAndOptionalScopes, req.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(contract.ErrNotFound, "session is expired")
		}
		return nil, contract.ErrInternalError
	}

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	if isExpired {
		log.Extract(ctx).Info("session is expired")

		return nil, errors.Wrapf(contract.ErrInvalidArguments, "session is expired")
	}

	if session.State != models.SessionStateUnclaimed {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "precondition failed session not unclaimed")
	}

	// Update accepted session to state claimed
	err = a.repo.SessionDB.UpdateSessionState(session, models.SessionStateClaimed)
	if err != nil {
		log.Extract(ctx).WithField("session_id", req.SessionID).WithError(err).Error("updating session state failed")
		return nil, contract.ErrInternalError
	}

	// Update accepted session subject
	err = a.repo.SessionDB.UpdateSessionSubject(session, subject)
	if err != nil {
		log.Extract(ctx).WithField("session_id", req.SessionID).WithError(err).Error("updating session subject failed")
		return nil, contract.ErrInternalError
	}
	session.State = models.SessionStateClaimed

	return sessionToResponse(session), nil
}

func (a *App) Reject(ctx context.Context, req *models.SessionRequest) error {
	session, err := a.repo.SessionDB.GetSessionByID(repository.NoPreload, req.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(contract.ErrNotFound, sessionNotFound)
		}
		return contract.ErrInternalError
	}

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	if isExpired {
		log.Extract(ctx).Info("session is expired")

		return errors.Wrapf(contract.ErrInvalidArguments, "session is expired")
	}

	// Check if session is claimed
	if session.State != models.SessionStateClaimed {
		return errors.Wrapf(contract.ErrInvalidArguments, "precondition failed session not unclaimed")
	}

	// Update accepted session to state rejected
	err = a.repo.SessionDB.UpdateSessionState(session, models.SessionStateRejected)
	if err != nil {
		log.Extract(ctx).WithField("session_id", session.ID).WithError(err).Error("updating session state failed")
		return contract.ErrInternalError
	}

	return nil
}

func (a *App) GetSessionDetails(ctx context.Context, req *models.SessionRequest) (*models.SessionResponse, error) {
	session, err := a.repo.SessionDB.GetSessionByID(repository.PreloadAll, req.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(contract.ErrNotFound, sessionNotFound)
		}
		return nil, contract.ErrInternalError
	}

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	if isExpired {
		log.Extract(ctx).Info("session is expired")

		return nil, errors.Wrapf(contract.ErrInvalidArguments, "session is expired")
	}

	if session.State == models.SessionStateUnclaimed || session.State == models.SessionStateClaimed {
		return sessionToResponse(session), nil
	}

	return nil, errors.Wrapf(contract.ErrInvalidArguments, "session has wrong state")
}

// GenerateSessionFinaliseToken Generates a password for the session, which can be reused to check if the session is valid.
// The session token expires if the session's state becomes claimed.
func (a *App) GenerateSessionFinaliseToken(ctx context.Context, req *models.SessionRequest) (*models.SessionAuthorization, error) {
	session, err := a.repo.SessionDB.GetSessionByID(repository.NoPreload, req.SessionID)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("could not find the session")

		return nil, contract.ErrInternalError
	}

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	if isExpired {
		log.Extract(ctx).Info("session is expired")

		return nil, errors.Wrap(contract.ErrInvalidArguments, "session is expired")
	}

	if session.State != models.SessionStateUnclaimed {
		return nil, errors.Wrap(contract.ErrInvalidArguments, "precondition failed session not accepted")
	}

	token, err := authtoken.NewToken(a.conf.AuthorizationCodeLength)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error signing token when swapping")

		return nil, contract.ErrInternalError
	}

	hash, err := a.passwordManager.GenerateHash(token)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error generating hash from token")

		return nil, contract.ErrInternalError
	}

	err = a.repo.SessionDB.SetSessionFinaliseToken(session, hash)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("could not set the finalise token for the session")

		return nil, contract.ErrInternalError
	}

	return &models.SessionAuthorization{FinaliseToken: token}, nil
}

func (a *App) Status(ctx context.Context, req *models.SessionRequest) (*models.StatusResponse, error) {
	session, err := a.repo.SessionDB.GetSessionByID(repository.NoPreload, req.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(contract.ErrNotFound, sessionNotFound)
		}
		return nil, contract.ErrInternalError
	}

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	if isExpired {
		log.Extract(ctx).Info("session is expired")

		return nil, errors.Wrapf(contract.ErrInvalidArguments, "session is expired")
	}

	return &models.StatusResponse{
		State: session.State,
	}, nil
}

// Finalise updates the session state to code_granted and redirects to the redirect target with included authorization_code
func (a *App) Finalise(ctx context.Context, req *models.FinaliseRequest) (*models.FinaliseResponse, error) {
	session, err := a.repo.SessionDB.GetSessionByID(repository.PreloadRequiredAndOptionalScopes, req.SessionID)
	if err != nil {
		log.WithError(err).Error("getting session by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(contract.ErrNotFound, sessionNotFound)
		}
		return nil, contract.ErrInternalError
	}

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	if isExpired {
		log.Extract(ctx).Info("session is expired")

		return nil, errors.Wrapf(contract.ErrInvalidArguments, "session is expired")
	}

	logger := log.Extract(ctx).WithField("session_id", session.ID)

	// Check if session is accepted
	if session.State != models.SessionStateAccepted {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "precondition failed session not accepted")
	}

	// Check if the given password matches with the session
	if req.SessionFinaliseToken == "" {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "did not provide a pass for session")
	}

	ok, err := a.passwordManager.ComparePassword(req.SessionFinaliseToken, session.FinaliseToken)
	if err != nil {
		return nil, contract.ErrInternalError
	}
	if !ok {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "Given Finalise Token did not match")
	}

	code, err := a.setSessionAuthenticationCode(session)
	if err != nil {
		return nil, contract.ErrInternalError
	}

	// Update accepted session
	err = a.repo.SessionDB.UpdateSessionAuthorizationCode(session, *session.AuthorizationCode)
	if err != nil {
		logger.WithError(err).Error("updating session status failed")
		return nil, contract.ErrInternalError
	}
	err = a.repo.SessionDB.UpdateSessionState(session, models.SessionStateCodeGranted)
	if err != nil {
		logger.WithError(err).Error("updating session status failed")
		return nil, contract.ErrInternalError
	}

	// Redirect to redirect location
	location := fmt.Sprintf("%s/?%s=%s", session.RedirectTarget.RedirectTarget, authCodeParamKey, code)

	return &models.FinaliseResponse{RedirectLocation: location}, nil
}

type createTokenClaims struct {
	ConsentID      *uuid.UUID
	ClientID       string
	Scopes         interface{}
	ClientMetadata map[string]interface{}
	Subject        string
	Subjects       map[string]interface{}
	Audience       jwt.ClaimStrings
}

// TokenClientFlow creates a token using credentials.
func (a *App) TokenClientFlow(ctx context.Context, req *models.TokenClientFlowRequest) (*models.TokenResponse, error) {
	clientID, err := a.identityProvider.GetIdentity(ctx, &req.Metadata)
	if err != nil {
		log.WithError(err).Error("validating identity")
		return nil, err
	}

	// Parse scopes.
	scopesClaim := a.parseScopesString(req.Scope)
	scopes, err := a.validateScopes(ctx, scopesClaim)
	if err != nil {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "validating scopes in client flow")
	}

	// Check if incorrect scopes are given.
	if len(scopesClaim) != len(scopes) && req.Scope != "" {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "invalid scope")
	}

	audiences, err := a.audienceProvider.GetAudience(ctx, req, scopes)
	if err != nil {
		return nil, err
	}

	claims := &createTokenClaims{
		ClientID: clientID,
		Audience: audiences,
		Subject:  clientID,
		Scopes:   scopesClaim,
	}

	token, err := a.createToken(ctx, claims)
	if err != nil {
		if errors.Is(err, contract.ErrMultipleAudiencesNotAllowed) {
			return nil, errors.Wrapf(contract.ErrInvalidArguments, "a scope you requested has multiple audiences, please specify the audience via the audience parameter")
		}

		log.WithError(err).Error("generating jwt token")
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken: token,
	}, nil
}

// Token creates a token using a refresh token or authorization code.
func (a *App) Token(ctx context.Context, req *models.TokenRequest) (*models.TokenResponse, error) {
	clientID, err := a.identityProvider.GetIdentity(ctx, &req.Metadata)
	if err != nil {
		log.WithError(err).Error("validating identity")
		return nil, err
	}

	if req.GrantType == "refresh_token" {
		return a.RefreshToken(ctx, req.RefreshToken)
	}

	code := req.AuthorizationCode
	logger := log.Extract(ctx).WithFields(log.Fields{
		"authorization_code": code,
		"client_id":          clientID,
	})

	hash, err := authtoken.Hash(code)
	if err != nil {
		logger.WithError(err).Error("generating hash for authorization code failed")
		return nil, contract.ErrInternalError
	}
	session, err := a.repo.SessionDB.GetSessionByCodeAndClientID(repository.PreloadAll, hash, clientID)
	if err != nil {
		if errors.Is(err, contract.ErrUnableToRetrieveTokenExpiration) {
			return nil, errors.Wrapf(contract.ErrDeadlineExceeded, err.Error())
		}
		return nil, errors.Wrapf(contract.ErrNotFound, sessionNotFound)
	}

	isExpired := session.IsSessionExpired(a.conf.AuthorizationCodeExpirationTime)
	if isExpired {
		log.Extract(ctx).Info("session is expired")

		return nil, errors.Wrapf(contract.ErrInvalidArguments, "session is expired")
	}

	if session.State != models.SessionStateCodeGranted {
		// Do not give away information about this session -- since code is not known to be equal to ID
		return nil, errors.Wrapf(contract.ErrInvalidArguments, contract.ErrUnableToRetrieveTokenInvalidState.Error())
	}

	// Set Scopes
	scopes, err := a.getScopesFromSession(ctx, session)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("getting scopes from session")
		return nil, contract.ErrInternalError
	}

	nidSubject, audienceSubject, err := a.getSubjects(ctx, session, scopes)
	if err != nil {
		if errors.Is(err, contract.ErrUnableToRetrieveTokenExpiration) {
			return nil, errors.Wrapf(contract.ErrDeadlineExceeded, err.Error())
		}
		return nil, contract.ErrUnauthenticated
	}

	refreshToken, err := a.createRefreshToken(ctx, session.ID)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("creating refresh token")
		return nil, contract.ErrUnauthenticated
	}

	token, err := a.createTokenForSession(ctx, session, nidSubject, scopes, audienceSubject)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("creating token")
		return nil, contract.ErrUnauthenticated
	}

	tokenParts := strings.Split(token, ".")
	// sanity check. a JWT should always have 3 parts
	// nolint: gomnd
	if len(tokenParts) == 3 {
		log.Extract(ctx).WithField("token", tokenParts[0]+"."+tokenParts[1]).Info("token created")
	}

	return &models.TokenResponse{AccessToken: token, RefreshToken: refreshToken, TokenType: tokenType}, nil
}

func (a *App) SwapToken(ctx context.Context, req *models.SwapTokenRequest) (*models.TokenResponse, error) {
	inTokenClaims := &models.TokenClaims{}
	err := a.jwtClient.ValidateAndParseClaims(req.CurrentToken, inTokenClaims)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to parse token")

		return nil, errors.Wrapf(contract.ErrInvalidArguments, "unable to parse token")
	}

	// FIXME we should probably add some checks instead of just flikkering the input in the token https://lab.weave.nl/twi/core/-/issues/49
	audience, err := a.repo.AudienceDB.GetAudienceByURI(req.Audience)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(contract.ErrInvalidArguments, "audience \"%s\" does not exist", req.Audience)
		}

		log.Extract(ctx).WithError(err).Error("error getting audience")
		return nil, contract.ErrInternalError
	}
	accessModels, err := a.repo.AccessModelDB.GetAccessModelsByAudience(true, audience)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(contract.ErrNotFound, "error getting access_models")
		}

		log.Extract(ctx).WithError(err).Error("error getting access_models")

		return nil, contract.ErrInternalError
	}

	// Set Queries
	scopes := make(map[string]interface{})
	err = addAccessModelsToScopes(scopes, accessModels)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error adding accepted access models to queries")

		return nil, contract.ErrInternalError
	}

	defaultClaims := utilJWT.NewDefaultClaims(time.Duration(a.conf.JWTExpirationHours))
	defaultClaims.Issuer = a.conf.Issuer
	defaultClaims.Audience = []string{req.Audience}
	defaultClaims.Subject = inTokenClaims.Subject

	outTokenClaims := &models.TokenClaims{
		DefaultClaims:  defaultClaims,
		ClientID:       inTokenClaims.ClientID,
		Scopes:         scopes,
		Subjects:       inTokenClaims.Subjects,
		ClientMetadata: inTokenClaims.ClientMetadata,
	}

	newToken, err := a.jwtClient.SignToken(outTokenClaims)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error signing token when swapping")

		return nil, contract.ErrInternalError
	}

	a.stats.TokenSwapped.With(prometheus.Labels{"audience": req.Audience}).Inc()

	tokenParts := strings.Split(newToken, ".")
	// sanity check. a JWT should always have 3 parts
	// nolint: gomnd
	if len(tokenParts) == 3 {
		log.Extract(ctx).WithField("token", tokenParts[0]+"."+tokenParts[1]).Info("token swapped")
	}

	return &models.TokenResponse{AccessToken: newToken, TokenType: tokenType}, nil
}

// RefreshToken creates a new token with a new refresh token
func (a *App) RefreshToken(ctx context.Context, tokenString string) (*models.TokenResponse, error) {
	_, claims, err := a.jwtClient.ParseWithClaims(tokenString)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("parsing string token with claims")
		return nil, errors.Wrapf(contract.ErrUnauthenticated, "invalid refresh token")
	}

	dbRefreshToken, err := a.repo.RefreshTokenDB.GetWithClaims(claims)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("getting refresh token with client id")
		return nil, errors.Wrapf(contract.ErrUnauthenticated, "invalid refresh token")
	}

	err = a.repo.RefreshTokenDB.Delete(ctx, dbRefreshToken.ID)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("deleting all refresh tokens of client")
		return nil, errors.Wrapf(contract.ErrUnauthenticated, "invalid refresh token")
	}

	dbRefreshToken.Session, err = a.repo.SessionDB.GetSessionByID(repository.PreloadAll, dbRefreshToken.Session.ID.String())
	if err != nil {
		log.Extract(ctx).WithError(err).Error("getting and preloading session with session id")
		return nil, errors.Wrapf(contract.ErrUnauthenticated, "invalid refresh token")
	}

	scopes, err := a.getScopesFromSession(ctx, dbRefreshToken.Session)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("getting scopes from session")
		return nil, errors.Wrapf(contract.ErrUnauthenticated, "invalid refresh token")
	}

	nidSubject, audienceSubject, err := a.getSubjects(ctx, dbRefreshToken.Session, scopes)
	if err != nil {
		if errors.Is(err, contract.ErrUnableToRetrieveTokenExpiration) {
			return nil, errors.Wrapf(contract.ErrDeadlineExceeded, err.Error())
		}
		return nil, contract.ErrUnauthenticated
	}

	newToken, err := a.createTokenForSession(ctx, dbRefreshToken.Session, nidSubject, scopes, audienceSubject)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("creating new access token")
		return nil, errors.Wrapf(contract.ErrUnauthenticated, "invalid refresh token")
	}

	refreshToken, err := a.createRefreshToken(ctx, dbRefreshToken.SessionID)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("creating new refresh token")
		return nil, errors.Wrapf(contract.ErrUnauthenticated, "invalid refresh token")
	}

	return &models.TokenResponse{AccessToken: newToken, RefreshToken: refreshToken, TokenType: tokenType}, nil
}

func (a *App) RegisterAccessModel(ctx context.Context, req *models.AccessModelRequest) error {
	if req.ScopeName == "" {
		return errors.Wrapf(contract.ErrInvalidArguments, "no scope name specified")
	}

	if err := a.validateQueryModelForAudience(ctx, req.QueryModelJSON, req.Audience); err != nil {
		return err
	}

	// Get Audience and insert AccessModel
	audience, err := a.repo.AudienceDB.GetAudienceByURI(req.Audience)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(contract.ErrNotFound, "audience \"%s\" does not exist", req.Audience)
		}

		log.Extract(ctx).WithError(err).WithField("audience", req.Audience).Error("getting audience")
		return contract.ErrInternalError
	}

	accessModel := &models.AccessModel{
		AudienceID:  audience.ID,
		Description: req.Description,
		Hash:        fmt.Sprintf("%x", sha256.Sum256([]byte(req.QueryModelJSON))),
		JSONModel:   req.QueryModelJSON,
		Name:        req.ScopeName,
	}

	err = a.repo.AccessModelDB.CreateAccessModel(accessModel)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error inserting access_model")
		return contract.ErrInternalError
	}

	return nil
}

func (a *App) fillSessionFromRequest(ctx context.Context, session *models.Session, req models.CreateSessionRequest) error {
	client, err := a.repo.ClientDB.GetClientByID(req.GetClientID())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(contract.ErrInvalidArguments, "client with id \"%s\" does not exist", req.GetClientID())
		}

		log.Extract(ctx).WithError(err).WithField("id", req.GetClientID()).Error("getting client")
		return contract.ErrInternalError
	}
	session.ClientID = client.ID

	audience, err := a.repo.AudienceDB.GetAudienceByURI(req.GetAudience())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(contract.ErrInvalidArguments, "audience \"%s\" does not exist", req.GetAudience())
		}

		log.Extract(ctx).WithError(err).WithField("audience", req.GetAudience()).Error("getting audience")
		return contract.ErrInternalError
	}
	session.Audience = audience
	session.AudienceID = audience.ID

	redirectTarget, err := a.repo.RedirectTargetDB.GetRedirectTarget(req.GetRedirectURI(), req.GetClientID())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(contract.ErrInvalidArguments, "client with id \"%s\" does not have redirect URI \"%s\"", req.GetClientID(), req.GetRedirectURI())
		}
		log.Extract(ctx).WithError(err).WithField("redirect_target", req.GetRedirectURI()).
			WithField("client_id", req.GetClientID()).Error("getting redirect target")
		return contract.ErrInternalError
	}
	session.RedirectTarget = redirectTarget
	session.RedirectTargetID = redirectTarget.ID

	return nil
}

func (a *App) getRequiredAccessModels(ctx context.Context, scopes string, audience *models.Audience) ([]*models.AccessModel, error) {
	openIDScopeSpecified := false
	requiredAccessModels := make([]*models.AccessModel, 0)
	for _, specifiedScope := range strings.Split(scopes, " ") {
		scope := strings.Trim(specifiedScope, " ")
		if scope == "openid" {
			openIDScopeSpecified = true
			continue
		}
		parts := strings.Split(scope, "@")
		if len(parts) != accessModelIdentifierPartAmount {
			return nil, errors.Wrapf(contract.ErrInvalidArguments, "scope \"%s\" is of invalid format, should be \"name@hash\"", scope)
		}
		name := parts[0]
		hash := parts[1]
		accessModel, err := a.repo.AccessModelDB.GetAccessModelByAudienceWithScope(name, hash, audience)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.Wrapf(contract.ErrInvalidArguments, "audience \"%s\" does not support access model \"%s\"", audience.Audience, scope)
			}

			log.Extract(ctx).WithError(err).Errorf("getting access model %s", scope)
			return nil, contract.ErrInternalError
		}
		requiredAccessModels = append(requiredAccessModels, accessModel)
	}

	if !openIDScopeSpecified {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "the \"openid\" scope must be specified")
	}
	if len(requiredAccessModels) == 0 {
		return nil, errors.Wrapf(contract.ErrInvalidArguments, "at least one access model-scope must be specified")
	}

	return requiredAccessModels, nil
}

func (a *App) getOptionalAccessModels(ctx context.Context, optionalScopes string, audience *models.Audience) ([]*models.AccessModel, error) {
	optionalAccessModels := make([]*models.AccessModel, 0)
	trimmedOptionalScopes := strings.Trim(optionalScopes, " ")

	if trimmedOptionalScopes == "" {
		return optionalAccessModels, nil
	}

	for _, specifiedScope := range strings.Split(trimmedOptionalScopes, " ") {
		scope := strings.Trim(specifiedScope, " ")
		parts := strings.Split(scope, "@")
		if len(parts) != accessModelIdentifierPartAmount {
			return nil, errors.Wrapf(contract.ErrInvalidArguments, "scope \"%s\" is of invalid format, should be \"name@hash\"", scope)
		}

		name := parts[0]
		hash := parts[1]
		accessModel, err := a.repo.AccessModelDB.GetAccessModelByAudienceWithScope(name, hash, audience)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.Wrapf(contract.ErrInvalidArguments, "audience \"%s\" does not support access model \"%s\"", audience.Audience, scope)
			}

			log.Extract(ctx).WithError(err).Errorf("getting access model %s", scope)
			return nil, contract.ErrInternalError
		}
		optionalAccessModels = append(optionalAccessModels, accessModel)
	}

	return optionalAccessModels, nil
}

func (a *App) getSubjectFromJWTPayload(payload string) (string, error) {
	claims := struct {
		Subject string `json:"sub"`

		jwt.Claims
	}{}

	// parse Claims
	var claimBytes []byte
	var err error

	if claimBytes, err = a.jwtParserUnverified.DecodeSegment(payload); err != nil {
		return "", errors.Wrapf(jwt.ErrTokenMalformed, err.Error())
	}
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))
	err = dec.Decode(&claims)
	// Handle decode error
	if err != nil {
		return "", errors.Wrapf(jwt.ErrTokenMalformed, err.Error())
	}

	return claims.Subject, nil
}

func (a *App) setSessionAuthenticationCode(session *models.Session) (string, error) {
	// Generate random authorization code for session
	code, err := authtoken.NewToken(a.conf.AuthorizationCodeLength)
	if err != nil {
		log.WithError(err).Error("generating authorization code failed")

		return "", contract.ErrInternalError
	}

	hash, err := authtoken.Hash(code)
	if err != nil {
		log.WithError(err).Error("generating hash for authorization code failed")

		return "", contract.ErrInternalError
	}

	session.AuthorizationCode = &hash

	return code, nil
}

func (a *App) getScopesFromSession(ctx context.Context, session *models.Session) (map[string]interface{}, error) {
	scopes := make(map[string]interface{})
	err := addAccessModelsToScopes(scopes, session.AcceptedAccessModels)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error adding accepted access models to scopes")

		return nil, err
	}
	err = addAccessModelsToScopes(scopes, session.RequiredAccessModels)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error adding required access models to scopes")

		return nil, err
	}

	return scopes, nil
}

// Ensure that gql query adheres schema defined by service
func (a *App) validateQueryModelForAudience(ctx context.Context, queryModelJSON string, audience string) error {
	schema, err := a.schemaFetcher.FetchSchema(ctx, audience)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error fetching graphql schema")

		return contract.ErrInternalError
	}

	var queryModel map[string]*gqlutil.AccessModel
	if err := json.Unmarshal([]byte(queryModelJSON), &queryModel); err != nil {
		return errors.Wrapf(contract.ErrInvalidArguments, "query model json incorrectly formatted")
	}
	if err := schema.ValidateQueryModel(queryModel); err != nil {
		return errors.Wrapf(contract.ErrInvalidArguments, cases.Title(language.English, cases.NoLower).String(err.Error()))
	}

	return nil
}

type gqlAccessModel struct {
	Type  string                 `json:"t"`
	Path  string                 `json:"p"`
	Model map[string]interface{} `json:"m"`
}

type restAccessModel struct {
	Type   string            `json:"t"`
	Path   string            `json:"p"`
	Query  map[string]string `json:"q"`
	Body   string            `json:"b"`
	Method string            `json:"m"`
}

func addAccessModelsToScopes(scopes map[string]interface{}, accessModels []*models.AccessModel) error {
	for _, accessModel := range accessModels {
		switch accessModel.Type {
		case models.AccessModelTypeGQL:
			if accessModel.GqlAccessModel == nil {
				return errors.Wrapf(contract.ErrInvalidQueryModelType, "GQL type but no GQL access model related (%s)", accessModel.ID.String())
			}
			accessModelAsMap := make(map[string]interface{})
			err := json.Unmarshal([]byte(accessModel.GqlAccessModel.JSONModel), &accessModelAsMap)
			if err != nil {
				return errors.Wrapf(err, "parsing query model json of %s@%s", accessModel.Name, accessModel.Hash)
			}
			scopes[accessModel.Name] = gqlAccessModel{
				Model: accessModelAsMap,
				Path:  accessModel.GqlAccessModel.Path,
				Type:  "GQL",
			}
		case models.AccessModelTypeREST:
			if accessModel.RestAccessModel == nil {
				return errors.Wrapf(contract.ErrInvalidQueryModelType, "REST type but no REST access model related (%s)", accessModel.ID.String())
			}

			var query map[string]string
			err := json.Unmarshal([]byte(accessModel.RestAccessModel.Query), &query)
			if err != nil {
				return errors.Wrapf(err, "marshalling query for REST access model with ID %s", accessModel.RestAccessModel.ID.String())
			}

			scopes[accessModel.Name] = restAccessModel{
				Type:   "REST",
				Path:   accessModel.RestAccessModel.Path,
				Query:  query,
				Body:   accessModel.RestAccessModel.Body,
				Method: accessModel.RestAccessModel.Method,
			}
		default:
			return fmt.Errorf("%w %v", contract.ErrInternalError, errors.Wrap(contract.ErrInvalidQueryModelType, "unable to add access model to scopes"))
		}
	}

	return nil
}

func (a *App) getSubjects(ctx context.Context, session *models.Session, scopes map[string]interface{}) (string, string, error) {
	if session == nil {
		return "", "", errors.New("session is nil")
	}

	if session.Subject == "" {
		if scopeContainsString(scopes, "$$nid:subject$$", "$$nid:bsn$$") {
			log.Extract(ctx).Error("scope requires subject")
			return "", "", contract.ErrInternalError
		}

		return "", "", nil
	}

	nidSubject, err := a.pseudonymizer.GetPseudonym(ctx, session.Subject, a.conf.Namespace)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error getting pseudonym")

		return "", "", contract.ErrInternalError
	}

	audienceSubject, err := a.pseudonymizer.GetPseudonym(ctx, session.Subject, session.Audience.Namespace)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error getting pseudonym")

		return "", "", contract.ErrInternalError
	}

	return nidSubject, audienceSubject, nil
}

func scopeContainsString(scopes map[string]interface{}, subs ...string) bool {
	for _, scope := range scopes {
		for _, sub := range subs {
			if strings.Contains(fmt.Sprint(scope), sub) {
				return true
			}
		}
	}

	return false
}

func (a *App) createRefreshToken(ctx context.Context, sessionID uuid.UUID) (string, error) {
	tokenID, err := uuid.NewV4()
	if err != nil {
		return "", contract.ErrInternalError
	}

	refreshTokenClaims := utilJWT.NewDefaultRefreshTokenClaims(sessionID.String(), tokenID.String(), time.Duration(a.conf.JWTRefreshExpirationHours))
	refreshToken, err := a.jwtClient.SignToken(&refreshTokenClaims)
	if err != nil {
		return "", err
	}

	err = a.repo.RefreshTokenDB.Add(ctx, &models.RefreshToken{
		ID:        tokenID,
		SessionID: sessionID,
	})
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (a *App) createTokenForSession(ctx context.Context, session *models.Session, nidSubject string, scopes map[string]interface{}, audienceSubject string) (string, error) {
	clientMetadata, err := setClientMetadataToClaims(session)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("marshalling client metadata")
		return "", contract.ErrInternalError
	}

	// Create consent id ourselves, so we can use it to create a consent in the wallet service.
	consentID, err := uuid.NewV4()
	if err != nil {
		log.WithError(err).Error("creating consent id")
		return "", contract.ErrInternalError
	}

	claims := &createTokenClaims{
		ConsentID: &consentID,
		ClientID:  session.Client.ID.String(),
		Scopes:    scopes,
		Audience:  []string{session.Audience.Audience},
		Subjects: map[string]interface{}{
			session.Audience.Namespace: audienceSubject,
		},
		ClientMetadata: clientMetadata,
		Subject:        nidSubject,
	}

	token, err := a.createToken(ctx, claims)
	if err != nil {
		return "", contract.ErrInternalError
	}

	scopeNames := make([]string, len(scopes))
	var i int
	for q := range scopes {
		scopeNames[i] = strings.Split(q, "@")[0]
		i++
	}

	// Call the wallet service to create consent.
	if nidSubject != "" {
		err = a.createConsentInUserWallet(ctx, consentID, session.Client.ID.String(), session.Subject, "name", strings.Join(scopeNames, ", "), token)
		if err != nil {
			return "", err
		}
	}

	// Update session state to state token granted.
	err = a.repo.SessionDB.UpdateSessionState(session, models.SessionStateTokenGranted)
	if err != nil {
		log.Extract(ctx).WithField("session_id", session.ID).WithError(err).Error("updating session state failed")
		return "", contract.ErrInternalError
	}

	return token, nil
}

// createToken creates a signed token with the given claims.
func (a *App) createToken(_ context.Context, tokenClaims *createTokenClaims) (string, error) {
	// Set default claims properties
	defaultClaims := utilJWT.NewDefaultClaims(time.Duration(a.conf.JWTExpirationHours))
	defaultClaims.Issuer = a.conf.Issuer
	defaultClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(a.conf.JWTExpirationHours)))
	defaultClaims.Subject = tokenClaims.Subject

	// Generate a new subject if the subject is empty.
	if defaultClaims.Subject == "" {
		subject, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("creating subject id")
			return "", contract.ErrInternalError
		}
		defaultClaims.Subject = subject.String()
	}

	// Set the audience if it is not empty.
	if len(tokenClaims.Audience) > 0 {
		defaultClaims.Audience = tokenClaims.Audience
	}

	// Generate a new consent id if the consent id is empty.
	if tokenClaims.ConsentID == nil {
		consentID, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("creating consent id")
			return "", contract.ErrInternalError
		}
		tokenClaims.ConsentID = &consentID
	}

	// Populate all gathered token claims.
	claims := &models.TokenClaims{
		DefaultClaims:  defaultClaims,
		ClientID:       tokenClaims.ClientID,
		Scopes:         tokenClaims.Scopes,
		Subjects:       tokenClaims.Subjects,
		ClientMetadata: tokenClaims.ClientMetadata,
		ConsentID:      tokenClaims.ConsentID.String(),
	}

	// Sign the token with the claims.
	return a.jwtClient.SignToken(claims)
}

// validateScope validates the scopes with the database and returns the valid scopes.
func (a *App) validateScopes(ctx context.Context, scopes []string) ([]*models.Scope, error) {
	matchingScopes, err := a.repo.ScopeDB.ListAllMatching(ctx, scopes)
	if err != nil {
		log.WithError(err).Error("matching scopes in the database")
		return nil, err
	}

	return matchingScopes, nil
}

func (a *App) parseScopesString(scope string) []string {
	scope = strings.TrimSpace(scope)

	// If no scope is provided.
	if scope == "" {
		return []string{}
	}

	scopes := strings.Split(scope, " ")
	scopes = sliceutil.RemoveDuplicates(scopes)

	return scopes
}

func (a *App) createConsentInUserWallet(ctx context.Context, id uuid.UUID, clientID, pseudo, name, desc, token string) error {
	pGranted := timestamppb.New(time.Now())
	if !pGranted.IsValid() {
		log.Error("unable to convert to proto timestamp")
		return contract.ErrInternalError
	}
	req := &walletPB.CreateConsentRequest{
		Id:          id.String(),
		Description: desc,
		Name:        name,
		GrantedAt:   pGranted,
		UserPseudo:  pseudo,
		ClientId:    clientID,
		AccessToken: token,
	}

	_, err := a.walletClient.CreateConsent(ctx, req)
	if err != nil {
		log.WithError(err).Error("unable to create consent")
		return contract.ErrInternalError
	}

	return nil
}

func accessModelExists(supplied, required, optional []*models.AccessModel) error {
	for _, s := range supplied {
		found := false
		for _, r := range required {
			if r.ID == s.ID {
				return errors.Wrapf(contract.ErrInvalidArguments, "payload should not contain id's of required_access_models these are accepted automatically")
			}
		}
		for _, o := range optional {
			if o.ID == s.ID {
				found = true

				break
			}
		}
		if !found {
			return errors.Wrapf(contract.ErrNotFound, "one ore more access_models from payload are not found in session's optional access_models")
		}
	}

	return nil
}

// sessionToResponse creates session_response message from model
func sessionToResponse(m *models.Session) *models.SessionResponse {
	s := models.SessionResponse{}

	s.ID = m.ID.String()
	s.State = m.State

	if m.Client != nil {
		c := models.Client{}
		c.ID = m.Client.ID
		c.Name = m.Client.Name
		c.Logo = m.Client.Logo
		c.Icon = m.Client.Icon
		c.Color = m.Client.Color
		s.Client = &c
	}

	if m.Audience != nil {
		a := models.Audience{}
		a.ID = m.Audience.ID
		a.Audience = m.Audience.Audience
		a.Namespace = m.Audience.Namespace
		s.Audience = &a
	}

	s.RequiredAccessModels = m.RequiredAccessModels
	s.OptionalAccessModels = m.OptionalAccessModels
	s.AcceptedAccessModels = m.AcceptedAccessModels

	return &s
}
