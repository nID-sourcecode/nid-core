package main

import (
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/metrics"
)

// Stats contains the prometheus stats for the service
type Stats struct{}

// CreateStats will initialise the prometheus stats
func CreateStats(_ metrics.Scope) *Stats {
	return &Stats{}
}
