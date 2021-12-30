package main

import (
	errgrpc "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/svc/auth/models"
	pb "lab.weave.nl/nid/nid-core/svc/auth/proto"
)

const (
	noPreload                        models.PreloadOption = 0
	preloadRequiredAndOptionalScopes models.PreloadOption = 1
	preloadAll                       models.PreloadOption = 2
)

// accessModelExists checks if supplied access_models are found in required and optional session access_models
func accessModelExists(supplied, required, optional []*models.AccessModel) error {
	for _, s := range supplied {
		found := false
		for _, r := range required {
			if r.ID == s.ID {
				return errgrpc.ErrFailedPrecondition("payload should not contain id's of required_access_models these are accepted automatically")
			}
		}
		for _, o := range optional {
			if o.ID == s.ID {
				found = true

				break
			}
		}
		if !found {
			return errgrpc.ErrNotFound("one ore more access_models from payload are not found in session's optional access_models")
		}
	}

	return nil
}

// sessionToResponse creates session_response message from model
func sessionToResponse(m *models.Session) *pb.SessionResponse {
	s := pb.SessionResponse{}

	s.Id = m.ID.String()
	s.State = getState(m.State)

	if m.Client != nil {
		c := pb.Client{}
		c.Id = m.Client.ID.String()
		c.Name = m.Client.Name
		c.Logo = m.Client.Logo
		c.Icon = m.Client.Icon
		c.Color = m.Client.Color
		s.Client = &c
	}

	if m.Audience != nil {
		a := pb.Audience{}
		a.Id = m.Audience.ID.String()
		a.Audience = m.Audience.Audience
		a.Namespace = m.Audience.Namespace
		s.Audience = &a
	}

	s.RequiredAccessModels = accessModelToResponse(m.RequiredAccessModels)
	s.OptionalAccessModels = accessModelToResponse(m.OptionalAccessModels)
	s.AcceptedAccessModels = accessModelToResponse(m.AcceptedAccessModels)

	return &s
}

// accessModelToResponse creates access_model message from model
func accessModelToResponse(m []*models.AccessModel) []*pb.AccessModel {
	var accessModels []*pb.AccessModel

	if len(m) > 0 {
		for _, accessModel := range m {
			accessModels = append(accessModels, &pb.AccessModel{
				Id:          accessModel.ID.String(),
				Name:        accessModel.Name,
				Hash:        accessModel.Hash,
				Description: accessModel.Description,
			})
		}

		return accessModels
	}

	return nil
}

// getState parses model enum to proto enum
func getState(s models.SessionState) pb.SessionState {
	switch s {
	case models.SessionStateClaimed:
		return pb.SessionState_CLAIMED
	case models.SessionStateAccepted:
		return pb.SessionState_ACCEPTED
	case models.SessionStateRejected:
		return pb.SessionState_REJECTED
	case models.SessionStateCodeGranted:
		return pb.SessionState_CODE_GRANTED
	case models.SessionStateTokenGranted:
		return pb.SessionState_TOKEN_GRANTED
	case models.SessionStateUnclaimed:
		return pb.SessionState_UNCLAIMED
	default:
		return pb.SessionState_UNSPECIFIED
	}
}
