// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	proto "lab.weave.nl/nid/nid-core/svc/auth/proto"
)

// WellKnownServer is an autogenerated mock type for the WellKnownServer type
type WellKnownServer struct {
	mock.Mock
}

// GetWellKnownOAuthAuthorizationServer provides a mock function with given fields: _a0, _a1
func (_m *WellKnownServer) GetWellKnownOAuthAuthorizationServer(_a0 context.Context, _a1 *proto.WellKnownRequest) (*proto.WellKnownResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *proto.WellKnownResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.WellKnownRequest) *proto.WellKnownResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.WellKnownResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.WellKnownRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWellKnownOpenIDConfiguration provides a mock function with given fields: _a0, _a1
func (_m *WellKnownServer) GetWellKnownOpenIDConfiguration(_a0 context.Context, _a1 *proto.WellKnownRequest) (*proto.WellKnownResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *proto.WellKnownResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.WellKnownRequest) *proto.WellKnownResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.WellKnownResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.WellKnownRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
