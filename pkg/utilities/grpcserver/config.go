package grpcserver

import (
	"crypto/rsa"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// DefaultGRPCPort default gRPC port value
const DefaultGRPCPort int = 3550

// NewDefaultConfig creates new default gRPC server config
func NewDefaultConfig() Config {
	return Config{
		Port:                   DefaultGRPCPort,
		LogLevel:               log.WarnLevel,
		LogFormatter:           &log.TextFormatter{},
		AdditionalInterceptors: []grpc.UnaryServerInterceptor{},
		RecoveryHandlerFunc:    getRecoveryFunction(),
	}
}

// Config grpc server config
type Config struct {
	Port                   int
	LogLevel               log.Level
	LogFormatter           log.Formatter
	AdditionalInterceptors []grpc.UnaryServerInterceptor
	RecoveryHandlerFunc    grpc_recovery.RecoveryHandlerFuncContext
	PubKey                 *rsa.PublicKey
}
