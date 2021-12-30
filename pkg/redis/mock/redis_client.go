// Code generated by mockery v2.6.0. DO NOT EDIT.

package mock

import (
	redis "github.com/go-redis/redis/v8"
	mock "github.com/stretchr/testify/mock"
)

// RedisClient is an autogenerated mock type for the RedisClient type
type RedisClient struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *RedisClient) Get(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Publish provides a mock function with given fields: _a0, _a1
func (_m *RedisClient) Publish(_a0 string, _a1 string) string {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Set provides a mock function with given fields: _a0, _a1
func (_m *RedisClient) Set(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: _a0
func (_m *RedisClient) Subscribe(_a0 string) *redis.PubSub {
	ret := _m.Called(_a0)

	var r0 *redis.PubSub
	if rf, ok := ret.Get(0).(func(string) *redis.PubSub); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.PubSub)
		}
	}

	return r0
}
