// Package errors provides different grpc error implementations
//
package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpcerrpb "lab.weave.nl/nid/nid-core/pkg/utilities/grpcserver/errors/proto"
)

// ErrorFunc error function interface
type ErrorFunc func(interface{}, ...interface{}) error

const (
	internalServerMessage = "internal server error"
)

// ErrCanceled indicates the operation was canceled (typically by the caller).
func ErrCanceled(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.Canceled, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_CANCELED})
}

// ErrUnknown error. An example of where this error may be returned is
// if a Status value received from another address space belongs to
// an error-space that is not known in this address space. Also
// errors raised by APIs that do not return enough error information
// may be converted to this error.
func ErrUnknown(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.Unknown, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_UNKNOWN})
}

// ErrInvalidArgument indicates client specified an invalid argument.
// ErrNote that this differs from FailedPrecondition. It indicates arguments
// that are problematic regardless of the state of the system
// (e.g., a malformed file name).
func ErrInvalidArgument(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.InvalidArgument, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_INVALID_ARGUMENT})
}

// ErrDeadlineExceeded means operation expired before completion.
// For operations that change the state of the system, this error may be
// returned even if the operation has completed successfully. For
// example, a successful response from a server could have been delayed
// long enough for the deadline to expire.
func ErrDeadlineExceeded(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.DeadlineExceeded, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_DEADLINE_EXCEEDED})
}

// ErrNotFound means some requested entity (e.g., file or directory) was
// not found.
func ErrNotFound(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.NotFound, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_NOT_FOUND})
}

// ErrAlreadyExists means an attempt to create an entity failed because one
// already exists.
func ErrAlreadyExists(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.AlreadyExists, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_ALREADY_EXISTS})
}

// ErrPermissionDenied indicates the caller does not have permission to
// execute the specified operation. It must not be used for rejections
// caused by exhausting some resource (use ResourceExhausted
// instead for those errors). It must not be
// used if the caller cannot be identified (use Unauthenticated
// instead for those errors).
func ErrPermissionDenied(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.PermissionDenied, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_PERMISSION_DENIED})
}

// ErrResourceExhausted indicates some resource has been exhausted, perhaps
// a per-user quota, or perhaps the entire file system is out of space.
func ErrResourceExhausted(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.ResourceExhausted, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_RESOURCE_EXHAUSTED})
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
func ErrFailedPrecondition(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.FailedPrecondition, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_FAILED_PRECONDITION})
}

// ErrAborted indicates the operation was aborted, typically due to a
// concurrency issue like sequencer check failures, transaction aborts,
// etc.
//
// See litmus test above for deciding between FailedPrecondition,
// Aborted, and Unavailable.
func ErrAborted(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.Aborted, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_ABORTED})
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
func ErrOutOfRange(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.OutOfRange, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_OUT_OF_RANGE})
}

// ErrUnimplemented indicates operation is not implemented or not
// supported/enabled in this service.
func ErrUnimplemented(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.Unimplemented, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_UNIMPLEMENTED})
}

// ErrInternal errors. Means some invariants expected by underlying
// system has been broken. If you see one of these errors,
// something is very broken.
//
// Deprecated: Use ErrInternalServer instead
func ErrInternal(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.Internal, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_INTERNAL})
}

// ErrInternalServer errors. Means some invariants expected by underlying
// system has been broken. If you see one of these errors,
// something is very broken.
func ErrInternalServer() error {
	return NewStatusWithDetails(codes.Internal, formatErrMessage(internalServerMessage), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_INTERNAL})
}

// ErrUnavailable indicates the service is currently unavailable.
// This is a most likely a transient condition and may be corrected
// by retrying with a backoff.
//
// See litmus test above for deciding between FailedPrecondition,
// Aborted, and Unavailable.
func ErrUnavailable(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.Unavailable, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_UNAVAILABLE})
}

// ErrDataLoss indicates unrecoverable data loss or corruption.
func ErrDataLoss(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.DataLoss, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_DATA_LOSS})
}

// ErrUnauthenticated indicates the request does not have valid
// authentication credentials for the operation.
func ErrUnauthenticated(message interface{}, keyvals ...interface{}) error {
	return NewStatusWithDetails(codes.Unauthenticated, formatErrMessage(message, keyvals...), &grpcerrpb.ErrDetails{ErrType: grpcerrpb.ErrorType_UNAUTHENTICATED})
}

func formatErrMessage(message interface{}, keyvals ...interface{}) string {
	var msg string
	switch actual := message.(type) {
	case string:
		msg = actual
	case error:
		msg = actual.Error()
	case fmt.Stringer:
		msg = actual.String()
	default:
		msg = fmt.Sprintf("%v", actual)
	}
	var meta map[string]interface{}
	l := len(keyvals)
	if l > 0 {
		meta = make(map[string]interface{})
	}
	for i := 0; i < l; i += 2 {
		k := keyvals[i]
		var v interface{} = "MISSING"
		if i+1 < l {
			v = keyvals[i+1]
		}
		meta[fmt.Sprintf("%v", k)] = v
	}
	return fmt.Sprintf(msg, keyvals...)
}

// NewStatusWithDetails create status error from given code msg and error type
func NewStatusWithDetails(code codes.Code, mesg string, details *grpcerrpb.ErrDetails) error {
	status := status.New(code, mesg)
	var err error
	status, err = status.WithDetails(details)
	if err != nil {
		panic(err)
	}
	return status.Err()
}

// GetDetails retrieve details from status
func GetDetails(st *status.Status) (*grpcerrpb.ErrDetails, bool) {
	details := st.Details()
	if len(details) != 1 {
		return nil, false
	}
	d := details[0]
	detail, ok := d.(*grpcerrpb.ErrDetails)
	if !ok {
		return nil, false
	}
	return detail, true
}
