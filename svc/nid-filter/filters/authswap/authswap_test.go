package authswap

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/extproc/filter"
	authpb "lab.weave.nl/nid/nid-core/svc/auth/proto"
	"lab.weave.nl/nid/nid-core/svc/auth/proto/mock"
)

type AuthSwapFilterTestSuite struct {
	suite.Suite
	authClient *mock.AuthClient
	filter     *Filter
}

func (s *AuthSwapFilterTestSuite) TestSetup() {
	s.authClient = &mock.AuthClient{}
	s.filter = &Filter{authClient: s.authClient}
}

func (s *AuthSwapFilterTestSuite) SetupSuite() {}

func (s *AuthSwapFilterTestSuite) Test_FilterSwap() {
	tests := []struct {
		Name         string
		Method       string
		TokenPayload string
		TokenReponse string
		Body         []byte
		Path         string
	}{
		{
			Name:         "TestPOSTRequest1",
			Path:         "/gql",
			Method:       http.MethodPost,
			Body:         []byte("{\"query\":\"\\n\\t\\tquery {\\n\\t\\t\\tusers(filter: {\\n\\t\\t\\t\\tpseudonym: {\\n\\t\\t\\t\\t\\teq: \\\"$$nid:subject$$\\\"\\n\\t\\t\\t\\t}\\n\\t\\t\\t}){\\n\\t\\t\\t\\tbankAccounts {\\n\\t\\t\\t\\t  amount\\n\\t\\t\\t\\t  savingsAccounts {\\n\\t\\t\\t\\t\\tamount\\n\\t\\t\\t\\t  }\\n\\t\\t\\t\\t}\\n\\t\\t\\t}\\n\\t\\t}\\n\\t\",\"variables\":null}"),
			TokenPayload: "testAwesomeToken",
			TokenReponse: "isThisTokenSwapped?", // Swap action adds Bearer prefix
		},
		{
			Name:         "TestPOSTRequest2",
			Path:         "/gql",
			Method:       http.MethodPost,
			Body:         []byte("{\"query\":\"\\n\\t\\tquery {\\n\\t\\t\\tusers(filter: {\\n\\t\\t\\t\\tpseudonym: {\\n\\t\\t\\t\\t\\teq: \\\"$$nid:subject$$\\\"\\n\\t\\t\\t\\t}\\n\\t\\t\\t}){\\n\\t\\t\\t\\tbankAccounts {\\n\\t\\t\\t\\t  amount\\n\\t\\t\\t\\t  savingsAccounts {\\n\\t\\t\\t\\t\\tamount\\n\\t\\t\\t\\t  }\\n\\t\\t\\t\\t}\\n\\t\\t\\t}\\n\\t\\t}\\n\\t\",\"variables\":null}"),
			TokenPayload: "testAwesomeToken",
			TokenReponse: "isThisTokenSwapped?", // Swap action adds Bearer prefix
		},
		{
			Name:         "TestGETRequest1",
			Path:         "/gql?query=%0A%09%09%7B%0A%09%09++users%28filter%3A+%7Bpseudonym%3A+%7Beq%3A+%22%24%24nid%3Asubject%24%24%22%7D%7D%29+%7B%0A%09%09%09contactDetails+%7B%0A%09%09%09++phone%0A%09%09%09%7D%0A%09%09++%7D%0A%09%09%7D%0A%09&variables=%7B%7D",
			Method:       http.MethodGet,
			TokenPayload: "randomAccessToken",
			TokenReponse: "swappedTokenY", // Swap action adds Bearer prefix
		},
		{
			Name:         "TestGETRequest2",
			Path:         "/gql?query=%0A%09%09%7B%0A%09%09++users%28filter%3A+%7Bpseudonym%3A+%7Beq%3A+%22%24%24nid%3Asubject%24%24%22%7D%7D%29+%7B%0A%09%09%09contactDetails+%7B%0A%09%09%09++phone%0A%09%09%09%7D%0A%09%09++%7D%0A%09%09%7D%0A%09&variables=%7B%7D",
			Method:       http.MethodGet,
			TokenPayload: "otherAccessToken",
			TokenReponse: "otherSwappedTokenYes", // Swap action adds Bearer prefix
		},
	}

	for _, t := range tests {
		s.Run(t.Name, func() {
			s.TestSetup()
			s.Require().NotEmpty(t.Method)

			// Set incoming request headers
			headers := map[string]string{
				"authorization":     "Bearer " + t.TokenPayload,
				":path":             t.Path,
				":method":           t.Method,
				"x-forwarded-proto": "http",
				":authority":        "databron.nid",
			}

			ctx := context.TODO()

			s.authClient.On("SwapToken", ctx, &authpb.SwapTokenRequest{
				CurrentToken: t.TokenPayload,
				Query:        "stub",
				Audience:     "http://databron.nid/gql",
			}).Return(&authpb.TokenResponse{
				AccessToken: t.TokenReponse,
				TokenType:   "Bearer",
			}, nil)

			res, err := s.filter.OnHTTPRequest(ctx, t.Body, headers)

			s.Require().NoError(err)
			s.Require().NotNil(res)

			expectedRes := &filter.ProcessingResponse{
				NewHeaders: map[string]string{
					"authorization":     "Bearer " + t.TokenReponse,
					":path":             t.Path,
					":method":           t.Method,
					"x-forwarded-proto": "http",
					":authority":        "databron.nid",
				},
				NewBody:           nil,
				ImmediateResponse: nil,
			}

			s.Require().Equal(expectedRes, res)
		})
	}
}

func (s *AuthSwapFilterTestSuite) Test_FilterIgnore() {
	tests := []struct {
		Name   string
		Method string
	}{
		{
			Name:   "TestPOSTRequest",
			Method: http.MethodPost,
		},
		{
			Name:   "TestGETRequest",
			Method: http.MethodGet,
		},
	}

	for _, t := range tests {
		s.Run(t.Name, func() {
			s.TestSetup()
			s.Require().NotEmpty(t.Method)

			// Set incoming request headers
			headers := map[string]string{
				":path":             "/gql",
				":method":           t.Method,
				"x-forwarded-proto": "http",
				":authority":        "databron.nid",
			}

			body := []byte("{\"query\":\"\\n\\t\\tquery {\\n\\t\\t\\tusers(filter: {\\n\\t\\t\\t\\tpseudonym: {\\n\\t\\t\\t\\t\\teq: \\\"$$nid:subject$$\\\"\\n\\t\\t\\t\\t}\\n\\t\\t\\t}){\\n\\t\\t\\t\\tbankAccounts {\\n\\t\\t\\t\\t  amount\\n\\t\\t\\t\\t  savingsAccounts {\\n\\t\\t\\t\\t\\tamount\\n\\t\\t\\t\\t  }\\n\\t\\t\\t\\t}\\n\\t\\t\\t}\\n\\t\\t}\\n\\t\",\"variables\":null}")

			ctx := context.TODO()
			res, err := s.filter.OnHTTPRequest(ctx, body, headers)

			s.Nil(res)
			s.Nil(err)
		})
	}
}

func TestAuthSwapFilterTestSuite(t *testing.T) {
	suite.Run(t, new(AuthSwapFilterTestSuite))
}
