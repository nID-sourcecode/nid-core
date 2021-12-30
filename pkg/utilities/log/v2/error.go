package log

import (
	"fmt"
)

var (
	// ErrInvalidFormat is thrown when given string doesn't conform to Format enum
	ErrInvalidFormat error = fmt.Errorf("invalid log format")
	// ErrInvalidLevel is thrown when given string doesn't conform to Level enum
	ErrInvalidLevel error = fmt.Errorf("invalid log level")
	// ErrIncorrectFormatter is thrown when format type is not recognised
	ErrIncorrectFormatter error = fmt.Errorf("unexpected format type")
)
