package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

// Logger embeds zap.Logger to inherit all standard methods.
type Logger struct {
	*zap.Logger
}

// SugaredLogger embeds zap.SugaredLogger with tracing support.
type SugaredLogger struct {
	*zap.SugaredLogger
	logger *Logger
}

// newLogger creates a new Logger instance with the given configuration.
func newLogger(opts ...Option) (*Logger, error) {
	cfg, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}
	return &Logger{
		Logger: cfg.buildZapLogger(),
	}, nil
}

// Trace extracts tracing fields from context and returns a logger with those fields.
// Usage: logger.Trace(ctx).Info("processing request")
func (l *Logger) Trace(ctx context.Context) *Logger {
	fields := ExtractTraceFields(ctx)
	return &Logger{
		Logger: l.With(fields...),
	}
}

// Sugar returns a SugaredLogger with tracing support.
func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{
		SugaredLogger: l.Logger.Sugar(),
		logger:        l,
	}
}

// Print functions (fmt replacement)
func (l *Logger) Print(args ...interface{}) {
	l.Info(fmt.Sprint(args...))
}

func (l *Logger) Println(args ...interface{}) {
	l.Info(fmt.Sprintln(args...))
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

// Trace extracts tracing fields from context and returns a sugared logger with those fields.
func (sl *SugaredLogger) Trace(ctx context.Context) *SugaredLogger {
	return sl.logger.Trace(ctx).Sugar()
}
