package auditlog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2/mock"
)

type AuditLogTestSuite struct {
	suite.Suite
	filter     *Filter
	loggerMock *mock.LoggerUtility
}

func (s *AuditLogTestSuite) SetupTest() {
	s.loggerMock = &mock.LoggerUtility{}
	s.filter = &Filter{logger: s.loggerMock}
}

func TestAuditLogTestSuite(t *testing.T) {
	suite.Run(t, &AuditLogTestSuite{})
}

func (s *AuditLogTestSuite) TestLogRequest_InvalidJWT() {
	headers := map[string]string{
		"authorization": "Bearer !NVAL!DBASE64",
		":method":       "GET",
		":path":         "/lekker/testen",
		":authority":    "something.nl",
		"x-request-id":  "ab23c",
	}

	loggerMock2 := &mock.LoggerUtility{}
	s.loggerMock.On("WithFields", log.Fields{
		"token":       "no valid token",
		"url":         "something.nl/lekker/testen",
		"body":        "",
		"http_method": "GET",
		"request_id":  "ab23c",
	}).Return(loggerMock2)

	loggerMock2.On("Info", "received request")

	ctx := context.TODO()

	res, err := s.filter.OnHTTPRequest(ctx, nil, headers)

	s.Require().NoError(err)
	s.Require().Nil(res)

	s.loggerMock.AssertExpectations(s.T())
	loggerMock2.AssertExpectations(s.T())
}

func (s *AuditLogTestSuite) TestLogRequest_ValidJWT() {
	ctx := context.TODO()

	loggerMock2 := &mock.LoggerUtility{}
	s.loggerMock.On("WithFields", log.Fields{
		"token": map[string]interface{}{
			"sub":   "1234567890",
			"name":  "John Doe",
			"iat":   float64(1516239022),
			"scope": "test",
		},
		"url":         "something.nl/lekker/testen",
		"body":        `{"contact":"nadya business"}`,
		"http_method": "POST",
		"request_id":  "ab23c",
	}).Return(loggerMock2)

	loggerMock2.On("Info", "received request")

	loggerMock3 := &mock.LoggerUtility{}
	s.loggerMock.On("WithFields", log.Fields{
		"status_code": "400",
		"request_id":  "ab23c",
	}).Return(loggerMock3)

	loggerMock3.On("Info", "received response")

	requestHeaders := map[string]string{
		"authorization": "Bearer header.ewogICJzdWIiOiAiMTIzNDU2Nzg5MCIsCiAgIm5hbWUiOiAiSm9obiBEb2UiLAogICJpYXQiOiAxNTE2MjM5MDIyLAogICJzY29wZSI6ICJ0ZXN0Igp9Cg.signature",
		":method":       "POST",
		":path":         "/lekker/testen",
		":authority":    "something.nl",
		"x-request-id":  "ab23c",
	}
	res, err := s.filter.OnHTTPRequest(ctx, []byte(`{"contact":"nadya business"}`), requestHeaders)

	s.Require().NoError(err)
	s.Require().Nil(res)

	responseHeaders := map[string]string{
		":status": "400",
	}

	res, err = s.filter.OnHTTPResponse(ctx, nil, responseHeaders)

	s.Require().NoError(err)
	s.Require().Nil(res)

	s.loggerMock.AssertExpectations(s.T())
	loggerMock2.AssertExpectations(s.T())
	loggerMock3.AssertExpectations(s.T())
}
