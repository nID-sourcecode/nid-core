package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xanzy/go-gitlab"
	"google.golang.org/grpc/codes"

	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
	mockobjectstorage "lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage/mock"
	gitmock "lab.weave.nl/nid/nid-core/svc/documentation/packages/git/mock"
	documentationPB "lab.weave.nl/nid/nid-core/svc/documentation/proto"
)

type DocumentationServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	documentationServiceClient *DocumentationServiceServer
}

func (s *DocumentationServiceTestSuite) SetupTest() {
	conf := &documentationConfig{
		GitlabBaseURL:           "https://lab.weave.nl/",
		GitlabAccessToken:       "myrandomaccesstoken",
		GitlabProjectIdentifier: "000",
		ObjectStorage: objectStorageConfig{
			Bucket: "nid-swagger-files",
		},
	}

	s.Ctx = context.Background()

	// Init mocked git client
	mockedGitClient := &gitmock.MockedClient{
		RepositoryFilesClient: &gitmock.MockedRepositoryFilesClient{},
		BranchesClient:        &gitmock.MockedBranchesClient{},
		TagsClient:            &gitmock.MockedTagsClient{},
		RepositoriesClient:    &gitmock.MockedRepoClient{},
	}

	// Default GetRawFile
	mockedGitClient.RepositoryFilesClient.On("GetRawFile", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]byte("# Hello, world!"), &gitlab.Response{
			Response: &http.Response{
				StatusCode: http.StatusOK,
			},
		}, nil)

	// Default ListTags
	mockedGitClient.TagsClient.On("ListTags", mock.Anything, mock.Anything, mock.Anything).
		Return([]*gitlab.Tag{
			{Name: "v1.0.0"},
			{Name: "v2.0.0"},
		}, &gitlab.Response{
			Response: &http.Response{
				StatusCode: http.StatusOK,
			},
		}, nil)

	// Default ListBranches
	mockedGitClient.BranchesClient.On("ListBranches", mock.Anything, mock.Anything, mock.Anything).
		Return([]*gitlab.Branch{
			{Name: "master"},
			{Name: "prod"},
			{Name: "pre-release"},
		}, &gitlab.Response{
			Response: &http.Response{
				StatusCode: http.StatusOK,
			},
		}, nil)

	// Default ListTree
	ref := "master"
	path := "documentation"
	mockedGitClient.RepositoriesClient.On("ListTree", mock.Anything, &gitlab.ListTreeOptions{
		Ref:  &ref,
		Path: &path,
	}, mock.Anything).
		Return([]*gitlab.TreeNode{
			{Name: "auth-service-flow.md", Path: "documentation/auth-service-flow.md", Type: "blob"},
			{Name: "other-service-flow.md", Path: "documentation/other-service-flow.md", Type: "blob"},
		}, &gitlab.Response{
			Response: &http.Response{
				StatusCode: http.StatusOK,
			},
		}, nil)

	// Init mocked bucket client
	mockedObjectStorage := &mockobjectstorage.Client{}

	// Default ListFiles
	mockedObjectStorage.On("List", mock.Anything, "auth/master").
		Return([]objectstorage.Object{
			{
				Key: "auth/master/auth.swagger.json",
			},
		}, nil)

	// Default ListFiles
	signed := "supersecretsignedurl"
	mockedObjectStorage.On("GetPresignedObjectURL", s.Ctx, "auth/master/auth.swagger.json", http.MethodGet, mock.Anything).Return(signed, nil)

	s.documentationServiceClient = &DocumentationServiceServer{
		conf:          conf,
		git:           mockedGitClient,
		storageClient: mockedObjectStorage,
	}
}

func (s *DocumentationServiceTestSuite) TestGetFileMarkdownFile() {
	response, err := s.documentationServiceClient.GetFile(s.Ctx, &documentationPB.GetFileRequest{
		Ref:         "master",
		FilePath:    "documentation/index.md",
		ServiceName: "auth",
	})
	s.NoError(err)
	s.Equal(response.Content, "# Hello, world!")
	s.Len(response.SwaggerFiles, 1)
	s.Equal(response.SwaggerFiles[0].Name, "auth/master/auth.swagger.json")
	s.Equal(response.SwaggerFiles[0].SignedUrl, "supersecretsignedurl")
}

func (s *DocumentationServiceTestSuite) TestGetFileNoService() {
	response, err := s.documentationServiceClient.GetFile(s.Ctx, &documentationPB.GetFileRequest{
		Ref:      "master",
		FilePath: "documentation/index.md",
	})

	s.NoError(err)
	s.Equal(response.Content, "# Hello, world!")
	s.Empty(response.SwaggerFiles)
}

func (s *DocumentationServiceTestSuite) TestErrorGetFileUnauthorized() {
	mockedUnit := &gitmock.MockedClient{
		RepositoryFilesClient: &gitmock.MockedRepositoryFilesClient{},
	}
	mockedUnit.RepositoryFilesClient.On("GetRawFile", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]byte("# Hello, world!"), &gitlab.Response{
			Response: &http.Response{
				StatusCode: 401,
			},
		}, nil)
	s.documentationServiceClient.git = mockedUnit

	_, err := s.documentationServiceClient.GetFile(s.Ctx, &documentationPB.GetFileRequest{
		Ref:      "master",
		FilePath: "documentation/index.md",
	})
	s.Error(err)
	s.VerifyStatusError(err, codes.Unauthenticated)
}

func (s *DocumentationServiceTestSuite) TestListRepositoryRefs() {
	response, err := s.documentationServiceClient.ListRepositoryRefs(s.Ctx, nil)
	s.Nil(err)
	s.NotNil(response)
	s.Len(response.Refs, 5)
}

func (s *DocumentationServiceTestSuite) TestListDirectoryFiles() {
	response, err := s.documentationServiceClient.ListDirectoryFiles(s.Ctx, &documentationPB.ListDirectoryFilesRequest{
		Ref:      "master",
		FilePath: "documentation",
	})
	s.Nil(err)
	s.NotNil(response)
	s.Len(response.Files, 2)
}

func TestDocumentationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentationServiceTestSuite))
}
