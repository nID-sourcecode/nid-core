package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	gqlClientMock "lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient/mock"
	metadataMock "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers/mock"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	onboardingPB "lab.weave.nl/nid/nid-core/svc/onboarding/proto"
	pseudoMock "lab.weave.nl/nid/nid-core/svc/pseudonymization/mock"
	pseudoPB "lab.weave.nl/nid/nid-core/svc/pseudonymization/proto"
)

var ErrAllesIsLek error = fmt.Errorf("alles is lek")

type DashboardServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	dataSourceServiceServer *DataSourceServiceServer
}

func (s *DashboardServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()
	walletClientMock := &gqlClientMock.Client{}
	walletClientMock.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		response := args.Get(2).(*Response)
		response.Users = []PseudoResponse{{Pseudonym: "somepseudonym"}}
	})

	pseudoClientMock := &pseudoMock.PseudonymizerClient{}
	convertResponse := &pseudoPB.ConvertResponse{
		Conversions: make(map[string][]byte),
	}
	convertResponse.Conversions["somepseudonym"] = []byte("resultpseudonym")
	pseudoClientMock.On("Convert", mock.Anything, mock.Anything).Return(convertResponse, nil)

	metadataMock := &metadataMock.GRPCMetadataHelperMock{}
	metadataMock.On("GetValFromCtx", s.Ctx, mock.Anything).Return("By=spiffe://cluster.local/ns/foo/sa/httpbin;Hash=<redacted>;Subject=\"\";URI=spiffe://cluster.local/ns/foo/sa/sleep", nil)

	s.dataSourceServiceServer = &DataSourceServiceServer{
		walletClient:           walletClientMock,
		pseudonimizationClient: pseudoClientMock,
		metadataHelper:         metadataMock,
	}
}

func (s *DashboardServiceTestSuite) TestTranslateBSN() {
	res, err := s.dataSourceServiceServer.ConvertBSNToPseudonym(s.Ctx, &onboardingPB.ConvertMessage{
		Bsn: "1234567890",
	})
	s.Require().NoError(err)
	s.Equal("resultpseudonym", fmt.Sprintf("%s", res.Pseudonym))
}

func (s *DashboardServiceTestSuite) TestCantGetPseudonymForBSN() {
	walletClientMock := &gqlClientMock.Client{}
	walletClientMock.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(ErrAllesIsLek)
	s.dataSourceServiceServer.walletClient = walletClientMock
	res, err := s.dataSourceServiceServer.ConvertBSNToPseudonym(s.Ctx, &onboardingPB.ConvertMessage{
		Bsn: "1234567890",
	})
	s.Require().Error(err)
	s.Nil(res)
}

func (s *DashboardServiceTestSuite) TestCantConvertPseudonym() {
	pseudoClientMock := &pseudoMock.PseudonymizerClient{}
	pseudoClientMock.On("Convert", mock.Anything, mock.Anything).Return(nil, ErrAllesIsLek)
	s.dataSourceServiceServer.pseudonimizationClient = pseudoClientMock
	res, err := s.dataSourceServiceServer.ConvertBSNToPseudonym(s.Ctx, &onboardingPB.ConvertMessage{
		Bsn: "1234567890",
	})
	s.Require().Error(err)
	s.Nil(res)
}

func (s *DashboardServiceTestSuite) TestPseudonymNotReturned() {
	pseudoClientMock := &pseudoMock.PseudonymizerClient{}
	convertResponse := &pseudoPB.ConvertResponse{
		Conversions: make(map[string][]byte),
	}
	pseudoClientMock.On("Convert", mock.Anything, mock.Anything).Return(convertResponse, nil)
	s.dataSourceServiceServer.pseudonimizationClient = pseudoClientMock
	res, err := s.dataSourceServiceServer.ConvertBSNToPseudonym(s.Ctx, &onboardingPB.ConvertMessage{
		Bsn: "1234567890",
	})
	s.Require().Error(err)
	s.Nil(res)
}

func TestDashboardServiceTestSuite(t *testing.T) {
	suite.Run(t, &DashboardServiceTestSuite{})
}
