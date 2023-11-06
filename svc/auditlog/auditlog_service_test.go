package main

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"                       //nolint:gomodguard
	logrustest "github.com/sirupsen/logrus/hooks/test" //nolint:gomodguard
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	pb "github.com/nID-sourcecode/nid-core/svc/auditlog/proto"
)

type AuditLogServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	auditLogServer pb.AuditlogServiceServer
	hook           *logrustest.Hook
}

func (s *AuditLogServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	logger, hook := logrustest.NewNullLogger()

	s.auditLogServer = &AuditLogServiceServer{logger: log.CustomLogrusUtility(logrus.NewEntry(logger))}
	s.hook = hook
}

func TestAuditLogServiceTestSuite(t *testing.T) {
	suite.Run(t, &AuditLogServiceTestSuite{})
}

func (s *AuditLogServiceTestSuite) TestInvalidJWT() {
	in := &pb.Request{
		Auth:       "Bearer !NVAL!DBASE64",
		Url:        "/lekker/testen",
		Body:       "{}",
		HttpMethod: "GET",
	}

	_, err := s.auditLogServer.LogRequest(context.Background(), in)
	s.NoError(err)

	// There is an error log but it should not be logged on the dedicated logger
	s.Require().Equal(1, len(s.hook.AllEntries()))
	entry := s.hook.LastEntry()

	s.Equal("token not parsable", entry.Data["token"])
}

func (s *AuditLogServiceTestSuite) TestValidJWT() {
	in := &pb.Request{
		Auth:       "ewogICJzdWIiOiAiMTIzNDU2Nzg5MCIsCiAgIm5hbWUiOiAiSm9obiBEb2UiLAogICJpYXQiOiAxNTE2MjM5MDIyLAogICJzY29wZSI6ICJ0ZXN0Igp9Cg",
		Url:        "/lekker/testen",
		Body:       "{}",
		HttpMethod: "GET",
		RequestId:  "ab23c",
	}

	_, err := s.auditLogServer.LogRequest(s.Ctx, in)
	s.Require().NoError(err)

	// There is an error log but it should not be logged on the dedicated logger
	s.Require().Equal(1, len(s.hook.AllEntries()))
	entry := s.hook.LastEntry()
	s.Equal(logrus.InfoLevel, entry.Level)
	s.Equal(in.HttpMethod, entry.Data["http_method"])
	s.Equal(in.Url, entry.Data["url"])
	s.Equal(in.Body, entry.Data["body"])
	s.Equal(in.RequestId, entry.Data["request_id"])
	s.Equal("test", (entry.Data["token"].(map[string]interface{}))["scope"])
}

func (s *AuditLogServiceTestSuite) TestValidateRequest() {
	valid := pb.Request{
		Auth:       "iets",
		Url:        "http://test",
		Body:       "{}",
		HttpMethod: "GET",
	}
	s.NoError(valid.Validate())

	emptyToken := pb.Request{
		Url:  "http://test",
		Body: "{}",
	}
	s.Error(emptyToken.Validate())

	emptyURL := pb.Request{
		Auth:       "iets",
		Body:       "{}",
		HttpMethod: "GET",
	}
	s.Error(emptyURL.Validate())
}

func (s *AuditLogServiceTestSuite) TestValidateResponse() {
	response := pb.Response{
		RequestId:  "aabbcc234",
		StatusCode: 403,
	}
	s.NoError(response.Validate())
}

func (s *AuditLogServiceTestSuite) TestLogResponse() {
	in := pb.Response{
		RequestId:  "aabbcc234",
		StatusCode: 403,
	}

	_, err := s.auditLogServer.LogResponse(s.Ctx, &in)
	s.Require().NoError(err)

	s.Require().Equal(1, len(s.hook.AllEntries()))
	entry := s.hook.LastEntry()
	s.Equal(logrus.InfoLevel, entry.Level)
	s.Equal(in.RequestId, entry.Data["request_id"])
	s.Equal(in.StatusCode, entry.Data["status_code"])
}
