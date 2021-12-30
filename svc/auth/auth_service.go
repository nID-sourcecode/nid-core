// Package authserviceserver provides all endpoint handlers for the auth service
package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"lab.weave.nl/nid/nid-core/pkg/authtoken"
	"lab.weave.nl/nid/nid-core/pkg/gqlutil"
	"lab.weave.nl/nid/nid-core/pkg/pseudonym"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	errgrpc "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v3"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	"lab.weave.nl/nid/nid-core/svc/auth/internal/callbackhandler"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
	pb "lab.weave.nl/nid/nid-core/svc/auth/proto"
	walletPB "lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
)

// AuthServiceServer auth server
type AuthServiceServer struct {
	db              *AuthDB
	stats           *Stats
	wk              *pb.WellKnownResponse
	pseudonymizer   pseudonym.Pseudonymizer
	jwtClient       *jwt.Client
	schemaFetcher   gqlutil.SchemaFetcher
	walletClient    walletPB.WalletClient
	conf            *AuthConfig
	metadataHelper  headers.MetadataHelper
	passwordManager password.IManager
	callbackhandler callbackhandler.CallbackHandler
}

// various error definitions
var (
	ErrSessionNotFound       = errors.New("session not found")
	ErrInvalidQueryModelType = errors.New("query model does not have a related specific query model")
)

const (
	accessModelIdentifierPartAmount = 2
	authCodeParamKey                = "authorization_code"
	accessModelGeneratedName        = "Generated"
)

// Accept implements pb.AuthServer interface
func (s *AuthServiceServer) Accept(ctx context.Context, req *pb.AcceptRequest) (*pb.SessionResponse, error) {
	subject, err := s.getSubjectFromJWT(ctx)
	if err != nil {
		return nil, err
	}

	session, err := s.db.SessionDB.GetSessionByIDAndSubject(ctx, preloadRequiredAndOptionalScopes, req.GetSessionId(), subject, s.conf.AuthorizationCodeExpirationTime)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound(ErrSessionNotFound)
		}
		return nil, errgrpc.ErrInternalServer()
	}

	if session.State != models.SessionStateClaimed {
		return nil, errgrpc.ErrFailedPrecondition("precondition failed session not claimed")
	}

	suppliedAccessModels := req.GetAccessModelIds()

	// Get access_models from payload ids
	accessModels, err := s.db.AccessModelDB.GetAccessModelsByIDs(suppliedAccessModels)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound("found %d but got %d access_model_ids", len(accessModels), len(suppliedAccessModels))
		}
		log.Extract(ctx).WithField("access_model_ids", suppliedAccessModels).WithError(err).Error("querying access_models failed")
		return nil, errgrpc.ErrInternalServer()
	}

	// Check if all found access_models are in sessions required or optional access_models
	err = accessModelExists(accessModels, session.RequiredAccessModels, session.OptionalAccessModels)
	if err != nil {
		return nil, err
	}

	// All supplied access models are valid save in accepted_access_models session association
	// Combine access_models with required access_models
	combinedAccessModels := append(accessModels, session.RequiredAccessModels...)
	err = s.db.SessionDB.UpdateAcceptedAccessModels(session, combinedAccessModels)
	if err != nil {
		log.Extract(ctx).WithField("access_model_ids", strings.Join(suppliedAccessModels, ",")).WithError(err).Error("updating accepted access models failed")
		return nil, errgrpc.ErrInternalServer()
	}

	// Update accepted session to state claimed
	err = s.db.SessionDB.UpdateSessionState(session, models.SessionStateAccepted)
	if err != nil {
		log.Extract(ctx).WithField("session_id", req.SessionId).WithError(err).Error("updating session state failed")
		return nil, errgrpc.ErrInternalServer()
	}

	session.State = models.SessionStateAccepted
	session.AcceptedAccessModels = combinedAccessModels

	return sessionToResponse(session), nil
}

