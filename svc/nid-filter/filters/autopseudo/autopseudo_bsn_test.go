package autopseudo

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	"lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto"
	wallet_mock "lab.weave.nl/nid/nid-core/svc/wallet-rpc/proto/mock"
)

type AutoBSNTestSuite struct {
	suite.Suite
	key              *rsa.PrivateKey
	authHeader       string
	filter           *Filter
	walletClientMock *wallet_mock.WalletClient
}

func (s *AutoBSNTestSuite) SetupSuite() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	s.key = key

	pseudonymData := []byte("ABCDEFG")
	encryptedPseudonym, err := rsa.EncryptPKCS1v15(rand.Reader, &s.key.PublicKey, pseudonymData)
	encryptedEncodedPseudonym := base64.StdEncoding.EncodeToString(encryptedPseudonym)
	s.Require().NoError(err)
	claims := map[string]interface{}{
		"subjects": map[string]string{
			"nid": encryptedEncodedPseudonym,
		},
	}
	claimsJSON, err := json.Marshal(claims)
	s.Require().NoError(err)
	s.authHeader = fmt.Sprintf(jwtFormat, base64.StdEncoding.EncodeToString(claimsJSON))
}

func (s *AutoBSNTestSuite) SetupTest() {
	s.walletClientMock = &wallet_mock.WalletClient{}
	s.filter = &Filter{
		config: &Config{
			Namespace:         "nid",
			Key:               s.key,
			SubjectIdentifier: "$$nid:bsn$$",
			TranslateToBSN:    true,
			WalletClient:      s.walletClientMock,
		},
	}
}

func (s *AutoBSNTestSuite) TestHeaders() {
	ctx := context.TODO()
	s.walletClientMock.On("GetBSNForPseudonym", ctx, &proto.GetBSNForPseudonymRequest{Pseudonym: "QUJDREVGRw=="}).Return(&proto.GetBSNForPseudonymResponse{Bsn: "123456789"}, nil)

	res, err := s.filter.OnHTTPRequest(context.TODO(), nil, map[string]string{
		"authorization": s.authHeader,
		":path":         "/something?apple=something+containing+%24%24nid%3Absn%24%24+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	})

	s.Require().NoError(err)
	s.Require().NotNil(res)

	expectedResponse := &filter.ProcessingResponse{
		NewHeaders: map[string]string{
			"authorization": s.authHeader,
			":path":         "/something?apple=something+containing+123456789+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
		},
		NewBody:           nil,
		ImmediateResponse: nil,
	}
	s.Require().EqualValues(expectedResponse, res)
}

func (s *AutoBSNTestSuite) TestBody() {
	ctx := context.TODO()
	s.walletClientMock.On("GetBSNForPseudonym", ctx, &proto.GetBSNForPseudonymRequest{Pseudonym: "QUJDREVGRw=="}).Return(&proto.GetBSNForPseudonymResponse{Bsn: "123456789"}, nil)

	headers := map[string]string{
		"authorization": s.authHeader,
		":path":         "/something",
	}

	res, err := s.filter.OnHTTPRequest(context.TODO(), []byte("Some random nonsense containing $$nid:bsn$$ indeed very !n$t4$$$$eresting"), headers)

	s.Require().NoError(err)
	expectedResponse := &filter.ProcessingResponse{
		NewHeaders: map[string]string{
			"authorization":  s.authHeader,
			"content-length": "71",
			":path":          "/something",
		},
		NewBody:           []byte("Some random nonsense containing 123456789 indeed very !n$t4$$$$eresting"),
		ImmediateResponse: nil,
	}

	s.Equal(expectedResponse, res)
}

func (s *AutoBSNTestSuite) TestHeadersAndBody() {
	ctx := context.TODO()
	s.walletClientMock.On("GetBSNForPseudonym", ctx, &proto.GetBSNForPseudonymRequest{Pseudonym: "QUJDREVGRw=="}).Return(&proto.GetBSNForPseudonymResponse{Bsn: "123456789"}, nil)

	res, err := s.filter.OnHTTPRequest(context.TODO(), []byte("Some random nonsense containing $$nid:bsn$$ indeed very !n$t4$$$$eresting"), map[string]string{
		"authorization": s.authHeader,
		":path":         "/something?apple=something+containing+%24%24nid%3Absn%24%24+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	})

	s.Require().NoError(err)
	s.Require().NotNil(res)

	expectedResponse := &filter.ProcessingResponse{
		NewHeaders: map[string]string{
			"authorization":  s.authHeader,
			"content-length": "71",
			":path":          "/something?apple=something+containing+123456789+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
		},
		NewBody:           []byte("Some random nonsense containing 123456789 indeed very !n$t4$$$$eresting"),
		ImmediateResponse: nil,
	}
	s.Require().EqualValues(expectedResponse, res)
}

func TestAutoBSNTestSuite(t *testing.T) {
	suite.Run(t, &AutoBSNTestSuite{})
}
