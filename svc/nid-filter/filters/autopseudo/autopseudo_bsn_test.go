package autopseudo

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/nID-sourcecode/nid-core/svc/nid-filter/contract"

	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
	walletMock "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto/mock"
)

type AutoBSNTestSuite struct {
	suite.Suite
	key              *rsa.PrivateKey
	authHeader       string
	filter           contract.AuthorizationRule
	walletClientMock *walletMock.WalletClient
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
	s.authHeader = fmt.Sprintf(jwtFormat, base64.RawURLEncoding.EncodeToString(claimsJSON))
}

func (s *AutoBSNTestSuite) SetupTest() {
	s.walletClientMock = &walletMock.WalletClient{}
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
	headers := map[string]string{
		"authorization": s.authHeader,
		":path":         "/something?apple=something+containing+%24%24nid%3Absn%24%24+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	}
	err := s.filter.Check(context.TODO(), returnAuthV3CheckRequest("", headers))

	s.Require().NoError(err)

	NewHeaders := map[string]string{
		"authorization": s.authHeader,
		":path":         "/something?apple=something+containing+123456789+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	}
	s.Require().EqualValues(NewHeaders, headers)
}

func (s *AutoBSNTestSuite) TestBody() {
	ctx := context.TODO()
	s.walletClientMock.On("GetBSNForPseudonym", ctx, &proto.GetBSNForPseudonymRequest{Pseudonym: "QUJDREVGRw=="}).Return(&proto.GetBSNForPseudonymResponse{Bsn: "123456789"}, nil)

	headers := map[string]string{
		"authorization": s.authHeader,
		":path":         "/something",
	}

	authRequest := returnAuthV3CheckRequest("Some random nonsense containing $$nid:bsn$$ indeed very !n$t4$$$$eresting", headers)
	err := s.filter.Check(context.TODO(), authRequest)

	s.Require().NoError(err)
	newHeaders := map[string]string{
		"authorization":  s.authHeader,
		"content-length": "71",
		":path":          "/something",
	}
	newBody := "Some random nonsense containing 123456789 indeed very !n$t4$$$$eresting"
	currentBody := authRequest.GetAttributes().GetRequest().GetHttp().GetBody()
	currentHeaders := authRequest.GetAttributes().GetRequest().GetHttp().GetHeaders()

	s.Equal(newHeaders, currentHeaders)
	s.Equal(newBody, currentBody)
}

func (s *AutoBSNTestSuite) TestHeadersAndBody() {
	ctx := context.TODO()
	s.walletClientMock.On("GetBSNForPseudonym", ctx, &proto.GetBSNForPseudonymRequest{Pseudonym: "QUJDREVGRw=="}).Return(&proto.GetBSNForPseudonymResponse{Bsn: "123456789"}, nil)

	authRequest := returnAuthV3CheckRequest("Some random nonsense containing $$nid:bsn$$ indeed very !n$t4$$$$eresting", map[string]string{
		"authorization": s.authHeader,
		":path":         "/something?apple=something+containing+%24%24nid%3Absn%24%24+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	})

	err := s.filter.Check(context.TODO(), authRequest)

	s.Require().NoError(err)

	newHeaders := map[string]string{
		"authorization":  s.authHeader,
		"content-length": "71",
		":path":          "/something?apple=something+containing+123456789+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	}
	newBody := []byte("Some random nonsense containing 123456789 indeed very !n$t4$$$$eresting")

	currentHeaders := authRequest.GetAttributes().GetRequest().GetHttp().GetHeaders()
	currentBody := authRequest.GetAttributes().GetRequest().GetHttp().GetBody()

	s.Require().EqualValues(newHeaders, currentHeaders)
	s.Require().EqualValues(newBody, currentBody)
}

func TestAutoBSNTestSuite(t *testing.T) {
	suite.Run(t, &AutoBSNTestSuite{})
}

func returnAuthV3CheckRequest(body string, headers map[string]string) *authv3.CheckRequest {
	return &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: headers,
					Body:    body,
				},
			},
		},
	}
}
