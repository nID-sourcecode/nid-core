package mutation

import (
	"testing"

	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/stretchr/testify/suite"
)

type MutationTestSuite struct {
	suite.Suite
	calculator Calculator
}

func TestMutationTestSuite(t *testing.T) {
	suite.Run(t, &MutationTestSuite{})
}

func (s *MutationTestSuite) SetupSuite() {
	s.calculator = &DefaultCalculator{}
}

func (s *MutationTestSuite) TestDefaultCalculator_CalculateHeaderMutations() {
	tests := []struct {
		Name                   string
		OriginalHeaders        map[string]string
		NewHeaders             map[string]string
		ExpectedHeaderMutation *ext_proc_pb.HeaderMutation
	}{
		{
			Name: "NewHeadersNil",
			OriginalHeaders: map[string]string{
				"some-header": "some-value",
			},
			NewHeaders:             nil,
			ExpectedHeaderMutation: nil,
		},
		{
			Name: "HeaderRemoved",
			OriginalHeaders: map[string]string{
				"header1": "value1",
				"header2": "value2",
			},
			NewHeaders: map[string]string{
				"header2": "value2",
			},
			ExpectedHeaderMutation: &ext_proc_pb.HeaderMutation{
				SetHeaders:    nil,
				RemoveHeaders: []string{"header1"},
			},
		},
		{
			Name: "HeaderChanged",
			OriginalHeaders: map[string]string{
				"header1": "value1",
				"header2": "value2",
			},
			NewHeaders: map[string]string{
				"header1": "value1",
				"header2": "new_value",
			},
			ExpectedHeaderMutation: &ext_proc_pb.HeaderMutation{
				SetHeaders: []*envoy_core_v3.HeaderValueOption{
					{
						Header: &envoy_core_v3.HeaderValue{
							Key:   "header2",
							Value: "new_value",
						},
					},
				},
				RemoveHeaders: nil,
			},
		},
		{
			Name: "HeaderAdded",
			OriginalHeaders: map[string]string{
				"header1": "value1",
				"header2": "value2",
			},
			NewHeaders: map[string]string{
				"header1": "value1",
				"header2": "value2",
				"header3": "value3",
			},
			ExpectedHeaderMutation: &ext_proc_pb.HeaderMutation{
				SetHeaders: []*envoy_core_v3.HeaderValueOption{
					{
						Header: &envoy_core_v3.HeaderValue{
							Key:   "header3",
							Value: "value3",
						},
					},
				},
				RemoveHeaders: nil,
			},
		},
		{
			Name: "Everything",
			OriginalHeaders: map[string]string{
				"header1": "value1",
				"header2": "value2",
			},
			NewHeaders: map[string]string{
				"header2": "new_value",
				"header3": "value3",
			},
			ExpectedHeaderMutation: &ext_proc_pb.HeaderMutation{
				SetHeaders: []*envoy_core_v3.HeaderValueOption{
					{
						Header: &envoy_core_v3.HeaderValue{
							Key:   "header3",
							Value: "value3",
						},
					},
					{
						Header: &envoy_core_v3.HeaderValue{
							Key:   "header2",
							Value: "new_value",
						},
					},
				},
				RemoveHeaders: []string{"header1"},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			mutation := s.calculator.CalculateHeaderMutations(test.OriginalHeaders, test.NewHeaders)
			if test.ExpectedHeaderMutation == nil {
				s.Nil(mutation)
			} else {
				s.ElementsMatch(test.ExpectedHeaderMutation.SetHeaders, mutation.SetHeaders)
				s.ElementsMatch(test.ExpectedHeaderMutation.RemoveHeaders, mutation.RemoveHeaders)
			}
		})
	}
}

func (s *MutationTestSuite) TestDefaultCalculator_CalculateBodyMutation() {
	tests := []struct {
		Name                 string
		OriginalBody         []byte
		NewBody              []byte
		ExpectedBodyMutation *ext_proc_pb.BodyMutation
	}{
		{
			Name:         "BodyFilled",
			OriginalBody: []byte{},
			NewBody:      []byte("summer body"),
			ExpectedBodyMutation: &ext_proc_pb.BodyMutation{
				Mutation: &ext_proc_pb.BodyMutation_Body{
					Body: []byte("summer body"),
				},
			},
		},
		{
			Name:         "BodyAdded",
			OriginalBody: nil,
			NewBody:      []byte("summer body"),
			ExpectedBodyMutation: &ext_proc_pb.BodyMutation{
				Mutation: &ext_proc_pb.BodyMutation_Body{
					Body: []byte("summer body"),
				},
			},
		},
		{
			Name:         "BodyChanged",
			OriginalBody: []byte("summer body"),
			NewBody:      []byte("autumn body"),
			ExpectedBodyMutation: &ext_proc_pb.BodyMutation{
				Mutation: &ext_proc_pb.BodyMutation_Body{
					Body: []byte("autumn body"),
				},
			},
		},
		{
			Name:         "BodyRemoved",
			OriginalBody: []byte("summer body"),
			NewBody:      []byte{},
			ExpectedBodyMutation: &ext_proc_pb.BodyMutation{
				Mutation: &ext_proc_pb.BodyMutation_ClearBody{ClearBody: true},
			},
		},
		{
			Name:                 "NoChanges",
			OriginalBody:         []byte("summer body"),
			NewBody:              nil,
			ExpectedBodyMutation: nil,
		},
	}

	for _, test := range tests {
		s.Run(test.Name, func() {
			mutation := s.calculator.CalculateBodyMutation(test.OriginalBody, test.NewBody)
			s.Equal(test.ExpectedBodyMutation, mutation)
		})
	}
}
