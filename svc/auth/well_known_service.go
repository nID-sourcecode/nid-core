package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nID-sourcecode/nid-core/svc/auth/internal/config"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/jwt/v3"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	pb "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
)

// WellKnownServiceServer well_kown server
type WellKnownServiceServer struct {
	wk        *pb.WellKnownResponse
	conf      *config.AuthConfig
	jwtClient *jwt.Client
}

// ErrDefinitions
var (
	ErrWellKnownAlreadyInitialised error = errors.New("tried to initialise wellknown, but is already initialised")
	ErrUnableToListClaims          error = errors.New("unable to list claim keys")
)

const (
	// AuthURITemplate is used for prefixing the well known endpoints
	AuthURITemplate = "http://auth.%v/v1/auth"
)

// Load /.well-known/openid-configuration endpoints from auth proto
func (s *WellKnownServiceServer) initWellKnown(fileDescriptor protoreflect.FileDescriptor) error {
	if s.wk != nil {
		return ErrWellKnownAlreadyInitialised
	}

	var options []string

	// Get all services in auth proto
	services := fileDescriptor.Services()
	for s := 0; s < services.Len(); s++ {
		service := services.Get(s)

		// Get methods in service
		methods := service.Methods()
		for m := 0; m < methods.Len(); m++ {
			found := methods.Get(m).Options().(*descriptorpb.MethodOptions).String()

			// Save method if options are found
			if len(found) > 0 {
				options = append(options, found)
			}
		}
	}

	s.wk = s.createWellKnownResponse(options)

	// Get the claims
	tokenClaims := &models.TokenClaims{}
	claimsKeyList, err := tokenClaims.ListKeys()
	if err != nil {
		return ErrUnableToListClaims
	}
	s.wk.ClaimTypesSupported = claimsKeyList

	return nil
}

func (s *WellKnownServiceServer) createWellKnownResponse(options []string) *pb.WellKnownResponse {
	handlerRegex := regexp.MustCompile(`\[auth\.well_known_openid_handler\]:([A-Z_]*)\s+(.*)($|\s)`)
	extractor := regexp.MustCompile(
		`\[google\.api\.http\]:{.*(delete|get|patch|post|put):"(.*?)".*}($|\s)`,
	)

	wk := &pb.WellKnownResponse{
		JwksUri:                          s.conf.JWKSURI,
		Issuer:                           s.conf.Issuer,
		IdTokenSigningAlgValuesSupported: []string{s.jwtClient.Opts.HeaderOpts.Alg},
		ResponseTypesSupported:           []string{CodeResponseType.String()},
		ScopesSupported:                  []string{"openid"},             // Only openid is shown within the well known. Other scopes will be kept private due to relation with audience
		GrantTypesSupported:              []string{"authorization_code"}, // Currently the only supported grant type as validated in the auth.proto
	}

	authBaseURI := fmt.Sprintf(AuthURITemplate, s.conf.ClusterHost)
	for _, o := range options {
		if len(handlerRegex.FindStringIndex(o)) > 0 {
			match := handlerRegex.FindAllStringSubmatch(o, -1)[0][1]

			if len(extractor.FindStringIndex(o)) > 0 {
				url := extractor.FindAllStringSubmatch(o, -1)[0][2]

				switch match {
				case "AUTHORIZATION_ENDPOINT":
					wk.AuthorizationEndpoint = authBaseURI + url
				case "TOKEN_ENDPOINT":
					wk.TokenEndpoint = authBaseURI + url
				case "REGISTRATION_ENDPOINT":
					wk.RegistrationEndpoint = authBaseURI + url
				case "SERVICE_DOCUMENTATION":
					wk.ServiceDocumentation = authBaseURI + url
				case "OP_POLICY_URI":
					wk.OpPolicyUri = authBaseURI + url
				case "OP_TOS_URI":
					wk.OpTosUri = authBaseURI + url
				case "REVOCATION_ENDPOINT":
					wk.RevocationEndpoint = authBaseURI + url
				case "INTROSPECTION_ENDPOINT":
					wk.IntrospectionEndpoint = authBaseURI + url
				case "USERINFO_ENDPOINT":
					wk.UserinfoEndpoint = authBaseURI + url
				case "CHECK_SESSION_IFRAME":
					wk.CheckSessionIframe = authBaseURI + url
				default:
					log.Info("Wellknown annotation: ", match)
				}
			}
		}
	}

	return wk
}

// GetWellKnownOpenIDConfiguration implements pb.AuthServer interface
func (s *WellKnownServiceServer) GetWellKnownOpenIDConfiguration(context.Context, *pb.WellKnownRequest) (*pb.WellKnownResponse, error) {
	return s.wk, nil
}

// GetWellKnownOAuthAuthorizationServer implements pb.AuthServer interface
func (s *WellKnownServiceServer) GetWellKnownOAuthAuthorizationServer(ctx context.Context, req *pb.WellKnownRequest) (*pb.WellKnownResponse, error) {
	return s.GetWellKnownOpenIDConfiguration(ctx, req)
}
