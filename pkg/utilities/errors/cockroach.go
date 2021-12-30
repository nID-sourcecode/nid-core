package errors

import (
	goErr "errors"

	cockroach "github.com/cockroachdb/errors"
)

// CockroachErrorUtility is the cockroachdb/errors implementation of ErrorUtility
type CockroachErrorUtility struct{}

// Wrap wraps an error with a message prefix.
// A stack trace is retained.
func (c *CockroachErrorUtility) Wrap(err error, msg string) error {
	return cockroach.WrapWithDepth(1, err, msg)
}

// WrapWithDepth is like Wrap except the depth to capture the stack
// trace is configurable.
// See the doc of `Wrap()` for more details.
func (c *CockroachErrorUtility) WrapWithDepth(depth int, err error, msg string) error {
	return cockroach.WrapWithDepth(depth+1, err, msg)
}

// Wrapf wraps an error with a formatted message prefix. A stack
// trace is also retained. If the format is empty, no prefix is added,
// but the extra arguments are still processed for reportable strings.
func (c *CockroachErrorUtility) Wrapf(err error, format string, args ...interface{}) error {
	return cockroach.WrapWithDepthf(1, err, format, args...)
}

// WrapWithDepthf is like Wrapf except the depth to capture the stack
// trace is configurable.
// The the doc of `Wrapf()` for more details.
func (c *CockroachErrorUtility) WrapWithDepthf(depth int, err error, format string, args ...interface{}) error {
	return cockroach.WrapWithDepthf(depth+1, err, format, args...)
}

// Errorf creates an error with a formatted error message. Using %w in the format wraps an error inside.
// A stack trace is retained.
func (c *CockroachErrorUtility) Errorf(format string, args ...interface{}) error {
	return cockroach.NewWithDepthf(1, format, args...)
}

// ErrorWithDepthf is like Errorf except the depth to capture the stack
// trace is configurable.
// The the doc of `Errorf()` for more details.
func (c *CockroachErrorUtility) ErrorWithDepthf(depth int, format string, args ...interface{}) error {
	return cockroach.NewWithDepthf(depth+1, format, args...)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// Use only if not wrapping. Otherwise, use Wrap.
func (c *CockroachErrorUtility) WithStack(err error) error {
	return cockroach.WithStackDepth(err, 1)
}

// WithStackDepth annotates err with a stack trace starting from the
// given call depth. The value zero identifies the caller
// of WithStackDepth itself.
// See the documentation of WithStack() for more details.
func (c *CockroachErrorUtility) WithStackDepth(err error, depth int) error {
	return cockroach.WithStackDepth(err, depth+1)
}

// WithSecondaryError enhances the error given as first argument with
// an annotation that carries the error given as second argument.  The
// second error does not participate in cause analysis (Is, etc) and
// is only revealed when printing out the error or collecting safe
// (PII-free) details for reporting.
//
// If additionalErr is nil, the first error is returned as-is.
func (c *CockroachErrorUtility) WithSecondaryError(err, otherErr error) error {
	return cockroach.WithSecondaryError(err, otherErr)
}

// CombineErrors returns err, or, if err is nil, otherErr.
// if err is non-nil, otherErr is attached as secondary error.
// See the documentation of `WithSecondaryError()` for details.
func (c *CockroachErrorUtility) CombineErrors(err, otherErr error) error {
	return cockroach.CombineErrors(err, otherErr)
}

// Is determines whether one of the causes of the given error or any
// of its causes is equivalent to some reference error.
func (c *CockroachErrorUtility) Is(err, reference error) bool {
	return cockroach.Is(err, reference)
}

// IsAny is like Is except that multiple references are compared.
func (c *CockroachErrorUtility) IsAny(err error, references ...error) bool {
	return cockroach.IsAny(err, references...)
}

// Cause returns the wrapped error. Most common use case is switches.
// In a single if statement, use `Is()` instead
func (c *CockroachErrorUtility) Cause(err error) error {
	return cockroach.Cause(err)
}

// New creates a new error with a message and no stacktrace.
func (c *CockroachErrorUtility) New(message string) error {
	return goErr.New(message) //nolint:goerr113
}
