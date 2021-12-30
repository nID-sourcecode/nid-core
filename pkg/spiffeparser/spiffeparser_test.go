package spiffeparser

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SpiffeParserTestSuite struct {
	suite.Suite
}

func TestSpiffeParserTestSuite(t *testing.T) {
	suite.Run(t, &SpiffeParserTestSuite{})
}

func (s *SpiffeParserTestSuite) TestParseSpiffeCert() {
	parser := NewDefaultSpiffeParser()
	cert, err := parser.Parse("By=spiffe://cluster.local/ns/foo/sa/httpbin;Hash=sfgiukdjhgrkevsdcn;Subject=\"\";URI=spiffe://cluster.local/ns/foo/sa/sleep")
	s.Require().NoError(err)

	s.Require().NotNil(cert.By)
	s.Require().NotNil(cert.URI)

	s.Require().Equal("cluster.local", cert.By.GetDomain())
	s.Require().Equal("foo", cert.By.GetNamespace())
	s.Require().Equal("httpbin", cert.By.GetServiceAccount())

	s.Require().Equal("cluster.local", cert.URI.GetDomain())
	s.Require().Equal("foo", cert.URI.GetNamespace())
	s.Require().Equal("sleep", cert.URI.GetServiceAccount())
}

func (s *SpiffeParserTestSuite) TestParseSpiffeCertWithDashes() {
	parser := NewDefaultSpiffeParser()
	cert, err := parser.Parse("By=spiffe://cluster-one.local/ns/nid/sa/default;Hash=2fc04659293bf4cb0384e0d1f65ce082a61fde425c5a0349e66c92cc44697736;Subject=\"\";URI=spiffe://cluster-two.local/ns/istio-system/sa/istio-ingressgateway-service-account")
	s.Require().NoError(err)

	s.Require().NotNil(cert.By)
	s.Require().NotNil(cert.URI)

	s.Require().Equal("cluster-one.local", cert.By.GetDomain())
	s.Require().Equal("nid", cert.By.GetNamespace())
	s.Require().Equal("default", cert.By.GetServiceAccount())

	s.Require().Equal("cluster-two.local", cert.URI.GetDomain())
	s.Require().Equal("istio-system", cert.URI.GetNamespace())
	s.Require().Equal("istio-ingressgateway-service-account", cert.URI.GetServiceAccount())
}
