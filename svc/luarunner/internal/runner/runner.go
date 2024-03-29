// Package runner provides a runner for running lua scripts
package runner

import (
	"context"

	"github.com/nID-sourcecode/nid-core/pkg/gqlclient"

	lua "github.com/yuin/gopher-lua"

	"github.com/nID-sourcecode/nid-core/pkg/luautil"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	auth "github.com/nID-sourcecode/nid-core/svc/auth/transport/grpc/proto"
)

// ScriptRunner interface to run scripts.
type ScriptRunner interface {
	RunScript(ctx context.Context, script string, input map[string]interface{}) error
}

// LuaRunner script runner implementation
type LuaRunner struct {
	authClient *auth.AuthClient
}

// NewLuaRunner instantiates an instance of LuaRunner
func NewLuaRunner(authClient *auth.AuthClient) LuaRunner {
	return LuaRunner{authClient: authClient}
}

// RunScript run given script
func (l *LuaRunner) RunScript(ctx context.Context, script string, input map[string]interface{}) error {
	state := lua.NewState()

	// setup
	l.setupFunctions(ctx, state)

	lInput := luautil.ToLuaValue(state, input)

	if err := state.DoString(script); err != nil {
		return errors.Wrap(err, "unable to execute lua script")
	}

	err := state.CallByParam(lua.P{
		Fn:      state.GetGlobal("handle"),
		NRet:    0,
		Protect: true,
	}, lInput)
	if err != nil {
		log.WithError(err).Error("call state by param")
		return err
	}

	return nil
}

// SetupFunctions adds global functions to the lua state. This makes it possible to run custom go functions within the lua script state.
func (l *LuaRunner) setupFunctions(ctx context.Context, state *lua.LState) {
	state.SetContext(ctx)
	state.SetGlobal("graphql", state.NewFunction(luautil.NewLuaGraphQLCaller(gqlclient.NewRestyGQLClientFactory()).Call))
	state.SetGlobal("headlessAuthorization", state.NewFunction(luautil.NewHeadlessAuthCaller(l.authClient).Call))
	luautil.AddAllLogFunctions(state, log.GetLogger())
}