// Authorize implements pb.AuthServer interface
func (s *AuthServiceServer) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*empty.Empty, error) {
	if strings.TrimSpace(req.ResponseType) != CodeResponseType.String() {
		return nil, errgrpc.ErrInvalidArgument("no response type other than \"code\" is supported")
	}

	session := &models.Session{
		State: models.SessionStateUnclaimed,
	}

	if err := s.fillSessionFromRequest(ctx, session, req); err != nil {
		return nil, err
	}

	requiredAccessModels, err := s.getRequiredAccessModels(ctx, req.GetScope(), session.Audience)
	if err != nil {
		return nil, err
	}
	session.RequiredAccessModels = requiredAccessModels

	optionalAccessModels, err := s.getOptionalAccessModels(ctx, req.GetOptionalScopes(), session.Audience)
	if err != nil {
		return nil, err
	}
	session.OptionalAccessModels = optionalAccessModels

	err = s.db.SessionDB.CreateSession(session)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("inserting session")
		return nil, errgrpc.ErrInternalServer()
	}

	location := fmt.Sprintf("%s#%s", s.conf.AuthRequestURI, session.ID.String())
	if err := grpc.SetHeader(ctx, metadata.Pairs("grpc-statuscode", "302", "location", location)); err != nil {
		log.Extract(ctx).WithError(err).Error("setting headers")

		return nil, errgrpc.ErrInternalServer()
	}

	return &empty.Empty{}, nil
}

// AuthorizeHeadless creates a session with granted state and calls the redirect URL with the generated authorization code
func (s *AuthServiceServer) AuthorizeHeadless(ctx context.Context, req *pb.AuthorizeHeadlessRequest) (*empty.Empty, error) {
	session := &models.Session{
		State: models.SessionStateCodeGranted,
	}

	if err := s.fillSessionFromRequest(ctx, session, req); err != nil {
		return nil, err
	}

	code, err := s.setSessionAuthenticationCode(session)
	if err != nil {
		return nil, err
	}

	accessModelHash := fmt.Sprintf("%x", sha256.Sum256([]byte(req.GetQueryModelJson())))
	accessModel, err := s.db.AccessModelDB.GetAccessModelByAudienceWithScope(accessModelGeneratedName, accessModelHash, session.Audience)
	if err != nil {
		// FIXME: Schema checker does not work with Vecozo's schema  (TODO: Create issue for this on lab.weave.nl)
		// if err := s.validateQueryModelForAudience(ctx, req.GetQueryModelJson(), req.GetAudience()); err != nil {
		// 	return nil, err
		// }

		if req.QueryModelPath == "" {
			req.QueryModelPath = "/gql"
		}
		accessModel = &models.AccessModel{
			AudienceID:  session.AudienceID,
			Description: "",
			Hash:        accessModelHash,
			GqlAccessModel: &models.GqlAccessModel{
				JSONModel: req.GetQueryModelJson(),
				Path:      req.QueryModelPath,
			},
			Name: accessModelGeneratedName,
			Type: models.AccessModelTypeGQL,
		}

		err = s.db.AccessModelDB.CreateAccessModel(accessModel)
		if err != nil {
			log.Extract(ctx).WithError(err).Error("error inserting access_model")
			return nil, errgrpc.ErrInternalServer()
		}
	}

	session.RequiredAccessModels = []*models.AccessModel{accessModel}
	session.AcceptedAccessModels = []*models.AccessModel{accessModel}

	err = s.db.SessionDB.CreateSession(session)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("inserting session")
		return nil, errgrpc.ErrInternalServer()
	}

	if err := s.callbackhandler.HandleCallback(ctx, session.RedirectTarget.RedirectTarget, code); err != nil {
		log.Extract(ctx).WithError(err).Error("error handling callback")
		return nil, errgrpc.ErrInternalServer()
	}

	return &empty.Empty{}, nil
}

