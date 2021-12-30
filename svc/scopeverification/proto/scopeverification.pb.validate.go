// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: scopeverification.proto

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
var _scopeverification_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on VerifyRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *VerifyRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for AuthHeader

	// no validation rules for Method

	// no validation rules for Path

	// no validation rules for Body

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