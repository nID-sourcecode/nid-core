// Package mock is the mock package of the parent git package
package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/xanzy/go-gitlab"

	"github.com/nID-sourcecode/nid-core/svc/documentation/packages/git"
)

// MockedClient mocked gitlab client
type MockedClient struct {
	mock.Mock

	RepositoriesClient    *MockedRepoClient
	BranchesClient        *MockedBranchesClient
	TagsClient            *MockedTagsClient
	RepositoryFilesClient *MockedRepositoryFilesClient
}

// Repositories returns the mocked repositories client
func (g *MockedClient) Repositories() git.IRepositories {
	return g.RepositoriesClient
}

// Branches returns the mocked branches client
func (g *MockedClient) Branches() git.IBranches {
	return g.BranchesClient
}

// Tags returns the mocked tags client
func (g *MockedClient) Tags() git.ITags {
	return g.TagsClient
}

// RepositoryFiles returns the mocked repository files client
func (g *MockedClient) RepositoryFiles() git.IRepositoryFiles {
	return g.RepositoryFilesClient
}

// MockedRepoClient mocked client struct
type MockedRepoClient struct {
	mock.Mock
}

// ListTree mocks the list tree call
func (m *MockedRepoClient) ListTree(pid interface{}, opt *gitlab.ListTreeOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.TreeNode, *gitlab.Response, error) {
	args := m.MethodCalled("ListTree", pid, opt, options)

	return args.Get(0).([]*gitlab.TreeNode), args.Get(1).(*gitlab.Response), args.Error(2) //nolint:gomnd
}

// MockedBranchesClient mocked branches client
type MockedBranchesClient struct {
	mock.Mock
}

// ListBranches mocks the list branches call
func (m *MockedBranchesClient) ListBranches(pid interface{}, opt *gitlab.ListBranchesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Branch, *gitlab.Response, error) {
	args := m.MethodCalled("ListBranches", pid, opt, options)

	return args.Get(0).([]*gitlab.Branch), args.Get(1).(*gitlab.Response), args.Error(2) //nolint:gomnd
}

// MockedTagsClient mocked tags client
type MockedTagsClient struct {
	mock.Mock
}

// ListTags mocks the list tags call
func (m *MockedTagsClient) ListTags(pid interface{}, opt *gitlab.ListTagsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Tag, *gitlab.Response, error) {
	args := m.MethodCalled("ListTags", pid, opt, options)

	return args.Get(0).([]*gitlab.Tag), args.Get(1).(*gitlab.Response), args.Error(2) //nolint:gomnd
}

// MockedRepositoryFilesClient mocked repository files client
type MockedRepositoryFilesClient struct {
	mock.Mock
}

// GetRawFile mocks get raw file call
func (m *MockedRepositoryFilesClient) GetRawFile(pid interface{}, _ string, opt *gitlab.GetRawFileOptions, options ...gitlab.RequestOptionFunc) ([]byte, *gitlab.Response, error) {
	args := m.MethodCalled("GetRawFile", pid, opt, options)

	return args.Get(0).([]byte), args.Get(1).(*gitlab.Response), args.Error(2) //nolint:gomnd
}
