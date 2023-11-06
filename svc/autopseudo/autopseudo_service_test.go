package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"github.com/nID-sourcecode/nid-core/pkg/environment"
)

const jwtFormat = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.%s.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

type AutopseudoTestSuite struct {
	suite.Suite
	key    *rsa.PrivateKey
	engine *gin.Engine
}

func (s *AutopseudoTestSuite) SetupSuite() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	s.key = key
	s.engine = initRouter(nil, key, &AutoPseudoConfig{
		BaseConfig: environment.BaseConfig{},
		Namespace:  "nid",
	})
}

func (s *AutopseudoTestSuite) TestDecryptAndApplyPOST() {
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
	authHeader := fmt.Sprintf(jwtFormat, base64.RawURLEncoding.EncodeToString(claimsJSON))

	gofight.New().POST("/decryptAndApply").
		SetHeader(map[string]string{"Authorization": authHeader}).
		SetBody("Some random nonsense containing $$nid:subject$$ indeed very !n$t4$$$$eresting").
		SetDebug(true).
		SetQuery(map[string]string{
			"apple": "something containing $$nid:subject$$ and $$nid:soobjact$$",
			"pie":   "made of pears",
		}).
		Run(s.engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			s.Require().Equal(http.StatusOK, r.Code, r.Body.String())
			res := decryptAndApplyResponse{}
			err := json.Unmarshal(r.Body.Bytes(), &res)
			s.Require().NoError(err)
			s.Equal("Some random nonsense containing QUJDREVGRw== indeed very !n$t4$$$$eresting", res.Body)

			query, err := url.ParseQuery(res.Query)
			s.Require().NoError(err)
			s.Equal("something containing QUJDREVGRw== and $$nid:soobjact$$", query.Get("apple"))
			s.Equal("made of pears", query.Get("pie"))
		},
		)
}

func (s *AutopseudoTestSuite) TestDecryptAndApplyGET() {
	pseudonymData := []byte("HIJKLMNOP")
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
	authHeader := fmt.Sprintf(jwtFormat, base64.RawURLEncoding.EncodeToString(claimsJSON))

	gofight.New().GET("/decryptAndApply").
		SetHeader(map[string]string{"Authorization": authHeader}).
		SetBody("Some random nonsense containing $$nid:subject$$ indeed very !n$t4$$$$eresting").
		SetDebug(true).
		SetQuery(map[string]string{
			"apple": "something containing $$nid:subject$$ and $$nid:soobjact$$",
			"pie":   "made of pears",
		}).
		Run(s.engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			s.Require().Equal(http.StatusOK, r.Code, r.Body.String())
			res := decryptAndApplyResponse{}
			err := json.Unmarshal(r.Body.Bytes(), &res)
			s.Require().NoError(err)
			s.Equal("", res.Body)

			query, err := url.ParseQuery(res.Query)
			s.Require().NoError(err)
			s.Equal("something containing SElKS0xNTk9Q and $$nid:soobjact$$", query.Get("apple"))
			s.Equal("made of pears", query.Get("pie"))
		},
		)
}

func (s *AutopseudoTestSuite) TestDecrypt() {
	pseudonymData := []byte("HIJKLMNOP")
	encryptedPseudonym, err := rsa.EncryptPKCS1v15(rand.Reader, &s.key.PublicKey, pseudonymData)
	s.Require().NoError(err)
	encryptedEncodedPseudonym := base64.StdEncoding.EncodeToString(encryptedPseudonym)
	gofight.New().GET("/decrypt").SetQuery(map[string]string{
		"pseudonym": encryptedEncodedPseudonym,
	}).Run(s.engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		s.Equal(http.StatusOK, r.Code)
		s.Equal(`{"decrypted_pseudonym":"SElKS0xNTk9Q"}`, fmt.Sprintf("%v", r.Body))
	})
}

func TestAutopseudoTestSuite(t *testing.T) {
	suite.Run(t, &AutopseudoTestSuite{})
}
