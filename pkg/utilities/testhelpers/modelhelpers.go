// Package testhelpers provides helper functions to be used during testing
package testhelpers

import (
	"encoding/json"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

// EqualModel checks if two models are equal, it uses custom compare functions for time.Time and postgres.JSON fields
// time.Time is considered equal if the output .Unix() of both structs is equal
// postgres.JSON is considered equal if a json result unmarshalling the json.Rawmessage values of the fields into a map[string]interface{}{} are equal
func EqualModel(s suite.Suite, expected, actual interface{}) {
	// this is a simple wrapper function to make the actual logic more testable
	s.True(equalModel(expected, actual))
}

func equalModel(expected, actual interface{}) (bool, string) {
	dateComparer := cmp.Comparer(func(a, b time.Time) bool {
		return a.Unix() == b.Unix()
	})

	metadataJSONComparer := cmp.Comparer(func(a, b json.RawMessage) bool {
		if a == nil && b == nil {
			return true
		}

		valueA := map[string]interface{}{}
		err := json.Unmarshal(a, &valueA)
		if err != nil {
			return false
		}

		valueB := map[string]interface{}{}
		err = json.Unmarshal(b, &valueB)
		if err != nil {
			return false
		}

		return cmp.Equal(valueA, valueB)
	})

	compareOpts := []cmp.Option{dateComparer, metadataJSONComparer}
	return cmp.Equal(expected, actual, compareOpts...), cmp.Diff(expected, actual, compareOpts...)
}
