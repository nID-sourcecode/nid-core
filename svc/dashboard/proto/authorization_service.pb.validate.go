// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: authorization_service.proto

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
var _authorization_service_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on SigninResponseMessage with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *SigninResponseMessage) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Bearer

	return nil
}

// SigninResponseMessageValidationError is the validation error returned by
// SigninResponseMessage.Validate if the designated constraints aren't met.
type SigninResponseMessageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SigninResponseMessageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SigninResponseMessageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SigninResponseMessageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SigninResponseMessageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SigninResponseMessageValidationError) ErrorName() string {
	return "SigninResponseMessageValidationError"
}

// Error satisfies the builtin error interface
func (e SigninResponseMessageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSigninResponseMessage.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SigninResponseMessageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SigninResponseMessageValidationError{}