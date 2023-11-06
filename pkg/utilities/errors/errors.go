// Package errors provides an interface for handling errors, to ensure consistency across the company
// and flexibility with regards to the error package used.
package errors

// ErrorUtility is an interface containing the methods deemed useful for proper error handling
type ErrorUtility interface {
	Wrap(err error, msg string) error
	WrapWithDepth(depth int, err error, msg string) error
	Wrapf(err error, format string, args ...interface{}) error
	WrapWithDepthf(depth int, err error, format string, args ...interface{}) error
	Errorf(format string, args ...interface{}) error
	ErrorWithDepthf(depth int, format string, args ...interface{}) error
	WithStack(err error) error
	WithStackDepth(err error, depth int) error
	WithSecondaryError(err, otherErr error) error
	CombineErrors(err, otherErr error) error
	Is(err, reference error) bool
	IsAny(err error, references ...error) bool
	Cause(err error) error
	New(message string) error
}

var errorUtility ErrorUtility //nolint:gochecknoglobals

//nolint:gochecknoinits
func init() {
	errorUtility = &CockroachErrorUtility{}
}

// SetErrorUtility sets the ErrorUtility implementation that is used in the methods exported by this package
func SetErrorUtility(utility ErrorUtility) {
	errorUtility = utility
}

// Wrap wraps an error with a message prefix.
// A stack trace is retained.
func Wrap(err error, msg string) error {
	return errorUtility.WrapWithDepth(1, err, msg)
}

// WrapWithDepth is like Wrap except the depth to capture the stack
// trace is configurable.
// See the doc of `Wrap()` for more details.
func WrapWithDepth(depth int, err error, msg string) error {
	return errorUtility.WrapWithDepth(depth+1, err, msg)
}

// Wrapf wraps an error with a formatted message prefix. A stack
// trace is also retained. If the format is empty, no prefix is added,
// but the extra arguments are still processed for reportable strings.
func Wrapf(err error, format string, args ...interface{}) error {
	return errorUtility.WrapWithDepthf(1, err, format, args...)
}

// WrapWithDepthf is like Wrapf except the depth to capture the stack
// trace is configurable.
// See the doc of `Wrapf()` for more details.
func WrapWithDepthf(depth int, err error, format string, args ...interface{}) error {
	return errorUtility.WrapWithDepthf(depth+1, err, format, args...)
}

// Errorf creates an error with a formatted error message. Using %w in the format wraps an error inside.
// A stack trace is retained.
func Errorf(format string, args ...interface{}) error {
	return errorUtility.ErrorWithDepthf(1, format, args...)
}

// ErrorWithDepthf is like Errorf except the depth to capture the stack
// trace is configurable.
// See the doc of `Errorf()` for more details.
func ErrorWithDepthf(depth int, format string, args ...interface{}) error {
	return errorUtility.ErrorWithDepthf(depth+1, format, args...)
}

// WithSecondaryError enhances the error given as first argument with
// an annotation that carries the error given as second argument.  The
// second error does not participate in cause analysis (Is, etc) and
// is only revealed when printing out the error or collecting safe
// (PII-free) details for reporting.
//
// If additionalErr is nil, the first error is returned as-is.
func WithSecondaryError(err, otherErr error) error {
	return errorUtility.WithSecondaryError(err, otherErr)
}

// CombineErrors returns err, or, if err is nil, otherErr.
// if err is non-nil, otherErr is attached as secondary error.
// See the documentation of `WithSecondaryError()` for details.
func CombineErrors(err, otherErr error) error {
	return errorUtility.CombineErrors(err, otherErr)
}

// Is determines whether one of the causes of the given error or any
// of its causes is equivalent to some reference error.
func Is(err, reference error) bool {
	return errorUtility.Is(err, reference)
}

// IsAny is like Is except that multiple references are compared.
func IsAny(err error, references ...error) bool {
	return errorUtility.IsAny(err, references...)
}

// Cause returns the wrapped error. Most common use case is switches.
// In a single if statement, use `Is()` instead
func Cause(err error) error {
	return errorUtility.Cause(err)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// Use only if not wrapping. Otherwise, use Wrap.
func WithStack(err error) error {
	return errorUtility.WithStackDepth(err, 1)
}

// WithStackDepth annotates err with a stack trace starting from the
// given call depth. The value zero identifies the caller
// of WithStackDepth itself.
// See the documentation of WithStack() for more details.
func WithStackDepth(err error, depth int) error {
	return errorUtility.WithStackDepth(err, depth)
}

// New creates a new error with a message and no stacktrace.
func New(message string) error {
	return errorUtility.New(message)
}
