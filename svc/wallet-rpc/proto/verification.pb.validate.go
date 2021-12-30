// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: verification.proto

package proto

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _verification_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on VerifyRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *VerifyRequest) Validate() error {
	if m == nil {
		return nil
	}

	if err := m._validateUuid(m.GetId()); err != nil {
		return VerifyRequestValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
	}

	// no validation rules for Code

	return nil
}

func (m *VerifyRequest) _validateUuid(uuid string) error {
	if matched := _verification_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// VerifyRequestValidationError is the validation error returned by
// VerifyRequest.Validate if the designated constraints aren't met.
type VerifyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e VerifyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e VerifyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e VerifyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e VerifyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e VerifyRequestValidationError) ErrorName() string { return "VerifyRequestValidationError" }

// Error satisfies the builtin error interface
func (e VerifyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sVerifyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = VerifyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = VerifyRequestValidationError{}

// Validate checks the field values on RetryVerifyRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RetryVerifyRequest) Validate() error {
	if m == nil {
		return nil
	}

	if err := m._validateUuid(m.GetId()); err != nil {
		return RetryVerifyRequestValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
	}

	return nil
}

func (m *RetryVerifyRequest) _validateUuid(uuid string) error {
	if matched := _verification_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// RetryVerifyRequestValidationError is the validation error returned by
// RetryVerifyRequest.Validate if the designated constraints aren't met.
type RetryVerifyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RetryVerifyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RetryVerifyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RetryVerifyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RetryVerifyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RetryVerifyRequestValidationError) ErrorName() string {
	return "RetryVerifyRequestValidationError"
}

// Error satisfies the builtin error interface
func (e RetryVerifyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRetryVerifyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RetryVerifyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RetryVerifyRequestValidationError{}

// Validate checks the field values on VerifyResponse with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *VerifyResponse) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	return nil
}

// VerifyResponseValidationError is the validation error returned by
// VerifyResponse.Validate if the designated constraints aren't met.
type VerifyResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e VerifyResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e VerifyResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e VerifyResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e VerifyResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e VerifyResponseValidationError) ErrorName() string { return "VerifyResponseValidationError" }

// Error satisfies the builtin error interface
func (e VerifyResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sVerifyResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = VerifyResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = VerifyResponseValidationError{}

// Validate checks the field values on RetryPhoneRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *RetryPhoneRequest) Validate() error {
	if m == nil {
		return nil
	}

	if err := m._validateUuid(m.GetId()); err != nil {
		return RetryPhoneRequestValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
	}

	// no validation rules for VerificationType

	return nil
}

func (m *RetryPhoneRequest) _validateUuid(uuid string) error {
	if matched := _verification_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// RetryPhoneRequestValidationError is the validation error returned by
// RetryPhoneRequest.Validate if the designated constraints aren't met.
type RetryPhoneRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RetryPhoneRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RetryPhoneRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RetryPhoneRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RetryPhoneRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RetryPhoneRequestValidationError) ErrorName() string {
	return "RetryPhoneRequestValidationError"
}

// Error satisfies the builtin error interface
func (e RetryPhoneRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRetryPhoneRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RetryPhoneRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RetryPhoneRequestValidationError{}
