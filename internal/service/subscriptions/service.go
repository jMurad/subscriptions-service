package subscriptions

import (
	"context"
	"time"

	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubService {
	return &SubService{
		repo: repo,
	}
}

func (s *SubService) Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error) {
	log := logger.FromContext(ctx)
	start := time.Now()

	log.Info("service create subscription started",
		zap.String("service_name", sub.ServiceName),
		zap.String("user_id", sub.UserID.String()),
	)

	id, err := s.repo.Create(ctx, sub)
	if err != nil {
		log.Error("service create subscription failed", zap.Error(err))
		return uuid.Nil, err
	}

	log.Info("service create subscription completed",
		zap.String("subscription_id", id.String()),
		zap.Duration("duration", time.Since(start)),
	)

	return id, nil
}

func (s *SubService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	log := logger.FromContext(ctx)
	start := time.Now()

	log.Info("service get subscription started", zap.String("subscription_id", id.String()))

	subscription, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("service get subscription failed",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("service get subscription completed",
		zap.String("subscription_id", id.String()),
		zap.Duration("duration", time.Since(start)),
	)

	return subscription, nil
}

func (s *SubService) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	log := logger.FromContext(ctx)
	start := time.Now()

	log.Info("service list subscriptions started",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	subscriptions, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		log.Error("service list subscriptions failed", zap.Error(err))
		return nil, err
	}

	log.Info("service list subscriptions completed",
		zap.Int("count", len(subscriptions)),
		zap.Duration("duration", time.Since(start)),
	)

	return subscriptions, nil
}

func (s *SubService) Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error {
	log := logger.FromContext(ctx)
	start := time.Now()

	log.Info("service update subscription started", zap.String("subscription_id", id.String()))

	err := s.repo.Update(ctx, id, update)
	if err != nil {
		log.Error("service update subscription failed",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	log.Info("service update subscription completed",
		zap.String("subscription_id", id.String()),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

func (s *SubService) Delete(ctx context.Context, id uuid.UUID) error {
	log := logger.FromContext(ctx)
	start := time.Now()

	log.Info("service delete subscription started", zap.String("subscription_id", id.String()))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Error("service delete subscription failed",
			zap.String("subscription_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	log.Info("service delete subscription completed",
		zap.String("subscription_id", id.String()),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

func (s *SubService) Total(ctx context.Context, userID *uuid.UUID, serviceName *string, from time.Time, to time.Time) (int, error) {
	return s.repo.Total(ctx, userID, serviceName, from, to)
}
