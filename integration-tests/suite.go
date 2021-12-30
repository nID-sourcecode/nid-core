//go:build integration || to || files
// +build integration to files

// Package integration runs integration tests
package integration

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	gql "lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/dial"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	auth "lab.weave.nl/nid/nid-core/svc/auth/proto"
	dashboard "lab.weave.nl/nid/nid-core/svc/dashboard/proto"
	documentationPB "lab.weave.nl/nid/nid-core/svc/documentation/proto"
	infomanagerPB "lab.weave.nl/nid/nid-core/svc/info-manager/proto"
	pseudonymization "lab.weave.nl/nid/nid-core/svc/pseudonymization/proto"
	walletPB "lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
)

type config struct {
	BackendURL         string `envconfig:"default=nid.example.n-id.network,BACKEND_URL"`
	BackendPort        int    `envconfig:"default=443,BACKEND_PORT"`
	IsTLS              bool   `envconfig:"default=true,IS_TLS"`
	Namespace          string `envconfig:"NAMESPACE"`
	CleanupCallbackURL string `envconfig:"optional,CLEANUP_CALLBACK_URL"`
	Service            struct {
		Dashboard        string `envconfig:"default=dashboard,SERVICE_DASHBOARD"`
		Auth             string `envconfig:"default=auth,SERVICE_AUTH"`
		AuthGQL          string `envconfig:"default=auth-gql,SERVICE_AUTH_GQL"`
		WalletGQL        string `envconfig:"default=wallet-gql,SERVICE_WALLET_GQL"`
		Pseudonymization string `envconfig:"default=pseudonymization,SERVICE_PSEUDONYMIZATION"`
		Databron         string `envconfig:"default=databron,SERVICE_DATABRON"`
		Documentation    string `envconfig:"default=documentation,SERVICE_DOCUMENTATION"`
		Information      string `envconfig:"default=information,SERVICE_INFORMATION"`
		WalletRPC        string `envconfig:"default=wallet-rpc,SERVICE_WALLET_RPC"`
		Info             string `envconfig:"default=info,SERVICE_INFO"`
		InfoManager      string `envconfig:"default=info-manager,SERVICE_INFO_MANAGER"`
		InfoManagerGQL   string `envconfig:"default=info-manager-gql,SERVICE_INFO_MANAGER_GQL"`
	}
}

type clients struct {
	dashboardClient        dashboard.DashboardClient
	dashboardAuthClient    dashboard.AuthorizationServiceClient
	authClient             auth.AuthClient
	authHTTPClient         *resty.Client
	authGQLClient          gql.Client
	databronClient         gql.Client
	walletGQLClient        gql.Client
	pseudonymizationClient pseudonymization.PseudonymizerClient
	documentationClient    documentationPB.DocumentationClient
	walletAuthClient       walletPB.AuthorizationClient
	walletClient           walletPB.WalletClient
	infoClient             gql.Client
	infoManagerGQLClient   gql.Client
	infoManagerClient      infomanagerPB.InfoManagerClient
}

const (
	schemeHTTPS = "https"
	schemeHTTP  = "http"
)

// BaseTestSuite provides integrationTestSuite helper functionality
type BaseTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	ctx            context.Context
	envConfig      config
	callbackClient *http.Client
	clients
}

// TearDownSuite after the tests have run tear down
func (t *BaseTestSuite) TearDownSuite() {
	if t.envConfig.CleanupCallbackURL == "" {
		log.Warn("No cleanup callback given")
	} else {
		req, err := http.NewRequestWithContext(t.ctx, http.MethodGet, t.envConfig.CleanupCallbackURL, nil)
		t.Require().NoError(err)

		resp, err := t.callbackClient.Do(req)
		t.Require().NoError(err)
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.WithError(err).Error("unable to close response body")
			}
		}()

		log.Info("callback response: ", resp)
	}
}

