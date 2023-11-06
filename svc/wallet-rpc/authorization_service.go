// Package wallet-rpc wallet implementing rpc client
package main

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lestrrat-go/jwx/jwt/openid"

	"github.com/nID-sourcecode/nid-core/pkg/authtoken"
	pw "github.com/nID-sourcecode/nid-core/pkg/password"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	grpcerrors "github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/headers"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/jwt/v2"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	authscopes "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto/scopes"
	"github.com/nID-sourcecode/nid-core/svc/wallet-gql/models"
	"github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
)

// AuthorizationServer handles authorization for the dashboard
type AuthorizationServer struct {
	stats          *Stats
	metadataHelper headers.MetadataHelper
	jwtClient      *jwt.Client
	db             *WalletDB
	pwManager      pw.IManager
}

var (
	// ErrIncorrectUsernameOrPassword is returned if the username or password is incorrect
	ErrIncorrectUsernameOrPassword = grpcerrors.ErrInvalidArgument("incorrect username or password")
	// ErrInternal is returned if something goes wrong internally
	ErrInternal = grpcerrors.ErrInternalServer()
)

const (
	walletGQLScope  = "wallet_gql"
	codeBitLength   = 16
	secretBitLength = 16
)

// SignIn signs in a device
func (a *AuthorizationServer) SignIn(ctx context.Context, _ *emptypb.Empty) (*proto.SignInResponse, error) {
	code, secret, err := a.metadataHelper.GetBasicAuth(ctx)
	if err != nil {
		return nil, grpcerrors.ErrInvalidArgument(errors.Wrap(err, "retrieving basic auth").Error())
	}

	device, err := a.db.DeviceDB.GetByCode(code, true)
	if err != nil {
		return nil, ErrIncorrectUsernameOrPassword
	}

	matches, err := a.pwManager.ComparePassword(secret, device.Secret)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to compare secret")

		return nil, ErrInternal
	}
	if !matches {
		return nil, ErrIncorrectUsernameOrPassword
	}

	customClaims := make(map[string]interface{})
	customClaims[openid.SubjectKey] = device.User.Pseudonym

	customClaims["scope"] = []string{
		authscopes.AuthClaim,
		authscopes.AuthAccept,
		authscopes.AuthReject,
		walletGQLScope,
	}

	bearer, err := a.jwtClient.SignToken(customClaims)
	if err != nil {
		log.Extract(ctx).WithError(err).WithField("pseudonym", device.User.Pseudonym).Error("unable to create signed token")

		return nil, ErrInternal
	}

	return &proto.SignInResponse{
		Bearer: bearer,
	}, nil
}

// RegisterDevice registers a device and generates a code and secret for it
func (a *AuthorizationServer) RegisterDevice(ctx context.Context, _ *emptypb.Empty) (*proto.RegisterDeviceResponse, error) {
	bsn, password, err := a.metadataHelper.GetBasicAuth(ctx)
	if err != nil {
		return nil, grpcerrors.ErrInvalidArgument(errors.Wrap(err, "retrieving basic auth").Error())
	}

	user, err := a.db.UserDB.GetByBsn(bsn)
	if err != nil {
		return nil, ErrIncorrectUsernameOrPassword
	}

	matches, err := a.pwManager.ComparePassword(password, user.Password)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to compare secret")

		return nil, ErrInternal
	}
	if !matches {
		return nil, ErrIncorrectUsernameOrPassword
	}

	code, err := authtoken.NewToken(codeBitLength)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("creating device code")

		return nil, ErrInternal
	}

	secret, err := authtoken.NewToken(secretBitLength)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("creating device secret")

		return nil, ErrInternal
	}
	hashedSecret, err := a.pwManager.GenerateHash(secret)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("generating device secret hash")

		return nil, ErrInternal
	}

	newDevice := models.Device{
		Code:   code,
		Secret: hashedSecret,
		UserID: user.ID,
	}

	err = a.db.db.Create(&newDevice).Error
	if err != nil {
		log.Extract(ctx).WithError(err).Error("inserting device")

		return nil, ErrInternal
	}

	return &proto.RegisterDeviceResponse{
		Code:   code,
		Secret: secret,
	}, nil
}
