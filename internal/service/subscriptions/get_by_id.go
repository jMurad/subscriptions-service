package subscriptions

import (
	"context"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
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
