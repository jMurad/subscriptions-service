package app

import (
	"net/http"
	"subscriptions-service/internal/middleware"
	"subscriptions-service/internal/transport/http/handler/subscriptions"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func NewRouter(logg *zap.Logger, h *subscriptions.Handler) http.Handler {
	r := chi.NewRouter()

	middleware.Register(r, logg)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/user/{user_id}", h.GetByUserID)
		r.Get("/{id}", h.GetByID)
		r.Get("/total", h.Total)
		r.Patch("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})

	return r
}
