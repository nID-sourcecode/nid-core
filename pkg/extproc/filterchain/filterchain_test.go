package filterchain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	"lab.weave.nl/nid/nid-core/pkg/extproc/filter/mock"
)

type FilterChainTestSuite struct {
	suite.Suite
}

func TestFilterChainTestSuite(t *testing.T) {
	suite.Run(t, &FilterChainTestSuite{})
}

func (s *FilterChainTestSuite) TestBuildFilterChain() {
	mockFilter1 := &mock.Filter{}
	// mockFilter2 := mock.Filter{}
	mockFilter3 := &mock.Filter{}

	mockFilterInitializer1 := &mock.FilterInitializer{}
	mockFilterInitializer1.On("Name").Return("somefilter")
	mockFilterInitializer1.On("NewFilter").Return(mockFilter1, nil)
	mockFilterInitializer2 := &mock.FilterInitializer{}
	mockFilterInitializer2.On("Name").Return("someotherfilter")
	// mockFilterInitializer2.On("InitializeFilter").Return(mockFilter2, nil)
	mockFilterInitializer3 := &mock.FilterInitializer{}
	mockFilterInitializer3.On("Name").Return("yetanotherfilter")
	mockFilterInitializer3.On("NewFilter").Return(mockFilter3, nil)

	initializers := []filter.Initializer{
		mockFilterInitializer1,
		mockFilterInitializer2,
		mockFilterInitializer3,
	}

	chain, err := BuildDefaultFilterChain([]string{"somefilter", "yetanotherfilter"}, initializers)

	s.Require().NoError(err)
	s.Require().NotNil(chain)
	s.Require().NotNil(chain.filters)
	s.Require().Len(chain.filters, 2)
	s.Require().Equal(mockFilter1, chain.filters[0])
	s.Require().Equal(mockFilter3, chain.filters[1])
}

func (s *FilterChainTestSuite) TestFilterChainGoesInReverseOnResponse() {
	mockFilter1 := &mock.Filter{}
	mockFilter2 := &mock.Filter{}
	mockFilter3 := &mock.Filter{}

	filters := []filter.Filter{
		mockFilter1,
		mockFilter2,
		mockFilter3,
	}

	ctx := context.TODO()
	mockFilter1.On("Name").Return("filter1")
	mockFilter1.On("OnHTTPResponse", ctx, []byte("appeltaart"), map[string]string{}).Return(&filter.ProcessingResponse{
		NewHeaders:        nil,
		NewBody:           []byte("appeltaart is lekker"),
		ImmediateResponse: nil,
	}, nil)
	mockFilter2.On("Name").Return("filter2")
	mockFilter2.On("OnHTTPResponse", ctx, []byte("appel"), map[string]string{}).Return(&filter.ProcessingResponse{
		NewHeaders:        nil,
		NewBody:           []byte("appeltaart"),
		ImmediateResponse: nil,
	}, nil)
	mockFilter3.On("Name").Return("filter3")
	mockFilter3.On("OnHTTPResponse", ctx, []byte("a"), map[string]string{}).Return(&filter.ProcessingResponse{
		NewHeaders:        nil,
		NewBody:           []byte("appel"),
		ImmediateResponse: nil,
	}, nil)

	chain := &DefaultChain{filters: filters}

	chain.ProcessResponseHeaders(map[string]string{})
	res, err := chain.ProcessResponseBody(ctx, []byte("a"))

	s.Require().NoError(err)
	s.Require().Equal([]byte("appeltaart is lekker"), res.Body)
}

func (s *FilterChainTestSuite) TestFilterChainRequest() {
	mockFilter1 := &mock.Filter{}
	mockFilter2 := &mock.Filter{}
	mockFilter3 := &mock.Filter{}

	filters := []filter.Filter{
		mockFilter1,
		mockFilter2,
		mockFilter3,
	}

	mapSomeValue := map[string]string{"coolheader": "some value"}
	mapSomeOtherValue := map[string]string{"coolheader": "some other value"}

	ctx := context.TODO()
	mockFilter1.On("Name").Return("filter1")
	mockFilter1.On("OnHTTPRequest", ctx, []byte("a"), map[string]string{}).Return(&filter.ProcessingResponse{
		NewHeaders:        mapSomeValue,
		NewBody:           []byte("appel"),
		ImmediateResponse: nil,
	}, nil)
	mockFilter2.On("Name").Return("filter2")
	mockFilter2.On("OnHTTPRequest", ctx, []byte("appel"), map[string]string{"coolheader": "some value"}).Return(&filter.ProcessingResponse{
		NewHeaders:        mapSomeOtherValue,
		NewBody:           []byte("appeltaart"),
		ImmediateResponse: nil,
	}, nil)
	mockFilter3.On("Name").Return("filter3")
	mockFilter3.On("OnHTTPRequest", ctx, []byte("appeltaart"), map[string]string{"coolheader": "some other value"}).Return(&filter.ProcessingResponse{
		NewHeaders:        map[string]string{"coolheader": "some other value"},
		NewBody:           []byte("appeltaart is beter dan chocoladetaart"),
		ImmediateResponse: nil,
	}, nil)

	chain := &DefaultChain{filters: filters}

	res, err := chain.ProcessRequestHeaders(ctx, map[string]string{}, false)
	s.Require().Equal(&ProcessingResponse{}, res)
	s.Require().NoError(err)

	res, err = chain.ProcessRequestBody(ctx, []byte("a"))

	s.Require().NoError(err)
	s.Equal(map[string]string{"coolheader": "some other value"}, res.Headers)
	s.Equal([]byte("appeltaart is beter dan chocoladetaart"), res.Body)
}
