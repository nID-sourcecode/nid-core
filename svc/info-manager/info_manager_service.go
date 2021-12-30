package main

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"

	"lab.weave.nl/nid/nid-core/pkg/authtoken"
	errgrpc "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
	"lab.weave.nl/nid/nid-core/svc/info-manager/inforestarter"
	"lab.weave.nl/nid/nid-core/svc/info-manager/models"
	"lab.weave.nl/nid/nid-core/svc/info-manager/proto"
)

// InfoManagerServiceServer infomanager server struct
type InfoManagerServiceServer struct {
	db            *InfoManagerDB
	storageClient objectstorage.Client
	infoRestarter inforestarter.InfoRestarter
	logger        log.LoggerUtility
}

// DefaultSignedURLExpTime is the time the url expires
const DefaultSignedURLExpTime = 1 * time.Minute

// ScriptsUpload is the endpoint for uploading scripts
func (s InfoManagerServiceServer) ScriptsUpload(ctx context.Context, req *proto.ScriptsUploadRequest) (*empty.Empty, error) {
	// Get existing object on script_id
	scriptIDPayload := req.GetScriptId()
	scriptID, err := uuid.FromString(scriptIDPayload)
	if err != nil {
		log.WithError(err).Info("error parsing string uuid to uuid")
		return nil, errgrpc.ErrInternalServer()
	}

	script, err := s.db.ScriptDB.Get(ctx, scriptID)
	if err != nil {
		log.WithError(err).Info("error fetching script from database")
		return nil, errgrpc.ErrInternalServer()
	}

	// Check for existing objects on script
	objects, err := s.storageClient.List(ctx, scriptIDPayload)
	if err != nil {
		log.WithError(err).Info("error listing existing objects from bucket")
		return nil, errgrpc.ErrInternalServer()
	}

	// Create new object for script
	scriptBytes := req.GetScript()
	reader := bytes.NewReader(scriptBytes)

	hash, err := authtoken.Hash(string(scriptBytes))
	if err != nil {
		log.WithError(err).Info("error creating hash for bytes")
		return nil, errgrpc.ErrInternalServer()
	}

	object := objectstorage.Object{
		Key:         strings.Join([]string{scriptIDPayload, hash}, "/"),
		Size:        int64(len(scriptBytes)),
		ContentType: "text/json",
	}

	// Write object to storage client
	err = s.storageClient.Write(ctx, &object, reader, true)
	if err != nil {
		log.WithError(err).Info("error writing new object to bucket")
		return nil, errgrpc.ErrInternalServer()
	}

	// Save new source
	source := &models.ScriptSource{
		ChangeDescription: req.ChangeDescription,
		Checksum:          hash,
		Version:           strconv.Itoa(len(objects) + 1),
		ScriptID:          script.ID,
	}
	err = s.db.ScriptSourceDB.Add(ctx, source)
	if err != nil {
		log.WithError(err).Info("error creating new script source in database")
		return nil, errgrpc.ErrInternalServer()
	}

	if script.Status == models.ScriptStatusActive {
		go s.restartInfoServices()
	}

	return &empty.Empty{}, nil
}

// ScriptsGet endpoint for getting script source on script_id and version
func (s InfoManagerServiceServer) ScriptsGet(ctx context.Context, req *proto.ScriptsGetRequest) (*proto.ScriptsGetRespone, error) {
	// Get existing object on script_id
	scriptIDPayload := req.GetScriptId()
	scriptID, err := uuid.FromString(scriptIDPayload)
	if err != nil {
		log.WithError(err).Info("error parsing string uuid to uuid")
		return nil, errgrpc.ErrInternalServer()
	}

	version := req.GetVersion()

	script, err := s.db.ScriptDB.GetWithScriptSourcesByVersion(ctx, scriptID, version)
	if err != nil {
		log.WithError(err).Info("error searching script with source in database")
		return nil, errgrpc.ErrInternalServer()
	}

	// Check if we found the script source by version
	if script.ScriptSources != nil && len(script.ScriptSources) == 0 {
		return nil, errgrpc.ErrInternalServer()
	}

	url, err := s.storageClient.GetPresignedObjectURL(ctx, strings.Join([]string{scriptIDPayload, script.ScriptSources[0].Checksum}, "/"), http.MethodGet, DefaultSignedURLExpTime)
	if err != nil {
		log.WithError(err).Info("error getting presigned object url from bucket")
		return nil, errgrpc.ErrInternalServer()
	}

	return &proto.ScriptsGetRespone{
		SignedUrl: url,
	}, nil
}

// ScriptsTest endpoint for testing and validating scripts
func (s InfoManagerServiceServer) ScriptsTest(ctx context.Context, req *proto.ScriptsTestRequest) (*empty.Empty, error) {
	return nil, errgrpc.ErrUnimplemented("")
}

func (s *InfoManagerServiceServer) restartInfoServices() {
	err := s.infoRestarter.RestartInfoServices()
	if err != nil {
		s.logger.WithError(err).Error("restarting info services")
	}
}
