package subscriptions

import (
	"context"
	apperrors "subscriptions-service/internal/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
    UPDATE subscriptions
    SET deleted_at = NOW(),
        status = 'canceled'
    WHERE id = $1
    AND deleted_at IS NULL
    `

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return apperrors.MapPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.MapPostgresError(pgx.ErrNoRows)
	}

	return nil
}
