package subscriptions

import (
	"context"
	"strconv"
	"strings"
	apperrors "subscriptions-service/internal/errors"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) Total(ctx context.Context, userID *uuid.UUID, serviceName *string, from time.Time, to time.Time) (int, error) {
	conditions := make([]string, 0)
	args := []any{from, to}

	argPos := 3

	if userID != nil {
		conditions = append(conditions, "user_id = $"+strconv.Itoa(argPos))
		args = append(args, *userID)

		argPos++
	}
	if serviceName != nil {
		conditions = append(conditions, "service_name = $"+strconv.Itoa(argPos))
		args = append(args, *serviceName)

		argPos++
	}

	query := `
    SELECT COALESCE(SUM(
		price * (
			(
				DATE_PART(
					'year',
					AGE(
						LEAST(
							COALESCE(end_date, $2),
							$2
						),
						GREATEST(start_date, $1)
					)
				) * 12
			)
			+
			DATE_PART(
				'month',
				AGE(
					LEAST(
						COALESCE(end_date, $2),
						$2
					),
					GREATEST(start_date, $1)
				)
			)
			+ 1
		)
	), 0)
	FROM subscriptions
	WHERE start_date <= $2
	AND (
		end_date IS NULL
		OR end_date >= $1
	)
    `

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var total int

	err := r.db.QueryRow(
		ctx,
		query,
		args...,
	).Scan(&total)

	if err != nil {
		return 0, apperrors.MapPostgresError(err)
	}

	return total, nil
}