func (s *AuthServiceServer) fillSessionFromRequest(ctx context.Context, session *models.Session, req pb.CreateSessionRequest) error {
	client, err := s.db.ClientDB.GetClientByID(req.GetClientId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errgrpc.ErrInvalidArgument("client with id \"%s\" does not exist", req.GetClientId())
		}

		log.Extract(ctx).WithError(err).WithField("id", req.GetClientId()).Error("getting client")
		return errgrpc.ErrInternalServer()
	}
	session.ClientID = client.ID

	audience, err := s.db.AudienceDB.GetAudienceByURI(req.GetAudience())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errgrpc.ErrInvalidArgument("audience \"%s\" does not exist", req.GetAudience())
		}

		log.Extract(ctx).WithError(err).WithField("audience", req.GetAudience()).Error("getting audience")
		return errgrpc.ErrInternalServer()
	}
	session.Audience = audience
	session.AudienceID = audience.ID

	redirectTarget, err := s.db.RedirectTargetDB.GetRedirectTarget(req.GetRedirectUri(), req.GetClientId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errgrpc.ErrInvalidArgument("client with id \"%s\" does not have redirect URI \"%s\"", req.GetClientId(), req.GetRedirectUri())
		}
		log.Extract(ctx).WithError(err).WithField("redirect_target", req.GetRedirectUri()).
			WithField("client_id", req.GetClientId()).Error("getting redirect target")
		return errgrpc.ErrInternalServer()
	}
	session.RedirectTarget = redirectTarget
	session.RedirectTargetID = redirectTarget.ID

	return nil
}

// GenerateSessionFinaliseToken Generates a password for the session, which can be reused to check if the session is valid.
// The session token expires if the session's state becomes claimed.
func (s *AuthServiceServer) GenerateSessionFinaliseToken(ctx context.Context, req *pb.SessionRequest) (*pb.SessionAuthorization, error) {
	session, err := s.db.SessionDB.GetSessionByID(ctx, noPreload, req.GetSessionId(), s.conf.AuthorizationCodeExpirationTime)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("could not find the session")

		return nil, errgrpc.ErrInternalServer()
	}

	if session.State != models.SessionStateUnclaimed {
		return nil, errgrpc.ErrFailedPrecondition("precondition failed session not accepted")
	}

	token, err := authtoken.NewToken(s.conf.AuthorizationCodeLength)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error signing token when swapping")

		return nil, errgrpc.ErrInternalServer()
	}

	hash, err := s.passwordManager.GenerateHash(token)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error generating hash from token")

		return nil, errgrpc.ErrInternalServer()
	}

	err = s.db.SessionDB.SetSessionFinaliseToken(session, hash)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("could not set the finalise token for the session")

		return nil, errgrpc.ErrInternalServer()
	}

	return &pb.SessionAuthorization{FinaliseToken: token}, nil
}

func (s *AuthServiceServer) getOptionalAccessModels(ctx context.Context, optionalScopes string, audience *models.Audience) ([]*models.AccessModel, error) {
	optionalAccessModels := make([]*models.AccessModel, 0)
	trimmedOptionalScopes := strings.Trim(optionalScopes, " ")
	if trimmedOptionalScopes != "" {
		for _, specifiedScope := range strings.Split(trimmedOptionalScopes, " ") {
			scope := strings.Trim(specifiedScope, " ")
			parts := strings.Split(scope, "@")
			if len(parts) != accessModelIdentifierPartAmount {
				return nil, errgrpc.ErrInvalidArgument("scope \"%s\" is of invalid format, should be \"name@hash\"", scope)
			}
			name := parts[0]
			hash := parts[1]
			accessModel, err := s.db.AccessModelDB.GetAccessModelByAudienceWithScope(name, hash, audience)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errgrpc.ErrInvalidArgument("audience \"%s\" does not support access model \"%s\"", audience.Audience, scope)
				}

				log.Extract(ctx).WithError(err).Errorf("getting access model %s", scope)
				return nil, errgrpc.ErrInternalServer()
			}
			optionalAccessModels = append(optionalAccessModels, accessModel)
		}
	}

	return optionalAccessModels, nil
}

