package main

import (
	"encoding/json"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v3"
)

// TokenClaims specifies the claims used for JWT's
type TokenClaims struct {
	*jwt.DefaultClaims
	ClientID       string                 `json:"client_id"`
	Subjects       map[string]interface{} `json:"subjects"`
	Scopes         map[string]interface{} `json:"scopes"`
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
