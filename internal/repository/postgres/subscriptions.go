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

func (r *SubRepo) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
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
		log.Error("repository list subscriptions query failed", zap.Error(err))
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
			log.Error("repository list subscriptions scan failed", zap.Error(err))
			return nil, err
		}

		subscriptions = append(subscriptions, sub)
	}

	log.Info("repository list subscriptions query completed",
		zap.Int("count", len(subscriptions)),
		zap.Duration("duration", time.Since(start)),
	)

	return subscriptions, nil
}

func (r *SubRepo) Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error {
	log := logger.FromContext(ctx)
	start := time.Now()

	query := `
	UPDATE subscriptions
	SET
		service_name = COALESCE($1, service_name),
		price = COALESCE($2, price),
		end_date = COALESCE($3, end_date)
	WHERE id = $4
	`

	_, err := r.db.Exec(
		ctx,
		query,
		update.ServiceName,
		update.Price,
		update.EndDate,
		id,
	)
	if err != nil {
		log.Error("repository update subscription query failed",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	log.Info("repository update subscription query completed",
		zap.String("subscription_id", id.String()),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

func (r *SubRepo) Delete(ctx context.Context, id uuid.UUID) error {
	log := logger.FromContext(ctx)
	start := time.Now()

	query := `
	DELETE FROM subscriptions
	WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		log.Error("repository delete subscription query failed",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	if result.RowsAffected() == 0 {
		err = errors.New("subscription not found")
		log.Warn("repository delete subscription affected zero rows",
			zap.String("subscription_id", id.String()),
		)
		return err
	}

	log.Info("repository delete subscription query completed",
		zap.String("subscription_id", id.String()),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

func (r *SubRepo) Total(ctx context.Context, userID *uuid.UUID, serviceName *string, from time.Time, to time.Time) (int, error) {
	query := `
	SELECT COALESCE(SUM(price), 0)
	FROM subscriptions
	WHERE start_date >= $1
	AND start_date <= $2
	`

	args := []interface{}{
		from,
		to,
	}

	argPos := 3

	if userID != nil {
		query += ` AND user_id = $` + strconv.Itoa(argPos)
		args = append(args, *userID)
		argPos++
	}

	if serviceName != nil {
		query += ` AND service_name = $` + strconv.Itoa(argPos)
		args = append(args, *serviceName)
	}

	var total int

	err := r.db.QueryRow(
		ctx,
		query,
		args...,
	).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
