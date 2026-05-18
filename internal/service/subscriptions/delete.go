package subscriptions

import (
	"context"
	"subscriptions-service/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
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
