package subscriptions

import (
	"context"

	"subscriptions-service/internal/model"
	"subscriptions-service/internal/repository"

	"github.com/google/uuid"
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
	return s.repo.Create(ctx, sub)
}

func (s *SubService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}
