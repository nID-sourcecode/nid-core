//go:build integration || to || files
// +build integration to files

package integration

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/grpckeys"
	documentationPB "lab.weave.nl/nid/nid-core/svc/documentation/proto"
)

type DocumentationIntegrationTestSuite struct {
	BaseTestSuite
	dashboardConfig
}

func (s *DocumentationIntegrationTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.Require().NoError(envconfig.Init(&s.dashboardConfig), "unable to initialise dashboard environment config")

	bearer, err := signinDashboard(s.ctx, s.dashboardAuthClient, s.dashboardConfig.DefaultUserEmail, s.dashboardConfig.DefaultUserPass)
	s.Require().NoError(err)

	s.ctx = metadata.AppendToOutgoingContext(s.ctx, grpckeys.AuthorizationKey.String(), fmt.Sprintf("Bearer %s", bearer))
}

// func TestDocumentationIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(DocumentationIntegrationTestSuite))
// }

func (s *DocumentationIntegrationTestSuite) TestGetRefs() {
	resp, err := s.documentationClient.ListRepositoryRefs(s.ctx, &empty.Empty{})
	s.Require().NoError(err)
	s.NotEmpty(resp.Refs)
}

func (s *DocumentationIntegrationTestSuite) TestGetRefsFindTags() {
	resp, err := s.documentationClient.ListRepositoryRefs(s.ctx, &empty.Empty{})
	s.Require().NoError(err)
	s.NotEmpty(resp.Refs)
	refs := resp.GetRefs()
	var tags []*documentationPB.Ref
	for _, r := range refs {
		if r.GetType() == documentationPB.RefType_TAG {
			tags = append(tags, r)
		}
	}
	s.NotEmpty(tags)
	s.Greater(len(tags), 0)
}

func (s *DocumentationIntegrationTestSuite) TestListDirectoryFiles() {
	resp, err := s.documentationClient.ListDirectoryFiles(s.ctx, &documentationPB.ListDirectoryFilesRequest{
		Ref:      "master",
		FilePath: "documentation",
	})
	s.Require().NoError(err)
	s.NotEmpty(resp.Files)
}

func (s *DocumentationIntegrationTestSuite) TestListDirectoryFilesFindDocumentation() {
	resp, err := s.documentationClient.ListDirectoryFiles(s.ctx, &documentationPB.ListDirectoryFilesRequest{
		Ref:      "master",
		FilePath: "documentation",
	})
	s.Require().NoError(err)
	s.NotEmpty(resp.Files)
	files := resp.GetFiles()
	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.GetName())
	}
	s.Contains(fileNames, "auth-service-flow.md")
}

func (s *DocumentationIntegrationTestSuite) TestGetAuthServiceFlowDoc() {
	resp, err := s.documentationClient.GetFile(s.ctx, &documentationPB.GetFileRequest{
		Ref:         "master",
		FilePath:    "documentation/auth-service-flow.md",
		ServiceName: "auth",
	})
	s.Require().NoError(err, "when calling GetFile on the documentationClient")
	s.NotEmpty(resp.Content)
}

func (s *DocumentationIntegrationTestSuite) TestGetFileErrorInvalidExtension() {
	_, err := s.documentationClient.GetFile(s.ctx, &documentationPB.GetFileRequest{
		Ref:         "master",
		FilePath:    "documentation/auth-service-flow.pdf",
		ServiceName: "auth",
	})
	s.Require().Error(err)
	s.EqualError(err, "rpc error: code = InvalidArgument desc = invalid GetFileRequest.FilePath: value does not have suffix \".md\"")
}

func (s *DocumentationIntegrationTestSuite) TestGetFileProtected() {
	invalidBearerCTX := metadata.AppendToOutgoingContext(context.Background(), grpckeys.AuthorizationKey.String(), "Bearer aa.bb.cc")

	_, err := s.documentationClient.GetFile(invalidBearerCTX, &documentationPB.GetFileRequest{
		Ref:         "master",
		FilePath:    "documentation/auth-service-flow.pdf",
		ServiceName: "auth",
	})
	s.VerifyStatusError(err, codes.Unauthenticated)
}

func (s *DocumentationIntegrationTestSuite) TestListDirectoryFilesRequestProtected() {
	invalidBearerCTX := metadata.AppendToOutgoingContext(context.Background(), grpckeys.AuthorizationKey.String(), "Bearer aa.bb.cc")

	_, err := s.documentationClient.ListDirectoryFiles(invalidBearerCTX, &documentationPB.ListDirectoryFilesRequest{
		Ref:      "master",
		FilePath: "documentation",
	})
	s.VerifyStatusError(err, codes.Unauthenticated)
}

func (s *DocumentationIntegrationTestSuite) TestListRepositoryFilesProtected() {
	invalidBearerCTX := metadata.AppendToOutgoingContext(context.Background(), grpckeys.AuthorizationKey.String(), "Bearer aa.bb.cc")

	_, err := s.documentationClient.ListRepositoryRefs(invalidBearerCTX, &empty.Empty{})
	s.VerifyStatusError(err, codes.Unauthenticated)
}
