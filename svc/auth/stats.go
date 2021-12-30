package main

import (
	"github.com/prometheus/client_golang/prometheus"

	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/metrics"
)

// Stats contains the prometheus stats for the auth service
type Stats struct {
	tokenSwapped *prometheus.CounterVec
}

// CreateStats will initialise the prometheus stats
func CreateStats(scope metrics.Scope) *Stats {
	return &Stats{
		tokenSwapped: scope.RegisterNewCounterVector("token_swapped", "Counter for the tokens swapped per audience", []string{"audience"}),
	}
}
