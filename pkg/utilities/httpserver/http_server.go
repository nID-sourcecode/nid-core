// Package httpserver provides the default HTTP server
// and a way to create custom http servers
package httpserver

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Options server options
type Options struct {
	SetLogFormatter func()
	LogMiddleware   func(c *gin.Context)
}

// DefaultServerOptions creates the default server options
func DefaultServerOptions() *Options {
	return &Options{
		SetLogFormatter: func() {
			log.SetFormatter(&log.TextFormatter{
				DisableColors: false,
				DisableQuote:  true,
				FullTimestamp: true,
			})
		},
		LogMiddleware: logrusLoggerMiddleware,
	}
}

// NewGinServer create a new gin server
func NewGinServer() *gin.Engine {
	return NewGinServerWithOpts(DefaultServerOptions())
}

// NewGinServerWithOpts create new gin server with specified options
func NewGinServerWithOpts(opts *Options) *gin.Engine {
	if opts == nil {
		return NewGinServer()
	}
	opts.SetLogFormatter()
	r := gin.New()
	r.Use(opts.LogMiddleware)
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

func logrusLoggerMiddleware(c *gin.Context) {
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
