package app

import (
	"subscriptions-service/internal/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func newRouter(logg *zap.Logger, deps *Dependencies) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(logg))
	r.Use(middleware.Recovery(logg))

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", deps.SubscriptionHandler.Create)
		r.Get("/{id}", deps.SubscriptionHandler.GetByID)
		r.Get("/", deps.SubscriptionHandler.List)
		r.Get("/total", deps.SubscriptionHandler.Total)
		r.Patch("/{id}", deps.SubscriptionHandler.Update)
		r.Delete("/{id}", deps.SubscriptionHandler.Delete)
	})

	return r
}
