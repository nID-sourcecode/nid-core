package main

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	grpcerrors "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers"
	"lab.weave.nl/nid/nid-core/pkg/utilities/jwt/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	"lab.weave.nl/nid/nid-core/pkg/utilities/password"
	"lab.weave.nl/nid/nid-core/svc/dashboard/proto"
	dashboardscopes "lab.weave.nl/nid/nid-core/svc/dashboard/proto/scopes"
	documentationscopes "lab.weave.nl/nid/nid-core/svc/documentation/proto/scopes"
)

// AuthorizationServiceServer handles authorization for the dashboard
type AuthorizationServiceServer struct {
	stats          *Stats
	metadataHelper headers.MetadataHelper
	jwtClient      *jwt.Client
	db             *DashboardDB
	pwManager      password.IManager
}

// Signin signin a dashboard user
func (a *AuthorizationServiceServer) Signin(ctx context.Context, in *empty.Empty) (*proto.SigninResponseMessage, error) {
	u, p, err := a.metadataHelper.GetBasicAuth(ctx)
	if err != nil {
		log.Extract(ctx).WithError(err).Info("error retrieving basic auth")

		return nil, grpcerrors.ErrInvalidArgument("no basic auth header provided")
	}

	user, err := a.db.UserDB.GetOnEmail(ctx, u)
	if err != nil {
		log.Extract(ctx).WithError(err).WithField("email", u).Info("unexisting user tried to sign in")

		return nil, grpcerrors.ErrInvalidArgument("incorrect username or password")
	}

	matches, err := a.pwManager.ComparePassword(p, user.Password)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to compare password")

		return nil, grpcerrors.ErrInvalidArgument("incorrect username or password")
	}
	if !matches {
		log.Extract(ctx).WithError(err).WithField("user_id", user.ID).Info("password does not match")

		return nil, grpcerrors.ErrInvalidArgument("incorrect username or password")
	}

	customClaims := make(map[string]interface{})
	customClaims["scope"] = append(dashboardscopes.GetAllScopes(), documentationscopes.GetAllScopes()...)
	bearer, err := a.jwtClient.SignToken(customClaims)
	if err != nil {
		log.Extract(ctx).WithError(err).WithField("user_id", user.ID).Error("unable to create signed token")

		return nil, grpcerrors.ErrInternalServer()
	}

	return &proto.SigninResponseMessage{
		Bearer: bearer,
	}, nil
}
