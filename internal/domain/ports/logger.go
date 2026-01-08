package ports

// LogLevel represents the severity of a log entry.
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// Logger defines the interface for structured logging.
type Logger interface {
	// Log writes a structured log entry.
	Log(level LogLevel, message string, fields map[string]interface{})
	// LogPluginRun writes a log entry for a plugin execution.
	LogPluginRun(pluginName string, args string, output []byte, exitCode int)
}
