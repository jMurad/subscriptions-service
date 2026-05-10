package repository

import (
	"context"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error)
}
