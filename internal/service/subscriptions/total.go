package subscriptions

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) Total(ctx context.Context, userID *uuid.UUID, serviceName *string, from *time.Time, to *time.Time) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, longTimeout)
	defer cancel()

	log := logger.FromContext(ctx)

	effectiveFrom := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	effectiveTo := time.Now()

	if from != nil {
		effectiveFrom = *from
	}

	if to != nil {
		effectiveTo = *to
	}

	if effectiveFrom.After(effectiveTo) {
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
		effectiveFrom,
		effectiveTo,
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
