package subscriptions

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubRepo struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubRepo {
	return &SubRepo{
		db: db,
	}
}
