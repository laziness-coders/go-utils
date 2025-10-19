package logger

// Log format constants
const (
	FormatJSON = "json"
	FormatText = "text"
	FormatFile = "file"
)

// Log level constants
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

const (
	defaultMaxSize    = 100 // MB
	defaultMaxAge     = 30  // days
	defaultMaxBackups = 3   // files
	defaultCallerSkip = 1   // skip 1 stack frame

	// Permissions for log directory
	// Means: rwxr-xr-x
	// Owner can read, write, execute; Group can read, execute; Others can read, execute
	logDirPermission = 0o755
)
