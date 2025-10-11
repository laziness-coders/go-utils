# Logger Package

A flexible and powerful logging package for Go applications built on top of [Uber's Zap](https://github.com/uber-go/zap). This package provides structured logging with context support for distributed tracing integration (Datadog, OpenTelemetry, etc.).

## Features

- **High Performance**: Built on Uber's Zap logger, one of the fastest structured logging libraries
- **Context Support**: Automatic trace and span ID extraction from context for distributed tracing
- **Flexible Output**: Support for console, JSON, and file outputs
- **Configurable**: Easy configuration through structs or config files
- **Base & Sugar Loggers**: Access both structured (base) and printf-style (sugar) loggers
- **Multiple Log Levels**: Debug, Info, Warn, Error, DPanic, Panic, Fatal
- **Global Logger**: Convenient global logger instance for easy usage
- **Tracing Integration**: Built-in support for Datadog and OpenTelemetry trace IDs

## Installation

```bash
go get github.com/laziness-coders/go-utils/logger
```

## Quick Start

### Basic Usage

```go
package main

import (
    "github.com/laziness-coders/go-utils/logger"
    "go.uber.org/zap"
)

func main() {
    // Use the default global logger
    logger.Info("Application started")
    logger.Debug("Debug information")
    logger.Warn("Warning message")
    logger.Error("Error occurred", zap.String("error", "something went wrong"))
    
    // Don't forget to sync before exit
    defer logger.Sync()
}
```

### Using Sugar Logger

```go
package main

import (
    "github.com/laziness-coders/go-utils/logger"
)

func main() {
    sugar := logger.Sugar()
    
    // Printf-style logging
    sugar.Infof("User %s logged in", "john")
    sugar.Debugw("Request processed",
        "method", "GET",
        "path", "/api/users",
        "duration", 42,
    )
    
    defer sugar.Sync()
}
```

## Configuration

### Using Config Struct

```go
package main

import (
    "github.com/laziness-coders/go-utils/logger"
)

func main() {
    cfg := &logger.Config{
        Level:            logger.LogLevelDebug,
        Format:           logger.LogFormatJSON,
        EnableCaller:     true,
        EnableStacktrace: true,
        Development:      true,
    }
    
    if err := logger.InitFromConfig(cfg); err != nil {
        panic(err)
    }
    
    logger.Info("Logger initialized with custom config")
}
```

### Configuration from File

```yaml
# config.yaml
level: "debug"
format: "json"
enable_caller: true
enable_stacktrace: true
development: false
```

```go
package main

import (
    "github.com/laziness-coders/go-utils/configs"
    "github.com/laziness-coders/go-utils/logger"
)

type AppConfig struct {
    Logger logger.Config `mapstructure:"logger"`
}

func main() {
    cfg := &AppConfig{}
    if err := configs.New(cfg).Load(configs.AppEnvironmentDev, "./configs"); err != nil {
        panic(err)
    }
    
    if err := logger.InitFromConfig(&cfg.Logger); err != nil {
        panic(err)
    }
    
    logger.Info("Logger initialized from config file")
}
```

## Configuration Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Level` | `LogLevel` | `info` | Minimum log level (debug, info, warn, error, dpanic, panic, fatal) |
| `Format` | `LogFormat` | `console` | Output format (console, json, text, file) |
| `FilePath` | `string` | `""` | Path to log file (required when format is "file") |
| `EnableCaller` | `bool` | `false` | Enable logging of caller information (file:line) |
| `EnableStacktrace` | `bool` | `false` | Enable automatic stacktrace for errors |
| `Development` | `bool` | `false` | Enable development mode (DPanic panics, more verbose) |
| `Encoding` | `string` | `console` | Log encoding (json or console) |

## Output Formats

### Console Format (Human-Readable)

```go
cfg := &logger.Config{
    Level:  logger.LogLevelInfo,
    Format: logger.LogFormatConsole,
}
logger.InitFromConfig(cfg)
logger.Info("Application started", zap.String("version", "1.0.0"))
```

Output:
```
2024-01-15T10:30:45.123Z	INFO	Application started	{"version": "1.0.0"}
```

### JSON Format (Machine-Readable)

```go
cfg := &logger.Config{
    Level:  logger.LogLevelInfo,
    Format: logger.LogFormatJSON,
}
logger.InitFromConfig(cfg)
logger.Info("User login", zap.String("user", "john"), zap.Int("user_id", 123))
```

Output:
```json
{"level":"info","time":"2024-01-15T10:30:45.123Z","msg":"User login","user":"john","user_id":123}
```

### File Output

```go
cfg := &logger.Config{
    Level:    logger.LogLevelInfo,
    Format:   logger.LogFormatFile,
    FilePath: "/var/log/myapp/app.log",
}
logger.InitFromConfig(cfg)
logger.Info("Logging to file")
```

Or specify a directory (defaults to `app.log`):

```go
cfg := &logger.Config{
    Level:    logger.LogLevelInfo,
    Format:   logger.LogFormatFile,
    FilePath: "/var/log/myapp/", // Will create app.log in this directory
}
```

## Context and Distributed Tracing

### Basic Context Usage

```go
package main

import (
    "context"
    "github.com/laziness-coders/go-utils/logger"
    "go.uber.org/zap"
)

func main() {
    ctx := context.Background()
    
    // Add trace ID to context
    ctx = logger.ContextWithTraceID(ctx, "trace-123-456")
    ctx = logger.ContextWithSpanID(ctx, "span-789-012")
    
    // Logger automatically extracts trace information from context
    logger.WithContext(ctx).Info("Processing request", 
        zap.String("endpoint", "/api/users"))
}
```

Output:
```json
{
  "level": "info",
  "time": "2024-01-15T10:30:45.123Z",
  "msg": "Processing request",
  "trace_id": "trace-123-456",
  "span_id": "span-789-012",
  "endpoint": "/api/users"
}
```

### Integration with Datadog

```go
import (
    "context"
    "github.com/laziness-coders/go-utils/logger"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func handleRequest(ctx context.Context) {
    // Extract Datadog trace information
    span, ctx := tracer.StartSpanFromContext(ctx, "operation")
    defer span.Finish()
    
    // The logger will automatically include dd.trace_id and dd.span_id
    logger.WithContext(ctx).Info("Handling request")
}
```

### Integration with OpenTelemetry

```go
import (
    "context"
    "github.com/laziness-coders/go-utils/logger"
    "go.opentelemetry.io/otel/trace"
)

func handleRequest(ctx context.Context) {
    span := trace.SpanFromContext(ctx)
    if span.SpanContext().IsValid() {
        traceID := span.SpanContext().TraceID().String()
        spanID := span.SpanContext().SpanID().String()
        
        ctx = logger.ContextWithTraceID(ctx, traceID)
        ctx = logger.ContextWithSpanID(ctx, spanID)
    }
    
    logger.WithContext(ctx).Info("Processing with OpenTelemetry trace")
}
```

## Advanced Usage

### Creating Instance Logger

```go
// Create a dedicated logger instance for a component
func NewUserService() *UserService {
    cfg := &logger.Config{
        Level:        logger.LogLevelDebug,
        Format:       logger.LogFormatJSON,
        EnableCaller: true,
    }
    
    log, err := logger.New(cfg)
    if err != nil {
        panic(err)
    }
    
    return &UserService{
        logger: log,
    }
}

type UserService struct {
    logger *logger.Logger
}

func (s *UserService) CreateUser(ctx context.Context, name string) error {
    s.logger.WithContext(ctx).Info("Creating user", 
        zap.String("name", name))
    return nil
}
```

### Adding Persistent Fields

```go
// Create a logger with persistent fields
serviceLogger := logger.With(
    zap.String("service", "user-service"),
    zap.String("version", "1.0.0"),
)

serviceLogger.Info("Service started") // Will include service and version fields
```

### Using Base Logger for Structured Logging

```go
import "go.uber.org/zap"

// Base logger provides structured logging
logger.Base().Info("User created",
    zap.String("user_id", "123"),
    zap.String("email", "user@example.com"),
    zap.Int("age", 25),
    zap.Bool("verified", true),
)
```

### Using Sugar Logger for Printf-Style

```go
// Sugar logger provides printf-style logging
logger.Sugar().Infof("User %s created with ID %d", "john", 123)

// Or with structured fields
logger.Sugar().Infow("User created",
    "user_id", 123,
    "email", "user@example.com",
    "age", 25,
)
```

## Log Levels

| Level | Description | Use Case |
|-------|-------------|----------|
| `Debug` | Verbose information for debugging | Development, troubleshooting |
| `Info` | General informational messages | Application lifecycle events |
| `Warn` | Warning messages for potentially harmful situations | Deprecated API usage, non-critical issues |
| `Error` | Error messages for serious problems | Failed operations, exceptions |
| `DPanic` | Critical errors (panics in development) | Serious bugs that should be fixed |
| `Panic` | Logs and then panics | Unrecoverable errors |
| `Fatal` | Logs and then calls os.Exit(1) | Catastrophic failures |

## Best Practices

1. **Always sync before exit**: Call `logger.Sync()` or `defer logger.Sync()` to flush buffered logs
2. **Use structured logging**: Prefer structured fields over formatted strings for better queryability
3. **Include context**: Use `WithContext()` to automatically include trace IDs
4. **Set appropriate log levels**: Use Debug in development, Info/Warn in production
5. **Use instance loggers**: Create dedicated logger instances for different components
6. **Configure from files**: Use configuration files for easy environment-specific settings

## Examples

### HTTP Middleware with Tracing

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // Generate or extract trace ID
        traceID := r.Header.Get("X-Trace-ID")
        if traceID == "" {
            traceID = generateTraceID()
        }
        ctx = logger.ContextWithTraceID(ctx, traceID)
        
        // Log request
        logger.WithContext(ctx).Info("HTTP request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
        )
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Error Handling

```go
func processOrder(ctx context.Context, orderID string) error {
    log := logger.WithContext(ctx)
    
    log.Info("Processing order", zap.String("order_id", orderID))
    
    if err := validateOrder(orderID); err != nil {
        log.Error("Order validation failed",
            zap.String("order_id", orderID),
            zap.Error(err),
        )
        return err
    }
    
    log.Info("Order processed successfully", zap.String("order_id", orderID))
    return nil
}
```

## Performance

This logger is built on Uber's Zap, which is one of the fastest structured logging libraries in Go:

- Zero allocation in most cases
- Reflection-free
- Optimized for high-throughput applications
- Minimal performance overhead

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
