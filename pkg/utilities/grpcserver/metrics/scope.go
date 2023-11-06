// Package metrics implements prometheus scopes
package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// Buckets holds all the buckets for prometheus Histogram
type Buckets []float64

// Scope interface for scope
type Scope interface {
	MustRegister(...prometheus.Collector)
	RegisterNewCounterVector(name, help string, labels []string) *prometheus.CounterVec
	RegisterNewGaugeVector(name, help string, labels []string) *prometheus.GaugeVec
	RegisterNewHistogramVector(buckets Buckets, name, help string, labels []string) *prometheus.HistogramVec
}

// PromScope scope implementation for prometheus
type PromScope struct {
	prefix     []string
	registerer prometheus.Registerer
}

// NewPromScope create a new prometheus scope
func NewPromScope(registerer prometheus.Registerer, scopes ...string) Scope {
	return &PromScope{
		prefix:     scopes,
		registerer: registerer,
	}
}

// MustRegister registers given prometheus collector on scope registerer
func (s *PromScope) MustRegister(collectors ...prometheus.Collector) {
	s.registerer.MustRegister(collectors...)
}

// RegisterNewCounterVector create and registers new counter vector
func (s *PromScope) RegisterNewCounterVector(name, help string, labels []string) *prometheus.CounterVec {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: s.statName(name),
			Help: help,
		},
		labels,
	)
	s.MustRegister(vec)
	return vec
}

// RegisterNewGaugeVector create and registers a new gauge vector
func (s *PromScope) RegisterNewGaugeVector(name, help string, labels []string) *prometheus.GaugeVec {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: s.statName(name),
			Help: help,
		},
		labels,
	)
	s.MustRegister(vec)
	return vec
}

// RegisterNewHistogramVector create and registers a new histogram
func (s *PromScope) RegisterNewHistogramVector(buckets Buckets, name, help string, labels []string) *prometheus.HistogramVec {
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    s.statName(name),
			Help:    help,
			Buckets: buckets,
		},
		labels,
	)
	s.MustRegister(vec)
	return vec
}

func (s *PromScope) statName(name string) string {
	if len(s.prefix) > 0 {
		return strings.Join(s.prefix, "_") + "_" + name
	}
	return name
}

// NopeScope scope mock
type NopeScope struct{}

// NewNopeScope get mocked scope
func NewNopeScope() Scope {
	return &NopeScope{}
}

// MustRegister registers given prometheus collector on scope registerer
func (s *NopeScope) MustRegister(_ ...prometheus.Collector) {
}

// RegisterNewCounterVector creates and registers new counter vector
func (s *NopeScope) RegisterNewCounterVector(name, help string, labels []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
}

// RegisterNewGaugeVector creates and registers a new gauge vector on the nopescope
func (s *NopeScope) RegisterNewGaugeVector(name, help string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, labels,
	)
}

// RegisterNewHistogramVector create and registers a new histogram
func (s *NopeScope) RegisterNewHistogramVector(_ Buckets, name, help string, labels []string) *prometheus.HistogramVec {
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name,
			Help: help,
		}, labels,
	)
	return vec
}
