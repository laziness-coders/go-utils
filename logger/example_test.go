package logger_test

import (
	"context"
	"github.com/laziness-coders/go-utils/logger"
	"go.uber.org/zap"
)

// Example_basicUsage demonstrates basic logger usage
func Example_basicUsage() {
	// Initialize logger with custom config
	cfg := &logger.Config{
		Level:  logger.LogLevelInfo,
		Format: logger.LogFormatConsole,
	}

	if err := logger.InitFromConfig(cfg); err != nil {
		panic(err)
	}

	// Use the global logger
	logger.Info("Application started")
	logger.Debug("This won't be shown (level is Info)")
	logger.Warn("Warning message")
	logger.Error("Error occurred", zap.String("error", "something went wrong"))

	// Don't forget to sync
	defer logger.Sync()
}

// Example_contextTracing demonstrates context-based tracing
func Example_contextTracing() {
	// Initialize logger
	cfg := &logger.Config{
		Level:  logger.LogLevelInfo,
		Format: logger.LogFormatJSON,
	}
	logger.InitFromConfig(cfg)

	// Create context with trace information
	ctx := context.Background()
	ctx = logger.ContextWithTraceID(ctx, "trace-123-456")
	ctx = logger.ContextWithSpanID(ctx, "span-789-012")

	// Logger automatically includes trace IDs
	logger.WithContext(ctx).Info("Processing request",
		zap.String("endpoint", "/api/users"),
		zap.Int("status", 200),
	)

	defer logger.Sync()
}

// Example_sugarLogger demonstrates using the sugar logger
func Example_sugarLogger() {
	cfg := &logger.Config{
		Level:  logger.LogLevelInfo,
		Format: logger.LogFormatConsole,
	}
	logger.InitFromConfig(cfg)

	sugar := logger.Sugar()

	// Printf-style logging
	sugar.Infof("User %s logged in at %s", "john", "2024-01-15")

	// Structured logging with sugar
	sugar.Infow("Request processed",
		"method", "GET",
		"path", "/api/users",
		"duration_ms", 42,
		"status", 200,
	)

	defer sugar.Sync()
}

// Example_fileLogging demonstrates logging to a file
func Example_fileLogging() {
	cfg := &logger.Config{
		Level:    logger.LogLevelInfo,
		Format:   logger.LogFormatFile,
		FilePath: "/tmp/app.log",
	}

	if err := logger.InitFromConfig(cfg); err != nil {
		panic(err)
	}

	logger.Info("This message will be written to /tmp/app.log")
	logger.Error("Error message in file", zap.String("error", "disk full"))

	defer logger.Sync()
}

// Example_instanceLogger demonstrates creating a dedicated logger instance
func Example_instanceLogger() {
	cfg := &logger.Config{
		Level:        logger.LogLevelDebug,
		Format:       logger.LogFormatJSON,
		EnableCaller: true,
	}

	// Create a dedicated logger for a component
	serviceLogger, err := logger.New(cfg)
	if err != nil {
		panic(err)
	}

	// Add persistent fields
	serviceLogger = serviceLogger.With(
		zap.String("service", "user-service"),
		zap.String("version", "1.0.0"),
	)

	// Use the instance logger
	serviceLogger.Info("Service started")
	serviceLogger.Debug("Debug information")

	defer serviceLogger.Sync()
}

// Example_differentLevels demonstrates different log levels
func Example_differentLevels() {
	cfg := &logger.Config{
		Level:  logger.LogLevelDebug,
		Format: logger.LogFormatConsole,
	}
	logger.InitFromConfig(cfg)

	logger.Debug("Detailed debug information")
	logger.Info("General information")
	logger.Warn("Warning: deprecated API usage")
	logger.Error("Error: failed to process request")

	// DPanic logs in development mode will panic
	// logger.DPanic("Critical error")

	defer logger.Sync()
}
