// Package servicebase is the base of a grpc service
package servicebase

// Registry represent the base of a service registry
type Registry struct {
	LogMode bool
	Port    int
}

// Service represent the base of a service
type Service struct{}
