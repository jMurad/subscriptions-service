package subscriptions

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	query := `
    SELECT
        id,
        user_id,
        service_name,
        price,
        start_date,
        end_date,
        created_at
    FROM subscriptions
    WHERE id = $1
    AND deleted_at IS NULL
    `

	var sub model.Subscription

	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
	)

	if err != nil {
		return nil, apperrors.MapPostgresError(err)
	}

	return &sub, nil
}
