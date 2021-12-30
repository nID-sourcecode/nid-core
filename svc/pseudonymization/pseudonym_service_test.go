package main

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	metadataMock "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/headers/mock"
	"lab.weave.nl/nid/nid-core/pkg/utilities/grpctesthelpers"
	"lab.weave.nl/nid/nid-core/svc/pseudonymization/keymanager"
	pb "lab.weave.nl/nid/nid-core/svc/pseudonymization/proto"
)

type PseudonymServiceTestSuite struct {
	grpctesthelpers.GrpcTestSuite

	PseudonymizerServer *PseudonymizerServer
	keyManager          keymanager.KeyManager
}

func (s *PseudonymServiceTestSuite) SetupTest() {
	s.GrpcTestSuite.SetupTest()

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	s.Require().NoError(err)

	jwkKey, err := jwk.New(privateKey.PublicKey)
	s.Require().NoError(err)

	jwkKey.Set(jwk.AlgorithmKey, "RSA1_5")
	jwkKey.Set(jwk.KeyUsageKey, string(jwk.ForEncryption))

	fetcherMock := &keymanager.JWKSFetcherMock{}
	fetcherMock.On("Fetch", mock.Anything).Return(&jwk.Set{
		Keys: []jwk.Key{jwkKey},
	}, nil)

	keyManager := keymanager.NewKeyManager("someurl", durationDay, fetcherMock)
	s.keyManager = keyManager

	metadataMock := &metadataMock.GRPCMetadataHelperMock{}
	metadataMock.On("GetValFromCtx", s.Ctx, mock.Anything).Return("By=spiffe://cluster.local/ns/foo/sa/httpbin;Hash=<redacted>;Subject=\"\";URI=spiffe://cluster.local/ns/foo/sa/sleep", nil)

	s.PseudonymizerServer = &PseudonymizerServer{
		KeyManager:     keyManager,
		metadataHelper: metadataMock,
	}
}

func (s *PseudonymServiceTestSuite) TearDownTest() {
	s.keyManager.Cleanup()
}

func (s *PseudonymServiceTestSuite) TestGeneratePseudonym() {
	res, err := s.PseudonymizerServer.Generate(s.Ctx, &pb.GenerateRequest{
		Amount: 1,
	})
	s.Require().NoError(err)
	s.Require().Equal(1, len(res.GetPseudonyms()))
	s.Equal(64, len(res.Pseudonyms[0]))
}

func (s *PseudonymServiceTestSuite) TestConvert() {
	pseudonym := "935oMsp1RiaMJEHaKBF+/4eL0mEQtQZENAOEbeO1f/YSJbzjyx8AKH1io2Z7L2WS"
	res, err := s.PseudonymizerServer.Convert(s.Ctx, &pb.ConvertRequest{
		NamespaceTo: "dummy",
		Pseudonyms:  []string{pseudonym},
	})
	s.Require().NoError(err)
	s.NotNil(res.Conversions[pseudonym])
}

func TestPseudonymServiceTestSuite(t *testing.T) {
	suite.Run(t, &PseudonymServiceTestSuite{})
}