func (s *AuthServiceServer) getRequiredAccessModels(ctx context.Context, scopes string, audience *models.Audience) ([]*models.AccessModel, error) {
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
			return nil, errgrpc.ErrInvalidArgument("scope \"%s\" is of invalid format, should be \"name@hash\"", scope)
		}
		name := parts[0]
		hash := parts[1]
		accessModel, err := s.db.AccessModelDB.GetAccessModelByAudienceWithScope(name, hash, audience)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errgrpc.ErrInvalidArgument("audience \"%s\" does not support access model \"%s\"", audience.Audience, scope)
			}

			log.Extract(ctx).WithError(err).Errorf("getting access model %s", scope)
			return nil, errgrpc.ErrInternalServer()
		}
		requiredAccessModels = append(requiredAccessModels, accessModel)
	}

	if !openIDScopeSpecified {
		return nil, errgrpc.ErrInvalidArgument("the \"openid\" scope must be specified")
	}
	if len(requiredAccessModels) == 0 {
		return nil, errgrpc.ErrInvalidArgument("at least one access model-scope must be specified")
	}

	return requiredAccessModels, nil
}

// Claim implements pb.AuthServer interface
func (s *AuthServiceServer) Claim(ctx context.Context, req *pb.SessionRequest) (*pb.SessionResponse, error) {
	subject, err := s.getSubjectFromJWT(ctx)
	if err != nil {
		return nil, err
	}

	if subject == "" {
		return nil, errgrpc.ErrFailedPrecondition("subject cannot be empty")
	}

	session, err := s.db.SessionDB.GetSessionByID(ctx, preloadRequiredAndOptionalScopes, req.GetSessionId(), s.conf.AuthorizationCodeExpirationTime)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound(ErrSessionNotFound)
		}
		return nil, errgrpc.ErrInternalServer()
	}

	if session.State != models.SessionStateUnclaimed {
		return nil, errgrpc.ErrFailedPrecondition("precondition failed session not unclaimed")
	}

	// Update accepted session to state claimed
	err = s.db.SessionDB.UpdateSessionState(session, models.SessionStateClaimed)
	if err != nil {
		log.Extract(ctx).WithField("session_id", req.SessionId).WithError(err).Error("updating session state failed")
		return nil, errgrpc.ErrInternalServer()
	}

	// Update accepted session subject
	err = s.db.SessionDB.UpdateSessionSubject(session, subject)
	if err != nil {
		log.Extract(ctx).WithField("session_id", req.SessionId).WithError(err).Error("updating session subject failed")
		return nil, errgrpc.ErrInternalServer()
	}
	session.State = models.SessionStateClaimed

	return sessionToResponse(session), nil
}

type subClaims struct {
	Subject string `json:"sub"`
}

func (s *AuthServiceServer) getSubjectFromJWT(ctx context.Context) (string, error) {
	md, err := s.metadataHelper.MetadataFromCtx(ctx)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("getting metadata")
		return "", errgrpc.ErrInternalServer()
	}

	b64claims, err := s.metadataHelper.GetMetadataValue(md, "claims")
	if err != nil {
		log.Extract(ctx).WithError(err).Error("getting claims header")
		return "", errgrpc.ErrInternalServer()
	}

	claims := subClaims{}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(b64claims)
	if err != nil {
		return "", errgrpc.ErrInvalidArgument("base64 decoding claims")
	}

	err = json.Unmarshal(claimsJSON, &claims)
	if err != nil {
		return "", errgrpc.ErrInvalidArgument("parsing claims")
	}

	return claims.Subject, nil
}

