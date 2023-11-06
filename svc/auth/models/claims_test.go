package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScopesMap(t *testing.T) {
	tests := []struct {
		Scenario string
		Claims   TokenClaims
		Expected map[string]interface{}
		Err      error
	}{
		{
			Scenario: "Valid map",
			Claims:   TokenClaims{Scopes: map[string]interface{}{"key": "value"}},
			Expected: map[string]interface{}{"key": "value"},
			Err:      nil,
		},
		{
			Scenario: "Invalid type (not a map)",
			Claims:   TokenClaims{Scopes: []interface{}{"item1", "item2"}},
			Expected: nil,
			Err:      errTypeAssertionFailed,
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			scopes, err := test.Claims.ScopesMap()
			if err != nil {
				assert.Equal(t, true, errors.Is(err, test.Err))
			}

			assert.Equal(t, test.Expected, scopes)
		})
	}
}

func TestScopesArray(t *testing.T) {
	tests := []struct {
		Scenario string
		Claims   TokenClaims
		Expected []string
		Err      error
	}{
		{
			Scenario: "Valid array",
			Claims:   TokenClaims{Scopes: []interface{}{"item1", "item2"}},
			Expected: []string{"item1", "item2"},
			Err:      nil,
		},
		{
			Scenario: "Invalid type (not an array)",
			Claims:   TokenClaims{Scopes: map[string]interface{}{"key": "value"}},
			Expected: nil,
			Err:      errTypeAssertionFailed,
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			scopes, err := test.Claims.ScopesArray()
			if err != nil {
				assert.Equal(t, true, errors.Is(err, test.Err))
			}

			assert.Equal(t, test.Expected, scopes)
		})
	}
}
