package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xanzy/go-gitlab"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/servicebase"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
	"lab.weave.nl/nid/nid-core/svc/documentation/packages/git"
	pb "lab.weave.nl/nid/nid-core/svc/documentation/proto"
)

// DefaultSignedURLExpTime default signed url expiration time
const DefaultSignedURLExpTime = 15 * time.Minute

// DocumentationServiceServer the documentation service client
type DocumentationServiceServer struct {
	servicebase.Service
	stats         *Stats
	conf          *documentationConfig
	git           git.IGitClient
	storageClient objectstorage.Client
}

// GetFile gets a file within the repository on given path
func (s *DocumentationServiceServer) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	// Search raw file on filepath
	ref := req.GetRef()
	service := req.GetServiceName()
	bytes, resp, err := s.git.RepositoryFiles().GetRawFile(s.conf.GitlabProjectIdentifier, req.GetFilePath(), &gitlab.GetRawFileOptions{
		Ref: &ref,
	}, nil)
	returnCode := s.DocumentationErrHandler(ctx, resp, err)
	if returnCode != codes.OK {
		return nil, status.New(returnCode, "unable to get raw file").Err()
	}

	response := &pb.GetFileResponse{Content: string(bytes)}

	if service != "" {
		files, err := s.GetSwaggerFiles(ctx, service, ref)
		if err != nil {
			log.Extract(ctx).WithError(err).WithFields(log.Fields{"ref": ref, "service": service}).Error("getting swagger files failed")
			return nil, grpcerrors.ErrInternalServer()
		}
		response.SwaggerFiles = files
	}

	return response, nil
}

// ListDirectoryFiles lists markdown files within given directory path
func (s *DocumentationServiceServer) ListDirectoryFiles(ctx context.Context, req *pb.ListDirectoryFilesRequest) (*pb.ListDirectoryFilesResponse, error) {
	// Fetch tree for filepath
	ref := req.GetRef()
	filePath := req.GetFilePath()
	treeItems, resp, err := s.git.Repositories().ListTree(s.conf.GitlabProjectIdentifier, &gitlab.ListTreeOptions{
		Path: &filePath,
		Ref:  &ref,
	}, nil)
	code := s.DocumentationErrHandler(ctx, resp, err)
	if code != codes.OK {
		return nil, status.New(code, "unable to list tree").Err()
	}

	// Create list of found markdown files in tree
	response := pb.ListDirectoryFilesResponse{}
	for _, treeFile := range treeItems {
		if strings.HasSuffix(treeFile.Name, ".md") {
			response.Files = append(response.Files, &pb.File{
				Name:      treeFile.Name,
				Extension: ".md",
				Path:      treeFile.Path,
				Type:      treeFile.Type,
			})
		}
	}

	return &response, nil
}

// ListRepositoryRefs fetches branches and tags for repository
func (s *DocumentationServiceServer) ListRepositoryRefs(ctx context.Context, req *empty.Empty) (*pb.ListRepositoryRefsResponse, error) {
	// Fetch tags for git project id
	tags, resp, err := s.git.Tags().ListTags(s.conf.GitlabProjectIdentifier, nil, nil)
	code := s.DocumentationErrHandler(ctx, resp, err)
	if code != codes.OK {
		return nil, status.New(code, "unable to list tags").Err()
	}

	// Fetch branches for git project id
	branches, resp, err := s.git.Branches().ListBranches(s.conf.GitlabProjectIdentifier, nil, nil)
	code = s.DocumentationErrHandler(ctx, resp, err)
	if code != codes.OK {
		return nil, errors.Wrap(err, "unable to list branches")
	}

	// Create ref response
	response := pb.ListRepositoryRefsResponse{}
	for _, tag := range tags {
		response.Refs = append(response.Refs, &pb.Ref{
			Name: tag.Name,
			Type: pb.RefType_TAG,
		})
	}
	for _, branch := range branches {
		response.Refs = append(response.Refs, &pb.Ref{
			Name: branch.Name,
			Type: pb.RefType_BRANCH,
		})
	}

	return &response, nil
}

// DocumentationErrHandler handles common errors for gitlab package calls
func (s *DocumentationServiceServer) DocumentationErrHandler(ctx context.Context, resp *gitlab.Response, err error) codes.Code {
	if resp != nil && resp.StatusCode >= http.StatusBadRequest {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return codes.NotFound
		case http.StatusUnauthorized:
			return codes.Unauthenticated
		}
	}

	if err != nil {
		log.Extract(ctx).WithFields(log.Fields{
			"project_id": s.conf.GitlabProjectIdentifier,
		}).WithError(err).Error("Getting requested doc file failed")

		return codes.Internal
	}

	return codes.OK
}

// GetSwaggerFiles search for service swagger documentation files for given service namespace
func (s *DocumentationServiceServer) GetSwaggerFiles(ctx context.Context, service, ref string) ([]*pb.SwaggerFile, error) {
	log := log.Extract(ctx).WithFields(log.Fields{"service": service, "ref": ref, "bucket": s.conf.ObjectStorage.Bucket})
	var files []*pb.SwaggerFile

	// Check if bucket has files with prefix
	prefix := fmt.Sprintf("%s/%s", service, ref)
	objects, err := s.storageClient.List(ctx, prefix)
	if err != nil {
		log.WithError(err).Error("listing swagger files from bucket failed")
		return nil, grpcerrors.ErrInternalServer()
	}

	// Sign url for files
	for i := range objects {
		if !strings.HasSuffix(objects[i].Key, ".json") {
			continue
		}
		singedURL, err := s.storageClient.GetPresignedObjectURL(ctx, objects[i].Key, http.MethodGet, DefaultSignedURLExpTime)
		if err != nil {
			log.WithError(err).WithField("object_key", objects[i].Key).Error("getting signed url for swagger doc failed")
			return nil, grpcerrors.ErrInternalServer()
		}
		files = append(files, &pb.SwaggerFile{
			Name:      objects[i].Key,
			SignedUrl: singedURL,
		})
	}

	return files, nil
}
