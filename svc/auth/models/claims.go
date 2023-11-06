package models

import (
	"encoding/json"
	"fmt"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/jwt/v3"
)

var errTypeAssertionFailed = errors.New("Type assertion failed")

// TokenClaims specifies the claims used for JWT's
type TokenClaims struct {
	*jwt.DefaultClaims
	ClientID       string                 `json:"client_id"`
	Subjects       map[string]interface{} `json:"subjects"`
	Scopes         interface{}            `json:"scopes"`
	ConsentID      string                 `json:"consent_id,omitempty"` // ConsentID is not set for swapping tokens
	ClientMetadata map[string]interface{} `json:"client_metadata"`
}

// ListKeys will list all the keys inside token claims
func (t *TokenClaims) ListKeys() ([]string, error) {
	claimsJSON, err := json.Marshal(t)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling claims")
	}

	var claimsMap map[string]interface{}
	err = json.Unmarshal(claimsJSON, &claimsMap)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling claims as map")
	}

	keyList := make([]string, len(claimsMap))
	var i int
	for key := range claimsMap {
		keyList[i] = key
		i++
	}

	return keyList, nil
}

// ScopesMap will return the scopes as a map. If the scopes are not a map, an error will be returned.
func (t *TokenClaims) ScopesMap() (map[string]interface{}, error) {
	convertedVal, ok := t.Scopes.(map[string]interface{})
	if !ok {
		return nil, errTypeAssertionFailed
	}

	return convertedVal, nil
}

// ScopesArray will return the scopes as an array. If the scopes are not an array, an error will be returned.
func (t *TokenClaims) ScopesArray() ([]string, error) {
	slice, ok := t.Scopes.([]interface{})
	if !ok {
		return nil, errTypeAssertionFailed
	}

	strSlice := make([]string, len(slice))
	for i, value := range slice {
		strSlice[i] = fmt.Sprintf("%v", value)
	}

	return strSlice, nil
}
