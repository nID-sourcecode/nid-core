// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pseudonymizer.proto

package pseudonymization

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
var _pseudonymizer_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on ConvertRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ConvertRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for NamespaceTo

	return nil
}

// ConvertRequestValidationError is the validation error returned by
// ConvertRequest.Validate if the designated constraints aren't met.
type ConvertRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConvertRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConvertRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConvertRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConvertRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConvertRequestValidationError) ErrorName() string { return "ConvertRequestValidationError" }

// Error satisfies the builtin error interface
func (e ConvertRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConvertRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConvertRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConvertRequestValidationError{}

// Validate checks the field values on ConvertResponse with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ConvertResponse) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Conversions

	return nil
}

// ConvertResponseValidationError is the validation error returned by
// ConvertResponse.Validate if the designated constraints aren't met.
type ConvertResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConvertResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConvertResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConvertResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConvertResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConvertResponseValidationError) ErrorName() string { return "ConvertResponseValidationError" }

// Error satisfies the builtin error interface
func (e ConvertResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConvertResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConvertResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConvertResponseValidationError{}

// Validate checks the field values on GenerateRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *GenerateRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Amount

	return nil
}

// GenerateRequestValidationError is the validation error returned by
// GenerateRequest.Validate if the designated constraints aren't met.
type GenerateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateRequestValidationError) ErrorName() string { return "GenerateRequestValidationError" }

// Error satisfies the builtin error interface
func (e GenerateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateRequestValidationError{}

// Validate checks the field values on GenerateResponse with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *GenerateResponse) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// GenerateResponseValidationError is the validation error returned by
// GenerateResponse.Validate if the designated constraints aren't met.
type GenerateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateResponseValidationError) ErrorName() string { return "GenerateResponseValidationError" }

// Error satisfies the builtin error interface
func (e GenerateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateResponseValidationError{}