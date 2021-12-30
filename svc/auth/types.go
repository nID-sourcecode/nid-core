package main

// ResponseType specifies the supported OAuth2 response types
type ResponseType int

const (
	// CodeResponseType specifies that the client wants to start the authorization code grant type flow
	CodeResponseType ResponseType = iota + 1
)

// String will represent the string equivalent of ResponseType value
func (rt ResponseType) String() string {
	return [...]string{"", "code"}[rt]
}
