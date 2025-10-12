package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

// Example demonstrates basic usage of the logger package
func Example() {
	// Initialize global logger with functional options
	err := Init(
		WithLevel(LevelInfo),
		WithFormat(FormatJSON),
	)
	if err != nil {
		panic(err)
	}

	// Basic logging
	Info("application started")
	Debug("debug message", zap.String("key", "value"))

	// Print functions (fmt replacement)
	Print("simple message")
	Println("message with newline")
	Printf("formatted message: %s", "value")

	// Level-specific printf functions
	Debugf("debug: %d", 42)
	Infof("info: %s", "information")
	Warnf("warning: %s", "warning message")
	Errorf("error: %s", "error message")

	// Sugar logger
	sugar := Sugar()
	sugar.Infow("user action",
		"action", "login",
		"user_id", "123",
	)

	// Context with tracing
	ctx := context.WithValue(context.Background(), "dd.trace_id", "12345")
	ctxWithValue := context.WithValue(ctx, "dd.span_id", "67890")

	Trace(ctx).Info("processing request")
	Trace(ctx).Warn("warning with trace", zap.Int("count", 5))

	// Sugar with tracing
	sugar.Trace(ctxWithValue).Infow("processing with trace",
		"request_id", "req-123",
	)
}

// ExampleFile demonstrates file output with rotation
func ExampleFile() {
	// File output with rotation using functional options
	err := Init(
		WithLevel(LevelInfo),
		WithFileOutput("/tmp/example.log"),
		WithMaxSize(1), // 1MB for testing
		WithMaxAge(1),  // 1 day
		WithMaxBackups(2),
		WithCompress(false),
	)
	if err != nil {
		panic(err)
	}

	// These logs will go to file
	Info("this goes to file")
	Warn("this also goes to file")
	Error("error message in file")
}

// ExampleTrace demonstrates tracing integration
func ExampleTrace() {
	// Context with OpenTelemetry-style trace
	ctx := context.Background()
	// In real usage, this would come from OpenTelemetry or Datadog

	// Context with custom trace
	ctx = context.WithValue(ctx, "trace_id", "custom-trace-789")
	ctx = context.WithValue(ctx, "span_id", "custom-span-012")

	// All these will include trace information
	Trace(ctx).Info("processing request")
	Trace(ctx).Debug("debug with trace")
	Trace(ctx).Error("error with trace")
}

// TestExample runs the examples (they don't actually test anything, just demonstrate usage)
func TestExample(t *testing.T) {
	Example()
	ExampleFile()
	ExampleTrace()
}
