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
)

const jwtFormat = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.%s.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

type AutopseudoTestSuite struct {
	suite.Suite
	key        *rsa.PrivateKey
	authHeader string
	filter     *Filter
}

func (s *AutopseudoTestSuite) SetupSuite() {
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

func (s *AutopseudoTestSuite) SetupTest() {
	s.filter = &Filter{
		config: &Config{
			Namespace:         "nid",
			Key:               s.key,
			SubjectIdentifier: "$$nid:subject$$",
			TranslateToBSN:    false,
		},
	}
}

func (s *AutopseudoTestSuite) TestHeaders() {
	authRequest := returnAuthV3CheckRequest("", map[string]string{
		"authorization": s.authHeader,
		":path":         "/something?apple=something+containing+%24%24nid%3Asubject%24%24+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	})
	err := s.filter.Check(context.TODO(), authRequest)

	s.Require().NoError(err)

	newHeader := map[string]string{
		"authorization": s.authHeader,
		":path":         "/something?apple=something+containing+QUJDREVGRw%3D%3D+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	}

	currentHeaders := authRequest.GetAttributes().GetRequest().GetHttp().GetHeaders()
	currentBody := authRequest.GetAttributes().GetRequest().GetHttp().GetBody()

	s.Equal(newHeader, currentHeaders)
	s.Empty(currentBody)
}

func (s *AutopseudoTestSuite) TestHeaders_ComplexGQLQuery() {
	authRequest := returnAuthV3CheckRequest("", map[string]string{
		"authorization": s.authHeader,
		":path":         "/gql?query=%0A%09%09%7B%0A%09%09++users%28filter%3A+%7Bpseudonym%3A+%7Beq%3A+%22%24%24nid%3Asubject%24%24%22%7D%7D%29+%7B%0A%09%09%09%09contactDetails+%7B%0A%09%09%09++phone%0A%09%09%09++address+%7B%0A%09%09%09%09houseNumber%0A%09%09%09++%7D%0A%09%09%09%7D%0A%09%09++%7D%0A%09%09%7D%0A%09&variables=%7B%7D",
	})
	err := s.filter.Check(context.TODO(), authRequest)

	s.Require().NoError(err)

	NewHeaders := map[string]string{
		"authorization": s.authHeader,
		":path":         "/gql?query=%0A%09%09%7B%0A%09%09++users%28filter%3A+%7Bpseudonym%3A+%7Beq%3A+%22QUJDREVGRw%3D%3D%22%7D%7D%29+%7B%0A%09%09%09%09contactDetails+%7B%0A%09%09%09++phone%0A%09%09%09++address+%7B%0A%09%09%09%09houseNumber%0A%09%09%09++%7D%0A%09%09%09%7D%0A%09%09++%7D%0A%09%09%7D%0A%09&variables=%7B%7D",
	}

	currentHeaders := authRequest.GetAttributes().GetRequest().GetHttp().GetHeaders()
	currentBody := authRequest.GetAttributes().GetRequest().GetHttp().GetBody()

	s.Equal(NewHeaders, currentHeaders)
	s.Empty(currentBody)
}

func (s *AutopseudoTestSuite) TestBody() {
	authRequest := returnAuthV3CheckRequest("Some random nonsense containing $$nid:subject$$ indeed very !n$t4$$$$eresting", map[string]string{
		"authorization": s.authHeader,
		":path":         "/something",
	})

	err := s.filter.Check(context.TODO(), authRequest)

	s.Require().NoError(err)
	NewHeaders := map[string]string{
		"authorization":  s.authHeader,
		"content-length": "74",
		":path":          "/something",
	}
	NewBody := "Some random nonsense containing QUJDREVGRw== indeed very !n$t4$$$$eresting"

	currentHeaders := authRequest.GetAttributes().GetRequest().GetHttp().GetHeaders()
	currentBody := authRequest.GetAttributes().GetRequest().GetHttp().GetBody()

	s.Equal(currentHeaders, NewHeaders)
	s.Equal(currentBody, NewBody)
}

func (s *AutopseudoTestSuite) TestHeadersAndBody() {
	authRequest := returnAuthV3CheckRequest("Some random nonsense containing $$nid:subject$$ indeed very !n$t4$$$$eresting",
		map[string]string{
			"authorization": s.authHeader,
			":path":         "/something?apple=something+containing+%24%24nid%3Asubject%24%24+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
		})

	err := s.filter.Check(context.TODO(), authRequest)

	s.Require().NoError(err)

	NewHeaders := map[string]string{
		"authorization":  s.authHeader,
		"content-length": "74",
		":path":          "/something?apple=something+containing+QUJDREVGRw%3D%3D+and+%24%24nid%3Asoobjact%24%24&pie=made+of+pears",
	}
	NewBody := []byte("Some random nonsense containing QUJDREVGRw== indeed very !n$t4$$$$eresting")

	currentHeaders := authRequest.GetAttributes().GetRequest().GetHttp().GetHeaders()
	currentBody := authRequest.GetAttributes().GetRequest().GetHttp().GetBody()

	s.Require().EqualValues(NewHeaders, currentHeaders)
	s.Require().EqualValues(NewBody, currentBody)
}

func TestAutopseudoTestSuite(t *testing.T) {
	suite.Run(t, &AutopseudoTestSuite{})
}
