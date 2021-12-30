// Package callbackhandler contains the interface for callback handlers
package callbackhandler

import "context"

// CallbackHandler is an interface for handling callbacks.
type CallbackHandler interface {
	HandleCallback(context.Context, string, string) error
}