// Finalise updates the session state to code_granted and redirects to the redirect target with included authorization_code
func (s *AuthServiceServer) Finalise(ctx context.Context, req *pb.FinaliseRequest) (*pb.FinaliseResponse, error) {
	session, err := s.db.SessionDB.GetSessionByID(ctx, preloadRequiredAndOptionalScopes, req.GetSessionId(), s.conf.AuthorizationCodeExpirationTime)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound(ErrSessionNotFound)
		}
		return nil, errgrpc.ErrInternalServer()
	}

	logger := log.Extract(ctx).WithField("session_id", session.ID)

	// Check if session is accepted
	if session.State != models.SessionStateAccepted {
		return nil, errgrpc.ErrFailedPrecondition("precondition failed session not accepted")
	}

	// Check if the given password matches with the session
	if req.GetSessionFinaliseToken() == "" {
		return nil, errgrpc.ErrFailedPrecondition("did not provide a pass for session")
	}

	ok, err := s.passwordManager.ComparePassword(req.GetSessionFinaliseToken(), session.FinaliseToken)
	if err != nil {
		return nil, errgrpc.ErrInternalServer()
	}
	if !ok {
		return nil, errgrpc.ErrFailedPrecondition("Given Finalise Token did not match")
	}

	code, err := s.setSessionAuthenticationCode(session)
	if err != nil {
		return nil, errgrpc.ErrInternalServer()
	}

	// Update accepted session
	err = s.db.SessionDB.UpdateSessionAuthorizationCode(session, *session.AuthorizationCode)
	if err != nil {
		logger.WithError(err).Error("updating session status failed")
		return nil, errgrpc.ErrInternalServer()
	}
	err = s.db.SessionDB.UpdateSessionState(session, models.SessionStateCodeGranted)
	if err != nil {
		logger.WithError(err).Error("updating session status failed")
		return nil, errgrpc.ErrInternalServer()
	}

	// Redirect to redirect location
	location := fmt.Sprintf("%s/?%s=%s", session.RedirectTarget.RedirectTarget, authCodeParamKey, code)

	return &pb.FinaliseResponse{RedirectLocation: location}, nil
}

func (s *AuthServiceServer) setSessionAuthenticationCode(session *models.Session) (string, error) {
	// Generate random authorization code for session
	code, err := authtoken.NewToken(s.conf.AuthorizationCodeLength)
	if err != nil {
		log.WithError(err).Error("generating authorization code failed")

		return "", errgrpc.ErrInternalServer()
	}

	hash, err := authtoken.Hash(code)
	if err != nil {
		log.WithError(err).Error("generating hash for authorization code failed")

		return "", errgrpc.ErrInternalServer()
	}

	session.AuthorizationCode = &hash

	return code, nil
}

// GetSessionDetails implements pb.AuthServer interface
func (s *AuthServiceServer) GetSessionDetails(ctx context.Context, req *pb.SessionRequest) (*pb.SessionResponse, error) {
	session, err := s.db.SessionDB.GetSessionByID(ctx, preloadAll, req.GetSessionId(), s.conf.AuthorizationCodeExpirationTime)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound(ErrSessionNotFound)
		}
		return nil, errgrpc.ErrInternalServer()
	}

	if session.State == models.SessionStateUnclaimed || session.State == models.SessionStateClaimed {
		return sessionToResponse(session), nil
	}

	return nil, errgrpc.ErrFailedPrecondition("session has wrong state")
}

// RegisterAccessModel registers the audiences accessModelGql1 if its validated by the gql schema
func (s *AuthServiceServer) RegisterAccessModel(ctx context.Context, req *pb.AccessModelRequest) (*empty.Empty, error) {
	if req.GetScopeName() == "" {
		return nil, errgrpc.ErrInvalidArgument("no scope name specified")
	}

	if err := s.validateQueryModelForAudience(ctx, req.GetQueryModelJson(), req.GetAudience()); err != nil {
		return nil, err
	}

	// Get Audience and insert AccessModel
	audience, err := s.db.AudienceDB.GetAudienceByURI(req.GetAudience())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound("audience \"%s\" does not exist", req.GetAudience())
		}

		log.Extract(ctx).WithError(err).WithField("audience", req.GetAudience()).Error("getting audience")
		return nil, errgrpc.ErrInternalServer()
	}

	accessModel := &models.AccessModel{
		AudienceID:  audience.ID,
		Description: req.GetDescription(),
		Hash:        fmt.Sprintf("%x", sha256.Sum256([]byte(req.GetQueryModelJson()))),
		JSONModel:   req.GetQueryModelJson(),
		Name:        req.GetScopeName(),
	}

	err = s.db.AccessModelDB.CreateAccessModel(accessModel)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error inserting access_model")
		return nil, errgrpc.ErrInternalServer()
	}

	return &empty.Empty{}, nil
}

