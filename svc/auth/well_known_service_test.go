package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vrischmann/envconfig"

	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v3"
	pb "lab.weave.nl/nid/nid-core/svc/auth/proto"
)

type WellKnownServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	srv  *WellKnownServiceServer
	conf *AuthConfig
}

func (s *WellKnownServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
}

func (s *WellKnownServiceTestSuite) SetupSuite() {
	c := &AuthConfig{}
	if err := envconfig.InitWithOptions(c, envconfig.Options{AllOptional: true}); err != nil {
		s.Failf("init conf failed", "%+v", err)
	}
	s.conf = c
	s.conf.AuthorizationCodeExpirationTime = time.Minute
	s.conf.AuthorizationCodeLength = 32
	s.conf.AuthRequestURI = "https://authrequest.com"
	s.conf.ClusterHost = "test.uri.com/jwks.json"
	s.conf.Issuer = "someissuer"

	s.srv = &WellKnownServiceServer{
		wk:   nil,
		conf: s.conf,
	}

	// Init the JWT
	priv, pub, err := jwt.GenerateTestKeys()
	s.Require().NoError(err, "error generating test keys")
	s.conf.Namespace = "nid"
	s.conf.Issuer = "issuer.nl"
	opts := jwt.DefaultOpts()
	s.srv.jwtClient = jwt.NewJWTClientWithOpts(priv, pub, opts)

	err = s.srv.initWellKnown(pb.File_auth_proto)
	s.Require().NoError(err, "error initialising wellknown")
}

// Well_known_openpID
func (s *WellKnownServiceTestSuite) TestGetWellKnown() {
	s.srv.conf.JWKSURI = "test.uri.com/jwks.json"

	s.srv.wk = nil
	err := s.srv.initWellKnown(File_auth_test_proto)
	s.Require().NoError(err, "error initialising wellknown")
	authURI := fmt.Sprintf(AuthURITemplate, s.srv.conf.ClusterHost)

	wk, err := s.srv.GetWellKnownOpenIDConfiguration(s.Ctx, &pb.WellKnownRequest{})
	s.Require().NoError(err)
	s.Equal(s.srv.conf.Issuer, wk.Issuer)
	s.Equal([]string{s.srv.jwtClient.Opts.HeaderOpts.Alg}, wk.IdTokenSigningAlgValuesSupported)
	s.Equal([]string{CodeResponseType.String()}, wk.ResponseTypesSupported)
	s.Equal([]string{"authorization_code"}, wk.GrantTypesSupported)
	s.Equal([]string{"openid"}, wk.ScopesSupported)
	s.Equal(s.srv.conf.JWKSURI, wk.JwksUri)
	s.Equal(authURI+"/token", wk.TokenEndpoint)
	s.Equal(authURI+"/oppolicy", wk.OpPolicyUri)
	s.Equal(authURI+"/authorize", wk.AuthorizationEndpoint)
	s.Equal(authURI+"/servicedocs", wk.ServiceDocumentation)
	s.Equal(authURI+"/op_tos_uri", wk.OpTosUri)
	s.Equal(authURI+"/revoke", wk.RevocationEndpoint)
	s.Equal(authURI+"/userinfo", wk.UserinfoEndpoint)
	s.Equal(authURI+"/introspect", wk.IntrospectionEndpoint)
	s.Equal(authURI+"/checksess_iframe", wk.CheckSessionIframe)
}

func (s *WellKnownServiceTestSuite) TestGetWellKnownMissingValue() {
	s.srv.conf.JWKSURI = "iets.json"

	s.srv.wk = nil
	err := s.srv.initWellKnown(File_auth_test_2_proto)
	s.Require().NoError(err, "error initialising wellknown")
	wk, err := s.srv.GetWellKnownOpenIDConfiguration(s.Ctx, &pb.WellKnownRequest{})
	authURI := fmt.Sprintf(AuthURITemplate, s.srv.conf.ClusterHost)

	s.Require().NoError(err)
	s.Equal(s.srv.conf.JWKSURI, wk.JwksUri)
	s.Equal(authURI+"/token2", wk.TokenEndpoint)
	s.Equal(authURI+"/oppolicy2", wk.OpPolicyUri)
	s.Equal(authURI+"/authorize2", wk.AuthorizationEndpoint)
	s.Equal(authURI+"/servicedocs2", wk.ServiceDocumentation)
	s.Equal(authURI+"/op_tos_uri2", wk.OpTosUri)
	s.Equal(authURI+"/revoke2", wk.RevocationEndpoint)
	s.Equal(authURI+"/userinfo2", wk.UserinfoEndpoint)
	s.Equal("", wk.IntrospectionEndpoint)
	s.Equal(authURI+"/checksess_iframe2", wk.CheckSessionIframe)
}

func TestWellKnownServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WellKnownServiceTestSuite))
}
