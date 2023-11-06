// Package dial implements the dialer service
package dial

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
)

// ErrNotConnected not connected
var ErrNotConnected = errors.New("not connected")

const (
	defaultMaxAttempts       = 4
	defaultBackoff           = 0.01
	defaultBackoffMultiplier = 1.0
)

// MethodConfig specifies the configuration for dialling a service
// https://github.com/grpc/proposal/blob/master/A6-client-retries.md
type MethodConfig struct {
	Services     []MethodService `json:"name,omitempty"` // Leaving this empty will result in applying the config to the whole service
	WaitForReady bool            `json:"waitForReady"`
	RetryPolicy  *RetryPolicy    `json:"retryPolicy"`
}

// MethodService specifies the grpc service or method the config should be applied to.
type MethodService struct {
	Service string `json:"service"`
}

// RetryPolicy specifies the policy to apply when trying to connect to a gRPC Service.
type RetryPolicy struct {
	MaxAttempts          int              `json:"MaxAttempts"`
	InitialBackoff       BackOffInSeconds `json:"InitialBackoff"`
	MaxBackoff           BackOffInSeconds `json:"MaxBackoff"`
	BackoffMultiplier    float64          `json:"BackoffMultiplier"`
	RetryableStatusCodes []codes.Code     `json:"RetryableStatusCodes"`
}

// BackOffInSeconds is the backoff time in seconds as float32.
type BackOffInSeconds float32

// MarshalJSON will convert the long int value to a string with a 's' appended
func (b BackOffInSeconds) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%fs", b))
}

// NewDefaultMethodConfig will initiate the default method config
func NewDefaultMethodConfig() *MethodConfig {
	return &MethodConfig{
		WaitForReady: true,
		RetryPolicy: &RetryPolicy{
			MaxAttempts:       defaultMaxAttempts,
			InitialBackoff:    defaultBackoff,
			MaxBackoff:        defaultBackoff,
			BackoffMultiplier: defaultBackoffMultiplier,
			RetryableStatusCodes: []codes.Code{
				codes.Unavailable,
			},
		},
	}
}

// ServiceWithRetry dials a grpc service with given retryPolicy
func ServiceWithRetry(service string, methodConfig *MethodConfig, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	bMethodConfig, err := json.Marshal(methodConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse method config")
	}
	dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(string(bMethodConfig)))
	return Service(service, dialOptions...)
}

// Service dials a grpc service.
func Service(service string, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(service, dialOptions...)
	if err != nil {
		return conn, errors.Wrap(err, "unable to dial the service")
	}

	maxRetries := 100
	for conn.GetState() == connectivity.Connecting && maxRetries > 0 {
		log.Infof("Dialling %s ....", service)
		maxRetries--

		time.Sleep(time.Second)
	}

	state := conn.GetState()
	// https://grpc.github.io/grpc/core/md_doc_connectivity-semantics-and-api.html
	if state != connectivity.Ready && state != connectivity.Idle {
		return conn, errors.Wrapf(ErrNotConnected, "%s for %s", conn.GetState().String(), service)
	}

	return conn, nil
}
