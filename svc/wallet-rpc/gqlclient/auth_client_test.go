package gqlclient

import (
	"testing"

	gqlMock "github.com/nID-sourcecode/nid-core/pkg/gqlclient/mock"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
)

type AuthClientTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	mockGqlClient *gqlMock.Client
	authClient    *AuthClient
}

func (s *AuthClientTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
	s.mockGqlClient = &gqlMock.Client{}
	s.authClient = &AuthClient{
		client: s.mockGqlClient,
	}
}

func (s *AuthClientTestSuite) TestFetchClient() {
	clientID := uuid.Must(uuid.NewV4())
	client := Client{
		Color: "blue",
		Icon:  "icon",
		Logo:  "logo",
		Name:  "name",
	}

	s.mockGqlClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		response := args.Get(2).(*map[string]Client)
		(*response)["client"] = client
	})
	out, err := s.authClient.FetchClient(s.Ctx, clientID)
	s.Require().NoError(err)

	s.Equal(client.Color, out.Color)
	s.Equal(client.Icon, out.Icon)
	s.Equal(client.Logo, out.Logo)
	s.Equal(client.Name, out.Name)
}

func (s *AuthClientTestSuite) TestFetchClient_EmptyRes() {
	s.mockGqlClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	_, err := s.authClient.FetchClient(s.Ctx, uuid.Must(uuid.NewV4()))
	s.Require().Error(err)
}

func (s *AuthClientTestSuite) TestFetchClient_Error() {
	s.mockGqlClient.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("alles is lek"))
	_, err := s.authClient.FetchClient(s.Ctx, uuid.Must(uuid.NewV4()))
	s.Require().Error(err)
}

func TestAuthClientTestSuite(t *testing.T) {
	suite.Run(t, &AuthClientTestSuite{})
}
