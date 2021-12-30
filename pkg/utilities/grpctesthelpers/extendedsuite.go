package grpctesthelpers

import (
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
)

// ExtendedTestSuite is an extension of the testify suite with additional methods
type ExtendedTestSuite struct {
	suite.Suite
}

// NoErrorWithFail checks if the error is nil and fails the test otherwise
// Deprecated: use .Require().NoError(...) instead https://godoc.org/github.com/stretchr/testify/suite#Suite.Require
func (s *ExtendedTestSuite) NoErrorWithFail(err error, args ...interface{}) {
	if err != nil {
		s.FailNow("Unexpected error:"+err.Error(), args...)
	}
}

// ProtoEqual checks if two proto messages are equal using the proto.Equal function
func (s *ExtendedTestSuite) ProtoEqual(expected, actual proto.Message) bool {
	return s.True(proto.Equal(expected, actual), "Not equal\nExpected:\t%+v\nActual:\t\t%+v\n", expected, actual)
}

// RequireProtoEqual checks if two proto messages are equal using the proto.Equal function and using the Require context for the suite
func (s *ExtendedTestSuite) RequireProtoEqual(expected, actual proto.Message) {
	s.Require().True(s.ProtoEqual(expected, actual))
}
