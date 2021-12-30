//go:build integration || to || files
// +build integration to files

package integration

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vrischmann/envconfig"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
	authgqlclient "lab.weave.nl/nid/nid-core/svc/wallet-rpc/gqlclient"
)

type WalletIntegrationTestSuite struct {
	BaseTestSuite
	walletConfig
}

type walletConfig struct {
	DeviceCode   string `envconfig:"DEVICE_CODE"`
	DeviceSecret string `envconfig:"DEVICE_SECRET"`
}

func (s *WalletIntegrationTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.Require().NoError(envconfig.Init(&s.walletConfig), "unable to initialise wallet environment config")
}

func TestWalletIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WalletIntegrationTestSuite))
}

func (s *WalletIntegrationTestSuite) TestGetAuthClient() {
	authClient := authgqlclient.NewAuthClient(fmt.Sprintf("http://%s.%s:%d/gql", s.envConfig.Service.AuthGQL, s.envConfig.BackendURL, s.envConfig.BackendPort))
	_, err := authClient.FetchClient(s.ctx, uuid.Must(uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")))
	// we expect the auth client to deny this request, since it's not from the wallet
	s.Error(err)
	s.True(errors.Is(err, gqlclient.ErrRemoteErrorResponse), "error was %v instead of ErrRemoteErrorResponse", err)
	s.Contains(err.Error(), "record not found")
}
