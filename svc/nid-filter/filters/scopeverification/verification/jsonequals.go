package verification

import (
	"bytes"
	"encoding/json"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

// JSONEquals checks whether the JSON of two objects is equal
func JSONEquals(a, b interface{}) (bool, error) {
	jsonA, err := json.Marshal(a)
	if err != nil {
		return false, errors.Wrap(err, "unable to marshall object a")
	}
	jsonB, err := json.Marshal(b)
	if err != nil {
		return false, errors.Wrap(err, "unable to marshall object b")
	}

	return bytes.Equal(jsonA, jsonB), nil
}
