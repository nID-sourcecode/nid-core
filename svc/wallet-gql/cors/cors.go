// Package cors provides different implementations of origin providerss
package cors

import (
	"net/http"
	"strconv"
	"strings"
)

// OriginProvider interface definition for an origin provider
type OriginProvider interface {
	GetOrigin(r *http.Request) *string
}

func originFromRequestHeader(r *http.Request) string {
	return r.Header.Get("Origin")
}

/**
 * Accept any origin absolutely.
 */
type reflectOrigin struct{}

// GetOrigin retrieves origin header for reflect origin
func (r reflectOrigin) GetOrigin(req *http.Request) *string {
	o := originFromRequestHeader(req)

	return &o
}

// ReflectOrigin Accept any origin. This will set the Access-Control-Allow-Origin to the same Origin as set in the request header.
func ReflectOrigin() OriginProvider {
	return reflectOrigin{}
}

// Options origin options
type Options struct {
	/**
	 * Configures the Access-Control-Allow-Credentials CORS header. Set to true to pass the header, otherwise it is
	 * omitted.
	 */
	AllowCredentials bool

	/**
	 * Configures the Access-Control-Allow-Methods CORS header. If not specified, defaults to ["GET", "HEAD", "PUT",
	 * "PATCH", "POST", "DELETE"].
	 */
	AllowedMethods []string

	/**
	 * Configures the Access-Control-Allow-Headers CORS header. If not specified, defaults to reflecting the headers
	 * specified in the request's Access-Control-Request-Headers header.
	 */
	AllowedHeaders []string

	/**
	 * Configures the Access-Control-Expose-Headers CORS header. If not specified, no custom headers are exposed.
	 */
	ExposedHeaders []string

	/**
	 * Configures the Access-Control-Allow-Credentials CORS header. Set to true to pass the header, otherwise it is
	 * omitted.
	 */
	Credentials bool

	/**
	 * Configures the Access-Control-Max-Age CORS header. Set to an integer to pass the header, otherwise it is omitted.
	 */
	MaxAge int64

	/**
	 * Provides a status code to use for successful OPTIONS requests, since some legacy browsers (IE11, various
	 * SmartTVs) choke on 204.
	 */
	OptionsSuccessStatus int

	/**
	 * Configures the Access-Control-Allow-Origin CORS header.
	 *
	 * Builtin origin handlers:
	 * - AnyOrigin()                    - Always returns "*"
	 * - Origin(origin, origin, origin) - Returns the request Origin if it was passed through Origin(...)
	 * - ReflectOrigin()                - Always returns the request Origin
	 */
	Origin OriginProvider

	/**
	 * Pass the CORS preflight response to the next handler.
	 */
	PreflightContinue bool

	/**
	 * Disable the CORS middleware.
	 */
	Disable bool
}

// WithOptions creates cors middleware for given options
func WithOptions(options *Options) func(http.Handler) http.Handler { // nolint: gocognit
	if options.AllowedMethods == nil {
		options.AllowedMethods = []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"}
	}

	if options.OptionsSuccessStatus <= 0 {
		options.OptionsSuccessStatus = http.StatusNoContent
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.Disable {
				next.ServeHTTP(w, r)

				return
			}
			if origin := options.Origin.GetOrigin(r); origin != nil {
				w.Header().Add("Access-Control-Allow-Origin", *origin)

				if *origin != "*" {
					w.Header().Add("Vary", "Origin")
				}
			}

			if options.ExposedHeaders != nil {
				if headers := strings.Join(options.ExposedHeaders, ","); len(headers) > 0 {
					w.Header().Add("Access-Control-Expose-Headers", headers)
				}
			}

			if options.Credentials {
				w.Header().Add("Access-Control-Allow-Credentials", "true")
			}

			// Check if preflight
			// nolint: nestif
			if r.Method == "OPTIONS" && len(r.Header.Get("Access-Control-Request-Method")) > 0 && len(r.Header.Get("Origin")) > 0 {
				// Set allowed methods
				w.Header().Add("Access-Control-Allow-Methods", strings.Join(options.AllowedMethods, ","))

				// Set allowed headers
				if options.AllowedHeaders != nil {
					if headers := strings.Join(options.AllowedHeaders, ","); len(headers) > 0 {
						w.Header().Add("Access-Control-Allow-Headers", headers)
					}
				} else {
					if headers := r.Header.Get("Access-Control-Request-Headers"); len(headers) > 0 {
						w.Header().Add("Access-Control-Allow-Headers", headers)
						w.Header().Add("Vary", "Access-Control-Request-Headers")
					}
				}

				// Set max age
				if options.MaxAge > 0 {
					w.Header().Add("Access-Control-Max-Age", strconv.FormatInt(options.MaxAge, 10))
				}

				// If PreflightContinue is set we run the next handler. Otherwise we return an empty body here.
				if !options.PreflightContinue {
					w.WriteHeader(options.OptionsSuccessStatus)
					w.Header().Add("Content-Length", "0")

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Cors creates cors middleware for given origin provider with options
func Cors(origin OriginProvider) func(http.Handler) http.Handler {
	return WithOptions(&Options{Origin: origin})
}
