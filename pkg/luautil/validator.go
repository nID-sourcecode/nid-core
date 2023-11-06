package luautil

import (
	lua "github.com/yuin/gopher-lua"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
)

var (
	// ErrVariableNotDefined returned when required variable is not found in script
	ErrVariableNotDefined = errors.New("validator could not find required variable in script")
	// ErrFunctionsNotDefined returned when required function is not found in script
	ErrFunctionsNotDefined = errors.New("validator could not find required function in script")
)

// FunctionConfig config for lua script functions
type FunctionConfig struct {
	Name                string
	MockFunc            *lua.LGFunction
	RequiredDefinitions RequiredDefinitions
}

// LuaValidator is the settings struct for the lua script validator
type LuaValidator struct {
	RequiredDefinitions RequiredDefinitions
}

// NewLuaValidator constructor
func NewLuaValidator(requiredDefinitions RequiredDefinitions) *LuaValidator {
	return &LuaValidator{RequiredDefinitions: requiredDefinitions}
}

// RequiredDefinitions config that stores function and variable rules for validator
type RequiredDefinitions struct {
	Functions []FunctionConfig
	Variables []string
}

// ValidateScript function validates a compiled lua script with the rules defined in the validator
func (l *LuaValidator) ValidateScript(proto *lua.FunctionProto) (err error) {
	// Base rule script must run without errors
	L := lua.NewState()
	defer L.Close()

	// Set global mock functions when defined
	for function := range l.RequiredDefinitions.Functions {
		funcObj := l.RequiredDefinitions.Functions[function]
		if funcObj.MockFunc != nil {
			L.SetGlobal(funcObj.Name, L.NewFunction(*funcObj.MockFunc))
		}
	}

	// Push proto and call
	L.Push(L.NewFunctionFromProto(proto))
	err = L.PCall(0, lua.MultRet, nil)
	if err != nil {
		return errors.Wrap(err, "failed to execute lua script")
	}

	// Check if we find all variables defined in config
	for _, variable := range l.RequiredDefinitions.Variables {
		luaValue := L.GetGlobal(variable)
		goValue, err := ToGoValue(luaValue)
		if err != nil {
			return errors.Errorf("%w: parsing lua value to go value failed for variable %s", ErrVariableNotDefined, variable)
		}
		if goValue == nil {
			return errors.Errorf("%w: %s", ErrVariableNotDefined, variable)
		}
	}

	// Check if we find all functions defined in config
	for function := range l.RequiredDefinitions.Functions {
		funcObj := l.RequiredDefinitions.Functions[function]
		funcSrc, ok := L.GetGlobal(funcObj.Name).(*lua.LFunction)
		if !ok {
			return errors.Errorf("%w: %s", ErrFunctionsNotDefined, function)
		}

		// Check if function has validator rules
		if funcSrc.Proto != nil && (len(funcObj.RequiredDefinitions.Functions) > 0 || len(funcObj.RequiredDefinitions.Variables) > 0) {
			// Create validator with new rules
			validator := &LuaValidator{
				RequiredDefinitions: funcObj.RequiredDefinitions,
			}
			// Call validator for function
			return validator.ValidateScript(funcSrc.Proto)
		}
	}

	return nil
}
