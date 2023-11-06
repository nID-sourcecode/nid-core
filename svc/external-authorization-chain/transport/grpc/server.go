// Package grpc deals with the transport layer for the grpc server.
package grpc

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/nID-sourcecode/nid-core/pkg/interceptor/xrequestid"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/servicebase"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/contract"
	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/internal"
	extauthz "github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/transport/grpc/ext_authz"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// Server Certificate authorization service server
type Server struct{}

// ServerRegistry is an implementation of grpc service registry
type ServerRegistry struct {
	servicebase.Registry

	server Server
	app    contract.App
}

// RegisterServices registers the external processor server
func (r ServerRegistry) RegisterServices(grpcServer *grpc.Server) {
	authv3.RegisterAuthorizationServer(grpcServer, extauthz.New(r.app))
}

// New creates a new grpc transport layer
func New(conf internal.GRPCConfig, app contract.App) {
	registry := ServerRegistry{server: Server{}, app: app}

	grpcConfig := grpcserver.NewDefaultConfig()
	grpcConfig.Port = conf.Port
	grpcConfig.AdditionalInterceptors = []grpc.UnaryServerInterceptor{
		xrequestid.AddXRequestID,
		otelgrpc.UnaryServerInterceptor(),
	}
	grpcConfig.AdditionalStreamServerInterceptor = []grpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
	}

	err := grpcserver.InitWithConf(registry, &grpcConfig)
	if err != nil {
		log.WithError(err).Fatal("Error initialising grpc server")
	}
}
