// Package prometheus creates metrics for auth service
package stats

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/metrics"
)

// Stats contains the prometheus stats for the auth service
type Stats struct {
	TokenSwapped *prometheus.CounterVec
}

// CreateStats will initialise the prometheus stats
func CreateStats(scope metrics.Scope) *Stats {
	return &Stats{
		TokenSwapped: scope.RegisterNewCounterVector("token_swapped", "Counter for the tokens swapped per audience", []string{"audience"}),
	}
}
