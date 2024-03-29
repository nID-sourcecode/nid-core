// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package mocks

import (
	context "context"
	"github.com/nID-sourcecode/nid-core/pkg/gqlclient"

	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, req, resp
func (_m *Client) Get(ctx context.Context, req gqlclient.Request, resp interface{}) error {
	ret := _m.Called(ctx, req, resp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, gqlclient.Request, interface{}) error); ok {
		r0 = rf(ctx, req, resp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Post provides a mock function with given fields: ctx, req, resp
func (_m *Client) Post(ctx context.Context, req gqlclient.Request, resp interface{}) error {
	ret := _m.Called(ctx, req, resp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, gqlclient.Request, interface{}) error); ok {
		r0 = rf(ctx, req, resp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: ctx, req, resp, method
func (_m *Client) Run(ctx context.Context, req gqlclient.Request, resp interface{}, method gqlclient.Method) error {
	ret := _m.Called(ctx, req, resp, method)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, gqlclient.Request, interface{}, gqlclient.Method) error); ok {
		r0 = rf(ctx, req, resp, method)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
