package middleware

import (
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Register(r chi.Router, log *zap.Logger) {
	r.Use(RequestID)
	r.Use(Logger(log))
	r.Use(Recovery)
	r.Use(Timeout(30 * time.Second))
}
