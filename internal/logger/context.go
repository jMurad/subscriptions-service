package logger

import (
	"context"
	"subscriptions-service/internal/contextkeys"

	"go.uber.org/zap"
)

type contextKey string

// Put logger into context
func ToContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, contextkeys.Logger, log)
}

// Get logger from context
func FromContext(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(contextkeys.Logger).(*zap.Logger)
	if !ok || log == nil {
		fallback, _ := zap.NewProduction()
		return fallback
	}

	return log
}
