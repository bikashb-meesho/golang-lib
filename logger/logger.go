package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
)

// Logger wraps a structured logger with common functionality
type Logger struct {
	log *zap.Logger
}

// Config holds logger configuration
type Config struct {
	Level       string // debug, info, warn, error
	Environment string // development, production
	Service     string // service name for log tagging
}

// New creates a new logger instance with the provided configuration
func New(cfg Config) (*Logger, error) {
	if cfg.Service == "" {
		return nil, fmt.Errorf("service name is required")
	}

	var zapConfig zap.Config

	switch cfg.Environment {
	case "production":
		zapConfig = zap.NewProductionConfig()
	case "development":
		zapConfig = zap.NewDevelopmentConfig()
	default:
		zapConfig = zap.NewProductionConfig()
	}

	// Set log level
	level := zap.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	// Add service name as a field to all logs
	zapLogger = zapLogger.With(zap.String("service", cfg.Service))

	return &Logger{log: zapLogger}, nil
}

// Default creates a logger with sensible defaults for development
func Default() *Logger {
	logger, err := New(Config{
		Level:       "info",
		Environment: "development",
		Service:     "app",
	})
	if err != nil {
		// Fallback to basic logger if configuration fails
		zapLogger, _ := zap.NewDevelopment()
		return &Logger{log: zapLogger}
	}
	return logger
}

// Info logs an info level message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

// Warn logs a warning level message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

// Error logs an error level message
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

// Fatal logs a fatal level message and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
	os.Exit(1)
}

// With creates a child logger with additional fields
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{log: l.log.With(fields...)}
}

// contextKey is a custom type for context keys
type contextKey string

const requestIDKey contextKey = "request_id"

// WithContext adds context fields to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract request ID or trace ID from context if available
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return l.With(zap.String("request_id", reqID))
	}
	return l
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.log.Sync()
}
