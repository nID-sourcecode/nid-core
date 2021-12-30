// Package http mocks a http client
package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// ClientMock mock http client
type ClientMock struct {
	mock.Mock
}

// Get mocks the httpclient get function
func (c *ClientMock) Get(location string) (*http.Response, error) {
	args := c.Called(location)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Do mocks the Do function
func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// GetMockHTTPResponse create an http response
func GetMockHTTPResponse(resp interface{}, statusCode int) *http.Response {
	b, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}
}
