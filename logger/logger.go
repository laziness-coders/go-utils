package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.Logger with context support
type Logger struct {
	base  *zap.Logger
	sugar *zap.SugaredLogger
}

// contextKey is the type used for context keys
type contextKey string

const (
	// traceIDKey is the context key for trace ID
	traceIDKey contextKey = "trace_id"
	// spanIDKey is the context key for span ID
	spanIDKey contextKey = "span_id"
)

// New creates a new Logger instance from the given configuration
func New(cfg *Config) (*Logger, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Set defaults
	cfg.SetDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Create zap logger
	zapLogger, err := buildZapLogger(cfg)
	if err != nil {
		return nil, err
	}

	return &Logger{
		base:  zapLogger,
		sugar: zapLogger.Sugar(),
	}, nil
}

// buildZapLogger builds a zap.Logger from the configuration
func buildZapLogger(cfg *Config) (*zap.Logger, error) {
	// Parse log level
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Use colored output for console format
	if cfg.Format == LogFormatConsole || cfg.Format == LogFormatText {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create encoder
	var encoder zapcore.Encoder
	if cfg.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create writer syncer
	var writeSyncer zapcore.WriteSyncer
	if cfg.Format == LogFormatFile {
		// Write to file
		file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	} else {
		// Write to stdout
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger options
	opts := []zap.Option{}

	if cfg.EnableCaller {
		opts = append(opts, zap.AddCaller())
	}

	if cfg.EnableStacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	return zap.New(core, opts...), nil
}

// parseLevel parses a LogLevel string into a zapcore.Level
func parseLevel(level LogLevel) (zapcore.Level, error) {
	switch level {
	case LogLevelDebug:
		return zapcore.DebugLevel, nil
	case LogLevelInfo:
		return zapcore.InfoLevel, nil
	case LogLevelWarn:
		return zapcore.WarnLevel, nil
	case LogLevelError:
		return zapcore.ErrorLevel, nil
	case LogLevelDPanic:
		return zapcore.DPanicLevel, nil
	case LogLevelPanic:
		return zapcore.PanicLevel, nil
	case LogLevelFatal:
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, nil
	}
}

// Base returns the underlying zap.Logger
func (l *Logger) Base() *zap.Logger {
	return l.base
}

// Sugar returns the underlying zap.SugaredLogger
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// WithContext returns a logger with fields extracted from context
// This supports tracing integration with Datadog, OpenTelemetry, etc.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	fields := extractFieldsFromContext(ctx)
	if len(fields) == 0 {
		return l
	}

	return &Logger{
		base:  l.base.With(fields...),
		sugar: l.sugar.With(fieldsToInterfaces(fields)...),
	}
}

// extractFieldsFromContext extracts trace and span IDs from context
func extractFieldsFromContext(ctx context.Context) []zap.Field {
	fields := []zap.Field{}

	// Extract trace ID
	if traceID := ctx.Value(traceIDKey); traceID != nil {
		if tid, ok := traceID.(string); ok && tid != "" {
			fields = append(fields, zap.String("trace_id", tid))
		}
	}

	// Extract span ID
	if spanID := ctx.Value(spanIDKey); spanID != nil {
		if sid, ok := spanID.(string); ok && sid != "" {
			fields = append(fields, zap.String("span_id", sid))
		}
	}

	// Support for dd.trace_id and dd.span_id (Datadog format)
	if traceID := ctx.Value("dd.trace_id"); traceID != nil {
		fields = append(fields, zap.Any("dd.trace_id", traceID))
	}

	if spanID := ctx.Value("dd.span_id"); spanID != nil {
		fields = append(fields, zap.Any("dd.span_id", spanID))
	}

	return fields
}

// fieldsToInterfaces converts zap.Field slice to interface slice for sugar logger
func fieldsToInterfaces(fields []zap.Field) []interface{} {
	result := make([]interface{}, 0, len(fields)*2)
	for _, field := range fields {
		result = append(result, field.Key, field.String)
	}
	return result
}

// With creates a child logger with the given fields
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		base:  l.base.With(fields...),
		sugar: l.sugar.With(fieldsToInterfaces(fields)...),
	}
}

// Debug logs a message at DebugLevel
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.base.Debug(msg, fields...)
}

// Info logs a message at InfoLevel
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.base.Info(msg, fields...)
}

// Warn logs a message at WarnLevel
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.base.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.base.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.base.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.base.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.base.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.base.Sync()
}

// ContextWithTraceID returns a context with trace ID
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// ContextWithSpanID returns a context with span ID
func ContextWithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDKey, spanID)
}
