package log

// Format specifies the format used for logging
type Format string

const (
	// TextFormat formats the logs as text
	TextFormat Format = "text"
	// JSONFormat formats the logs as JSON
	JSONFormat Format = "json"
	// LogstashFormat formats the logs conform Logstash
	LogstashFormat Format = "logstash"
	// JoonixFormat formats the logs conform Joonix
	JoonixFormat Format = "joonix"
)

func (format Format) isValid() error {
	switch format {
	case TextFormat, JSONFormat, LogstashFormat, JoonixFormat:
		return nil
	}
	return ErrInvalidFormat
}