// Ensure that gql query adheres schema defined by service
func (s *AuthServiceServer) validateQueryModelForAudience(ctx context.Context, queryModelJSON string, audience string) error {
	schema, err := s.schemaFetcher.FetchSchema(ctx, audience)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error fetching graphql schema")

		return errgrpc.ErrInternalServer()
	}

	var queryModel map[string]*gqlutil.AccessModel
	if err := json.Unmarshal([]byte(queryModelJSON), &queryModel); err != nil {
		return errgrpc.ErrInvalidArgument("query model json incorrectly formated")
	}
	if err := schema.ValidateQueryModel(queryModel); err != nil {
		return errgrpc.ErrInvalidArgument(strings.Title(err.Error()))
	}

	return nil
}

// Reject implements pb.AuthServer interface
func (s *AuthServiceServer) Reject(ctx context.Context, req *pb.SessionRequest) (*empty.Empty, error) {
	session, err := s.db.SessionDB.GetSessionByID(ctx, noPreload, req.GetSessionId(), s.conf.AuthorizationCodeExpirationTime)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound(ErrSessionNotFound)
		}
		return nil, errgrpc.ErrInternalServer()
	}

	// Check if session is claimed
	if session.State != models.SessionStateClaimed {
		return nil, errgrpc.ErrFailedPrecondition("precondition failed session not unclaimed")
	}

	// Update accepted session to state rejected
	err = s.db.SessionDB.UpdateSessionState(session, models.SessionStateRejected)
	if err != nil {
		log.Extract(ctx).WithField("session_id", session.ID).WithError(err).Error("updating session state failed")
		return nil, errgrpc.ErrInternalServer()
	}

	return &empty.Empty{}, nil
}

// Status implements pb.AuthServer interface
func (s *AuthServiceServer) Status(ctx context.Context, req *pb.SessionRequest) (*pb.StatusResponse, error) {
	session, err := s.db.SessionDB.GetSessionByID(ctx, noPreload, req.GetSessionId(), s.conf.AuthorizationCodeExpirationTime)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound(ErrSessionNotFound)
		}
		return nil, errgrpc.ErrInternalServer()
	}

	return &pb.StatusResponse{
		State: getState(session.State),
	}, nil
}

// SwapToken swaps a token with access to a specific service to a token with access to another service based on the roles
func (s *AuthServiceServer) SwapToken(ctx context.Context, in *pb.SwapTokenRequest) (*pb.TokenResponse, error) {
	inTokenClaims := &TokenClaims{}
	err := s.jwtClient.ValidateAndParseClaims(in.GetCurrentToken(), inTokenClaims)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to parse token")

		return nil, status.Errorf(codes.InvalidArgument, "unable to parse token")
	}

	// FIXME we should probably add some checks instead of just flikkering the input in the token https://lab.weave.nl/twi/core/-/issues/49
	audience, err := s.db.AudienceDB.GetAudienceByURI(in.GetAudience())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrInvalidArgument("audience \"%s\" does not exist", in.GetAudience())
		}

		log.Extract(ctx).WithError(err).Error("error getting audience")
		return nil, errgrpc.ErrInternalServer()
	}
	accessModels, err := s.db.AccessModelDB.GetAccessModelsByAudience(true, audience)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errgrpc.ErrNotFound("error getting access_models")
		}

		log.Extract(ctx).WithError(err).Error("error getting access_models")

		return nil, errgrpc.ErrInternalServer()
	}

	// Set Queries
	scopes := make(map[string]interface{})
	err = addAccessModelsToScopes(scopes, accessModels)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error adding accepted access models to queries")

		return nil, errgrpc.ErrInternalServer()
	}

	defaultClaims := jwt.NewDefaultClaims()
	defaultClaims.Issuer = s.conf.Issuer
	defaultClaims.Audience = []string{in.GetAudience()}
	defaultClaims.Subject = inTokenClaims.Subject

	outTokenClaims := &TokenClaims{
		DefaultClaims:  defaultClaims,
		ClientID:       inTokenClaims.ClientID,
		Scopes:         scopes,
		Subjects:       inTokenClaims.Subjects,
		ClientMetadata: inTokenClaims.ClientMetadata,
	}

	newToken, err := s.jwtClient.SignToken(outTokenClaims)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error signing token when swapping")

		return nil, errgrpc.ErrInternalServer()
	}

	s.stats.tokenSwapped.With(prometheus.Labels{"audience": in.GetAudience()}).Inc()

	tokenParts := strings.Split(newToken, ".")
	// sanity check. a JWT should always have 3 parts
	// nolint: gomnd
	if len(tokenParts) == 3 {
		log.Extract(ctx).WithField("token", tokenParts[0]+"."+tokenParts[1]).Info("token swapped")
	}

	return &pb.TokenResponse{AccessToken: newToken, TokenType: tokenType}, nil
}

