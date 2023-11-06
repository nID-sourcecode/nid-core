// Package spiffeparser parses spiffe certificate URIs
package spiffeparser

import (
	"regexp"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

// IstioSpiffeCert contains the data from an istio spiffe cert
type IstioSpiffeCert struct {
	By  *IstioSpiffeParts
	URI *IstioSpiffeParts
}

// IstioSpiffeParts contains the individual parts of an identifier in an istio spiffe cert
type IstioSpiffeParts struct {
	domain         string
	namespace      string
	serviceAccount string
}

// GetDomain returns the domain
func (p *IstioSpiffeParts) GetDomain() string {
	return p.domain
}

// GetNamespace returns the namespace
func (p *IstioSpiffeParts) GetNamespace() string {
	return p.namespace
}

// GetServiceAccount returns the service account
func (p *IstioSpiffeParts) GetServiceAccount() string {
	return p.serviceAccount
}

// SpiffeParser parses istio spiffe certs
type SpiffeParser interface {
	Parse(cert string) *IstioSpiffeCert
}

// NewDefaultSpiffeParser creates a new default spiffe parser
func NewDefaultSpiffeParser() *DefaultSpiffeParser {
	expr, err := regexp.Compile("([A-Za-z]+)=spiffe://([a-zA-Z.0-9-]+)/ns/([a-zA-Z0-9-]+)/sa/([a-zA-Z0-9-]+)")
	if err != nil {
		panic(err)
	}
	return &DefaultSpiffeParser{expr: expr}
}

// DefaultSpiffeParser is the default spiffe parser implementation
type DefaultSpiffeParser struct {
	expr *regexp.Regexp
}

// Error definitions
var (
	ErrInvalidNumberOfMatchSegments = errors.New("invalid number of match segments, should be 5")
	ErrMissingEntry                 = errors.New("missing cert entry")
)

const (
	uriIdentifier              = "URI"
	byIdentifier               = "By"
	expectedAmountOfMatchParts = 5
)

// Parse parses an istio spiffe cert
func (p *DefaultSpiffeParser) Parse(cert string) (*IstioSpiffeCert, error) {
	matches := p.expr.FindAllStringSubmatch(cert, 2)

	parsedCert := &IstioSpiffeCert{}

	for _, match := range matches {
		if len(match) != expectedAmountOfMatchParts {
			return nil, errors.Errorf("%w (was %d)", ErrInvalidNumberOfMatchSegments, len(match))
		}
		parts := &IstioSpiffeParts{
			domain:         match[2],
			namespace:      match[3],
			serviceAccount: match[4],
		}
		switch match[1] {
		case uriIdentifier:
			parsedCert.URI = parts
		case byIdentifier:
			parsedCert.By = parts
		}
	}
	if parsedCert.URI == nil {
		return nil, errors.Errorf("%w: %s", ErrMissingEntry, uriIdentifier)
	}
	if parsedCert.By == nil {
		return nil, errors.Errorf("%w: %s", ErrMissingEntry, byIdentifier)
	}
	return parsedCert, nil
}
