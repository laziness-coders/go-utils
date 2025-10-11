package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the logger configuration
type Config struct {
	// Level is the minimum enabled logging level
	Level LogLevel `mapstructure:"level" json:"level" yaml:"level"`

	// Format is the log output format (json, console, text, file)
	Format LogFormat `mapstructure:"format" json:"format" yaml:"format"`

	// FilePath is the path to the log file (used when Format is "file")
	// Can be a file path or directory path
	FilePath string `mapstructure:"file_path" json:"file_path" yaml:"file_path"`

	// EnableCaller enables logging the caller information (file and line number)
	EnableCaller bool `mapstructure:"enable_caller" json:"enable_caller" yaml:"enable_caller"`

	// EnableStacktrace enables automatic stacktrace capturing for logs at or above a certain level
	EnableStacktrace bool `mapstructure:"enable_stacktrace" json:"enable_stacktrace" yaml:"enable_stacktrace"`

	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `mapstructure:"development" json:"development" yaml:"development"`

	// Encoding sets the logger's encoding. Valid values are "json" and "console"
	Encoding string `mapstructure:"encoding" json:"encoding" yaml:"encoding"`
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:            LogLevelInfo,
		Format:           LogFormatConsole,
		EnableCaller:     false,
		EnableStacktrace: false,
		Development:      false,
		Encoding:         "console",
	}
}

// Validate validates the logger configuration
func (c *Config) Validate() error {
	// Validate log level
	switch c.Level {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelDPanic, LogLevelPanic, LogLevelFatal:
		// valid
	default:
		return fmt.Errorf("invalid log level: %s", c.Level)
	}

	// Validate log format
	switch c.Format {
	case LogFormatJSON, LogFormatConsole, LogFormatText, LogFormatFile:
		// valid
	default:
		return fmt.Errorf("invalid log format: %s", c.Format)
	}

	// If format is file, validate file path
	if c.Format == LogFormatFile {
		if c.FilePath == "" {
			return fmt.Errorf("file_path is required when format is 'file'")
		}

		// Check if path is a directory or file
		info, err := os.Stat(c.FilePath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat file path: %w", err)
		}

		// If path exists and is a directory, create a default log file name
		if err == nil && info.IsDir() {
			c.FilePath = filepath.Join(c.FilePath, "app.log")
		}

		// Ensure the directory exists
		dir := filepath.Dir(c.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	return nil
}

// SetDefaults sets default values for empty fields
func (c *Config) SetDefaults() {
	if c.Level == "" {
		c.Level = LogLevelInfo
	}

	if c.Format == "" {
		c.Format = LogFormatConsole
	}

	// Set encoding based on format if not explicitly set
	if c.Encoding == "" {
		switch c.Format {
		case LogFormatJSON, LogFormatFile:
			c.Encoding = "json"
		case LogFormatConsole, LogFormatText:
			c.Encoding = "console"
		default:
			c.Encoding = "console"
		}
	}
}
