// Package mutation contains logic to determine header and body mutations based on the changes made in filters.
package mutation

import (
	"bytes"

	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

// Calculator calculates header and body mutations from the differences between new and original.
type Calculator interface {
	CalculateHeaderMutations(originalHeaders map[string]string, newHeaders map[string]string) *ext_proc_pb.HeaderMutation
	CalculateBodyMutation(originalBody []byte, newBody []byte) *ext_proc_pb.BodyMutation
}

// DefaultCalculator is the default calculator implementation
type DefaultCalculator struct{}

// CalculateHeaderMutations calculates header mutations by comparing the original and new headers
func (d DefaultCalculator) CalculateHeaderMutations(originalHeaders map[string]string, newHeaders map[string]string) *ext_proc_pb.HeaderMutation {
	if newHeaders == nil {
		return nil
	}

	res := &ext_proc_pb.HeaderMutation{
		SetHeaders:    []*envoy_core_v3.HeaderValueOption{},
		RemoveHeaders: []string{},
	}
	for key := range originalHeaders {
		if _, keyExistsInNewHeaders := newHeaders[key]; !keyExistsInNewHeaders {
			res.RemoveHeaders = append(res.RemoveHeaders, key)
		}
	}

	for key, newValue := range newHeaders {
		oldValue, keyExistsInOldHeaders := originalHeaders[key]
		if !keyExistsInOldHeaders || oldValue != newValue {
			res.SetHeaders = append(res.SetHeaders, &envoy_core_v3.HeaderValueOption{
				Header: &envoy_core_v3.HeaderValue{
					Key:   key,
					Value: newValue,
				},
			})
		}
	}

	if len(res.SetHeaders) == 0 {
		if len(res.RemoveHeaders) == 0 {
			return nil
		}
		res.SetHeaders = nil
	} else if len(res.RemoveHeaders) == 0 {
		res.RemoveHeaders = nil
	}

	return res
}

// CalculateBodyMutation calculates the body mutation by comparing the original and new body
func (d DefaultCalculator) CalculateBodyMutation(originalBody []byte, newBody []byte) *ext_proc_pb.BodyMutation {
	if newBody == nil || bytes.Equal(originalBody, newBody) {
		return nil
	}

	if len(newBody) == 0 {
		return &ext_proc_pb.BodyMutation{Mutation: &ext_proc_pb.BodyMutation_ClearBody{ClearBody: true}}
	}

	return &ext_proc_pb.BodyMutation{Mutation: &ext_proc_pb.BodyMutation_Body{Body: newBody}}
}
