package testhelpers

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/suite"
)

type ModelHelpersTestSuite struct {
	suite.Suite
}

func TestModelHelpersTestSuite(t *testing.T) {
	suite.Run(t, &ModelHelpersTestSuite{})
}

func (s *ModelHelpersTestSuite) TestEqualModelsDate() {
	type iets struct {
		DateField time.Time
	}

	ietsValueA := iets{DateField: time.Now()}
	ietsValueADifferentNano := iets{time.Unix(ietsValueA.DateField.Unix(), 0)}
	ietsValueB := iets{DateField: time.Now().AddDate(-3, -2, 5)}

	s.Run("equal", func() {
		s.True(equalModel(ietsValueA, ietsValueA))
	})

	s.Run("completelyDifferent", func() {
		s.False(equalModel(ietsValueA, ietsValueB))
	})

	s.Run("nanosNotEqual", func() {
		s.True(equalModel(ietsValueA, ietsValueADifferentNano))
	})

	s.Run("bothNilTimes", func() {
		s.True(equalModel(iets{}, iets{}))
	})
}

func (s ModelHelpersTestSuite) TestEqualModelsJSON() {
	type testStruct struct {
		JSONField postgres.Jsonb
	}

	s.Run("whitespaceIgnored", func() {
		testStructValue := testStruct{JSONField: postgres.Jsonb{RawMessage: []byte(`{"test":   "yes"}`)}}
		testStructValueNoWhitespace := testStruct{JSONField: postgres.Jsonb{RawMessage: []byte(`{"test":"yes"}`)}}

		s.True(equalModel(testStructValue, testStructValueNoWhitespace))
	})

	s.Run("testequalNil", func() {
		testStructValue := testStruct{JSONField: postgres.Jsonb{}}
		otherTestStructValue := testStruct{JSONField: postgres.Jsonb{}}

		s.True(equalModel(testStructValue, otherTestStructValue))
	})

	s.Run("valueNotIgnored", func() {
		testStructValue := testStruct{JSONField: postgres.Jsonb{RawMessage: []byte(`{"test":   "yes"}`)}}
		testStructValueDifferent := testStruct{JSONField: postgres.Jsonb{RawMessage: []byte(`{"test":"yes"}`)}}

		s.True(equalModel(testStructValue, testStructValueDifferent))
	})

	s.Run("valueNestedEqual", func() {
		json := []byte(`
		{
  			"test": {
				"a": "something",
 				"b": "else"
			}
		}
		`)

		testStructValue := testStruct{JSONField: postgres.Jsonb{RawMessage: json}}
		testStructValueDifferent := testStruct{JSONField: postgres.Jsonb{RawMessage: json}}

		s.True(equalModel(testStructValue, testStructValueDifferent))
	})

	s.Run("valueNestedNotEqual", func() {
		jsonValue := []byte(`
		{
  			"test": {
				"a": "something",
 				"b": "else"
			}
		}
		`)

		otherJSONValue := []byte(`
		{
  			"test": {
				"a": 1223,
 				"b": "else"
			}
		}
		`)

		testStructValue := testStruct{JSONField: postgres.Jsonb{RawMessage: jsonValue}}
		testStructValueDifferent := testStruct{JSONField: postgres.Jsonb{RawMessage: otherJSONValue}}

		s.False(equalModel(testStructValue, testStructValueDifferent))
	})
}
