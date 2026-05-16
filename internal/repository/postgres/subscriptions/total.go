package subscriptions

import (
	"context"
	"strconv"
	"strings"
	apperrors "subscriptions-service/internal/errors"
	"time"

	"github.com/google/uuid"
)

func (r *SubRepo) Total(ctx context.Context, userID *uuid.UUID, serviceName *string, from *time.Time, to *time.Time) (int, error) {
	conditions := make([]string, 0)
	args := make([]any, 0)

	argPos := 1

	if from != nil {
		conditions = append(conditions, "start_date >= $"+strconv.Itoa(argPos))
		args = append(args, *from)

		argPos++
	}
	if to != nil {
		conditions = append(conditions, "start_date <= $"+strconv.Itoa(argPos))
		args = append(args, *to)

		argPos++
	}
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
    SELECT COALESCE(SUM(price), 0)
    FROM subscriptions
    WHERE deleted_at IS NULL
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
