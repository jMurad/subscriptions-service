package middleware

import (
	"net/http"
	"runtime/debug"
	"subscriptions-service/internal/logger"

	"go.uber.org/zap"
)

// Recover from panics
func Recovery(baseLogger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(
				r.Context(),
			)

			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						zap.Any("panic", rec),
						zap.ByteString("stack", debug.Stack()),
					)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
