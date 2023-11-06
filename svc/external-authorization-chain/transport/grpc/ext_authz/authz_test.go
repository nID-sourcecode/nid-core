package extauthz

import (
	"context"
	"testing"

	"github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/contract"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc/codes"

	"github.com/stretchr/testify/mock"

	contractMock "github.com/nID-sourcecode/nid-core/svc/external-authorization-chain/contract/mocks"

	"github.com/stretchr/testify/suite"
)

type AuthzTestSuite struct {
	suite.Suite
	mockApp *contractMock.App

	extAuthz *ExternalAuthorization
}

func (a *AuthzTestSuite) SetupTest() {
	a.mockApp = contractMock.NewApp(a.T())

	a.extAuthz = New(a.mockApp)
}

func (a *AuthzTestSuite) TestCheckEndPoints() {
	a.mockApp.On("CheckEndpoints", mock.Anything, mock.Anything).Return(nil)

	response, err := a.extAuthz.Check(context.Background(), &authv3.CheckRequest{})

	a.NoError(err)
	a.Equal(int32(codes.OK), response.Status.Code)
}

func (a *AuthzTestSuite) TestCheckEndPoint_Fails() {
	a.mockApp.On("CheckEndpoints", mock.Anything, mock.Anything).Return(errors.New("an error"))

	response, err := a.extAuthz.Check(context.Background(), &authv3.CheckRequest{})

	a.NoError(err)
	a.Equal(int32(codes.Unauthenticated), response.Status.Code)
}

func (a *AuthzTestSuite) TestCheckEndPoint_Fails_Missing_Header_Error() {
	a.mockApp.On("CheckEndpoints", mock.Anything, mock.Anything).Return(contract.ErrTargetServiceHeaderNotFound)

	response, err := a.extAuthz.Check(context.Background(), &authv3.CheckRequest{})

	a.NoError(err)
	a.Equal(int32(codes.Unauthenticated), response.Status.Code)
}

func TestAuthzSuite(t *testing.T) {
	suite.Run(t, new(AuthzTestSuite))
}
