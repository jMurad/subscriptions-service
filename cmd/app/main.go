package main

import (
	"context"
	"log"
	"net/http"

	"subscriptions-service/internal/config"
	"subscriptions-service/internal/handler"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/repository/postgres"
	"subscriptions-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	defer logg.Sync()

	// --- DB ---
	db, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// --- repositories ---
	repo := postgres.NewSubscriptionRepository(db)

	// --- services ---
	svc := service.NewSubscriptionService(repo)

	// --- handlers ---
	h := handler.NewSubscriptionHandler(svc)

	// --- router ---
	r := chi.NewRouter()

	r.Post("/subscriptions", h.Create)

	logg.Info("server started")

	err = http.ListenAndServe(":"+cfg.AppPort, r)
	if err != nil {
		log.Fatal(err)
	}
}
