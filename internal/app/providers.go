package app

import (
	repository "subscriptions-service/internal/repository/postgres/subscriptions"
	service "subscriptions-service/internal/service/subscriptions"
	handler "subscriptions-service/internal/transport/http/handler/subscriptions"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Providers struct {
	Handler *handler.Handler
}

func NewProviders(db *pgxpool.Pool) *Providers {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	return &Providers{
		Handler: h,
	}
}
