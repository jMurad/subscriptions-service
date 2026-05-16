package subscriptions

import (
	"context"
	"strings"
	"time"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	shortTimeout  = 3 * time.Second
	mediumTimeout = 5 * time.Second
	longTimeout   = 10 * time.Second
)

type SubService struct {
	repo Repository
}

func NewSubscriptionService(repo Repository) *SubService {
	return &SubService{
		repo: repo,
	}
}

func (s *SubService) Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, shortTimeout)
	defer cancel()

	log := logger.FromContext(ctx)

	sub.ServiceName = strings.TrimSpace(sub.ServiceName)

	if sub.ServiceName == "" {
		err := apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"service_name is required",
		)

		log.Warn("validation failed",
			zap.String("field", "service_name"),
			zap.Error(err),
		)

		return uuid.Nil, err
	}

	if sub.Price <= 0 {
		err := apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"price must be greater than 0",
		)

		log.Warn("validation failed",
			zap.String("field", "price"),
			zap.Error(err),
		)

		return uuid.Nil, err
	}

	if sub.EndDate != nil &&
		sub.EndDate.Before(sub.StartDate) {
		err := apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"end_date must be greater than or equal to start_date",
		)

		log.Warn(
			"validation failed",
			zap.String("field", "end_date"),
			zap.Error(err),
		)

		return uuid.Nil, err
	}

	id, err := s.repo.Create(ctx, sub)
	if err != nil {
		log.Error("failed to create subscription",
			zap.String("user_id", sub.UserID.String()),
			zap.String("service_name", sub.ServiceName),
			zap.Error(err),
		)

		return uuid.Nil, err
	}

	log.Info("subscription created",
		zap.String("subscription_id", id.String()),
		zap.String("user_id", sub.UserID.String()),
	)

	return id, nil
}

func (s *SubService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, shortTimeout)
	defer cancel()

	log := logger.FromContext(ctx)

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to get subscription",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)

		return nil, err
	}
	log.Info("subscription created",
		zap.String("subscription_id", id.String()),
		zap.String("user_id", sub.UserID.String()),
	)

	return sub, nil
}

func (s *SubService) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, mediumTimeout)
	defer cancel()

	log := logger.FromContext(ctx)

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {

		err := apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"offset must be greater than or equal to 0",
		)

		log.Warn("validation failed",
			zap.String("field", "offset"),
			zap.Error(err),
		)

		return nil, err
	}

	subscriptions, err := s.repo.GetByUserID(
		ctx,
		userID,
		limit,
		offset,
	)
	if err != nil {
		log.Error("failed to get subscriptions by user id",
			zap.String("user_id", userID.String()),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(err),
		)

		return nil, err
	}

	return subscriptions, nil
}

func (s *SubService) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, mediumTimeout)
	defer cancel()

	log := logger.FromContext(ctx)

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {

		err := apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"offset must be greater than or equal to 0",
		)

		log.Warn("validation failed",
			zap.String("field", "offset"),
			zap.Error(err),
		)

		return nil, err
	}

	subscriptions, err := s.repo.List(
		ctx,
		limit,
		offset,
	)
	if err != nil {
		log.Error("failed to list subscriptions",
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(err),
		)

		return nil, err
	}

	return subscriptions, nil
}

func (s *SubService) Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, shortTimeout)
	defer cancel()

	log := logger.FromContext(ctx)

	if update.ServiceName == nil && update.Price == nil && update.EndDate == nil {
		err := apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"at least one field is required",
		)

		log.Warn("validation failed",
			zap.Error(err),
		)

		return err
	}

	if update.ServiceName != nil {
		trimmed := strings.TrimSpace(*update.ServiceName)

		if trimmed == "" {
			err := apperrors.WrapMessage(
				nil,
				apperrors.ErrValidation,
				"service_name cannot be empty",
			)

			log.Warn("validation failed",
				zap.String("field", "service_name"),
				zap.Error(err),
			)

			return err
		}

		update.ServiceName = &trimmed
	}

	if update.Price != nil && *update.Price <= 0 {
		err := apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"price must be greater than 0",
		)

		log.Warn("validation failed",
			zap.String("field", "price"),
			zap.Error(err),
		)

		return err
	}

	err := s.repo.Update(
		ctx,
		id,
		update,
	)
	if err != nil {
		log.Error("failed to update subscription",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)

		return err
	}

	log.Info(
		"subscription updated",
		zap.String("subscription_id", id.String()),
	)

	return nil
}

func (s *SubService) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, shortTimeout)
	defer cancel()

	log := logger.FromContext(ctx)

	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Error("failed to delete subscription",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)

		return err
	}

	log.Info("subscription deleted",
		zap.String("subscription_id", id.String()),
	)

	return nil
}

func (s *SubService) Total(ctx context.Context, userID *uuid.UUID, serviceName *string, from *time.Time, to *time.Time) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, longTimeout)
	defer cancel()

	log := logger.FromContext(ctx)

	if from != nil && to != nil && from.After(*to) {
		err := apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"from date must be before to date",
		)

		log.Warn("validation failed",
			zap.String("field", "date_range"),
			zap.Error(err),
		)

		return 0, err
	}

	total, err := s.repo.Total(
		ctx,
		userID,
		serviceName,
		from,
		to,
	)
	if err != nil {
		fields := []zap.Field{zap.Error(err)}

		if userID != nil {
			fields = append(
				fields,
				zap.String("user_id", userID.String()),
			)
		}

		if serviceName != nil {
			fields = append(
				fields,
				zap.String("service_name", *serviceName),
			)
		}

		log.Error("failed to calculate subscriptions total",
			fields...,
		)

		return 0, err
	}

	return total, nil
}
