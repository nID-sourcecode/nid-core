package luautil

import (
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	auth "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
	lua "github.com/yuin/gopher-lua"
)

// HeadlessAuthCaller handles the authentication for the lua script.
type HeadlessAuthCaller struct {
	authClient *auth.AuthClient
}

// NewHeadlessAuthCaller returns a new instance of NewHeadlessAuthCaller
func NewHeadlessAuthCaller(authClient *auth.AuthClient) *HeadlessAuthCaller {
	return &HeadlessAuthCaller{authClient: authClient}
}

const (
	clientIDIndex = iota + 1
	redirectURLIndex
	audienceIndex
	modelJSON
	modelPath
)

// Call implements the Call method for lua state. HeadlessAuthCaller handles the headless authorization.
func (h *HeadlessAuthCaller) Call(state *lua.LState) int {
	_, err := (*h.authClient).AuthorizeHeadless(state.Context(), &auth.AuthorizeHeadlessRequest{
		ResponseType:   "code",
		ClientId:       state.ToString(clientIDIndex),
		RedirectUri:    state.ToString(redirectURLIndex),
		Audience:       state.ToString(audienceIndex),
		QueryModelJson: state.ToString(modelJSON),
		QueryModelPath: state.ToString(modelPath),
	})
	if err != nil {
		log.WithError(err).Errorln("running authorize headless")
		return 0
	}

	return 1
}
