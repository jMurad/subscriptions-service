package middleware

import (
	"context"
	"net/http"
	"subscriptions-service/internal/contextkeys"

	"github.com/google/uuid"
)

type contextKey string

// Add request ID to request context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		ctx := context.WithValue(
			r.Context(),
			contextkeys.RequestID,
			requestID,
		)

		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(contextkeys.RequestID).(string)

	if !ok {
		return "unknown"
	}

	return requestID
}
