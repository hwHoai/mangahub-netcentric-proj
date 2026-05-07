package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

// Init initializes the global logger.
// If isJSON is true, it uses JSONHandler, otherwise TextHandler.
func Init(isJSON bool, level slog.Level) {
	once.Do(func() {
		opts := &slog.HandlerOptions{
			Level: level,
		}

		var handler slog.Handler
		if isJSON {
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			handler = slog.NewTextHandler(os.Stdout, opts)
		}

		defaultLogger = slog.New(handler)
		slog.SetDefault(defaultLogger)
	})
}

// Info logs at LevelInfo.
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Error logs at LevelError.
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// Warn logs at LevelWarn.
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// Debug logs at LevelDebug.
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// With returns a Logger that includes the given attributes in each output operation.
func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

// InfoContext logs at LevelInfo with the given context.
func InfoContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// ErrorContext logs at LevelError with the given context.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}
