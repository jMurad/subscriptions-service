package app

import (
	"context"

	"subscriptions-service/internal/config"
	"subscriptions-service/internal/handler"
	"subscriptions-service/internal/repository/postgres"
	"subscriptions-service/internal/service/subscriptions"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Dependencies struct {
	SubscriptionHandler *handler.SubscriptionHandler
}

func loadConfig(logg *zap.Logger) (*config.Config, error) {
	logg.Info("loading config")

	cfg, err := config.Load()
	if err != nil {
		logg.Error("failed to load config",
			zap.Error(err),
		)
		return nil, err
	}

	logg.Info("config loaded")

	return cfg, nil
}

func newLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func newPostgres(logg *zap.Logger, cfg *config.Config) (*pgxpool.Pool, error) {
	logg.Info("connecting to postgres")

	db, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		logg.Error("failed to connect postgres",
			zap.Error(err),
		)
		return nil, err
	}

	logg.Info("postgres connected")

	return db, nil
}

func initDependencies(logg *zap.Logger, db *pgxpool.Pool) *Dependencies {
	logg.Info("initializing dependencies")

	repo := postgres.NewSubscriptionRepository(db)

	logg.Info("repository initialized")

	svc := subscriptions.NewSubscriptionService(repo)

	logg.Info("service initialized")

	h := handler.NewSubscriptionHandler(svc)

	logg.Info("handler initialized")

	return &Dependencies{
		SubscriptionHandler: h,
	}
}
