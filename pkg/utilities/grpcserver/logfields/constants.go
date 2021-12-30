// Package logfields contains the fields to be used in the grpc withfield function
package logfields

// Different log field identifiers
const (
	Bearer         = "bearer"
	BearerError    = "bearer_error"
	Context        = "context"
	LogLevel       = "log_level"
	MetadataError  = "metadata_error"
	Port           = "port"
	RecoveredPanic = "recovered_panic"
	Stack          = "stack"
)
