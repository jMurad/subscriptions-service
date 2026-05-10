package postgres

import (
	"context"

	"subscriptions-service/internal/model"

	"github.com/google/uuid"
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

func (r *SubRepo) Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error) {
	query := `
	INSERT INTO subscriptions
	(service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	var id uuid.UUID

	err := r.db.QueryRow(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&id)

	return id, err
}
