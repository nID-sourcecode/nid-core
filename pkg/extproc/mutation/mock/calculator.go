// Code generated by mockery v2.4.0-beta. DO NOT EDIT.

package mock

import (
	envoy_service_ext_proc_v3alpha "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3alpha"
	mock "github.com/stretchr/testify/mock"
)

// Calculator is an autogenerated mock type for the Calculator type
type Calculator struct {
	mock.Mock
}

// CalculateBodyMutation provides a mock function with given fields: originalBody, newBody
func (_m *Calculator) CalculateBodyMutation(originalBody []byte, newBody []byte) *envoy_service_ext_proc_v3alpha.BodyMutation {
	ret := _m.Called(originalBody, newBody)

	var r0 *envoy_service_ext_proc_v3alpha.BodyMutation
	if rf, ok := ret.Get(0).(func([]byte, []byte) *envoy_service_ext_proc_v3alpha.BodyMutation); ok {
		r0 = rf(originalBody, newBody)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*envoy_service_ext_proc_v3alpha.BodyMutation)
		}
	}

	return r0
}

// CalculateHeaderMutations provides a mock function with given fields: originalHeaders, newHeaders
func (_m *Calculator) CalculateHeaderMutations(originalHeaders map[string]string, newHeaders map[string]string) *envoy_service_ext_proc_v3alpha.HeaderMutation {
	ret := _m.Called(originalHeaders, newHeaders)

	var r0 *envoy_service_ext_proc_v3alpha.HeaderMutation
	if rf, ok := ret.Get(0).(func(map[string]string, map[string]string) *envoy_service_ext_proc_v3alpha.HeaderMutation); ok {
		r0 = rf(originalHeaders, newHeaders)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*envoy_service_ext_proc_v3alpha.HeaderMutation)
		}
	}

	return r0
}
