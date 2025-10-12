package logger

import (
	"context"
	"go.uber.org/zap"
)

var (
	// globalLogger is the global logger instance
	globalLogger *Logger
)

// Global logger functions

// Init initializes the global logger with functional options.
// This should be called once at application startup.
func Init(opts ...Option) error {
	logger, err := newLogger(opts...)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	globalLogger.Panic(msg, fields...)
}

// Trace extracts tracing fields from context using the global logger.
func Trace(ctx context.Context) *Logger {
	if globalLogger != nil {
		return globalLogger.Trace(ctx)
	}
	panic("global logger is not initialized")
}

// Sugar returns a sugared logger from the global logger.
func Sugar() *SugaredLogger {
	if globalLogger != nil {
		return globalLogger.Sugar()
	}
	panic("global logger is not initialized")
}

// With creates a child logger with additional fields using the global logger.
func With(fields ...zap.Field) *Logger {
	if globalLogger != nil {
		return &Logger{
			Logger: globalLogger.With(fields...),
		}
	}
	panic("global logger is not initialized")
}

// Package-level print functions
func Print(args ...interface{}) {
	globalLogger.Print(args...)
}

func Println(args ...interface{}) {
	globalLogger.Println(args...)
}

func Printf(format string, args ...interface{}) {
	globalLogger.Printf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}