// SetupSuite sets up an integration test suite
func (t *BaseTestSuite) SetupSuite() {
	t.Require().NoError(envconfig.Init(&t.envConfig), "unable to initialise environment config")
	t.ctx = context.Background()
	t.callbackClient = http.DefaultClient
	// Create dashboard client connection
	conn, err := t.GetGRPCClient(t.envConfig.Service.Dashboard)
	t.Require().NoError(err, "unable to start grpc service for dashboard:")
	t.clients.dashboardClient = dashboard.NewDashboardClient(conn)
	// Create dashboard Auth client connection
	conn, err = t.GetGRPCClient(t.envConfig.Service.Dashboard)
	t.Require().NoError(err, "unable to start auth grpc service for dashboard")
	t.clients.dashboardAuthClient = dashboard.NewAuthorizationServiceClient(conn)
	// Create auth client connection
	conn, err = t.GetGRPCClient(t.envConfig.Service.Auth)
	t.Require().NoError(err, "unable to start grpc service for auth")
	t.clients.authClient = auth.NewAuthClient(conn)
	// Create pseudonymization client connection
	conn, err = t.GetGRPCClient(t.envConfig.Service.Pseudonymization)
	t.Require().NoError(err, "unable to start grpc service for pseudonymization")
	t.clients.pseudonymizationClient = pseudonymization.NewPseudonymizerClient(conn)
	// Create documentation client connection
	conn, err = t.GetGRPCClient(t.envConfig.Service.Documentation)
	t.Require().NoError(err, "unable to start grpc service for documentation")
	t.clients.documentationClient = documentationPB.NewDocumentationClient(conn)

	// Create documentation client connection
	conn, err = t.GetGRPCClient(t.envConfig.Service.WalletRPC)
	t.Require().NoError(err, "unable to start grpc service for documentation")
	t.clients.walletAuthClient = walletPB.NewAuthorizationClient(conn)
	t.clients.walletClient = walletPB.NewWalletClient(conn)

	// Create auth gql client connection
	t.clients.authGQLClient = t.GetGQLClient(t.envConfig.Service.AuthGQL)
	// Create auth gql client connection
	t.clients.walletGQLClient = t.GetGQLClient(t.envConfig.Service.WalletGQL)
	// Create databron gql client connection
	t.clients.databronClient = t.GetGQLClient(t.envConfig.Service.Databron)

	t.clients.infoClient = t.GetGQLClient(t.envConfig.Service.Info)

	t.clients.infoManagerGQLClient = t.GetGQLClient(t.envConfig.Service.InfoManagerGQL)

	t.clients.authHTTPClient = resty.New().SetHostURL(t.httpURL(t.envConfig.Service.Auth, ""))
	t.clients.authHTTPClient.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // Do not follow redirects
	}))

	conn, err = t.GetGRPCClient(t.envConfig.Service.InfoManager)
	t.Require().NoError(err, "unable to start grpc service for infomanager")
	t.clients.infoManagerClient = infomanagerPB.NewInfoManagerClient(conn)
}

// GetGRPCClient get a grpc client for a service
func (t *BaseTestSuite) GetGRPCClient(service string) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s.%s:%d", service, t.envConfig.BackendURL, t.envConfig.BackendPort)

	var err error
	var connection *grpc.ClientConn
	if t.envConfig.IsTLS {
		// nolint:gosec
		connection, err = dial.Service(addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		// temp fix for http weave-cluster
		connection, err = dial.Service(addr, grpc.WithInsecure())
	}
	t.Require().NoErrorf(err, "unable to connect trough \"%s\" to service \"%s\"", addr, service)

	return connection, nil
}

// GetGQLClient get a GraphQL client for a service
func (t *BaseTestSuite) GetGQLClient(service string) gql.Client {
	scheme := schemeHTTPS
	if !t.envConfig.IsTLS {
		scheme = schemeHTTP
	}
	addr := fmt.Sprintf("%s://%s.%s:%d/gql", scheme, service, t.envConfig.BackendURL, t.envConfig.BackendPort)

	return gql.NewClient(addr)
}

func (t *BaseTestSuite) httpURL(service, endpoint string) string {
	endpoint = strings.TrimLeft(endpoint, "/")
	if t.envConfig.IsTLS {
		return fmt.Sprintf("https://%s.%s:443/%s", service, t.envConfig.BackendURL, endpoint)
	}
	return fmt.Sprintf("http://%s.%s:80/%s", service, t.envConfig.BackendURL, endpoint)
}

// GetCtx returns returns a test context
func (t *BaseTestSuite) GetCtx() context.Context {
	return t.ctx
}

// GetResponse parses a response from http response to given value
func (t *BaseTestSuite) GetResponse(response *http.Response, v interface{}) {
	if response.Status != strconv.Itoa(http.StatusOK) &&
		response.Status != "200 OK" {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Errorf("status code: %s", response.Status)
		if len(bodyBytes) == 0 {
			log.Errorf("empty body returned")
		} else {
			log.Errorf("response http request: %s", string(bodyBytes))
		}
		t.FailNow("no status ok retrieved")
	}
	t.Require().NoError(json.NewDecoder(response.Body).Decode(&v))
}
