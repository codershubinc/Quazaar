package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Use JSON handler for production-like environments or text for development
	// For now, default to text as per request/preference
	handler := slog.NewTextHandler(os.Stdout, opts)
	Log = slog.New(handler)
	slog.SetDefault(Log)
}

func Info(msg string, args ...any) {
	if Log == nil {
		Init()
	}
	Log.Info(msg, args...)
}

func Error(msg string, args ...any) {
	if Log == nil {
		Init()
	}
	Log.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	if Log == nil {
		Init()
	}
	Log.Warn(msg, args...)
}

func Debug(msg string, args ...any) {
	if Log == nil {
		Init()
	}
	Log.Debug(msg, args...)
}
