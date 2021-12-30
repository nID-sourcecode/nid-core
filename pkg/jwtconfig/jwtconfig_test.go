package jwtconfig

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
)

type JWTConfigTestSuite struct {
	suite.Suite
}

func TestJWTConfigTestTestSuite(t *testing.T) {
	suite.Run(t, &JWTConfigTestSuite{})
}

func (s *JWTConfigTestSuite) TestRead() {
	path := s.T().TempDir()
	id := uuid.Must(uuid.NewV4()).String()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	pem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	s.Require().NoError(ioutil.WriteFile(path+"/kid", []byte(id), os.ModePerm))
	s.Require().NoError(ioutil.WriteFile(path+"/jwtkey.pem", pem, os.ModePerm))

	actualKey, err := Read(path)
	s.Require().NoError(err)

	expectedKey := &JWTKey{
		PrivateKey: key,
		ID:         id,
	}

	s.Equal(expectedKey, actualKey)
}

func (s *JWTConfigTestSuite) TestReadNoKID() {
	path := s.T().TempDir()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	pem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	s.Require().NoError(ioutil.WriteFile(path+"/jwtkey.pem", pem, os.ModePerm))

	actualKey, err := Read(path)
	s.Error(err)
	s.Nil(actualKey)
}

func (s *JWTConfigTestSuite) TestReadNoKey() {
	path := s.T().TempDir()
	id := uuid.Must(uuid.NewV4()).String()
	s.Require().NoError(ioutil.WriteFile(path+"/kid", []byte(id), os.ModePerm))

	actualKey, err := Read(path)
	s.Error(err)
	s.Nil(actualKey)
}

func (s *JWTConfigTestSuite) TestInvalidKey() {
	path := s.T().TempDir()
	id := uuid.Must(uuid.NewV4()).String()
	s.Require().NoError(ioutil.WriteFile(path+"/kid", []byte(id), os.ModePerm))

	s.Require().NoError(ioutil.WriteFile(path+"/jwtkey.pem", []byte("SOMEINVALID NOT A PEM, if this does not return an error everything is leaking"), os.ModePerm))

	actualKey, err := Read(path)
	s.Error(err)
	s.Nil(actualKey)
}
