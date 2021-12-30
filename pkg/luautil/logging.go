package luautil

import (
	lua "github.com/yuin/gopher-lua"

	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
)

// AddAllLogFunctions registers all LuaLogger functions in a lua state so that they can be called from a script.
func AddAllLogFunctions(state *lua.LState, logger log.LoggerUtility) {
	luaLogger := NewLuaLogger(logger)
	state.SetGlobal("logWarn", state.NewFunction(luaLogger.LogWarn))
	state.SetGlobal("logInfo", state.NewFunction(luaLogger.LogInfo))
	state.SetGlobal("logError", state.NewFunction(luaLogger.LogError))
	state.SetGlobal("logDebug", state.NewFunction(luaLogger.LogDebug))
}

// NewLuaLogger creates a new lua logger
func NewLuaLogger(logger log.LoggerUtility) *LuaLogger {
	return &LuaLogger{logger: logger.With("source", "lua")}
}

// LuaLogger logs from lua scripts
type LuaLogger struct {
	logger log.LoggerUtility
}

// LogWarn logs with warning level. This should be called from a lua script.
func (l *LuaLogger) LogWarn(state *lua.LState) int {
	message := state.ToString(1)
	l.logger.Warn(message)

	return 0
}

// LogInfo logs with info level. This should be called from a lua script.
func (l *LuaLogger) LogInfo(state *lua.LState) int {
	message := state.ToString(1)
	l.logger.Info(message)

	return 0
}

// LogError logs with error level. This should be called from a lua script.
func (l *LuaLogger) LogError(state *lua.LState) int {
	message := state.ToString(1)
	l.logger.Error(message)

	return 0
}

// LogDebug logs with debug level. This should be called from a lua script.
func (l *LuaLogger) LogDebug(state *lua.LState) int {
	message := state.ToString(1)
	l.logger.Debug(message)

	return 0
}
