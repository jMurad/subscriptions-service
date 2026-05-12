package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subscriptions-service/internal/config"
	"subscriptions-service/internal/handler"
	"subscriptions-service/internal/repository/postgres"
	"subscriptions-service/internal/service/subscriptions"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// --- load environment ---
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// --- log ---
	logg, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// --- DB ---
	db, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		logg.Fatal("failed to connect db", zap.Error(err))
	}

	defer db.Close()

	// --- repositories ---
	repo := postgres.NewSubscriptionRepository(db)

	// --- services ---
	svc := subscriptions.NewSubscriptionService(repo)

	// --- handlers ---
	h := handler.NewSubscriptionHandler(svc)

	// --- router ---
	r := chi.NewRouter()

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
		r.Get("/", h.List)
	})

	// Configure HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	// Start HTTP server in background goroutine
	go func() {
		logg.Info("server started", zap.String("port", cfg.AppPort))

		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logg.Fatal("server failed", zap.Error(err))
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	// Gracefully shutdown server
	logg.Info("shutdown server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logg.Fatal("shutdown failed", zap.Error(err))
	}

	logg.Info("server stopped")
}
