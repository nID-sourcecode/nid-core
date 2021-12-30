// Package git is a wrapper for the gitlab package
package git

import (
	"github.com/xanzy/go-gitlab"
)

// IGitClient wrapper of gitlab client
type IGitClient interface {
	Repositories() IRepositories
	Branches() IBranches
	Tags() ITags
	RepositoryFiles() IRepositoryFiles
}

// IRepositories wrapper interface
type IRepositories interface {
	ListTree(pid interface{}, opt *gitlab.ListTreeOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.TreeNode, *gitlab.Response, error)
}

// IBranches wrapper interface
type IBranches interface {
	ListBranches(pid interface{}, opt *gitlab.ListBranchesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Branch, *gitlab.Response, error)
}

// ITags wrapper interface
type ITags interface {
	ListTags(pid interface{}, opt *gitlab.ListTagsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Tag, *gitlab.Response, error)
}

// IRepositoryFiles wrapper interface
type IRepositoryFiles interface {
	GetRawFile(pid interface{}, fileName string, opt *gitlab.GetRawFileOptions, options ...gitlab.RequestOptionFunc) ([]byte, *gitlab.Response, error)
}

// Client is a gitlab client
type Client struct {
	gitlabClient *gitlab.Client
}

// NewGitClient constructs the gitlab client struct
func NewGitClient(c *gitlab.Client) IGitClient {
	return &Client{
		gitlabClient: c,
	}
}

// Repositories returns the repositories
func (g *Client) Repositories() IRepositories {
	return g.gitlabClient.Repositories
}

// Branches returns the branches type
func (g *Client) Branches() IBranches {
	return g.gitlabClient.Branches
}

// Tags returns the tags type
func (g *Client) Tags() ITags {
	return g.gitlabClient.Tags
}

// RepositoryFiles returns the repositoryFiles type
func (g *Client) RepositoryFiles() IRepositoryFiles {
	return g.gitlabClient.RepositoryFiles
}
