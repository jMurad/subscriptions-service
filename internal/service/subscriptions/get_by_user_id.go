package subscriptions

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Subscription, error) {
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
