package subscriptions

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
)

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

	if err != nil {
		return uuid.Nil, apperrors.MapPostgresError(err)
	}

	return id, nil
}
