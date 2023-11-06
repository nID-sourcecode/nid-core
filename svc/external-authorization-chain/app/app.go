// Package app implements business logic for external-authorization-chain
package app

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"

	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/internal"

	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/contract"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/backoff"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/dial"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// App handles the logic for this service
type App struct {
	endPointsClients    map[string][]authv3.AuthorizationClient
	denyByDefaultConfig *internal.DenyByDefault
}

// New returns new instance of App
func New(config internal.AppConfig) (*App, error) {
	app := &App{
		denyByDefaultConfig: config.DenyByDefault,
	}

	err := app.connectAuthorizationClients(config.Endpoints)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// CheckEndpoints checks each endpoint if the given request is authorized.
func (a *App) CheckEndpoints(ctx context.Context, request *authv3.CheckRequest) error {
	host := request.GetAttributes().GetRequest().GetHttp().GetHost()

	authzClients := a.endPointsClients[host]
	if len(authzClients) == 0 && a.denyByDefaultConfig.Enabled {
		// if length is zero, then the host is not known for authz clients,
		// so the host has to be applied to the allow lists in the DenyByDefault variabele of App.
		if !a.isHostAllowed(host) {
			return errors.Wrapf(contract.ErrHostIsDeniedByDefault, " denied request for host: %s", host)
		}
	}

	for i := 0; i < len(authzClients); i++ {
		client := authzClients[i]

		log.Debugf("authz checked for: %s", host)
		checkResponse, err := client.Check(ctx, request)
		if err != nil {
			return errors.Wrap(err, "tried checking authorization client")
		}

		if checkResponse.Status.Code != int32(codes.OK) {
			return errors.Wrapf(contract.ErrUnauthorized, " for hostname: %s", host)
		}
	}

	return nil
}

// connectAuthorizationClients connects with grpc endpoints and maps with the config the services that are linked with the authorization clients.
func (a *App) connectAuthorizationClients(endpoints map[string][]string) error {
	a.endPointsClients = make(map[string][]authv3.AuthorizationClient)

	clients, err := a.createClientsFromEndpoints(endpoints)
	if err != nil {
		return err
	}

	a.addServicesToAuthzClients(endpoints, clients)

	return nil
}

// createClientsFromEndpoints will connect with grpc servers and create client out of this using the keys of the endpoints map.
func (a *App) createClientsFromEndpoints(endpoints map[string][]string) (map[string]authv3.AuthorizationClient, error) {
	clients := make(map[string]authv3.AuthorizationClient)

	for key := range endpoints {
		clientConn, err := dial.Service(
			key, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithConnectParams(
				grpc.ConnectParams{
					Backoff:           backoff.Config{},
					MinConnectTimeout: time.Second * 5,
				},
			),
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		)
		if err != nil {
			return nil, err
		}

		clients[key] = authv3.NewAuthorizationClient(clientConn)
	}

	return clients, nil
}

// addServicesToAuthzClients will map the destination (key) services to the authz ([]value) services and save it in App struct
func (a *App) addServicesToAuthzClients(endpoints map[string][]string, clients map[string]authv3.AuthorizationClient) {
	for key, client := range clients {
		destinationsForAuthzService := endpoints[key]

		for i := 0; i < len(destinationsForAuthzService); i++ {
			endpoint := destinationsForAuthzService[i]

			clientsSlice := a.endPointsClients[endpoint]

			clientsSlice = append(clientsSlice, client)
			a.endPointsClients[endpoint] = clientsSlice
		}
	}
}

func (a *App) isHostAllowed(host string) bool {
	log.Debugf("checking if host is allowed: %s", host)
	host, _, err := net.SplitHostPort(host)
	if err != nil {
		log.WithError(err).Info("could not split host and port")
		return false
	}

	for i := 0; i < len(a.denyByDefaultConfig.Allow); i++ {
		allowHost := a.denyByDefaultConfig.Allow[i]
		if allowHost == host {
			return true
		}
	}

	log.Debugf("Denied host: %s", host)
	log.Debug(fmt.Sprintf("%s", a.denyByDefaultConfig.Allow))
	return false
}
