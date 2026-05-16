package subscriptions

import (
	"context"
	"strconv"
	"strings"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error {
	conditions := make([]string, 0)
	args := make([]any, 0)

	argPos := 1

	if update.ServiceName != nil {
		conditions = append(conditions, "service_name = $"+strconv.Itoa(argPos))
		args = append(args, *update.ServiceName)

		argPos++
	}

	if update.Price != nil {
		conditions = append(conditions, "price = $"+strconv.Itoa(argPos))
		args = append(args, *update.Price)

		argPos++
	}

	if update.EndDate != nil {
		conditions = append(conditions, "end_date = $"+strconv.Itoa(argPos))
		args = append(args, *update.EndDate)

		argPos++
	}

	if len(conditions) == 0 {
		return nil
	}

	query := `
	UPDATE subscriptions
	SET ` + strings.Join(conditions, ", ") + `
	WHERE id = $` + strconv.Itoa(argPos) + `
	  AND deleted_at IS NULL
	`

	args = append(args, id)

	result, err := r.db.Exec(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return apperrors.MapPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.MapPostgresError(pgx.ErrNoRows)
	}

	return nil
}