const (
	tokenType = "Bearer"
)

// Token implements pb.AuthServer interface
func (s *AuthServiceServer) Token(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	clientID, err := s.authenticateClient(ctx)
	if err != nil {
		return nil, err
	}

	code := req.GetAuthorizationCode()
	logger := log.Extract(ctx).WithFields(log.Fields{
		"authorization_code": code,
		"client_id":          clientID,
	})

	hash, err := authtoken.Hash(code)
	if err != nil {
		logger.WithError(err).Error("generating hash for authorization code failed")
		return nil, errgrpc.ErrInternalServer()
	}
	session, err := s.db.SessionDB.GetSessionByCodeAndClientID(ctx, preloadAll, hash, clientID, s.conf.AuthorizationCodeExpirationTime)
	if err != nil {
		if errors.Is(err, models.ErrUnableToRetrieveTokenExpiration) {
			return nil, errgrpc.ErrDeadlineExceeded(err)
		}
		return nil, errgrpc.ErrNotFound(ErrSessionNotFound)
	}

	if session.State != models.SessionStateCodeGranted {
		// Do not give away information about this session -- since code is not known to be equal to ID
		return nil, errgrpc.ErrFailedPrecondition("error getting token, session in incorrect state")
	}

	// Set Scopes
	scopes := make(map[string]interface{})
	err = addAccessModelsToScopes(scopes, session.AcceptedAccessModels)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error adding accepted access models to scopes")

		return nil, errgrpc.ErrInternalServer()
	}
	err = addAccessModelsToScopes(scopes, session.RequiredAccessModels)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("error adding required access models to scopes")

		return nil, errgrpc.ErrInternalServer()
	}

	// Set Subjects
	nidSubject := ""
	audienceSubject := ""
	if session.Subject != "" {
		nidSubject, err = s.pseudonymizer.GetPseudonym(ctx, session.Subject, s.conf.Namespace)
		if err != nil {
			log.Extract(ctx).WithError(err).Error("error getting pseudonym")

			return nil, errgrpc.ErrInternalServer()
		}

		audienceSubject, err = s.pseudonymizer.GetPseudonym(ctx, session.Subject, session.Audience.Namespace)
		if err != nil {
			log.Extract(ctx).WithError(err).Error("error getting pseudonym")

			return nil, errgrpc.ErrInternalServer()
		}
	} else if s.scopeContainsString(scopes, "$$nid:subject$$", "$$nid:bsn$$") {
		log.Extract(ctx).Error("scope requires subject")
		return nil, errgrpc.ErrInternalServer()
	}

	token, err := s.createToken(ctx, session, nidSubject, scopes, audienceSubject)
	if err != nil {
		return nil, err
	}

	tokenParts := strings.Split(token, ".")
	// sanity check. a JWT should always have 3 parts
	// nolint: gomnd
	if len(tokenParts) == 3 {
		log.Extract(ctx).WithField("token", tokenParts[0]+"."+tokenParts[1]).Info("token created")
	}

	return &pb.TokenResponse{AccessToken: token, TokenType: tokenType}, nil
}

func (s *AuthServiceServer) scopeContainsString(scopes map[string]interface{}, subs ...string) bool {
	for _, scope := range scopes {
		for _, sub := range subs {
			if strings.Contains(fmt.Sprint(scope), sub) {
				return true
			}
		}
	}

	return false
}

