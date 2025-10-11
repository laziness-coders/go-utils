package logger

// LogLevel represents the logging level
type LogLevel string

const (
	// LogLevelDebug logs are typically voluminous, and are usually disabled in production.
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo is the default logging priority.
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn logs are more important than Info, but don't need individual human review.
	LogLevelWarn LogLevel = "warn"
	// LogLevelError logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs.
	LogLevelError LogLevel = "error"
	// LogLevelDPanic logs are particularly important errors. In development the logger panics after writing the message.
	LogLevelDPanic LogLevel = "dpanic"
	// LogLevelPanic logs a message, then panics.
	LogLevelPanic LogLevel = "panic"
	// LogLevelFatal logs a message, then calls os.Exit(1).
	LogLevelFatal LogLevel = "fatal"
)

// LogFormat represents the log output format
type LogFormat string

const (
	// LogFormatJSON outputs logs in JSON format
	LogFormatJSON LogFormat = "json"
	// LogFormatConsole outputs logs in human-readable console format
	LogFormatConsole LogFormat = "console"
	// LogFormatText is an alias for console format
	LogFormatText LogFormat = "text"
	// LogFormatFile outputs logs to a file (uses JSON format by default)
	LogFormatFile LogFormat = "file"
)
