package main

import (
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/metrics"
)

// Stats contains the prometheus stats for the service
type Stats struct{}

// CreateStats will initialise the prometheus stats
func CreateStats(scope metrics.Scope) *Stats {
	return &Stats{}
}
