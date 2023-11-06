// Package environment provides the base config for all services
package environment

import (
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus" //nolint:gomodguard

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

// ConfigInterface needs to be implemented in Service Environments
type ConfigInterface interface {
	Init() error
	GetBaseConfig() *BaseConfig
}

// BaseConfig represents the base environment variables
type BaseConfig struct {
	Port        int    `envconfig:"default=8081,PORT"`
	LogLevel    string `envconfig:"default=info,LOG_LEVEL"`
	LogMode     bool   `envconfig:"default=false,LOG_MODE"`
	LogFormat   string `envconfig:"default=text,LOG_FORMAT"`
	Environment string `envconfig:"default=LOCAL,ENVIRONMENT"`
	PGHost      string `envconfig:"default=localhost,PG_HOST"`
	PGPort      int    `envconfig:"default=5432,PG_PORT"`
	PGUser      string `envconfig:"default=postgres,PG_USER"`
	PGPass      string `envconfig:"default=postgres,PG_PASS"`
	Namespace   string `envconfig:"NAMESPACE"`
}

// GetLogFormatter returns the logrus logformatter that should be used based on the environment
func (c *BaseConfig) GetLogFormatter() logrus.Formatter {
	if c.LogFormat == "logstash" {
		return logrustash.DefaultFormatter(logrus.Fields{})
	} else if c.LogFormat == "fluentd" {
		return joonix.NewFormatter()
	}

	return &logrus.TextFormatter{}
}

// GetLogLevel returns the logrus loglevel that should be used based on the environment
func (c *BaseConfig) GetLogLevel() logrus.Level {
	loglevel, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		log.WithError(err).WithField("loglevel", c.LogLevel).Error("Unable to parse loglevel, using info")

		return logrus.InfoLevel
	}

	return loglevel
}
