// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: autobsn.proto

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
var _autobsn_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on ReplacePlaceholderWithBSNRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *ReplacePlaceholderWithBSNRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Body

	// no validation rules for Query

	// no validation rules for Method

	// no validation rules for AuthorizationHeader

	return nil
}

// ReplacePlaceholderWithBSNRequestValidationError is the validation error
// returned by ReplacePlaceholderWithBSNRequest.Validate if the designated
// constraints aren't met.
type ReplacePlaceholderWithBSNRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReplacePlaceholderWithBSNRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReplacePlaceholderWithBSNRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReplacePlaceholderWithBSNRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReplacePlaceholderWithBSNRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReplacePlaceholderWithBSNRequestValidationError) ErrorName() string {
	return "ReplacePlaceholderWithBSNRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ReplacePlaceholderWithBSNRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReplacePlaceholderWithBSNRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReplacePlaceholderWithBSNRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReplacePlaceholderWithBSNRequestValidationError{}

// Validate checks the field values on ReplacePlaceholderWithBSNResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *ReplacePlaceholderWithBSNResponse) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Body

	// no validation rules for Query

	return nil
}

// ReplacePlaceholderWithBSNResponseValidationError is the validation error
// returned by ReplacePlaceholderWithBSNResponse.Validate if the designated
// constraints aren't met.
type ReplacePlaceholderWithBSNResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReplacePlaceholderWithBSNResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReplacePlaceholderWithBSNResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReplacePlaceholderWithBSNResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReplacePlaceholderWithBSNResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReplacePlaceholderWithBSNResponseValidationError) ErrorName() string {
	return "ReplacePlaceholderWithBSNResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ReplacePlaceholderWithBSNResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReplacePlaceholderWithBSNResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReplacePlaceholderWithBSNResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReplacePlaceholderWithBSNResponseValidationError{}
