package logger

import (
	"context"
	"go.opentelemetry.io/otel"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "valid json config",
			opts: []Option{
				WithLevel(LevelInfo),
				WithFormat(FormatJSON),
			},
			wantErr: false,
		},
		{
			name: "valid text config",
			opts: []Option{
				WithLevel(LevelDebug),
				WithFormat(FormatText),
			},
			wantErr: false,
		},
		{
			name: "valid file config",
			opts: []Option{
				WithLevel(LevelWarn),
				WithFormat(FormatFile),
				WithFilePath("/tmp/test.log"),
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			opts: []Option{
				WithLevel("invalid"),
				WithFormat(FormatJSON),
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			opts: []Option{
				WithLevel(LevelInfo),
				WithFormat("invalid"),
			},
			wantErr: true,
		},
		{
			name: "file format without path",
			opts: []Option{
				WithLevel(LevelInfo),
				WithFormat(FormatFile),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := newLogger(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("newLogger() returned nil logger when no error expected")
			}
		})
	}
}

func TestLogger_Trace(t *testing.T) {
	logger, _ := newLogger(
		WithLevel(LevelDebug),
		WithFormat(FormatJSON),
	)

	// Test with context containing trace info
	ctx := context.WithValue(context.Background(), "dd.trace_id", "12345")
	ctx = context.WithValue(ctx, "dd.span_id", "67890")

	tracedLogger := logger.Trace(ctx)
	if tracedLogger == nil {
		t.Error("Trace() returned nil logger")
	}

	// Test that the traced logger is different from the original
	if tracedLogger == logger {
		t.Error("Trace() should return a new logger instance")
	}
}

func TestLogger_PrintFunctions(t *testing.T) {
	logger, _ := newLogger(
		WithLevel(LevelDebug),
		WithFormat(FormatJSON),
	)

	// Test Print functions
	logger.Print("test print")
	logger.Println("test println")
	logger.Printf("test printf: %s", "value")
	logger.Debugf("test debugf: %d", 42)
	logger.Infof("test infof: %s", "info")
	logger.Warnf("test warnf: %s", "warning")
	logger.Errorf("test errorf: %s", "error")
}

func TestSugaredLogger_Trace(t *testing.T) {
	logger, _ := newLogger(
		WithLevel(LevelDebug),
		WithFormat(FormatJSON),
	)

	sugar := logger.Sugar()
	if sugar == nil {
		t.Error("Sugar() returned nil")
	}

	// Test trace with sugar logger
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace")
	tracedSugar := sugar.Trace(ctx)
	if tracedSugar == nil {
		t.Error("Trace() on sugar logger returned nil")
	}
}

func TestGlobalLogger(t *testing.T) {
	// Initialize global logger for testing
	err := Init(
		WithLevel(LevelInfo),
		WithFormat(FormatJSON),
	)
	if err != nil {
		t.Fatalf("Failed to initialize global logger: %v", err)
	}

	// Test global logger functions
	Info("test global info")
	Debug("test global debug")
	Warn("test global warn")
	Error("test global error")

	// Test global print functions
	Print("test global print")
	Println("test global println")
	Printf("test global printf: %s", "value")
	Debugf("test global debugf: %d", 42)
	Infof("test global infof: %s", "info")
	Warnf("test global warnf: %s", "warning")
	Errorf("test global errorf: %s", "error")
}

func TestGlobalTrace(t *testing.T) {
	// Initialize global logger for testing
	err := Init(
		WithLevel(LevelInfo),
		WithFormat(FormatJSON),
	)
	if err != nil {
		t.Fatalf("Failed to initialize global logger: %v", err)
	}

	ctx := context.WithValue(context.Background(), "dd.trace_id", "global-trace")
	ctx = context.WithValue(ctx, "dd.span_id", "global-span")

	tracedLogger := Trace(ctx)
	if tracedLogger == nil {
		t.Error("Global Trace() returned nil")
	}
}

func TestGlobalSugar(t *testing.T) {
	// Initialize global logger for testing
	err := Init(
		WithLevel(LevelInfo),
		WithFormat(FormatJSON),
	)
	if err != nil {
		t.Fatalf("Failed to initialize global logger: %v", err)
	}

	sugar := Sugar()
	if sugar == nil {
		t.Error("Global Sugar() returned nil")
	}

	sugar.Infow("test global sugar",
		"key1", "value1",
		"key2", "value2",
	)
}

func TestGlobalWith(t *testing.T) {
	// Initialize global logger for testing
	err := Init(
		WithLevel(LevelInfo),
		WithFormat(FormatJSON),
	)
	if err != nil {
		t.Fatalf("Failed to initialize global logger: %v", err)
	}

	withLogger := With(zap.String("service", "test"))
	if withLogger == nil {
		t.Error("Global With() returned nil")
	}

	withLogger.Info("test with fields")
}

