// Package errors provides different grpc error implementations
//
// nolint: grpcerr
package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
)

// ErrorFunc error function interface
type ErrorFunc func(string, ...protoiface.MessageV1) error

const (
	internalServerMessage = "internal server error"
)

// ErrCanceled indicates the operation was canceled (typically by the caller).
func ErrCanceled(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.Canceled, message, details...)
}

// ErrUnknown error. An example of where this error may be returned is
// if a Status value received from another address space belongs to
// an error-space that is not known in this address space. Also
// errors raised by APIs that do not return enough error information
// may be converted to this error.
func ErrUnknown(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.Unknown, message, details...)
}

// ErrInvalidArgument indicates client specified an invalid argument.
// ErrNote that this differs from FailedPrecondition. It indicates arguments
// that are problematic regardless of the state of the system
// (e.g., a malformed file name).
func ErrInvalidArgument(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.InvalidArgument, message, details...)
}

// ErrDeadlineExceeded means operation expired before completion.
// For operations that change the state of the system, this error may be
// returned even if the operation has completed successfully. For
// example, a successful response from a server could have been delayed
// long enough for the deadline to expire.
func ErrDeadlineExceeded(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.DeadlineExceeded, message, details...)
}

// ErrNotFound means some requested entity (e.g., file or directory) was
// not found.
func ErrNotFound(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.NotFound, message, details...)
}

// ErrAlreadyExists means an attempt to create an entity failed because one
// already exists.
func ErrAlreadyExists(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.AlreadyExists, message, details...)
}

// ErrPermissionDenied indicates the caller does not have permission to
// execute the specified operation. It must not be used for rejections
// caused by exhausting some resource (use ResourceExhausted
// instead for those errors). It must not be
// used if the caller cannot be identified (use Unauthenticated
// instead for those errors).
func ErrPermissionDenied(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.PermissionDenied, message, details...)
}

// ErrResourceExhausted indicates some resource has been exhausted, perhaps
// a per-user quota, or perhaps the entire file system is out of space.
func ErrResourceExhausted(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.ResourceExhausted, message, details...)
}

// ErrFailedPrecondition indicates operation was rejected because the
// system is not in a state required for the operation's execution.
// For example, directory to be deleted may be non-empty, an rmdir
// operation is applied to a non-directory, etc.
//
// A litmus test that may help a service implementor in deciding
// between FailedPrecondition, Aborted, and Unavailable:
// (a) Use Unavailable if the client can retry just the failing call.
// (b) Use Aborted if the client should retry at a higher-level
// (e.g., restarting a read-modify-write sequence).
// (c) Use FailedPrecondition if the client should not retry until
// the system state has been explicitly fixed. E.g., if an "rmdir"
// fails because the directory is non-empty, FailedPrecondition
// should be returned since the client should not retry unless
// they have first fixed up the directory by deleting files from it.
// (d) Use FailedPrecondition if the client performs conditional
// REST Get/Update/Delete on a resource and the resource on the
// server does not match the condition. E.g., conflicting
// read-modify-write on the same resource.
func ErrFailedPrecondition(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.FailedPrecondition, message, details...)
}

// ErrAborted indicates the operation was aborted, typically due to a
// concurrency issue like sequencer check failures, transaction aborts,
// etc.
//
// See litmus test above for deciding between FailedPrecondition,
// Aborted, and Unavailable.
func ErrAborted(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.Aborted, message, details...)
}

// ErrOutOfRange means operation was attempted past the valid range.
// E.g., seeking or reading past end of file.
//
// ErrUnlike InvalidArgument, this error indicates a problem that may
// be fixed if the system state changes. For example, a 32-bit file
// system will generate InvalidArgument if asked to read at an
// offset that is not in the range [0,2^32-1], but it will generate
// OutOfRange if asked to read from an offset past the current
// file size.
//
// There is a fair bit of overlap between FailedPrecondition and
// OutOfRange. We recommend using OutOfRange (the more specific
// error) when it applies so that callers who are iterating through
// a space can easily look for an OutOfRange error to detect when
// they are done.
func ErrOutOfRange(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.OutOfRange, message, details...)
}

// ErrUnimplemented indicates operation is not implemented or not
// supported/enabled in this service.
func ErrUnimplemented(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.Unimplemented, message, details...)
}

// ErrInternal errors. Means some invariants expected by underlying
// system has been broken. If you see one of these errors,
// something is very broken.
//
// Deprecated: Use ErrInternalServer instead
func ErrInternal(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.Internal, message, details...)
}

// ErrInternalServer errors. Means some invariants expected by underlying
// system has been broken. If you see one of these errors,
// something is very broken.
func ErrInternalServer() error {
	return NewStatusWithDetails(codes.Internal, internalServerMessage)
}

// ErrUnavailable indicates the service is currently unavailable.
// This is a most likely a transient condition and may be corrected
// by retrying with a backoff.
//
// See litmus test above for deciding between FailedPrecondition,
// Aborted, and Unavailable.
func ErrUnavailable(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.Unavailable, message, details...)
}

// ErrDataLoss indicates unrecoverable data loss or corruption.
func ErrDataLoss(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.DataLoss, message, details...)
}

// ErrUnauthenticated indicates the request does not have valid
// authentication credentials for the operation.
func ErrUnauthenticated(message string, details ...protoiface.MessageV1) error {
	return NewStatusWithDetails(codes.Unauthenticated, message, details...)
}

// NewStatusWithDetails create status error from given code msg and error type
func NewStatusWithDetails(code codes.Code, mesg string, details ...protoiface.MessageV1) error {
	status := status.New(code, mesg)
	var err error
	status, err = status.WithDetails(details...)
	if err != nil {
		panic(err)
	}
	return status.Err()
}

// GetDetails retrieve details from status
func GetDetails(st *status.Status) (*protoiface.MessageV1, bool) {
	details := st.Details()
	if len(details) != 1 {
		return nil, false
	}
	d := details[0]
	detail, ok := d.(protoiface.MessageV1)
	if !ok {
		return nil, false
	}
	return &detail, true
}
