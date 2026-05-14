package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const LoggerKey contextKey = "logger"

// Put logger into context
func ToContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, log)
}

// Get logger from context
func FromContext(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(LoggerKey).(*zap.Logger)
	if !ok || log == nil {
		fallback, _ := zap.NewProduction()
		return fallback
	}

	return log
}
