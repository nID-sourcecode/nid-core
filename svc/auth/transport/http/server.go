package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	"github.com/pkg/errors"
)

// Server handles the http server side for auth service
type Server struct {
	app  contract.App
	conf *HttpConfig
}

// New returns new instance of app struct
func New(app contract.App, conf *HttpConfig) *Server {
	return &Server{
		app:  app,
		conf: conf,
	}
}

// Run execute gin http server for auth service
func (s *Server) Run(port string) error {
	app := gin.Default()
	app.Use(
		gin.Recovery(),
	)

	app.GET("/authorize", s.authorize).
		POST("/authorize-headless", s.authorizeHeadless).
		POST("/claim", s.claim).
		POST("/accept", s.accept).
		POST("/reject", s.reject).
		POST("/generate-session-finalise-token", s.generateSessionFinaliseToken).
		POST("/details", s.getSessionDetails).
		POST("/status", s.status).
		POST("/auth.Auth/Status", s.status).
		POST("/finalise", s.finalise).
		GET("/token", s.token).
		POST("/token", s.tokenClientFlow).
		GET("/register", s.registerAccessModel).
		GET("/swap-token", s.swapToken).
		GET("/v1/health", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

	return app.Run(fmt.Sprintf(":%s", port))
}

func (s *Server) authorize(c *gin.Context) {
	authorizeRequest := &models.AuthorizeRequest{
		Scope:          c.Query("scope"),
		ResponseType:   c.Query("response_type"),
		ClientID:       c.Query("client_id"),
		RedirectURI:    c.Query("redirect_uri"),
		Audience:       c.Query("audience"),
		OptionalScopes: c.Query("optional_scopes"),
	}

	url, err := s.app.Authorize(c, authorizeRequest)
	if err != nil {
		_ = checkAndReturnGinError(err, c) //nolint:errcheck
		return
	}

	c.Status(http.StatusFound) // respond with 302: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/302
	c.Header("location", url)
}

func (s *Server) authorizeHeadless(c *gin.Context) {
	authorizeHeadlessRequest, err := parseBody[models.AuthorizeHeadlessRequest](c)
	if err != nil {
		return
	}
	err = s.app.AuthorizeHeadless(c, authorizeHeadlessRequest)

	_ = checkAndReturnGinError(err, c) //nolint:errcheck
}

func (s *Server) claim(c *gin.Context) {
	sessionRequest, err := parseBody[models.SessionRequest](c)
	if err != nil {
		log.WithError(err).Info("could not parse request to sessionrequest")
		c.AbortWithStatusJSON(401, contract.ErrInvalidArguments)
		return
	}

	sessionResponse, err := s.app.Claim(c, c.GetHeader("Claims"), sessionRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.JSON(200, sessionResponse)
}

func (s *Server) accept(c *gin.Context) {
	acceptRequest, err := parseBody[models.AcceptRequest](c)
	if err != nil {
		return
	}

	sessionResponse, err := s.app.Accept(c, c.GetHeader("Claims"), acceptRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.JSON(200, sessionResponse)
}

func (s *Server) reject(c *gin.Context) {
	sessionRequest, err := parseBody[models.SessionRequest](c)
	if err != nil {
		return
	}

	err = s.app.Reject(c, sessionRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.Status(200)
}

func (s *Server) generateSessionFinaliseToken(c *gin.Context) {
	sessionRequest, err := parseBody[models.SessionRequest](c)
	if err != nil {
		return
	}

	sessionAuthz, err := s.app.GenerateSessionFinaliseToken(c, sessionRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.JSON(200, sessionAuthz)
}

func (s *Server) getSessionDetails(c *gin.Context) {
	sessionRequest, err := parseBody[models.SessionRequest](c)
	if err != nil {
		return
	}

	sessionResponse, err := s.app.GetSessionDetails(c, sessionRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.JSON(200, sessionResponse)
}

func (s *Server) status(c *gin.Context) {
	sessionRequest, err := parseBody[models.SessionRequest](c)
	if err != nil {
		return
	}

	statusResponse, err := s.app.Status(c, sessionRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.JSON(200, statusResponse)
}

func (s *Server) finalise(c *gin.Context) {
	finaliseRequest, err := parseBody[models.FinaliseRequest](c)
	if err != nil {
		return
	}

	sessionAuthz, err := s.app.Finalise(c, finaliseRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.JSON(200, sessionAuthz)
}

func (s *Server) token(c *gin.Context) {
	tokenRequest := &models.TokenRequest{
		GrantType:         c.Query("grant_type"),
		AuthorizationCode: c.Query("authorization_code"),
		RefreshToken:      c.Query("refresh_token"),
	}

	if tokenRequest.GrantType == "" {
		_ = c.AbortWithError(401, errors.Wrapf(contract.ErrUnauthenticated, "grant_type is empty"))
		return
	}

	if tokenRequest.GrantType == "refresh_token" && tokenRequest.RefreshToken == "" {
		_ = c.AbortWithError(401, errors.Wrapf(contract.ErrUnauthenticated, "refresh token is empty for grant_type: refresh_token"))
		return
	}

	if tokenRequest.GrantType == "authorization_code" && tokenRequest.AuthorizationCode == "" {
		_ = c.AbortWithError(401, errors.Wrapf(contract.ErrUnauthenticated, "authorization code is empty for grant_type: authorization_code"))
		return
	}

	username, password, err := getBasicAuth(c)
	if err != nil {
		_ = c.AbortWithError(401, err)
		return
	}

	tokenRequest.Metadata = models.TokenRequestMetadata{
		Username:          username,
		Password:          password,
		CertificateHeader: c.GetHeader(s.conf.CertificateHeader),
	}

	tokenResponse, err := s.app.Token(c, tokenRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.JSON(200, tokenResponse)
}

func (s *Server) tokenClientFlow(c *gin.Context) {
	tokenClientFlowRequest, err := parseBody[models.TokenClientFlowRequest](c)
	if err != nil {
		return
	}

	username, password, err := getBasicAuth(c)
	if err != nil {
		_ = c.AbortWithError(401, err) //nolint:errcheck
		return
	}

	tokenClientFlowRequest.Metadata = models.TokenRequestMetadata{
		Username:          username,
		Password:          password,
		CertificateHeader: c.GetHeader(s.conf.CertificateHeader),
	}

	tokenResponse, err := s.app.TokenClientFlow(c, tokenClientFlowRequest)
	if checkAndReturnGinError(err, c) != nil {
		return
	}

	c.JSON(200, tokenResponse)
}

func (s *Server) registerAccessModel(c *gin.Context) {
	accessModelRequest := &models.AccessModelRequest{
		Audience:       c.Query("audience"),
		QueryModelJSON: c.Query("query_model_json"),
		ScopeName:      c.Query("scope_name"),
		Description:    c.Query("description"),
	}

	err := s.app.RegisterAccessModel(c, accessModelRequest)

	_ = checkAndReturnGinError(err, c)
}

func (s *Server) swapToken(c *gin.Context) {
	swapTokenRequest := &models.SwapTokenRequest{
		CurrentToken: c.Query("current_token"),
		Query:        c.Query("query"),
		Audience:     c.Query("audience"),
	}

	_, err := s.app.SwapToken(c, swapTokenRequest)

	_ = checkAndReturnGinError(err, c)
}

func parseBody[T any](c *gin.Context) (*T, error) {
	newBody := new(T)

	if err := c.ShouldBind(&newBody); err != nil {
		return nil, err
	}

	return newBody, nil
}

// checkAndReturnGinError handles the error. If error is nil, then nothing will happen and no error will be sent to the client.
func checkAndReturnGinError(err error, c *gin.Context) error {
	if err == nil {
		return nil
	}

	_ = c.Error(err) //nolint:errcheck
	switch {
	case errors.Is(err, contract.ErrInvalidArguments):
		_ = c.AbortWithError(400, contract.ErrInvalidArguments) //nolint:errcheck
		return err
	case errors.Is(err, contract.ErrInternalError):
		_ = c.AbortWithError(500, contract.ErrInternalError) //nolint:errcheck
		return err
	case errors.Is(err, contract.ErrNotFound):
		_ = c.AbortWithError(400, err)
		return err
	}

	c.AbortWithStatusJSON(500, "{ \"error\": \"error not handled\" }")
	return err
}
