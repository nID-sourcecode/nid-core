package integration

import (
	"context"
	"testing"

	"gopkg.in/yaml.v3"

	externalauthorization "github.com/nID-sourcecode/nid-core/svc/nid-filter/transport/external_authorization"

	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/internal"

	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/app"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/nID-sourcecode/nid-core/pkg/environment"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/dial"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/servicebase"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/contract"
	nidFilterContract "github.com/nID-sourcecode/nid-core/svc/nid-filter/contract"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const appConfigYaml = `
    deny_by_default:
      enabled: true
      allow:
        - luarunner.test.dev.domain.nl
        - auth.test.dev.domain.nl
        - nid.test.dev.domain.nl

    endpoints:
      localhost:3550:
        - luarunner.test.dev.domain.nl
        - auth.test.dev.domain.nl
        - nid.test.dev.domain.nl
`

type AppTestSuite struct {
	suite.Suite

	app contract.App

	serverConfig *environment.BaseConfig
	client       authv3.AuthorizationClient
}

func (a *AppTestSuite) TestCheckOK() {
	response, err := a.client.Check(
		context.Background(), returnCheckRequest(
			&authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{headerReturnCode: string(rune(OK))},
				},
			},
		),
	)

	a.NoError(err)
	a.NotNil(response.Status)
	a.Equal(0, int(response.Status.Code))
}

func (a *AppTestSuite) TestCheckFail() {
	response, err := a.client.Check(
		context.Background(), returnCheckRequest(
			&authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{headerReturnCode: string(rune(Fail))},
				},
			},
		),
	)

	a.NoError(err)
	a.NotNil(response.Status)
	a.Equal(16, int(response.Status.Code))
}

func (a *AppTestSuite) TestCheckEndpoints() {
	checkRequest := returnCheckRequest(
		&authv3.AttributeContext_Request{
			Http: &authv3.AttributeContext_HttpRequest{
				Host:    "luarunner.test.dev.domain.nl",
				Headers: map[string]string{headerReturnCode: string(rune(OK))},
			},
		},
	)

	err := a.app.CheckEndpoints(context.Background(), checkRequest)
	a.NoError(err)
}

func (a *AppTestSuite) TestCheckEndpoints_Fail() {
	checkRequest := returnCheckRequest(
		&authv3.AttributeContext_Request{
			Http: &authv3.AttributeContext_HttpRequest{
				Host:    "luarunner.test.dev.domain.nl",
				Headers: map[string]string{headerReturnCode: string(rune(Fail))},
			},
		},
	)

	err := a.app.CheckEndpoints(context.Background(), checkRequest)
	a.Error(err)
	a.ErrorIs(err, contract.ErrUnauthorized)
}

func (a *AppTestSuite) SetupSuite() {
	config := &environment.BaseConfig{
		Port:        3550,
		LogLevel:    "info",
		LogMode:     true,
		LogFormat:   "text",
		Environment: "LOCAL",
		PGHost:      "localhost",
		PGPort:      5432,
		PGUser:      "postgres",
		PGPass:      "postgres",
		Namespace:   "nid",
	}

	setupNIDFilter(config)
	a.serverConfig = config

	client, err := dial.Service("localhost:3550", grpc.WithTransportCredentials(insecure.NewCredentials()))
	a.Require().NoError(err)

	a.client = authv3.NewAuthorizationClient(client)

	var appConfig internal.AppConfig
	err = yaml.Unmarshal([]byte(appConfigYaml), &appConfig)
	a.Require().NoError(err)

	appHandler, err := app.New(
		internal.AppConfig{
			Endpoints: appConfig.Endpoints, DenyByDefault: appConfig.DenyByDefault,
		},
	)
	a.Require().NoError(err)

	a.app = appHandler
}

