package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/grpctesthelpers"
	"github.com/nID-sourcecode/nid-core/svc/autobsn/proto"
	walletPB "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto"
	walletmock "github.com/nID-sourcecode/nid-core/svc/wallet-rpc/proto/mock"
)

const jwtFormat = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.%s.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

type AutoBSNTestSuite struct {
	grpctesthelpers.GrpcTestSuite
	key *rsa.PrivateKey
}

func (s *AutoBSNTestSuite) SetupSuite() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	s.key = key
}

func (s *AutoBSNTestSuite) TestReplacePlaceholderWithBSN_WithBody() {
	pseudonymData := []byte("ABCDEFG")
	encryptedPseudonym, err := rsa.EncryptPKCS1v15(rand.Reader, &s.key.PublicKey, pseudonymData)
	encryptedEncodedPseudonym := base64.StdEncoding.EncodeToString(encryptedPseudonym)
	s.Require().NoError(err)
	claims := map[string]interface{}{
		"sub": encryptedEncodedPseudonym,
	}
	claimsJSON, err := json.Marshal(claims)
	s.Require().NoError(err)
	authHeader := fmt.Sprintf(jwtFormat, base64.StdEncoding.EncodeToString(claimsJSON))

	walletClientMock := walletmock.WalletClient{}
	server := NewAutoBSNServer(s.key, &walletClientMock)

	walletClientMock.On("GetBSNForPseudonym", mock.Anything, &walletPB.GetBSNForPseudonymRequest{Pseudonym: "QUJDREVGRw=="}).
		Return(&walletPB.GetBSNForPseudonymResponse{Bsn: "13243546"}, nil)

	res, err := server.ReplacePlaceholderWithBSN(context.TODO(), &proto.ReplacePlaceholderWithBSNRequest{
		Body:                "Some random nonsense containing $$nid:bsn$$ indeed very !n$t4$$$$eresting",
		Query:               "apple=something%20containing%20%24%24nid%3Absn%24%24%20and%20%24%24nid%3Abeesn%24%24&pie=made%20of%20pears",
		Method:              "POST",
		AuthorizationHeader: authHeader,
	})

	s.Require().NoError(err)
	s.Equal("Some random nonsense containing 13243546 indeed very !n$t4$$$$eresting", res.Body)

	query, err := url.ParseQuery(res.Query)
	s.Require().NoError(err)
	s.Equal("something containing 13243546 and $$nid:beesn$$", query.Get("apple"))
	s.Equal("made of pears", query.Get("pie"))
}

func (s *AutoBSNTestSuite) TestReplacePlaceholderWithBSN_WithoutBody() {
	pseudonymData := []byte("HIJKLMNOP")
	encryptedPseudonym, err := rsa.EncryptPKCS1v15(rand.Reader, &s.key.PublicKey, pseudonymData)
	encryptedEncodedPseudonym := base64.StdEncoding.EncodeToString(encryptedPseudonym)
	s.Require().NoError(err)
	claims := map[string]interface{}{
		"sub": encryptedEncodedPseudonym,
	}
	claimsJSON, err := json.Marshal(claims)
	s.Require().NoError(err)
	authHeader := fmt.Sprintf(jwtFormat, base64.StdEncoding.EncodeToString(claimsJSON))

	walletClientMock := walletmock.WalletClient{}
	server := NewAutoBSNServer(s.key, &walletClientMock)

	walletClientMock.On("GetBSNForPseudonym", mock.Anything, &walletPB.GetBSNForPseudonymRequest{Pseudonym: "SElKS0xNTk9Q"}).
		Return(&walletPB.GetBSNForPseudonymResponse{Bsn: "08978675"}, nil)

	res, err := server.ReplacePlaceholderWithBSN(context.TODO(), &proto.ReplacePlaceholderWithBSNRequest{
		Body:                "",
		Query:               "apple=something%20containing%20%24%24nid%3Absn%24%24%20and%20%24%24nid%3Abeesn%24%24&pie=made%20of%20pears",
		Method:              "GET",
		AuthorizationHeader: authHeader,
	})

	s.Require().NoError(err)
	query, err := url.ParseQuery(res.Query)
	s.Require().NoError(err)
	s.Equal("something containing 08978675 and $$nid:beesn$$", query.Get("apple"))
	s.Equal("made of pears", query.Get("pie"))
}

func TestAutoBSNTestSuite(t *testing.T) {
	suite.Run(t, &AutoBSNTestSuite{})
}
