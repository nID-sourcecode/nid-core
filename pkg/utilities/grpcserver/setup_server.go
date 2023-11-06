// Package grpcserver implements the grpcserver
package grpcserver

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"

	//nolint:gomodguard //needed for backwards compatibility
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/logfields"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/loggrpc"
)

// ServiceRegistry service registry can be implemented to register custom
// services to the grpc server
type ServiceRegistry interface {
	RegisterServices(*grpc.Server)
}

// DefaultServiceRegistry default service registry implementation
type DefaultServiceRegistry struct{}

// RegisterServices by default no services are registered
func (d DefaultServiceRegistry) RegisterServices(_ *grpc.Server) {}

// InitDefault initialises the grpc server with default config
func InitDefault(serviceRegistry ServiceRegistry) error {
	conf := NewDefaultConfig()
	return InitWithConf(serviceRegistry, &conf)
}

// InitWithConf initialises the grpc server with custom interceptor
func InitWithConf(serviceRegistry ServiceRegistry, conf *Config) error {
	log.SetFormatter(conf.LogFormatter)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.WithError(err).WithField(logfields.Port, conf.Port).Error("failed to listen")
		return err
	}

	if conf.PubKey != nil {
		SetPubKey(conf.PubKey)
	}

	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	contextLogger := log.New()
	contextLogger.SetLevel(conf.LogLevel)
	logrusEntry := log.NewEntry(contextLogger)
	// Shared options for the logger, with a custom gRPC code to log level function.
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(loggrpc.LogDecider),
	}

	log.WithField(logfields.LogLevel, conf.LogLevel).Info("loglevel set")
	// Shared options for the recovery middleware.
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(conf.RecoveryHandlerFunc),
	}

	// Make sure that log statements internal to gRPC library are logged using the logrus Logger as well.
	grpclogger := log.New()
	grpclogger.SetLevel(log.WarnLevel)
	grpclogger.SetFormatter(conf.LogFormatter)
	grpc_logrus.ReplaceGrpcLogger(log.NewEntry(grpclogger))

	// Create the chain of interceptors. Evaluated from left to right
	interceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
		grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		unaryContextLogInterceptor,
		grpc_validator.UnaryServerInterceptor(),
	}
	// Chain interceptors, evaluated from left to right
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			append(interceptors, conf.AdditionalInterceptors...)...,
		)),

		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(conf.AdditionalStreamServerInterceptor...)))

	go HandleSigGracefulShutdown(grpcServer)
	healthpb.RegisterHealthServer(grpcServer, health.NewServer())
	// Register services to grpc server
	serviceRegistry.RegisterServices(grpcServer)
	reflection.Register(grpcServer)
	log.Infof("Starting grpc server at port: %d", conf.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.WithError(err).Error("failed to serve")
		return err
	}
	return nil
}