func (a *AppTestSuite) Test_AppConfig() {
	allowedHostYaml := `
    deny_by_default:
      enabled: true
      allow: 
        - someother.url.dev.domain.nl
    endpoints:
      localhost:3550:
        - luarunner.test.dev.domain.nl
        - auth.test.dev.domain.nl
        - proxy-server.test.dev.domain.nl
`

	var appConfig internal.AppConfig
	err := yaml.Unmarshal([]byte(allowedHostYaml), &appConfig)
	a.Require().NoError(err)

	hosts, ok := appConfig.Endpoints["localhost:3550"]
	a.True(ok)
	a.Equal([]string{"luarunner.test.dev.domain.nl", "auth.test.dev.domain.nl", "proxy-server.test.dev.domain.nl"},
		hosts)
	a.True(appConfig.DenyByDefault.Enabled)
	a.Equal("someother.url.dev.domain.nl", appConfig.DenyByDefault.Allow[0])
}

func (a *AppTestSuite) Test_HostDeny() {
	checkRequest := returnCheckRequest(&authv3.AttributeContext_Request{
		Http: &authv3.AttributeContext_HttpRequest{
			Host: "denied.host.dev.domain.nl",
		},
	})

	err := a.app.CheckEndpoints(context.Background(), checkRequest)
	a.Error(err)
	a.ErrorContains(err, contract.ErrHostIsDeniedByDefault.Error())
}

func (a *AppTestSuite) Test_HostAccept_DenyByDefault_Disabled() {
	allowedHostYaml := `
    deny_by_default:
      enabled: false
      allow: []
    endpoints:
      localhost:3550:
        - luarunner.test.dev.domain.nl
        - auth.test.dev.domain.nl
        - proxy-server.test.dev.domain.nl
`

	var appConfig internal.AppConfig
	err := yaml.Unmarshal([]byte(allowedHostYaml), &appConfig)
	a.Require().NoError(err)

	newApp, err := app.New(appConfig)
	a.Require().NoError(err)

	checkRequest := returnCheckRequest(&authv3.AttributeContext_Request{
		Http: &authv3.AttributeContext_HttpRequest{
			Host: "random.host.dev.domain.nl",
		},
	})

	err = newApp.CheckEndpoints(context.TODO(), checkRequest)
	a.NoError(err)
}

func (a *AppTestSuite) Test_HostAccept_DenyByDefault_Enabled() {
	allowedHostYaml := `
    deny_by_default:
      enabled: true
      allow: []
    endpoints:
      localhost:3550:
        - luarunner.test.dev.domain.nl
        - auth.test.dev.domain.nl
        - proxy-server.test.dev.domain.nl
`

	var appConfig internal.AppConfig
	err := yaml.Unmarshal([]byte(allowedHostYaml), &appConfig)
	a.Require().NoError(err)

	newApp, err := app.New(appConfig)
	a.Require().NoError(err)

	checkRequest := returnCheckRequest(&authv3.AttributeContext_Request{
		Http: &authv3.AttributeContext_HttpRequest{
			Host: "random.host.dev.domain.nl",
		},
	})

	err = newApp.CheckEndpoints(context.TODO(), checkRequest)
	a.Error(err)
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, &AppTestSuite{})
}

func returnCheckRequest(request *authv3.AttributeContext_Request) *authv3.CheckRequest {
	return &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Source:            nil,
			Destination:       nil,
			Request:           request,
			ContextExtensions: nil,
			MetadataContext:   nil,
		},
	}
}

func setupNIDFilter(conf *environment.BaseConfig) {
	registry := initialise()

	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = conf.Port
	grpcConfig.LogLevel = conf.GetLogLevel()
	grpcConfig.LogFormatter = conf.GetLogFormatter()
	go func() {
		err := grpcserver.InitWithConf(registry, &grpcConfig)
		if err != nil {
			log.WithError(err).Fatal("Error initialising grpc server")
		}
	}()
}

func initialise() *NIDFilterServiceRegistry {
	authorizationRules := []nidFilterContract.AuthorizationRule{
		NewAuthzImpl(),
	}

	registry := &NIDFilterServiceRegistry{
		authorizationRules: authorizationRules,
	}

	return registry
}

// NIDFilterServiceRegistry is an implementation of grpc service registry
type NIDFilterServiceRegistry struct {
	servicebase.Registry

	authorizationRules []nidFilterContract.AuthorizationRule
}

// RegisterServices registers the external processor server
func (r *NIDFilterServiceRegistry) RegisterServices(grpcServer *grpc.Server) {
	authv3.RegisterAuthorizationServer(grpcServer, externalauthorization.New(r.authorizationRules))
}
