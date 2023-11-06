// Package luautil contains utility functions for lua runners
package luautil

import (
	"encoding/json"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	luajson "layeh.com/gopher-json"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

func convertDateStringsToTicks(obj map[string]interface{}) {
	for k, v := range obj {
		if value, ok := v.(string); ok {
			timestamp := time.Time{}
			err := json.Unmarshal([]byte("\""+value+"\""), &timestamp)
			if err == nil {
				obj[k] = float64(timestamp.Unix())
			}
		}
		if value, ok := v.(map[string]interface{}); ok {
			convertDateStringsToTicks(value)
		}
	}
}

// CompileScript parses string lua script to compiled script
func CompileScript(script string) (*lua.FunctionProto, error) {
	chunks, err := parse.Parse(strings.NewReader(script), "testScript")
	if err != nil {
		return nil, err
	}

	compiledScript, err := lua.Compile(chunks, "testScript")
	if err != nil {
		return nil, err
	}

	return compiledScript, nil
}

// ToLuaValue converts any go value to a lua value
func ToLuaValue(state *lua.LState, goValue interface{}) lua.LValue {
	if table, ok := goValue.(map[string]interface{}); ok {
		convertDateStringsToTicks(table)
	}
	luaValue := luajson.DecodeValue(state, goValue)

	return luaValue
}

// Error definitions
var (
	ErrValueTypeNotSupported = errors.New("unsupported kind of value read from lua")
	ErrGoValueNotAMap        = errors.New("LTable parsing result result was not a map")
	ErrNested                = errors.New("cannot encode recursively nested tables to JSON")
	ErrSparseArray           = errors.New("cannot encode sparse array")
	ErrInvalidKeys           = errors.New("cannot encode mixed or invalid key types")
)

// ToGoValue converts any lua value to a go value
func ToGoValue(lv lua.LValue) (interface{}, error) {
	return toGoValue(lv, make(map[lua.LValue]bool))
}

//nolint:gocognit
func toGoValue(lv lua.LValue, visited map[lua.LValue]bool) (interface{}, error) {
	switch converted := lv.(type) {
	case *lua.LNilType:
		return nil, nil
	case lua.LBool:
		return bool(converted), nil
	case lua.LNumber:
		return float64(converted), nil
	case lua.LString:
		return string(converted), nil
	case *lua.LTable:
		if visited[converted] {
			return nil, ErrNested
		}
		visited[converted] = true

		key, value := converted.Next(lua.LNil)

		switch key.Type() { // nolint
		case lua.LTNil: // empty array
			return make([]interface{}, 0), nil
		case lua.LTNumber: // array -> slice
			array := make([]interface{}, 0, converted.Len())
			expectedKey := lua.LNumber(1)
			for key != lua.LNil {
				if key.Type() != lua.LTNumber {
					return nil, ErrInvalidKeys
				}
				if expectedKey != key {
					return nil, ErrSparseArray
				}
				goValue, err := toGoValue(value, visited)
				if err != nil {
					return nil, errors.Wrap(err, "parsing array value")
				}
				array = append(array, goValue)
				expectedKey++
				key, value = converted.Next(key)
			}
			return array, nil
		case lua.LTString: // table -> map
			obj := make(map[string]interface{})
			for key != lua.LNil {
				if key.Type() != lua.LTString {
					return nil, ErrInvalidKeys
				}
				goValue, err := toGoValue(value, visited)
				if err != nil {
					return nil, errors.Wrap(err, "parsing table value")
				}
				obj[key.String()] = goValue
				key, value = converted.Next(key)
			}
			return obj, nil
		}
	}
	return nil, ErrValueTypeNotSupported
}

// ToGoMap converts a lua table to a go map.
func ToGoMap(luaTable *lua.LTable) (map[string]interface{}, error) {
	res, err := ToGoValue(luaTable)
	if err != nil {
		return nil, errors.Wrap(err, "parsing table")
	}

	resMap, ok := res.(map[string]interface{})
	if !ok {
		return nil, ErrGoValueNotAMap
	}
	return resMap, nil
}
