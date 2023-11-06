package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/headers"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	pb "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
)

// Server is the grpc server for auth service.
type Server struct {
	app            contract.App
	metadataHelper headers.MetadataHelper
	conf           *GrpcConfig
}

// New returns a new instance of GRPC server
func New(app contract.App, metadataHelper headers.MetadataHelper, conf *GrpcConfig) *Server {
	return &Server{
		app:            app,
		metadataHelper: metadataHelper,
		conf:           conf,
	}
}

// Accept implements pb.AuthServer interface
func (s *Server) Accept(ctx context.Context, req *pb.AcceptRequest) (*pb.SessionResponse, error) {
	acceptRequest := &models.AcceptRequest{
		SessionID:      req.SessionId,
		AccessModelIds: req.AccessModelIds,
	}

	authzHeader, err := s.metadataHelper.GetValFromCtx(ctx, "claims")
	if err != nil {
		return nil, errors.Wrapf(err, "getting authorization header from metadata")
	}

	sessionResponse, err := s.app.Accept(ctx, authzHeader, acceptRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "running accept method")
	}

	return sessionToResponse(sessionResponse), nil
}

// Authorize implements pb.AuthServer interface
func (s *Server) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*emptypb.Empty, error) {
	authorizeRequest := &models.AuthorizeRequest{
		Scope:          req.Scope,
		ResponseType:   req.ResponseType,
		ClientID:       req.ClientId,
		RedirectURI:    req.RedirectUri,
		Audience:       req.Audience,
		OptionalScopes: req.OptionalScopes,
	}

	location, err := s.app.Authorize(ctx, authorizeRequest)
	if err != nil {
		return &emptypb.Empty{}, errors.Wrapf(err, "authorizing the request")
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("grpc-statuscode", "302", "location", location)); err != nil {
		log.Extract(ctx).WithError(err).Error("setting headers")

		return nil, contract.ErrInternalError
	}

	return &emptypb.Empty{}, nil
}

// AuthorizeHeadless creates a session with granted state and calls the redirect URL with the generated authorization code
func (s *Server) AuthorizeHeadless(ctx context.Context, req *pb.AuthorizeHeadlessRequest) (*emptypb.Empty, error) {
	authorizeHeadlessReq := &models.AuthorizeHeadlessRequest{
		ResponseType:   req.ResponseType,
		ClientID:       req.ClientId,
		RedirectURI:    req.RedirectUri,
		Audience:       req.Audience,
		QueryModelJSON: req.QueryModelJson,
		QueryModelPath: req.QueryModelPath,
	}

	return &emptypb.Empty{}, s.app.AuthorizeHeadless(ctx, authorizeHeadlessReq)
}

func (s *Server) GenerateSessionFinaliseToken(ctx context.Context, req *pb.SessionRequest) (*pb.SessionAuthorization, error) {
	sessionRequest := &models.SessionRequest{
		SessionID: req.SessionId,
	}

	sessionAuthz, err := s.app.GenerateSessionFinaliseToken(ctx, sessionRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "generating session finalise token")
	}

	return &pb.SessionAuthorization{FinaliseToken: sessionAuthz.FinaliseToken}, nil
}

func (s *Server) Claim(ctx context.Context, req *pb.SessionRequest) (*pb.SessionResponse, error) {
	sessionRequest := &models.SessionRequest{
		SessionID: req.SessionId,
	}

	authzHeader, err := s.metadataHelper.GetValFromCtx(ctx, "claims")
	if err != nil {
		return nil, errors.Wrapf(err, "getting authorization header from metadata")
	}

	sessionResponse, err := s.app.Claim(ctx, authzHeader, sessionRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "running claim method")
	}

	return sessionToResponse(sessionResponse), nil
}

// Finalise updates the session state to code_granted and redirects to the redirect target with included authorization_code
func (s *Server) Finalise(ctx context.Context, req *pb.FinaliseRequest) (*pb.FinaliseResponse, error) {
	finaliseRequest := &models.FinaliseRequest{
		SessionID:            req.SessionId,
		SessionFinaliseToken: req.SessionFinaliseToken,
	}

	finalise, err := s.app.Finalise(ctx, finaliseRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "finalising the request")
	}

	return &pb.FinaliseResponse{RedirectLocation: finalise.RedirectLocation}, nil
}

