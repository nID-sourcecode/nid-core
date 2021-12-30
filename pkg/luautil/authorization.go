package luautil

import (
	"context"

	lua "github.com/yuin/gopher-lua"

	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
	auth "lab.weave.nl/nid/nid-core/svc/auth/proto"
)

// HeadlessAuthCaller handles the authentication for the lua script.
type HeadlessAuthCaller struct {
	authClient auth.AuthClient
	ctx        context.Context
}

// NewHeadlessAuthCaller returns a new instance of NewHeadlessAuthCaller
func NewHeadlessAuthCaller(ctx context.Context, authClient auth.AuthClient) *HeadlessAuthCaller {
	return &HeadlessAuthCaller{authClient: authClient, ctx: ctx}
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
	_, err := h.authClient.AuthorizeHeadless(h.ctx, &auth.AuthorizeHeadlessRequest{
		ResponseType:   "code",
		ClientId:       state.ToString(clientIDIndex),
		RedirectUri:    state.ToString(redirectURLIndex),
		Audience:       state.ToString(audienceIndex),
		QueryModelJson: state.ToString(modelJSON),
		QueryModelPath: state.ToString(modelPath),
	})

	log.WithError(err).Errorln("running authorize headless")

	if err != nil {
		return 0
	}

	return 1
}
