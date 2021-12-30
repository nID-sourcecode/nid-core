package luautil

import (
	"context"

	lua "github.com/yuin/gopher-lua"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
	"lab.weave.nl/nid/nid-core/pkg/utilities/gqlclient"
)

// LuaGraphQLCaller is the default implementation of the lua graphql caller
type LuaGraphQLCaller struct {
	ctx           context.Context
	clientFactory gqlclient.ClientFactory
}

// NewLuaGraphQLCaller creates a new default lua graphql caller
func NewLuaGraphQLCaller(ctx context.Context, clientFactory gqlclient.ClientFactory) *LuaGraphQLCaller {
	return &LuaGraphQLCaller{ctx: ctx, clientFactory: clientFactory}
}

const (
	luaEndpointIndex = iota + 1
	luaQueryIndex
	luaVariablesIndex
)

// Call is called from a lua script and makes a grapqhl call
func (c *LuaGraphQLCaller) Call(state *lua.LState) int {
	endpoint := state.ToString(luaEndpointIndex)
	query := state.ToString(luaQueryIndex)
	variableTable := state.ToTable(luaVariablesIndex)
	if variableTable == nil {
		state.ArgError(luaVariablesIndex, "variables must be a table")
	}
	variableMap, err := ToGoMap(variableTable)
	if err != nil {
		state.RaiseError("%v", errors.Wrap(err, "converting variable table to go map"))
	}

	client := c.clientFactory.NewClient(endpoint)
	req := gqlclient.NewRequest(query)
	req.Variables = variableMap

	response := make(map[string]interface{})
	if err := client.Get(c.ctx, req, &response); err != nil {
		state.RaiseError("%v", errors.Wrap(err, "doing graphql request from lua"))
	}

	state.Push(ToLuaValue(state, response))

	return 1
}
