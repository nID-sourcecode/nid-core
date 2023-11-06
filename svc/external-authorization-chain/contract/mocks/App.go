// Code generated by mockery v2.30.1. DO NOT EDIT.

package mock

import (
	context "context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	mock "github.com/stretchr/testify/mock"
)

// App is an autogenerated mock type for the App type
type App struct {
	mock.Mock
}

// CheckEndpoints provides a mock function with given fields: ctx, request
func (_m *App) CheckEndpoints(ctx context.Context, request *authv3.CheckRequest) error {
	ret := _m.Called(ctx, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *authv3.CheckRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewApp creates a new instance of App. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApp(t interface {
	mock.TestingT
	Cleanup(func())
}) *App {
	mock := &App{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}