// GetSessionDetails implements pb.AuthServer interface
func (s *Server) GetSessionDetails(ctx context.Context, req *pb.SessionRequest) (*pb.SessionResponse, error) {
	sessionRequest := &models.SessionRequest{
		SessionID: req.SessionId,
	}

	sessionResponse, err := s.app.GetSessionDetails(ctx, sessionRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "getting session details")
	}

	return sessionToResponse(sessionResponse), nil
}

// RegisterAccessModel registers the audiences accessModelGql1 if its validated by the gql schema
func (s *Server) RegisterAccessModel(ctx context.Context, req *pb.AccessModelRequest) (*emptypb.Empty, error) {
	accessModelReq := &models.AccessModelRequest{
		Audience:       req.Audience,
		QueryModelJSON: req.QueryModelJson,
		ScopeName:      req.ScopeName,
		Description:    req.Description,
	}

	return &emptypb.Empty{}, s.app.RegisterAccessModel(ctx, accessModelReq)
}

// Reject implements pb.AuthServer interface
func (s *Server) Reject(ctx context.Context, req *pb.SessionRequest) (*emptypb.Empty, error) {
	sessionRequest := &models.SessionRequest{
		SessionID: req.SessionId,
	}

	return &emptypb.Empty{}, s.app.Reject(ctx, sessionRequest)
}

// Status implements pb.AuthServer interface
func (s *Server) Status(ctx context.Context, req *pb.SessionRequest) (*pb.StatusResponse, error) {
	sessionRequest := &models.SessionRequest{
		SessionID: req.SessionId,
	}

	statusResponse, err := s.app.Status(ctx, sessionRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "getting status for session")
	}

	return &pb.StatusResponse{State: pb.SessionState(statusResponse.State)}, nil
}

// SwapToken swaps a token with access to a specific service to a token with access to another service based on the roles
func (s *Server) SwapToken(ctx context.Context, req *pb.SwapTokenRequest) (*pb.TokenResponse, error) {
	swapTokenReq := &models.SwapTokenRequest{
		CurrentToken: req.CurrentToken,
		Query:        req.Query,
		Audience:     req.Audience,
	}

	token, err := s.app.SwapToken(ctx, swapTokenReq)
	if err != nil {
		return nil, errors.Wrapf(err, "swapping token")
	}

	return &pb.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
	}, nil
}

// TokenClientFlow Client credentials flow to return an authorization code by client_id + client_secret
func (s *Server) TokenClientFlow(ctx context.Context, req *pb.TokenClientFlowRequest) (*pb.TokenResponse, error) {
	tokenClientFlowReq := &models.TokenClientFlowRequest{
		GrantType: req.GrantType,
		Scope:     req.Scope,
		Audience:  req.Audience,
	}

	username, password, err := s.metadataHelper.GetBasicAuth(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving basic auth")
	}

	md, err := s.metadataHelper.MetadataFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving accept header")
	}

	certificateHeader, err := s.metadataHelper.GetMetadataValue(md, s.conf.CertificateHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving certificate header")
	}

	tokenClientFlowReq.Metadata = models.TokenRequestMetadata{
		Username:          username,
		Password:          password,
		CertificateHeader: certificateHeader,
	}

	response, err := s.app.TokenClientFlow(ctx, tokenClientFlowReq)
	if err != nil {
		return nil, errors.Wrapf(err, "creating token for client flow")
	}

	return &pb.TokenResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		TokenType:    response.TokenType,
	}, nil
}

// Token implements pb.AuthServer interface
func (s *Server) Token(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	tokenRequest := &models.TokenRequest{
		GrantType:         req.GrantType,
		AuthorizationCode: req.GetAuthorizationCode(),
		RefreshToken:      req.GetRefreshToken(),
	}

	username, password, err := s.metadataHelper.GetBasicAuth(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving basic auth")
	}

	md, err := s.metadataHelper.MetadataFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving accept header")
	}

	certificateHeader, err := s.metadataHelper.GetMetadataValue(md, s.conf.CertificateHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving certificate header")
	}

	tokenRequest.Metadata = models.TokenRequestMetadata{
		Username:          username,
		Password:          password,
		CertificateHeader: certificateHeader,
	}

	tokenResponse, err := s.app.Token(ctx, tokenRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "creating token")
	}

	return &pb.TokenResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		TokenType:    tokenResponse.TokenType,
	}, nil
}
