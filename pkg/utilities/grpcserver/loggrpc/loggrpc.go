// Package loggrpc implements a logger for the grpc services
package loggrpc

import (
	"context"
	"io/ioutil"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

// LoggerFromContext returns the logger from the given context
//
// Deprecated: Please use lab.weave.nl/weave/utilities/log/v2.Extract insteads
func LoggerFromContext(ctx context.Context) *logrus.Entry {
	entry := ctxlogrus.Extract(ctx)
	// If we have a logger with these properties, it is a nullLogger, which logs nothing. We want to replace this logger with an actual logger.
	if entry.Logger.Out == ioutil.Discard && entry.Logger.Level == logrus.PanicLevel {
		return logrus.NewEntry(logrus.New())
	}
	return entry
}

// LogDecider determines whether the current call will be logged in the middleware
func LogDecider(fullMethodName string, err error) bool {
	// will not log gRPC calls if it was a call to healthcheck and no error was raised
	if err == nil && fullMethodName == "/grpc.health.v1.Health/Check" {
		return false
	}
	// Don't log calls that succeed
	if err == nil {
		return false
	}

	return true
}
