// Package log is the default weave logger and wraps logrus
package log

import (
	"context"
	"io"
)

// nolint: gochecknoglobals
var loggerUtility LoggerUtility

// Level specifies the log levels used
type Level string

const (
	// FatalLevel to os.Exit(1) after logging
	FatalLevel Level = "fatal"
	// ErrorLevel to specify something failed
	ErrorLevel Level = "error"
	// WarningLevel to specify something needs attention
	WarningLevel Level = "warning"
	// InfoLevel to specify something important happened
	InfoLevel Level = "info"
	// DebugLevel to show useful debugging information
	DebugLevel Level = "debug"
	// TraceLevel to show low level traces
	TraceLevel Level = "trace"
)

func (level Level) isValid() error {
	switch level {
	case FatalLevel, ErrorLevel, WarningLevel, InfoLevel, DebugLevel, TraceLevel:
		return nil
	}
	return ErrInvalidLevel
}

// nolint: gochecknoinits
func init() {
	loggerUtility = newLogrusLogger()
}

// SetLoggerUtility sets LoggerUtility implementation that is used in the methods exported by this package
func SetLoggerUtility(utility LoggerUtility) {
	loggerUtility = utility
}

// GetLogger returns the logger
func GetLogger() LoggerUtility {
	return loggerUtility
}

// Fields is a map for setting multiple log fields at once
type Fields map[string]interface{}

// LoggerUtility is the interface for logger
type LoggerUtility interface {
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})

	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})

	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})

	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})

	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	WithField(key string, value interface{}) LoggerUtility
	WithContext(ctx context.Context) LoggerUtility
	With(key string, value interface{}) LoggerUtility

	WithFields(Fields) LoggerUtility
	WithError(err error) LoggerUtility

	SetLevel(Level) error
	SetFormatter(formatter interface{}) error
	SetFormat(format Format) error

	Extract(ctx context.Context) LoggerUtility

	SetOut(writer io.Writer)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	loggerUtility.Debug(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	loggerUtility.Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	loggerUtility.Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	loggerUtility.Info(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	loggerUtility.Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	loggerUtility.Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	loggerUtility.Warn(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	loggerUtility.Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	loggerUtility.Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	loggerUtility.Error(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	loggerUtility.Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	loggerUtility.Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	loggerUtility.Fatal(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	loggerUtility.Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	loggerUtility.Fatalf(format, args...)
}

// WithField logs a message with a field
func WithField(key string, value interface{}) LoggerUtility {
	return loggerUtility.WithField(key, value)
}

// WithContext logs a message with a context
func WithContext(ctx context.Context) LoggerUtility {
	return loggerUtility.WithContext(ctx)
}

// With adds a field to the logger.
func With(key string, value interface{}) LoggerUtility {
	return loggerUtility.With(key, value)
}

// WithFields logs a message with fields
func WithFields(fields Fields) LoggerUtility {
	return loggerUtility.WithFields(fields)
}

// WithError logs a message with an error
func WithError(err error) LoggerUtility {
	return loggerUtility.WithError(err)
}

// SetLevel sets the log level
func SetLevel(level Level) error {
	if err := level.isValid(); err != nil {
		return err
	}
	return loggerUtility.SetLevel(level)
}

// SetFormatter sets the formatter
//
// Deprecated: please use SetFormat for setting the log format
func SetFormatter(formatter interface{}) error {
	return loggerUtility.SetFormatter(formatter)
}

// SetFormat set the log formatter based on format enum
func SetFormat(format Format) error {
	if err := format.isValid(); err != nil {
		return err
	}
	return loggerUtility.SetFormat(format)
}

// Extract extract given context into logger utility
func Extract(ctx context.Context) LoggerUtility {
	return loggerUtility.Extract(ctx)
}

// SetOut sets the output for the logger
func SetOut(writer io.Writer) {
	loggerUtility.SetOut(writer)
}
