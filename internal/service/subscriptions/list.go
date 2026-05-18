package subscriptions

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"

	"go.uber.org/zap"
)

func (s *service) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
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
