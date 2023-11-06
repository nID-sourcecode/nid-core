package log

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus" //nolint:gomodguard

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	joonix "github.com/joonix/log"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

type logrusLoggerUtility struct {
	entry *logrus.Entry
}

func newLogrusLogger() LoggerUtility {
	return &logrusLoggerUtility{entry: logrus.NewEntry(logrus.New())}
}

// CustomLogrusUtility creates logrus log utility for custom entry
func CustomLogrusUtility(entry *logrus.Entry) LoggerUtility {
	return &logrusLoggerUtility{entry: entry}
}

// Debug logs a message at level Debug on the standard logger.
func (l *logrusLoggerUtility) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Debugln logs a line message at level Debug on the standard logger.
func (l *logrusLoggerUtility) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (l *logrusLoggerUtility) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func (l *logrusLoggerUtility) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Infoln logs a message at level Info on the standard logger.
func (l *logrusLoggerUtility) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func (l *logrusLoggerUtility) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l *logrusLoggerUtility) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func (l *logrusLoggerUtility) Warnln(args ...interface{}) {
	l.entry.Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func (l *logrusLoggerUtility) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func (l *logrusLoggerUtility) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Errorln logs a message at level Error on the standard logger.
func (l *logrusLoggerUtility) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (l *logrusLoggerUtility) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l *logrusLoggerUtility) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func (l *logrusLoggerUtility) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (l *logrusLoggerUtility) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// WithField takes into consideration the field
func (l *logrusLoggerUtility) WithField(key string, value interface{}) LoggerUtility {
	return &logrusLoggerUtility{l.entry.WithField(key, value)}
}

// WithContext takes into consideration the field
func (l *logrusLoggerUtility) WithContext(ctx context.Context) LoggerUtility {
	return &logrusLoggerUtility{l.entry.WithContext(ctx)}
}

// With takes into consideration the field
func (l *logrusLoggerUtility) With(key string, value interface{}) LoggerUtility {
	return l.WithField(key, value)
}

// WithFields takes into consideration the fields
func (l *logrusLoggerUtility) WithFields(fields Fields) LoggerUtility {
	temp := logrus.Fields{}
	for k, v := range fields {
		temp[k] = v
	}

	return &logrusLoggerUtility{l.entry.WithFields(temp)}
}

// WithError takes into consideration the error
func (l *logrusLoggerUtility) WithError(err error) LoggerUtility {
	return &logrusLoggerUtility{l.entry.WithError(err)}
}

// SetLevel sets the log level
func (l *logrusLoggerUtility) SetLevel(level Level) error {
	lvl, err := logrus.ParseLevel(string(level))
	if err != nil {
		return errors.Wrap(err, "failed to parse level")
	}

	l.entry.Logger.Level = lvl
	return nil
}

// SetFormatter sets the formatter
func (l *logrusLoggerUtility) SetFormatter(formatter interface{}) error {
	switch t := formatter.(type) {
	case *logrus.TextFormatter:
		l.entry.Logger.SetFormatter(t)
		return nil
	case *logrus.JSONFormatter:
		l.entry.Logger.SetFormatter(t)
		return nil
	case logrustash.LogstashFormatter:
		l.entry.Logger.SetFormatter(t)
		return nil
	case *joonix.Formatter:
		l.entry.Logger.SetFormatter(t)
		return nil
	}

	return errors.Wrap(ErrIncorrectFormatter, fmt.Sprintf("unknown format type: %T", formatter))
}

// SetFormat will set the formatter by specified format enum
func (l *logrusLoggerUtility) SetFormat(format Format) error {
	formatMapper := map[Format]logrus.Formatter{
		TextFormat:     &logrus.TextFormatter{},
		JSONFormat:     &logrus.JSONFormatter{},
		LogstashFormat: logrustash.DefaultFormatter(logrus.Fields{}),
		JoonixFormat:   joonix.NewFormatter(),
	}
	f, ok := formatMapper[format]
	if !ok {
		return ErrIncorrectFormatter
	}

	l.entry.Logger.SetFormatter(f)
	return nil
}

// Extract extracts context into the logger
func (l *logrusLoggerUtility) Extract(ctx context.Context) LoggerUtility {
	entry := ctxlogrus.Extract(ctx)
	// If we have a logger with these properties, it is a nullLogger, which logs nothing. We want to replace this logger with an actual logger.
	if entry.Logger.Out == ioutil.Discard && entry.Logger.Level == logrus.PanicLevel {
		return newLogrusLogger()
	}
	return &logrusLoggerUtility{entry: entry}
}

// SetOut sets the output for the logger
func (l *logrusLoggerUtility) SetOut(writer io.Writer) {
	l.entry.Logger.Out = writer
}
