package app

import (
	"net/http"
	"subscriptions-service/internal/config"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func newHTTPServer(cfg *config.Config, router *chi.Mux) *http.Server {
	return &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}
}

func runServer(logg *zap.Logger, srv *http.Server, port string) {
	logg.Info("server started",
		zap.String("port", port),
	)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logg.Fatal("server failed",
			zap.Error(err),
		)
	}
}
