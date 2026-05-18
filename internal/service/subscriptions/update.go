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

func (s *service) Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error {
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
