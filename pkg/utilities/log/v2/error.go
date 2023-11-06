package log

import (
	"fmt"
)

var (
	// ErrInvalidFormat is thrown when given string doesn't conform to Format enum
	ErrInvalidFormat = fmt.Errorf("invalid log format")
	// ErrInvalidLevel is thrown when given string doesn't conform to Level enum
	ErrInvalidLevel = fmt.Errorf("invalid log level")
	// ErrIncorrectFormatter is thrown when format type is not recognised
	ErrIncorrectFormatter = fmt.Errorf("unexpected format type")
)
