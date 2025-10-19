package logger

import (
	"fmt"
	"github.com/laziness-coders/go-utils/generic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

// Config defines the logger configuration.
type Config struct {
	Level      string // debug, info, warn, error (default: info)
	Format     string // json, text, file (default: json)
	FilePath   string // file or directory path (only for file format)
	MaxSize    int    // max size in MB before rotation (default: 100)
	MaxAge     int    // max days to retain old logs (default: 30)
	MaxBackups int    // max number of old log files (default: 3)
	Compress   bool   // compress rotated files (default: false)
	IsDev      *bool  // use zap's development config for human-readable output (default: false)
	CallerSkip int    // caller skip for accurate logging (default: 1)

	AtomicLevel zap.AtomicLevel // atomic level for dynamic level changes
}

// newConfig creates a new Config with the given options applied.
func newConfig(opts ...Option) (*Config, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	cfg = cfg.setDefaults().setAtomicLevel()
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}

// setDefaults applies default values to the configuration.
func (c *Config) setDefaults() *Config {
	if c.Level == "" {
		c.Level = LevelDebug
	}
	if c.Format == "" {
		c.Format = FormatText
	}
	if c.MaxSize <= 0 {
		c.MaxSize = defaultMaxSize
	}
	if c.MaxAge <= 0 {
		c.MaxAge = defaultMaxAge
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = defaultMaxBackups
	}
	if c.CallerSkip <= 0 {
		c.CallerSkip = defaultCallerSkip
	}

	// Default IsDev to true if not set
	if c.IsDev == nil {
		c.IsDev = generic.ToPointer(true)
	}

	return c
}

func (c *Config) setAtomicLevel() *Config {
	var level zapcore.Level
	switch c.Level {
	case LevelDebug:
		level = zapcore.DebugLevel
	case LevelInfo:
		level = zapcore.InfoLevel
	case LevelWarn:
		level = zapcore.WarnLevel
	case LevelError:
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	c.AtomicLevel = zap.NewAtomicLevelAt(level)
	return c
}

// validate checks if the configuration is valid.
func (c Config) validate() error {
	// validate level
	switch c.Level {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		// valid
	default:
		return fmt.Errorf("invalid log level: %s", c.Level)
	}

	// validate format
	switch c.Format {
	case FormatJSON, FormatText, FormatFile:
		// valid
	default:
		return fmt.Errorf("invalid log format: %s", c.Format)
	}

	// If format is file, FilePath is required
	if c.Format == FormatFile && c.FilePath == "" {
		return fmt.Errorf("file path is required for file format")
	}

	// If FilePath is provided, ensure the directory exists or can be created
	if c.FilePath != "" {
		dir := filepath.Dir(c.FilePath)
		if err := os.MkdirAll(dir, logDirPermission); err != nil {
			return fmt.Errorf("failed to create log directory %s: %w", dir, err)
		}
	}

	return nil
}

func (c Config) IsDevelopment() bool {
	return c.IsDev != nil && *c.IsDev
}

func (c Config) buildZapEncoder() zapcore.Encoder {
	// Create encoder config
	var encoderConfig zapcore.EncoderConfig
	var encoder zapcore.Encoder

	// Use zap's development config when IsDev is true
	if c.IsDevelopment() {
		// Use zap's development config
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		return encoder
	}

	// Use production config for other cases
	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create encoder based on format
	switch c.Format {
	case FormatJSON, FormatFile:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case FormatText:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return encoder
}

func (c Config) buildZapWriteSyncer() zapcore.WriteSyncer {
	// Create writer sync based on format
	var writeSyncer zapcore.WriteSyncer
	switch c.Format {
	case FormatFile:
		writeSyncer = createFileWriter(c)
	default:
		writeSyncer = zapcore.AddSync(os.Stdout)
	}
	return writeSyncer
}

func (c Config) buildZapLogger() *zap.Logger {
	core := zapcore.NewCore(
		c.buildZapEncoder(),
		c.buildZapWriteSyncer(),
		c.AtomicLevel,
	)

	// Use zap's development options when IsDev is true
	var (
		zapLogger *zap.Logger
		zapOpts   = []zap.Option{
			zap.AddCaller(),
			zap.AddCallerSkip(c.CallerSkip),
			zap.AddStacktrace(zapcore.ErrorLevel),
		}
	)

	if c.IsDevelopment() {
		// Enable development options
		zapOpts = append(zapOpts, zap.Development())
	}

	zapLogger = zap.New(
		core,
		zapOpts...,
	)
	return zapLogger
}

// createFileWriter creates a file writer with rotation support.
func createFileWriter(cfg Config) zapcore.WriteSyncer {
	// If FilePath is a directory, create a default filename
	filePath := cfg.FilePath
	if stat, err := os.Stat(filePath); err == nil && stat.IsDir() {
		filePath = filepath.Join(filePath, "app.log")
	}

	// Create lumberjack logger for rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	}

	return zapcore.AddSync(lumberjackLogger)
}
