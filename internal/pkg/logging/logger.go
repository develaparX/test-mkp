package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

var Logger *slog.Logger

func Init() {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: getLogLevel(),
	}

	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Logger = slog.New(handler)
}

func getLogLevel() slog.Level {
	env := os.Getenv("ENV")
	if env == "production" {
		return slog.LevelWarn
	}
	return slog.LevelDebug
}

func WithRequestID(ctx context.Context) context.Context {
	requestID := uuid.New().String()
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}

func InfoCtx(ctx context.Context, msg string, args ...interface{}) {
	requestID := GetRequestID(ctx)
	allArgs := append([]interface{}{"request_id", requestID}, args...)
	Logger.Info(msg, allArgs...)
}

func ErrorCtx(ctx context.Context, msg string, args ...interface{}) {
	requestID := GetRequestID(ctx)
	allArgs := append([]interface{}{"request_id", requestID}, args...)
	Logger.Error(msg, allArgs...)
}

func DebugCtx(ctx context.Context, msg string, args ...interface{}) {
	requestID := GetRequestID(ctx)
	allArgs := append([]interface{}{"request_id", requestID}, args...)
	Logger.Debug(msg, allArgs...)
}

func WarnCtx(ctx context.Context, msg string, args ...interface{}) {
	requestID := GetRequestID(ctx)
	allArgs := append([]interface{}{"request_id", requestID}, args...)
	Logger.Warn(msg, allArgs...)
}

func Info(msg string, args ...interface{}) {
	Logger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Logger.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	Logger.Debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Logger.Warn(msg, args...)
}

func With(args ...interface{}) *slog.Logger {
	return Logger.With(args...)
}
