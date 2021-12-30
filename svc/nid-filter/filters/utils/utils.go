// Package utils contains utilities that are common to many filters.
package utils

import (
	"fmt"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3alpha"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"

	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
)

// GraphqlError creates a new immediateresponse containing a simple error that adheres to the graphql spec.
func GraphqlError(message string, code envoy_type_v3.StatusCode) *filter.ProcessingResponse {
	body := fmt.Sprintf(`{
    "errors": [
        {"message": "%s"}
    ]
}`,
		message)

	return &filter.ProcessingResponse{
		ImmediateResponse: &ext_proc_pb.ImmediateResponse{
			Status: &envoy_type_v3.HttpStatus{Code: code},
			Headers: &ext_proc_pb.HeaderMutation{
				SetHeaders: []*envoy_config_core_v3.HeaderValueOption{
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "Content-Type",
							Value: "application/json",
						},
					},
				},
			},
			Body: body,
		},
	}
}
