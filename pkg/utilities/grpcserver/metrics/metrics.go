package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	//nolint:gomodguard // needed for backwards compatibility
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/promrus"
)

const (
	// PrometheusStatPort default port for exposing prometheus stats
	PrometheusStatPort int = 8080
)

// ExposeProm starts http server on port 8080 exposing metrics route
func ExposeProm() {
	// Create the Prometheus hook:
	hook := promrus.MustNewPrometheusHook()

	// Configure logrus to use the Prometheus hook:
	log.AddHook(hook)
	ExposePromWithPort(PrometheusStatPort)
}

// ExposePromWithPort starts http server on specified port exposing metrics route
func ExposePromWithPort(port int) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Infof("Started http server for prometheus metrics on port 8080")
		// nolint:gosec
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
	}()
}
