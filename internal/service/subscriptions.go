package service

import (
	"context"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
}
