package repository

import (
	"context"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	List(ctx context.Context, limit, offset int) ([]model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error
}
