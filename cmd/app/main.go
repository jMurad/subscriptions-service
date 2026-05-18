// @title Subscriptions API
// @version 1.0
// @description REST API for managing subscriptions
// @host localhost:8080
// @BasePath /
package main

import (
	"context"

	"subscriptions-service/internal/app"
	"subscriptions-service/internal/config"
	"subscriptions-service/internal/logger"

	_ "subscriptions-service/docs"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New()

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	providers := app.NewProviders(db)

	router := app.NewRouter(log, providers.Handler)

	server := app.NewServer(router)

	if err := app.Run(server, log, db.Close); err != nil {
		log.Fatal(err.Error())
	}
}
