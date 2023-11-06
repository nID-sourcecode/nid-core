package internal

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/suite"
)

type LuaRunnerSuite struct {
	suite.Suite
	luaRunnerService *LuaRunnerService
}

func (s *LuaRunnerSuite) SetupTest() {
	s.luaRunnerService = NewLuaRunnerService(nil, nil, nil)
}

func (s *LuaRunnerSuite) TestCallbackInvalidBody() {
	httpRequest := &http.Request{
		Method: "POST",
		Body:   io.NopCloser(strings.NewReader(`{"organisatieId" : "5000","organisatieIdType" : "string","timestamp" : "string","abonnementId" : "string ","eventType": "string","recordId": "securatest\"},\"m\":{\"SecuraW\":\"SecuraW"}`)),
	}
	context := gin.Context{
		Writer: &responseWriter{ResponseWriter: httptest.NewRecorder()},
	}
	s.luaRunnerService.HTTPCallback(context.Writer, httpRequest)
	s.Equal("recordID is invalid", context.Writer.(*responseWriter).GetBody())
	s.Equal(400, context.Writer.Status())
}

func (s *LuaRunnerSuite) TestRunValidateJsonBodyFail() {
	httpRequest := &http.Request{
		Method: "POST",
		Body: io.NopCloser(strings.NewReader(`
{
  "örganisatieId": "5000_",
  "eenNull": null,
  "eenGetal": 2,
  "eenObject": {
    "objectBool": true
  },
  "eenArray": [
    "elem1",
    "elem2"
  ],
  "organisatieIdType": "string",
  "timestamp": "string",
  "abonnementId": "string ",
  "eventType": "string",
  "recordId": "securatest\"},\"m\":{\"SecuraW\":\"SecuraW"
}
`)),
	}
	err := s.luaRunnerService.ValidateJSONBody(httpRequest.Body)
	s.Error(err)
	s.True(errors.Is(err, ErrJSONBodyIsInvalid))
}

func (s *LuaRunnerSuite) TestRunValidateJsonBody() {
	httpRequest := &http.Request{
		Method: "POST",
		Body: io.NopCloser(strings.NewReader(`
{
  "örganisatieId": "5000",
  "eenNull": null,
  "eenGetal": 2,
  "eenObject": {
    "objectBool": true
  },
  "eenArray": [
    "elem1",
    "elem2"
  ],
  "organisatieIdType": "string",
  "timestamp": "string",
  "abonnementId": "string ",
  "eventType": "string",
  "recordId": "eenID/NogEenID"
}
`)),
	}
	err := s.luaRunnerService.ValidateJSONBody(httpRequest.Body)
	s.NoError(err)
}

func TestLuaRunnerSuite(t *testing.T) {
	suite.Run(t, new(LuaRunnerSuite))
}

type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			log.Info("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

func (w *responseWriter) GetBody() string {
	return w.ResponseWriter.(*httptest.ResponseRecorder).Body.String()
}

func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return n, err
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.size != -1
}

// Hijack implements the http.Hijacker interface.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify interface.
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify() //nolint
}

// Flush implements the http.Flush interface.
func (w *responseWriter) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *responseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
