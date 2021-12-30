package verification

import (
	"net/url"

	"lab.weave.nl/nid/nid-core/pkg/accessmodel"
	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

// Request corresponds to a request to be verified by a Verifier
type Request struct {
	Scopes map[string]*accessmodel.AccessModel
	Method string
	Path   string
	Query  url.Values
	Body   string
}

var (
	// ErrNotValid is returned by a verifier if the request is well-formed but access is not granted
	ErrNotValid = errors.New("request does not match scopes")
	// ErrBadRequest is returned by a verifier if the request is malformed
	ErrBadRequest = errors.New("bad request")
)

// Verifier verifies whether the request matches its scopes
type Verifier interface {
	Verify(req *Request) error
}