func (s *AuthServiceServer) createToken(ctx context.Context, session *models.Session, nidSubject string, scopes map[string]interface{}, audienceSubject string) (string, error) {
	consentID, err := uuid.NewV4()
	if err != nil {
		log.WithError(err).Error("creating consent id")
		return "", errgrpc.ErrInternalServer()
	}

	defaultClaims := jwt.NewDefaultClaims()
	defaultClaims.Audience = []string{session.Audience.Audience}
	defaultClaims.Subject = nidSubject
	defaultClaims.Issuer = s.conf.Issuer

	clientMetadataJSON := session.Client.Metadata.RawMessage
	clientMetadata := make(map[string]interface{})
	if clientMetadataJSON != nil && len(clientMetadataJSON) > 0 {
		err = json.Unmarshal(clientMetadataJSON, &clientMetadata)
		if err != nil {
			log.Extract(ctx).WithError(err).Error("marshalling client metadata")
			return "", errgrpc.ErrInternalServer()
		}
	}

	claims := &TokenClaims{
		DefaultClaims: defaultClaims,
		ClientID:      session.Client.ID.String(),
		Scopes:        scopes,
		Subjects: map[string]interface{}{
			session.Audience.Namespace: audienceSubject,
		},
		ClientMetadata: clientMetadata,
		ConsentID:      consentID.String(),
	}
	if err != nil {
		log.Extract(ctx).WithError(err).Error("converting token claims to map")
		return "", errgrpc.ErrInternalServer()
	}

	token, err := s.jwtClient.SignToken(claims)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("signing token")
		return "", errgrpc.ErrInternalServer()
	}

	scopeNames := make([]string, len(scopes))
	var i int
	for q := range scopes {
		scopeNames[i] = strings.Split(q, "@")[0]
		i++
	}

	// Call the wallet service to create consent
	if nidSubject != "" {
		err = s.createConsentInUserWallet(ctx, consentID, session.Client.ID.String(), session.Subject, "name", strings.Join(scopeNames, ", "), token)
		if err != nil {
			return "", err
		}
	}
	// Update session state to state token granted
	err = s.db.SessionDB.UpdateSessionState(session, models.SessionStateTokenGranted)
	if err != nil {
		log.Extract(ctx).WithField("session_id", session.ID).WithError(err).Error("updating session state failed")
		return "", errgrpc.ErrInternalServer()
	}

	return token, nil
}

func (s *AuthServiceServer) createConsentInUserWallet(ctx context.Context, id uuid.UUID, clientID, pseudo, name, desc, token string) error {
	pGranted := timestamppb.New(time.Now())
	if !pGranted.IsValid() {
		log.Error("unable to convert to proto timestamp")
		return errgrpc.ErrInternalServer()
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

	_, err := s.walletClient.CreateConsent(ctx, req)
	if err != nil {
		log.WithError(err).Error("unable to create create consent")
		return errgrpc.ErrInternalServer()
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

func (s *AuthServiceServer) authenticateClient(ctx context.Context) (string, error) {
	username, password, err := s.metadataHelper.GetBasicAuth(ctx)
	if err != nil {
		return "", errgrpc.ErrUnauthenticated("retrieving basic auth")
	}

	client, err := s.db.ClientDB.GetClientByID(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errgrpc.ErrUnauthenticated("invalid username")
		}
		log.Extract(ctx).WithField("id", username).WithError(err).Error("unable to get client")
		return "", errgrpc.ErrInternalServer()
	}

	passwordMatches, err := s.passwordManager.ComparePassword(password, client.Password)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to compare password")

		return "", errgrpc.ErrInternalServer()
	}
	if !passwordMatches {
		return "", errgrpc.ErrUnauthenticated("incorrect password")
	}

	return username, nil
}

func addAccessModelsToScopes(scopes map[string]interface{}, accessModels []*models.AccessModel) error {
	for _, accessModel := range accessModels {
		switch accessModel.Type {
		case models.AccessModelTypeGQL:
			if accessModel.GqlAccessModel == nil {
				return errors.Wrapf(ErrInvalidQueryModelType, "GQL type but no GQL access model related (%s)", accessModel.ID.String())
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
				return errors.Wrapf(ErrInvalidQueryModelType, "REST type but no REST access model related (%s)", accessModel.ID.String())
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
			return errors.Wrap(ErrInvalidQueryModelType, "unable to add access model to scopes")
		}
	}

	return nil
}
