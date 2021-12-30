package luautil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	lua "github.com/yuin/gopher-lua"
)

type LuaUtilTestSuite struct {
	suite.Suite
}

func TestLuaUtilTestSuite(t *testing.T) {
	suite.Run(t, &LuaUtilTestSuite{})
}

func (s *LuaUtilTestSuite) TestToGoValue() {
	tests := []struct {
		name     string
		luaValue string
		goValue  interface{}
	}{
		{
			name:     "Nil",
			luaValue: "nil",
			goValue:  nil,
		},
		{
			name:     "Integer",
			luaValue: "1",
			goValue:  float64(1),
		},
		{
			name:     "Float",
			luaValue: "3.1415926535897",
			goValue:  3.1415926535897,
		},
		{
			name:     "Bool",
			luaValue: "true",
			goValue:  true,
		},
		{
			name: "Complex table",
			luaValue: `{
    fruit = "banaan",
    somethingElse = {
        appel = "taart",
        citroen = false
    }
}`,
			goValue: map[string]interface{}{
				"fruit": "banaan",
				"somethingElse": map[string]interface{}{
					"appel":   "taart",
					"citroen": false,
				},
			},
		},
		{
			name:     "Array",
			luaValue: `{"Amsterdam", "Berlin", "London", "Paris", "Copenhagen"}`,
			goValue:  []interface{}{"Amsterdam", "Berlin", "London", "Paris", "Copenhagen"},
		},
	}

	testCounter := 0

	state := lua.NewState()
	state.SetGlobal("gocheck", state.NewFunction(func(L *lua.LState) int {
		test := tests[testCounter]
		s.Run(test.name, func() {
			luaValue := L.Get(1)
			goValue, err := ToGoValue(luaValue)
			s.Require().NoError(err)

			s.Require().Equal(test.goValue, goValue)
		})
		testCounter++

		return 0
	}))

	luaValues := make([]string, 0)
	for _, test := range tests {
		luaValues = append(luaValues, test.luaValue)
	}

	err := state.DoString("gocheck(" + strings.Join(luaValues, ")\ngocheck(") + ")")
	s.NoError(err)
}

func (s *LuaUtilTestSuite) TestToGoValueError() {
	tests := []struct {
		name         string
		luaValue     string
		errorMessage string
	}{
		{
			name:         "ArrayTable",
			luaValue:     `{"Amsterdam", "Berlin", "London", france = "Paris"}`,
			errorMessage: "cannot encode mixed or invalid key types",
		},
		{
			name: "ArrayTable2",
			luaValue: `{
    [1] = "banaan",
    somethingElse = {
        appel = "taart",
        citroen = false
    }
}`,
			errorMessage: "cannot encode mixed or invalid key types",
		},
		{
			name:         "SparseArray",
			luaValue:     `{"Amsterdam", "Berlin", "London", "Paris", [7] = "Copenhagen"}`,
			errorMessage: "cannot encode sparse array",
		},
	}

	testCounter := 0

	state := lua.NewState()
	state.SetGlobal("gocheck", state.NewFunction(func(L *lua.LState) int {
		test := tests[testCounter]
		s.Run(test.name, func() {
			luaValue := L.Get(1)
			_, err := ToGoValue(luaValue)
			s.Require().Error(err)

			s.Contains(err.Error(), test.errorMessage)
		})

		testCounter++

		return 0
	}))

	luaValues := make([]string, 0)
	for _, test := range tests {
		luaValues = append(luaValues, test.luaValue)
	}

	err := state.DoString("gocheck(" + strings.Join(luaValues, ")\ngocheck(") + ")")
	s.NoError(err)
}

func (s *LuaUtilTestSuite) TestToGoValueNestedError() {
	state := lua.NewState()
	state.SetGlobal("gocheck", state.NewFunction(func(L *lua.LState) int {
		luaValue := L.Get(1)
		_, err := ToGoValue(luaValue)
		s.Require().Error(err)

		s.ErrorIs(err, ErrNested)

		return 0
	}))

	err := state.DoString(`table = {}
table2 = { somekey = table }
table.someotherkey = table2
gocheck(table)`)
	s.Require().NoError(err)
}
