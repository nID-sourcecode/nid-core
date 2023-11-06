package grpctesthelpers

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcTestSuite is testsuite with a context that can be used for grpc calls
type GrpcTestSuite struct {
	ExtendedTestSuite
	Ctx                       context.Context // nolint:containedctx
	ServerTransportStreamMock *ServerTransportStreamMock
}

// SetupTest sets up tests in the GRPContextSuite
func (suite *GrpcTestSuite) SetupTest() {
	expectedStream := &ServerTransportStreamMock{
		Headers: make(map[string][]string),
	}
	suite.ServerTransportStreamMock = expectedStream
	suite.ServerTransportStreamMock.On("SetHeader", mock.AnythingOfType("metadata.MD")).Return(nil)
	suite.Ctx = grpc.NewContextWithServerTransportStream(context.Background(), grpc.ServerTransportStream(expectedStream))
}

// AddValToCtx sets given val for given key on the context
func (suite *GrpcTestSuite) AddValToCtx(key, val interface{}) {
	suite.Ctx = context.WithValue(suite.Ctx, key, val)
}

// VerifyStatusError verify if given err has given status
func (suite *GrpcTestSuite) VerifyStatusError(err error, expected codes.Code) {
	suite.Error(err, "Expected error not to be nil")
	statusErr, ok := status.FromError(err)
	if !ok {
		suite.Fail("Unknown statuscode: %d", status.Code(err))
	}
	suite.Equal(expected, statusErr.Code(), "Expected: %v but got %v(\"%v\")", expected, statusErr.Code(), statusErr.Err())
}

// VerifyCreatedHeader verify if given header has given value
func (suite *GrpcTestSuite) VerifyCreatedHeader(header, val string) {
	curHeader, ok := suite.ServerTransportStreamMock.Headers[header]
	suite.True(ok)
	suite.Equal(1, len(curHeader), "Expected status created")
	suite.Equal(val, curHeader[0], "Incorrect status for header")
}

// VerifyHeaderPresent verify if given header is present
func (suite *GrpcTestSuite) VerifyHeaderPresent(header string) {
	curHeader, ok := suite.ServerTransportStreamMock.Headers[header]
	suite.True(ok)
	suite.Equal(1, len(curHeader), "Expected status created")
}

// RetrieveHeader retrieve the header from a context
func (suite *GrpcTestSuite) RetrieveHeader(header string) string {
	curHeader, ok := suite.ServerTransportStreamMock.Headers[header]
	suite.True(ok)
	suite.Equal(1, len(curHeader), "Expected status created")
	return curHeader[0]
}

// VerifyHeaderNotPresent verifies that a given header is not present
func (suite *GrpcTestSuite) VerifyHeaderNotPresent(header string) {
	test, ok := suite.ServerTransportStreamMock.Headers[header]
	suite.False(ok)
	suite.Nil(test)
}
