// Package grpctesthelpers provides functionality to test grpc services
package grpctesthelpers

import (
	"fmt"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

// Error definitions
var (
	ErrFailHeader error = fmt.Errorf("failing header")
)

// ServerTransportStreamMock is a mock for the grpc.ServerTransportStreamMock interface
type ServerTransportStreamMock struct {
	mock.Mock
	Headers        map[string][]string
	FailingHeaders []string
}

// Method is a mock for the Method function
func (m *ServerTransportStreamMock) Method() string {
	args := m.Called()
	return args.String(0)
}

// SetHeader is a mock for the SetHeader function
func (m *ServerTransportStreamMock) SetHeader(md metadata.MD) error {
	for k, v := range md {
		if m.isFailingHeader(k) {
			return ErrFailHeader
		}
		m.Headers[k] = v
	}
	args := m.Called(md)
	return args.Error(0)
}

// SendHeader is a mock for the SendHeader function
func (m *ServerTransportStreamMock) SendHeader(md metadata.MD) error {
	args := m.Called(md)
	return args.Error(0)
}

// SetTrailer is a mock for the SetTrailer function
func (m *ServerTransportStreamMock) SetTrailer(md metadata.MD) error {
	args := m.Called(md)
	return args.Error(0)
}

// ResetHeaders reset current header map
func (m *ServerTransportStreamMock) ResetHeaders() {
	m.Headers = make(map[string][]string)
}

// AddFailingHeaders add failing headers to list of failing headers
func (m *ServerTransportStreamMock) AddFailingHeaders(header []string) {
	m.FailingHeaders = append(m.FailingHeaders, header...)
}

// ResetFailingHeaders reset failing header list
func (m *ServerTransportStreamMock) ResetFailingHeaders() {
	m.FailingHeaders = []string{}
}

func (m *ServerTransportStreamMock) isFailingHeader(header string) bool {
	for _, failingHeader := range m.FailingHeaders {
		if failingHeader == header {
			return true
		}
	}
	return false
}
