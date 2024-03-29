package server

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

type Logger struct {
	innerLogger *slog.Logger
	context     context.Context
}

type LoggerConfig struct {
	LoggerType string
	Level      slog.Level
	AddSource  bool
}

// NewLogger creates a new Logger instance with the provided configuration.
// It supports two types of loggers: "json" and "text" (passed as the LoggerType).
// The function initializes a context (which is used internally only), sets up a handler based on the logger type,
// and returns a pointer to the Logger struct.
// The handler is responsible for formatting and outputting log messages.
// If the configuration specifies to add source information to log messages,
// the source file name will be extracted and only the base name will be used.
// The function returns an error if an unsupported logger type is specified.
func NewLogger(config LoggerConfig) (*Logger, error) {
	ctx := context.Background()
	var handler slog.Handler
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "function" {
			return slog.Attr{}
		}

		if config.AddSource && a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}

		// if used in AWS or GCP (and Azure?), they'll add the timestamp themselves. Else uncomment to get timestamps
		// if a.Key == slog.TimeKey {
		// 	return slog.Attr{}
		// }
		return a
	}

	switch config.LoggerType {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       config.Level,
			AddSource:   config.AddSource,
			ReplaceAttr: replaceAttr,
		})
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:       config.Level,
			AddSource:   config.AddSource,
			ReplaceAttr: replaceAttr,
		})
	default:
		return nil, errors.New("unsupported logger type")
	}

	return &Logger{
		innerLogger: slog.New(handler),
		context:     ctx,
	}, nil
}

func (l *Logger) Info(msg string, attrs ...slog.Attr) {
	l.innerLogger.LogAttrs(l.context, slog.LevelInfo, msg, attrs...)
}

func (l *Logger) Error(msg string, attrs ...slog.Attr) {
	l.innerLogger.LogAttrs(l.context, slog.LevelError, msg, attrs...)
}

func (l *Logger) Debug(msg string, attrs ...slog.Attr) {
	l.innerLogger.LogAttrs(l.context, slog.LevelDebug, msg, attrs...)
}

func (l *Logger) Warn(msg string, attrs ...slog.Attr) {
	l.innerLogger.LogAttrs(l.context, slog.LevelWarn, msg, attrs...)
}

// Add more methods as needed
