package errors

import (
	"strings"
)

// Is checks if a grpc error contains the reference error.
// This is a placeholder method, so we all use the same method. In the future this should change to a better implementation
func Is(err, reference error) bool {
	return strings.Contains(err.Error(), reference.Error())
}
