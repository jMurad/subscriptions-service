package middleware

import (
	"net/http"
	"runtime/debug"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/transport/http/response"

	"go.uber.org/zap"
)

// Recover from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				log := logger.FromContext(r.Context())

				log.Error("panic recovered",
					zap.Any("panic", rec),
					zap.ByteString("stacktrace", debug.Stack()),
				)
				response.Error(w, apperrors.ErrInternal)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
