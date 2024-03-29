// Package httpserver provides the default HTTP server
// with logging middleware and health check endpoint
// and a way to create custom http servers
package httpserver

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

// Options server Options
type Options struct {
	LoggerUtility    log.LoggerUtility
	UseLogMiddleware bool
}

// DefaultServerOptions creates the default server Options
// nolint: golint
func DefaultServerOptions() *Options {
	return &Options{
		LoggerUtility:    log.GetLogger(),
		UseLogMiddleware: true,
	}
}

// NewGinServer create a new gin server
func NewGinServer() *gin.Engine {
	return NewGinServerWithOpts(DefaultServerOptions())
}

// NewGinServerWithOpts create new gin server with specified Options
func NewGinServerWithOpts(opts *Options) *gin.Engine {
	if opts == nil {
		return NewGinServer()
	}

	if opts.LoggerUtility == nil {
		opts.LoggerUtility = log.GetLogger()
	}

	r := gin.New()
	if opts.UseLogMiddleware {
		r.Use(loggerMiddleware)
	}
	r.GET("/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func loggerMiddleware(c *gin.Context) {
	if c.Request.RequestURI == "/v1/health" {
		return
	}
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	t := time.Now()

	c.Next()
	statusCode := c.Writer.Status()
	if statusCode >= http.StatusBadRequest {
		log.WithFields(log.Fields{
			"body":   blw.body.String(),
			"status": statusCode,
		}).Errorf("unable to complete request")
	}

	if time.Since(t) > time.Second {
		log.Warnf("slow request detected: %s took: %s", c.Request.RequestURI, time.Since(t))
	}
}
