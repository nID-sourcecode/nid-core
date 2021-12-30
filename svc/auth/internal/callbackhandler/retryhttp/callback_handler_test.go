package retryhttp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	httpClientMock "lab.weave.nl/nid/nid-core/pkg/utilities/http"
)

const testDomain = "url.com"

var errTest = errors.New("http error")

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
func (nopCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

type CallbackHandlerTestSuite struct {
	suite.Suite
}

func (s *CallbackHandlerTestSuite) TestHandleCallbackInvalidUrl() {
	callbackHandler := &CallbackHandler{
		Client: &httpClientMock.ClientMock{},
	}

	err := callbackHandler.HandleCallback(context.Background(), "@#$%^&*(", "authcode===")

	s.Error(err)
}

func (s *CallbackHandlerTestSuite) TestHandleCallbackDoError() {
	clientMock := &httpClientMock.ClientMock{}

	callbackHandler := &CallbackHandler{
		Client: clientMock,
	}

	clientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Host == testDomain
	})).Return(&http.Response{}, errTest)

	err := callbackHandler.HandleCallback(context.Background(), "http://url.com", "authcode===")

	s.Equal("executing request for url http://url.com: http error", err.Error())
}

func (s *CallbackHandlerTestSuite) TestHandleCallbackErrorStatusCode() {
	clientMock := &httpClientMock.ClientMock{}

	callbackHandler := &CallbackHandler{
		Client: clientMock,
	}
	clientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Host == testDomain
	})).Return(&http.Response{StatusCode: http.StatusNotFound}, nil)

	err := callbackHandler.HandleCallback(context.Background(), "http://url.com", "authcode===")

	s.Equal("request returned error status code: 404", err.Error())
}

func (s *CallbackHandlerTestSuite) TestHandleCallback() {
	clientMock := &httpClientMock.ClientMock{}

	callbackHandler := &CallbackHandler{
		Client: clientMock,
	}
	clientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Host == testDomain
	})).Return(&http.Response{Body: &nopCloser{}, StatusCode: http.StatusAccepted}, nil)

	err := callbackHandler.HandleCallback(context.Background(), "http://url.com", "authcode===")

	s.NoError(err)
}

func TestCallbackHandlerTestSuite(t *testing.T) {
	suite.Run(t, &CallbackHandlerTestSuite{})
}
