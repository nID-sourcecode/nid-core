// Code generated by mockery v2.9.4. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ScriptRunner is an autogenerated mock type for the ScriptRunner type
type ScriptRunner struct {
	mock.Mock
}

// RunScript provides a mock function with given fields: ctx, script, input
func (_m *ScriptRunner) RunScript(ctx context.Context, script string, input map[string]interface{}) error {
	ret := _m.Called(ctx, script, input)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) error); ok {
		r0 = rf(ctx, script, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
