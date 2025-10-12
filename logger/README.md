# Logger Package

A simple, concise zapper-based logger package with configurable output formats, log levels, file rotation, and context-based tracing integration for Datadog and OpenTelemetry.

## Features

- **Embedded zap.Logger**: Inherit all standard zap methods
- **Context-aware tracing**: Automatic trace/span ID extraction from context
- **Global logger support**: Use `logger.Info()` anywhere in your project
- **Multiple output formats**: JSON, text, and file output
- **File rotation**: Automatic log rotation with configurable retention
- **Print functions**: Drop-in replacement for `fmt.Print*` functions
- **Thread-safe**: Atomic level changes and thread-safe operations

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "github.com/laziness-coders/go-utils/logger"
    "go.uber.org/zap"
)

func main() {
    // Initialize global logger with functional options
    err := logger.Init(
        logger.WithLevel(logger.LevelInfo),
        logger.WithFormat(logger.FormatJSON),
    )
    if err != nil {
        panic(err)
    }

    // Use global logger
    logger.Info("application started")
    logger.Debug("debug message", zap.String("key", "value"))

    // Use with context for tracing
    ctx := context.Background() // context with trace info
    logger.Trace(ctx).Info("processing request")
    logger.Trace(ctx).Warn("warning message", zap.Int("count", 5))
}
```

### Print Functions (fmt replacement)

```go
// Instead of fmt.Print*, use logger functions
logger.Print("simple message")
logger.Println("message with newline")
logger.Printf("formatted message: %s", "value")

// Level-specific printf functions
logger.Debugf("debug: %d", 42)
logger.Infof("info: %s", "information")
logger.Warnf("warning: %s", "warning message")
logger.Errorf("error: %s", "error message")
```

### Sugar Logger

```go
// Get sugar logger
sugar := logger.Sugar()

// Use sugar methods
sugar.Infow("user action",
    "action", "login",
    "user_id", "123",
)

// Sugar with tracing
sugar.Trace(ctx).Infow("processing with trace",
    "request_id", "req-123",
)
```

## Configuration

### Config Struct

```go
type Config struct {
    Level      string // debug, info, warn, error (default: debug)
    Format     string // json, text, file (default: text)
    FilePath   string // file or directory path (only for file format)
    MaxSize    int    // max size in MB before rotation (default: 100)
    MaxAge     int    // max days to retain old logs (default: 30)
    MaxBackups int    // max number of old log files (default: 3)
    Compress   bool   // compress rotated files (default: false)
    IsDev      bool   // use zap's development config for human-readable output (default: false)
}
```

### Configuration Examples

#### Console Output (JSON)
```go
logger.Init(
    logger.WithLevel(logger.LevelInfo),
    logger.WithFormat(logger.FormatJSON),
)
```

#### Console Output (Text)
```go
logger.Init(
    logger.WithLevel(logger.LevelDebug),
    logger.WithFormat(logger.FormatText),
)
```

#### File Output with Rotation
```go
// Using functional options (recommended)
logger.Init(
    logger.WithLevel(logger.LevelInfo),
    logger.WithFileOutput("/var/log/myapp"),
    logger.WithMaxSize(100),
    logger.WithMaxAge(30),
    logger.WithMaxBackups(3),
    logger.WithCompress(true),
)

```

#### Single File Output
```go
// Using functional options (recommended)
logger.Init(
    logger.WithLevel(logger.LevelInfo),
    logger.WithFileOutput("/var/log/myapp.log"),
    logger.WithMaxSize(50),
    logger.WithMaxAge(7),
    logger.WithMaxBackups(5),
    logger.WithCompress(false),
)

```

#### Development Mode
```go
// Use zap's development config for human-readable output
// This is useful during local development
logger.Init(
    logger.WithLevel(logger.LevelDebug),
    logger.WithFormat(logger.FormatText),
    logger.WithDev(true),
)

// Or use the predefined development defaults
logger.Init(logger.WithDevelopmentDefaults())
```

#### Production Debug Mode
```go
// Sometimes you need debug level in production for troubleshooting
// but still want structured JSON logs (not development format)
logger.Init(
    logger.WithLevel(logger.LevelDebug),
    logger.WithFormat(logger.FormatJSON),
    logger.WithDev(false), // Keep production format
)

// This gives you debug-level logs in production-friendly JSON format
```

## Advanced Usage

### Custom Logger Instance

```go
// Create custom logger with functional options
log, err := logger.New(
    logger.WithLevel(logger.LevelDebug),
    logger.WithFormat(logger.FormatJSON),
)
if err != nil {
    log.Fatal("Failed to create logger", zap.Error(err))
}

