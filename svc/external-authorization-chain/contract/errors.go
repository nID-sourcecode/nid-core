package contract

import "errors"

var (
	// ErrTargetServiceHeaderNotFound error for when the target service name is not found in the headers.
	ErrTargetServiceHeaderNotFound = errors.New("header target service could not be found in the headers")

	// ErrUnauthorized error for when the request is not authorized by the services
	ErrUnauthorized = errors.New("authorization client did not respond with code OK")

	// ErrHostIsDeniedByDefault error for when the internal.Config.DenyByDefault is enabled and the given host is not in the allowed list.
	ErrHostIsDeniedByDefault = errors.New("host is not allowed and denied by default")
)
