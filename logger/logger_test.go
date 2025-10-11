package logger

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Level != LogLevelInfo {
		t.Errorf("Expected default level to be info, got %s", cfg.Level)
	}

	if cfg.Format != LogFormatConsole {
		t.Errorf("Expected default format to be console, got %s", cfg.Format)
	}

	if cfg.EnableCaller {
		t.Error("Expected EnableCaller to be false by default")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid console config",
			config: &Config{
				Level:  LogLevelInfo,
				Format: LogFormatConsole,
			},
			wantErr: false,
		},
		{
			name: "valid json config",
			config: &Config{
				Level:  LogLevelDebug,
				Format: LogFormatJSON,
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			config: &Config{
				Level:  LogLevel("invalid"),
				Format: LogFormatJSON,
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			config: &Config{
				Level:  LogLevelInfo,
				Format: LogFormat("invalid"),
			},
			wantErr: true,
		},
		{
			name: "file format without path",
			config: &Config{
				Level:  LogLevelInfo,
				Format: LogFormatFile,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSetDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.SetDefaults()

	if cfg.Level != LogLevelInfo {
		t.Errorf("Expected default level to be info, got %s", cfg.Level)
	}

	if cfg.Format != LogFormatConsole {
		t.Errorf("Expected default format to be console, got %s", cfg.Format)
	}

	if cfg.Encoding != "console" {
		t.Errorf("Expected default encoding to be console, got %s", cfg.Encoding)
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "nil config uses defaults",
			config:  nil,
			wantErr: false,
		},
		{
			name: "console logger",
			config: &Config{
				Level:  LogLevelInfo,
				Format: LogFormatConsole,
			},
			wantErr: false,
		},
		{
			name: "json logger",
			config: &Config{
				Level:  LogLevelDebug,
				Format: LogFormatJSON,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("Expected logger to be non-nil")
			}
			if logger != nil {
				if logger.Base() == nil {
					t.Error("Expected base logger to be non-nil")
				}
				if logger.Sugar() == nil {
					t.Error("Expected sugar logger to be non-nil")
				}
			}
		})
	}
}

func TestLoggerWithContext(t *testing.T) {
	logger, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test with trace ID
	ctx := ContextWithTraceID(context.Background(), "test-trace-123")
	loggerWithCtx := logger.WithContext(ctx)

	if loggerWithCtx == nil {
		t.Error("Expected logger with context to be non-nil")
	}

	// Test with span ID
	ctx = ContextWithSpanID(ctx, "test-span-456")
	loggerWithCtx = logger.WithContext(ctx)

	if loggerWithCtx == nil {
		t.Error("Expected logger with context to be non-nil")
	}

	// Test with nil context
	loggerWithCtx = logger.WithContext(nil)
	if loggerWithCtx != logger {
		t.Error("Expected logger with nil context to return same logger")
	}
}

func TestLoggerWith(t *testing.T) {
	logger, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	childLogger := logger.With(zap.String("key", "value"))
	if childLogger == nil {
		t.Error("Expected child logger to be non-nil")
	}

	if childLogger == logger {
		t.Error("Expected child logger to be different from parent")
	}
}

func TestLoggerLevels(t *testing.T) {
	logger, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test all log levels (these shouldn't panic or error)
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	logger.DPanic("dpanic message")

	// Note: We don't test Panic and Fatal as they would terminate the test
}

func TestFileLogger(t *testing.T) {
	// Create a temporary directory for test logs
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := &Config{
		Level:    LogLevelInfo,
		Format:   LogFormatFile,
		FilePath: logFile,
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	// Log some messages
	logger.Info("test message 1")
	logger.Warn("test message 2")

	// Sync to ensure writes are flushed
	if err := logger.Sync(); err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}

	// Read file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Check if log messages are in the file
	if len(content) == 0 {
		t.Error("Log file is empty")
	}
}

func TestFileLoggerWithDirectory(t *testing.T) {
	// Create a temporary directory for test logs
	tmpDir := t.TempDir()

	cfg := &Config{
		Level:    LogLevelInfo,
		Format:   LogFormatFile,
		FilePath: tmpDir, // Pass directory instead of file
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	// Log a message
	logger.Info("test message")

	// Sync to ensure writes are flushed
	if err := logger.Sync(); err != nil {
		t.Errorf("Failed to sync logger: %v", err)
	}

	// Check if default log file was created in the directory
	logFile := filepath.Join(tmpDir, "app.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Default log file (app.log) was not created in directory")
	}
}

func TestGlobalLogger(t *testing.T) {
	// Test global logger initialization
	if GetGlobalLogger() == nil {
		t.Error("Global logger should be initialized")
	}

	// Test setting global logger
	cfg := &Config{
		Level:  LogLevelDebug,
		Format: LogFormatJSON,
	}

	if err := InitFromConfig(cfg); err != nil {
		t.Fatalf("Failed to initialize from config: %v", err)
	}

	// Test global logger functions
	Info("test info")
	Debug("test debug")
	Warn("test warn")
	Error("test error")

	// Test Base and Sugar
	if Base() == nil {
		t.Error("Base() should return non-nil")
	}

	if Sugar() == nil {
		t.Error("Sugar() should return non-nil")
	}

	// Test WithContext
	ctx := ContextWithTraceID(context.Background(), "test-trace")
	loggerWithCtx := WithContext(ctx)
	if loggerWithCtx == nil {
		t.Error("WithContext() should return non-nil")
	}

	// Test With
	childLogger := With(zap.String("key", "value"))
	if childLogger == nil {
		t.Error("With() should return non-nil")
	}
}

func TestSugarLogger(t *testing.T) {
	logger, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	sugar := logger.Sugar()
	if sugar == nil {
		t.Fatal("Sugar logger should not be nil")
	}

	// Test sugar logger methods
	sugar.Infow("test message", "key", "value")
	sugar.Debugf("test debug: %s", "formatted")
	sugar.Warnw("test warning", "count", 42)
}

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()

	// Test adding trace ID
	ctx = ContextWithTraceID(ctx, "trace-123")
	if traceID := ctx.Value(traceIDKey); traceID != "trace-123" {
		t.Errorf("Expected trace ID to be 'trace-123', got %v", traceID)
	}

	// Test adding span ID
	ctx = ContextWithSpanID(ctx, "span-456")
	if spanID := ctx.Value(spanIDKey); spanID != "span-456" {
		t.Errorf("Expected span ID to be 'span-456', got %v", spanID)
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
	}{
		{"debug", LogLevelDebug},
		{"info", LogLevelInfo},
		{"warn", LogLevelWarn},
		{"error", LogLevelError},
		{"dpanic", LogLevelDPanic},
		{"panic", LogLevelPanic},
		{"fatal", LogLevelFatal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseLevel(tt.level)
			if err != nil {
				t.Errorf("parseLevel() error = %v", err)
			}
		})
	}
}
