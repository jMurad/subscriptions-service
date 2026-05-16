package postgres

import (
	"context"
	"errors"
	"strconv"
	"time"

	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
	log := logger.FromContext(ctx)
	start := time.Now()

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
		log.Error("repository create subscription query failed", zap.Error(err))
		return uuid.Nil, err
	}

	log.Info("repository create subscription query completed",
		zap.String("subscription_id", id.String()),
		zap.Duration("duration", time.Since(start)),
	)

	return id, nil
}

func (r *SubRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	log := logger.FromContext(ctx)
	start := time.Now()

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
		log.Error("repository get subscription query failed",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("repository get subscription query completed",
		zap.String("subscription_id", id.String()),
		zap.Duration("duration", time.Since(start)),
	)

	return &sub, nil
}

func (r *SubRepo) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Subscription, error) {
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

func (r *SubRepo) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
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
WHERE deleted_at IS NULL
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

func (r *SubRepo) Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error {
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

func (r *SubRepo) Delete(ctx context.Context, id uuid.UUID) error {
		query := `
	UPDATE subscriptions
SET deleted_at = NOW(),
        status = 'canceled'
	WHERE id = $1
AND deleted_at IS NULL
	`

	result, err := r.db.Exec(		ctx, 		query, id	)
	if err != nil {
		return apperrors.MapPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.MapPostgresError(pgx.ErrNoRows)
	}

	return nil
}

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
