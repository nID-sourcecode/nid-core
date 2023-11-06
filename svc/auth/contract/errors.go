package contract

import (
	"github.com/pkg/errors"
)

var (
	// ErrInvalidArguments error when wrong arguments were used.
	ErrInvalidArguments = errors.New("Invalid Argument")
	// ErrInternalError internal error.
	ErrInternalError = errors.New("Internal Error")
	// ErrNotFound resource not found error.
	ErrNotFound = errors.New("Not Found")
	// ErrUnauthenticated client is unauthenticated
	ErrUnauthenticated = errors.New("unauthenticated")
	// ErrDeadlineExceeded context deadline exceeded
	ErrDeadlineExceeded = errors.New("deadline exceeded")

	ErrUnableToRetrieveTokenExpiration   = errors.New("error getting token, token request outside expiration time")
	ErrUnableToRetrieveTokenInvalidState = errors.New("error getting token, session in incorrect state")

	ErrInvalidQueryModelType = errors.New("query model does not have a related specific query model")
	ErrAudienceNotDefined    = errors.New("audience not defined")

	ErrIncorrectBasicAuthPrefix = errors.New("incorrect basic auth prefix")

	ErrIncorrectBasicAuthFormat = errors.New("incorrect basic auth format")

	ErrSigningToken = errors.New("error signing token")

	ErrMultipleAudiencesNotAllowed = errors.New("multiple audiences are not allowed")

	ErrIncorrectEnvironmentConfig = errors.New("incorrect environment config")

	ErrInvalidAudienceProvider = errors.New("invalid audience provider")

	ErrInvalidIdentityProvider = errors.New("invalid identity provider")
)
