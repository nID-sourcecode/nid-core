// Package contract definition of interfaces of the app and adapters, to provide encapsulation of the service's layers.
package contract

import (
	"context"

	"github.com/nID-sourcecode/nid-core/svc/auth/models"
)

// App defines the method definitions for auth app.
type App interface {
	Authorize(context.Context, *models.AuthorizeRequest) (string, error)
	AuthorizeHeadless(context.Context, *models.AuthorizeHeadlessRequest) error
	Claim(ctx context.Context, jwtPayload string, acceptRequest *models.SessionRequest) (*models.SessionResponse, error)
	Accept(ctx context.Context, jwtPayload string, acceptRequest *models.AcceptRequest) (*models.SessionResponse, error)
	Reject(context.Context, *models.SessionRequest) error
	GenerateSessionFinaliseToken(context.Context, *models.SessionRequest) (*models.SessionAuthorization, error)
	GetSessionDetails(context.Context, *models.SessionRequest) (*models.SessionResponse, error)
	Status(context.Context, *models.SessionRequest) (*models.StatusResponse, error)
	Finalise(context.Context, *models.FinaliseRequest) (*models.FinaliseResponse, error)
	Token(ctx context.Context, req *models.TokenRequest) (*models.TokenResponse, error)
	TokenClientFlow(ctx context.Context, req *models.TokenClientFlowRequest) (*models.TokenResponse, error)
	RegisterAccessModel(context.Context, *models.AccessModelRequest) error
	SwapToken(context.Context, *models.SwapTokenRequest) (*models.TokenResponse, error)
}
