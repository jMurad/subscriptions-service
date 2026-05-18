package subscriptions

import (
	"context"
	"strings"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error) {
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
