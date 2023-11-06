package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicates_Strings(t *testing.T) {
	tests := []struct {
		scenario string
		input    []string
		expected []string
	}{
		{
			scenario: "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			scenario: "slice without duplicates",
			input:    []string{"apple", "banana", "orange"},
			expected: []string{"apple", "banana", "orange"},
		},
		{
			scenario: "slice with duplicates",
			input:    []string{"apple", "banana", "apple", "orange", "banana"},
			expected: []string{"apple", "banana", "orange"},
		},
	}

	for _, test := range tests {
		result := RemoveDuplicates(test.input)
		t.Run(test.scenario, func(t *testing.T) {
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestRemoveDuplicates_Ints(t *testing.T) {
	tests := []struct {
		scenario string
		input    []int
		expected []int
	}{
		{
			scenario: "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			scenario: "slice without duplicates",
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			scenario: "slice with duplicates",
			input:    []int{1, 2, 1},
			expected: []int{1, 2},
		},
	}

	for _, test := range tests {
		result := RemoveDuplicates(test.input)
		t.Run(test.scenario, func(t *testing.T) {
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestRemoveDuplicates_Structs(t *testing.T) {
	type testStruct struct {
		Name string
		Age  int
	}

	tests := []struct {
		scenario string
		input    []testStruct
		expected []testStruct
	}{
		{
			scenario: "empty slice",
			input:    []testStruct{},
			expected: []testStruct{},
		},
		{
			scenario: "slice without duplicates",
			input:    []testStruct{{Name: "John", Age: 20}, {Name: "Jane", Age: 21}},
			expected: []testStruct{{Name: "John", Age: 20}, {Name: "Jane", Age: 21}},
		},
		{
			scenario: "slice with duplicates",
			input:    []testStruct{{Name: "John", Age: 20}, {Name: "Jane", Age: 21}, {Name: "John", Age: 20}},
			expected: []testStruct{{Name: "John", Age: 20}, {Name: "Jane", Age: 21}},
		},
	}

	for _, test := range tests {
		result := RemoveDuplicates(test.input)
		t.Run(test.scenario, func(t *testing.T) {
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIntersect_Strings(t *testing.T) {
	tests := []struct {
		scenario string
		inputA   []string
		inputB   []string
		expected []string
	}{
		{
			scenario: "empty slice",
			inputA:   []string{},
			inputB:   []string{},
			expected: []string{},
		},
		{
			scenario: "slice without intersection",
			inputA:   []string{"apple", "banana", "orange"},
			inputB:   []string{"pear", "grape"},
			expected: []string{},
		},
		{
			scenario: "slice with one intersection",
			inputA:   []string{"apple", "banana", "orange"},
			inputB:   []string{"pear", "grape", "apple"},
			expected: []string{"apple"},
		},
		{
			scenario: "slice with multiple intersections",
			inputA:   []string{"apple", "banana", "orange", "pear", "grape"},
			inputB:   []string{"pear", "grape", "apple", "banana", "kiwi", "mango", "orange"},
			expected: []string{"apple", "banana", "orange", "pear", "grape"},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := Intersect(test.inputA, test.inputB)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIntersect_Structs(t *testing.T) {
	type testStruct struct {
		Name string
		Age  int
	}

	tests := []struct {
		scenario string
		inputA   []testStruct
		inputB   []testStruct
		expected []testStruct
	}{
		{
			scenario: "empty slice",
			inputA:   []testStruct{},
			inputB:   []testStruct{},
			expected: []testStruct{},
		},
		{
			scenario: "slice without intersection",
			inputA:   []testStruct{{Name: "John", Age: 20}, {Name: "Jane", Age: 21}},
			inputB:   []testStruct{{Name: "Joe", Age: 20}, {Name: "Kate", Age: 21}},
			expected: []testStruct{},
		},
		{
			scenario: "slice with one intersection",
			inputA:   []testStruct{{Name: "John", Age: 20}, {Name: "Jane", Age: 21}},
			inputB:   []testStruct{{Name: "Joe", Age: 20}, {Name: "Kate", Age: 21}, {Name: "John", Age: 20}},
			expected: []testStruct{{Name: "John", Age: 20}},
		},
		{
			scenario: "slice with multiple intersections",
			inputA:   []testStruct{{Name: "John", Age: 20}, {Name: "Jane", Age: 21}, {Name: "Joe", Age: 22}, {Name: "Kate", Age: 23}},
			inputB:   []testStruct{{Age: 20, Name: "John"}, {Name: "Kate", Age: 23}, {Name: "Pete", Age: 22}, {Name: "Jane", Age: 23}},
			expected: []testStruct{{Name: "John", Age: 20}, {Name: "Kate", Age: 23}},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := Intersect(test.inputA, test.inputB)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestContainsAll(t *testing.T) {
	tests := []struct {
		scenario  string
		inputAll  []string
		inputPart []string
		expected  bool
	}{
		{
			scenario:  "empty slice",
			inputAll:  []string{},
			inputPart: []string{},
			expected:  true,
		},
		{
			scenario:  "slice contains all elements",
			inputAll:  []string{"apple", "banana", "orange"},
			inputPart: []string{"apple", "banana"},
			expected:  true,
		},
		{
			scenario:  "slice doesn't contain all elements",
			inputAll:  []string{"apple", "banana", "orange"},
			inputPart: []string{"apple", "banana", "pear"},
			expected:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := ContainsAll(test.inputAll, test.inputPart)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		scenario string
		input    []string
		element  string
		expected bool
	}{
		{
			scenario: "empty slice",
			input:    []string{},
			element:  "apple",
			expected: false,
		},
		{
			scenario: "slice contains element",
			input:    []string{"apple", "banana", "orange"},
			element:  "apple",
			expected: true,
		},
		{
			scenario: "slice doesn't contain element",
			input:    []string{"apple", "banana", "orange"},
			element:  "pear",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := Contains(test.input, test.element)
			assert.Equal(t, test.expected, result)
		})
	}
}
