package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

var (
	// globalLogger is the global logger instance
	globalLogger *Logger
	// mu protects globalLogger
	mu sync.RWMutex
)

// init initializes the global logger with default configuration
func init() {
	// Create a default logger
	logger, err := New(DefaultConfig())
	if err != nil {
		panic(err)
	}
	globalLogger = logger
}

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	return globalLogger
}

// InitFromConfig initializes the global logger from configuration
func InitFromConfig(cfg *Config) error {
	logger, err := New(cfg)
	if err != nil {
		return err
	}
	SetGlobalLogger(logger)
	return nil
}

// Global logger convenience functions

// Base returns the base zap.Logger from the global logger
func Base() *zap.Logger {
	return GetGlobalLogger().Base()
}

// Sugar returns the sugar logger from the global logger
func Sugar() *zap.SugaredLogger {
	return GetGlobalLogger().Sugar()
}

// WithContext returns a logger with fields extracted from context
func WithContext(ctx context.Context) *Logger {
	return GetGlobalLogger().WithContext(ctx)
}

// With creates a child logger with the given fields
func With(fields ...zap.Field) *Logger {
	return GetGlobalLogger().With(fields...)
}

// Debug logs a message at DebugLevel using the global logger
func Debug(msg string, fields ...zap.Field) {
	GetGlobalLogger().Debug(msg, fields...)
}

// Info logs a message at InfoLevel using the global logger
func Info(msg string, fields ...zap.Field) {
	GetGlobalLogger().Info(msg, fields...)
}

// Warn logs a message at WarnLevel using the global logger
func Warn(msg string, fields ...zap.Field) {
	GetGlobalLogger().Warn(msg, fields...)
}

// Error logs a message at ErrorLevel using the global logger
func Error(msg string, fields ...zap.Field) {
	GetGlobalLogger().Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel using the global logger
func DPanic(msg string, fields ...zap.Field) {
	GetGlobalLogger().DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel using the global logger
func Panic(msg string, fields ...zap.Field) {
	GetGlobalLogger().Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel using the global logger
func Fatal(msg string, fields ...zap.Field) {
	GetGlobalLogger().Fatal(msg, fields...)
}

// Sync flushes any buffered log entries from the global logger
func Sync() error {
	return GetGlobalLogger().Sync()
}
