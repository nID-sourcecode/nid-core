// Package dial implements the dialer service
package dial

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
)

// Error definitions
var (
	ErrNotConnected            = errors.New("not connected")
	ErrInvalidFloatValue       = errors.New("unable to convert value")
	ErrUnkownConnectivityState = errors.New("unknown connectivity state")
)

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

// MethodService specifies the grpc service or method the config should be applied to
type MethodService struct {
	Service string `json:"service"`
}

// RetryPolicy specifies the policy to apply when trying to connect to a gRPC Service
type RetryPolicy struct {
	MaxAttempts          int              `json:"MaxAttempts"`
	InitialBackoff       BackOffInSeconds `json:"InitialBackoff"`
	MaxBackoff           BackOffInSeconds `json:"MaxBackoff"`
	BackoffMultiplier    float64          `json:"BackoffMultiplier"`
	RetryableStatusCodes []codes.Code     `json:"RetryableStatusCodes"`
}

// BackOffInSeconds is the backoff time in seconds as float32
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

	switch conn.GetState() {
	case connectivity.Ready:
		log.Warnf("Connected with %s, state: %s which should not be the initial state", service, conn.GetState().String())
		return conn, nil
	case connectivity.Shutdown, connectivity.Connecting, connectivity.TransientFailure:
		log.Warnf("Not connected with %s, state: %s", service, conn.GetState().String())
		return conn, ErrNotConnected
	case connectivity.Idle:
		return conn, nil
	default:
		return conn, errors.Wrapf(ErrUnkownConnectivityState, "%s for %s", conn.GetState().String(), service)
	}
}
