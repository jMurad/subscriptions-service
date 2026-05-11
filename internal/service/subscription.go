package service

import (
	"context"

	"subscriptions-service/internal/model"
	"subscriptions-service/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error) {
	return s.repo.Create(ctx, sub)
}

func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}
