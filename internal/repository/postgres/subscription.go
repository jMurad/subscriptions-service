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

func (r *SubRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	query := `
	SELECT
		id,
		service_name,
		price,
		user_id,
		start_date,
		end_date,
		created_at
	FROM subscriptions
	WHERE id = $1
	`

	var sub model.Subscription

	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *SubRepo) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	query := `
	SELECT
		id,
		service_name,
		price,
		user_id,
		start_date,
		end_date,
		created_at
	FROM subscriptions
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(
		ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription

	for rows.Next() {

		var sub model.Subscription

		err = rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}