func TestConfig_WithDefaults(t *testing.T) {
	// Create a completely fresh config to avoid interference from other tests
	config, err := newConfig()
	if err != nil {
		t.Fatalf("Failed to create new config: %v", err)
	}

	if config.Level != LevelDebug {
		t.Errorf("Expected default level to be %s, got %s", LevelDebug, config.Level)
	}
	if config.Format != FormatText {
		t.Errorf("Expected default format to be %s, got %s", FormatText, config.Format)
	}
	if config.MaxSize != defaultMaxSize {
		t.Errorf("Expected default MaxSize to be 100, got %d", config.MaxSize)
	}
	if config.MaxAge != defaultMaxAge {
		t.Errorf("Expected default MaxAge to be 30, got %d", config.MaxAge)
	}
	if config.MaxBackups != defaultMaxBackups {
		t.Errorf("Expected default MaxBackups to be 3, got %d", config.MaxBackups)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Level:  LevelInfo,
				Format: FormatJSON,
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			config: Config{
				Level:  "invalid",
				Format: FormatJSON,
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			config: Config{
				Level:  LevelInfo,
				Format: "invalid",
			},
			wantErr: true,
		},
		{
			name: "file format without path",
			config: Config{
				Level:  LevelInfo,
				Format: FormatFile,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractTraceFields(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected int // expected number of fields
	}{
		{
			name:     "nil context",
			ctx:      nil,
			expected: 0,
		},
		{
			name:     "empty context",
			ctx:      context.Background(),
			expected: 0,
		},
		{
			name: "context with trace_id and span_id",
			ctx: func() context.Context {
				ctx := context.Background()
				tracer := otel.Tracer("test-tracer")
				ctx, span := tracer.Start(ctx, "test-span")
				ctx = context.WithValue(ctx, "trace_id", "abc123")
				ctx = context.WithValue(ctx, "span_id", "def456")
				_ = span
				return ctx
			}(),
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := ExtractTraceFields(tt.ctx)
			if len(fields) != tt.expected {
				t.Errorf("ExtractTraceFields() returned %d fields, expected %d", len(fields), tt.expected)
			}
		})
	}
}

func TestFileRotation(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	logger, err := newLogger(
		WithLevel(LevelDebug),
		WithFormat(FormatFile),
		WithFilePath(logFile),
		WithMaxSize(1), // 1MB for testing
		WithMaxAge(1),  // 1 day for testing
		WithMaxBackups(2),
		WithCompress(false),
	)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Write some logs
	for i := 0; i < 100; i++ {
		logger.Info("test log message", zap.Int("iteration", i))
	}

	// Check if log file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

func TestDirectoryPath(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir := t.TempDir()

	logger, err := newLogger(
		WithLevel(LevelDebug),
		WithFormat(FormatFile),
		WithFilePath(tempDir), // Directory path
		WithMaxSize(1),
		WithMaxAge(1),
		WithMaxBackups(2),
		WithCompress(false),
	)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Write some logs
	logger.Info("test log message")

	// Check if log file was created in the directory
	expectedLogFile := filepath.Join(tempDir, "app.log")
	if _, err := os.Stat(expectedLogFile); os.IsNotExist(err) {
		t.Error("Log file was not created in directory")
	}
}

func TestInit(t *testing.T) {
	// Test successful initialization
	err := Init(
		WithLevel(LevelInfo),
		WithFormat(FormatJSON),
	)
	if err != nil {
		t.Errorf("Init() error = %v, want nil", err)
	}

	// Test that we can use global logger
	Info("test message")
}

func TestInitWithOptions(t *testing.T) {
	// Test functional options
	err := Init(
		WithLevel(LevelDebug),
		WithFormat(FormatText),
	)
	if err != nil {
		t.Errorf("InitWithOptions() error = %v, want nil", err)
	}

	// Test that we can use global logger
	Debug("test debug message")
}

func TestNewWithOptions(t *testing.T) {
	// Test creating logger with functional options
	logger, err := newLogger(
		WithLevel(LevelWarn),
		WithFormat(FormatJSON),
	)
	if err != nil {
		t.Errorf("NewWithOptions() error = %v, want nil", err)
	}
	if logger == nil {
		t.Error("NewWithOptions() returned nil logger")
	}

	// Test that logger works
	logger.Warn("test warn message")
}

func TestPredefinedOptions(t *testing.T) {
	// Test production defaults
	logger, err := newLogger(WithProductionDefaults())
	if err != nil {
		t.Errorf("newLogger(WithProductionDefaults()) error = %v", err)
	}
	if logger == nil {
		t.Error("newLogger(WithProductionDefaults()) returned nil logger")
	}

	// Test development defaults
	logger, err = newLogger()
	if err != nil {
		t.Errorf("newLogger() error = %v", err)
	}
	if logger == nil {
		t.Error("newLogger() returned nil logger")
	}

	// Test console output
	logger, err = newLogger(WithConsoleOutput())
	if err != nil {
		t.Errorf("newLogger(WithConsoleOutput()) error = %v", err)
	}
	if logger == nil {
		t.Error("newLogger(WithConsoleOutput()) returned nil logger")
	}

	// Test text output
	logger, err = newLogger(WithTextOutput())
	if err != nil {
		t.Errorf("newLogger(WithTextOutput()) error = %v", err)
	}
	if logger == nil {
		t.Error("newLogger(WithTextOutput()) returned nil logger")
	}
}

func TestNewWithoutOptions(t *testing.T) {
	// Test creating logger without any options (should use defaults)
	logger, err := newLogger()
	if err != nil {
		t.Errorf("newLogger() error = %v", err)
	}
	if logger == nil {
		t.Error("newLogger() returned nil logger")
	}

	// Test that logger works
	logger.Info("test info message")
	logger.Debug("test debug message")
	logger.Error("test error message")
	logger.Warn("test warn message")
}

func TestLogWithProductionDefaults(t *testing.T) {
	// Initialize logger with production defaults
	logger, err := newLogger(WithProductionDefaults())
	if err != nil {
		t.Fatalf("Failed to create logger with production defaults: %v", err)
	}

	// Test logging at different levels
	logger.Debug("This is a debug message") // Should not appear
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
}

func TestLogWithJSONFormat(t *testing.T) {
	// Initialize logger with JSON format
	logger, err := newLogger(
		WithLevel(LevelDebug),
		WithFormat(FormatJSON),
		WithDev(false),
	)

	if err != nil {
		t.Fatalf("Failed to create logger with JSON format: %v", err)
	}

	// Test logging at different levels
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
	logger.Print("This is a print message")
	logger.Printf("This is a formatted print message: %s", "formatted")
}
