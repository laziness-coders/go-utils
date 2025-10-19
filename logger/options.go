package logger

import "github.com/laziness-coders/go-utils/generic"

// Option is a functional option for configuring the logger.
type Option func(*Config)

// WithLevel sets the log level.
func WithLevel(level string) Option {
	return func(cfg *Config) {
		cfg.Level = level
	}
}

// WithFormat sets the log format.
func WithFormat(format string) Option {
	return func(cfg *Config) {
		cfg.Format = format
	}
}

// WithFilePath sets the file path for file output.
func WithFilePath(filePath string) Option {
	return func(cfg *Config) {
		cfg.FilePath = filePath
	}
}

// WithMaxSize sets the maximum size in MB before rotation.
func WithMaxSize(maxSize int) Option {
	return func(cfg *Config) {
		cfg.MaxSize = maxSize
	}
}

// WithMaxAge sets the maximum age in days to retain old logs.
func WithMaxAge(maxAge int) Option {
	return func(cfg *Config) {
		cfg.MaxAge = maxAge
	}
}

// WithMaxBackups sets the maximum number of old log files to retain.
func WithMaxBackups(maxBackups int) Option {
	return func(cfg *Config) {
		cfg.MaxBackups = maxBackups
	}
}

// WithCompress enables compression of rotated files.
func WithCompress(compress bool) Option {
	return func(cfg *Config) {
		cfg.Compress = compress
	}
}

// WithDev sets whether to use zap's development config for human-readable output.
func WithDev(isDev bool) Option {
	return func(cfg *Config) {
		cfg.IsDev = generic.ToPointer(isDev)
	}
}

// Predefined option sets for common configurations

// WithConsoleOutput configures console output (JSON format).
func WithConsoleOutput() Option {
	return func(cfg *Config) {
		cfg.Format = FormatJSON
		cfg.FilePath = ""
	}
}

// WithTextOutput configures text output to console.
func WithTextOutput() Option {
	return func(cfg *Config) {
		cfg.Format = FormatText
		cfg.FilePath = ""
	}
}

// WithFileOutput configures file output with rotation.
func WithFileOutput(filePath string) Option {
	return func(cfg *Config) {
		cfg.Format = FormatFile
		cfg.FilePath = filePath
	}
}

// WithDebugLevel configures debug level logging.
func WithDebugLevel() Option {
	return func(cfg *Config) {
		cfg.Level = LevelDebug
	}
}

// WithInfoLevel configures info level logging.
func WithInfoLevel() Option {
	return func(cfg *Config) {
		cfg.Level = LevelInfo
	}
}

// WithWarnLevel configures warn level logging.
func WithWarnLevel() Option {
	return func(cfg *Config) {
		cfg.Level = LevelWarn
	}
}

// WithErrorLevel configures error level logging.
func WithErrorLevel() Option {
	return func(cfg *Config) {
		cfg.Level = LevelError
	}
}

// WithCallerSkip sets the number of stack frames to skip for caller info.
func WithCallerSkip(skip int) Option {
	return func(cfg *Config) {
		cfg.CallerSkip = skip
	}
}

// WithProductionDefaults sets production-ready defaults.
func WithProductionDefaults() Option {
	return func(cfg *Config) {
		cfg.Level = LevelInfo
		cfg.Format = FormatFile
		cfg.FilePath = "logs/app.log"
		cfg.MaxSize = 100
		cfg.MaxAge = 30
		cfg.MaxBackups = 3
		cfg.Compress = true
	}
}
