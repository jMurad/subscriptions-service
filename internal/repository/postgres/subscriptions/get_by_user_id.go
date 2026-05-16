package subscriptions

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
)

func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Subscription, error) {
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
    WHERE user_id = $1
        AND deleted_at IS NULL
    ORDER BY created_at DESC
    LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, apperrors.MapPostgresError(err)
	}

	defer rows.Close()

	subscriptions := make([]model.Subscription, 0, limit)

	for rows.Next() {

		var sub model.Subscription

		err = rows.Scan(
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

		subscriptions = append(subscriptions, sub)
	}

	if rows.Err() != nil {
		return nil, apperrors.MapPostgresError(rows.Err())
	}

	return subscriptions, nil
}
