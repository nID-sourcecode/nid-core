package verification

import (
	"net/url"

	"github.com/nID-sourcecode/nid-core/pkg/accessmodel"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

// RESTVerifier verifies that the request matches any of the REST scopes in the token.
type RESTVerifier struct{}

// Verify verifies that the request matches any of the REST scopes in the token.
func (*RESTVerifier) Verify(req *Request) error {
	restScopes := accessmodel.FilterByType(req.Scopes, accessmodel.RESTType)
	if len(restScopes) == 0 {
		return errors.Errorf("%w: no REST scopes found", ErrNotValid)
	}

	for _, scope := range restScopes {
		model := scope.RESTAccessModel
		if model.Path == req.Path &&
			model.Method == req.Method &&
			model.Body == req.Body &&
			queryEquals(req.Query, model.Query) {
			return nil
		}
	}

	return ErrNotValid
}

func queryEquals(reqQuery url.Values, scopeQuery map[string]string) bool {
	for key, scopeValue := range scopeQuery {
		if reqQuery.Get(key) != scopeValue {
			return false
		}
	}

	for key := range reqQuery {
		_, scopeHasKey := scopeQuery[key]
		if !scopeHasKey {
			return false
		}
	}

	return true
}