// Use custom logger
log.Info("custom logger message")
log.Trace(ctx).Debug("debug with trace")

// Change log level dynamically
log.SetLevel(zapcore.ErrorLevel)
```

### Context Integration

The logger automatically extracts tracing information from context:

#### OpenTelemetry
```go
import "go.opentelemetry.io/otel/trace"

// OpenTelemetry trace context is automatically extracted
span := trace.SpanFromContext(ctx)
logger.Trace(ctx).Info("processing request") // Includes trace_id and span_id
```

#### Datadog
```go
// Datadog trace context is automatically extracted
ctx := context.WithValue(context.Background(), "trace_id", "12345")
ctx = context.WithValue(ctx, "span_id", "67890")
logger.Trace(ctx).Info("processing request") // Includes dd.trace_id and dd.span_id
```

#### Custom Trace Fields
```go
// Custom trace fields
ctx := context.WithValue(context.Background(), "trace_id", "custom-trace")
ctx = context.WithValue(ctx, "span_id", "custom-span")
logger.Trace(ctx).Info("processing request") // Includes trace_id and span_id
```

### Functional Options

The logger supports functional configuration options for cleaner, more readable setup:

```go
// Initialize with functional options
logger.Init(
    logger.WithLevel(logger.LevelDebug),
    logger.WithFormat(logger.FormatText),
    logger.WithDev(true), // Use development config
)

// Predefined option sets
logger.Init(logger.WithProductionDefaults())
logger.Init(logger.WithDevelopmentDefaults())

// Console output
logger.Init(logger.WithConsoleOutput())

// Text output
logger.Init(logger.WithTextOutput())

// File output with rotation
logger.Init(
    logger.WithFileOutput("/var/log/myapp"),
    logger.WithMaxSize(100),
    logger.WithMaxAge(30),
    logger.WithMaxBackups(3),
    logger.WithCompress(true),
)
```

#### Available Options

- `WithLevel(level string)` - Set log level (debug, info, warn, error)
- `WithFormat(format string)` - Set log format (json, text, file)
- `WithFilePath(path string)` - Set file path or directory for logs
- `WithMaxSize(mb int)` - Set max size in MB before rotation
- `WithMaxAge(days int)` - Set max days to retain old logs
- `WithMaxBackups(n int)` - Set max number of old log files
- `WithCompress(compress bool)` - Enable/disable compression of rotated files
- `WithDev(isDev bool)` - Use zap's development config for human-readable output
- `WithProductionDefaults()` - Apply production-friendly defaults
- `WithDevelopmentDefaults()` - Apply development-friendly defaults
- `WithConsoleOutput()` - Configure JSON console output
- `WithTextOutput()` - Configure text console output
- `WithFileOutput(path string)` - Configure file output with rotation

### Dynamic Level Changes

```go
// Change log level at runtime
logger.Global().SetLevel(zapcore.DebugLevel)

// Check current level
level := logger.Global().GetLevel()
```

### Child Loggers

```go
// Create child logger with additional fields
childLogger := logger.With(
    zap.String("service", "api"),
    zap.String("version", "1.0.0"),
)

childLogger.Info("service event") // Includes service and version fields
```

## Testing

```go
// Use no-op logger for testing
testLogger := logger.NewNop()
testLogger.Info("this won't actually log anything")

// Or set global logger for tests
logger.SetGlobal(logger.NewNop())
logger.Info("test message") // Won't output anything
```

## Constants

```go
// Log levels
logger.LevelDebug
logger.LevelInfo
logger.LevelWarn
logger.LevelError

// Log formats
logger.FormatJSON
logger.FormatText
logger.FormatFile
```

## Thread Safety

The logger is thread-safe and can be used concurrently from multiple goroutines. The global logger uses a read-write mutex for safe access.

## Dependencies

- `go.uber.org/zap` - Core logging functionality
- `gopkg.in/natefinch/lumberjack.v2` - Log rotation
- `go.opentelemetry.io/otel/trace` - OpenTelemetry integration (already in your project)

## Migration from fmt

Replace `fmt` calls with logger calls:

```go
// Before
fmt.Print("message")
fmt.Println("message")
fmt.Printf("formatted: %s", value)

// After
logger.Print("message")
logger.Println("message")
logger.Printf("formatted: %s", value)
```

This provides structured logging with proper levels and formatting while maintaining the same API.
