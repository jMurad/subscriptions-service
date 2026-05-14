package service

import (
	"context"
	"subscriptions-service/internal/model"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	List(ctx context.Context, limit, offset int) ([]model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
	Total(ctx context.Context, userID *uuid.UUID, serviceName *string, from time.Time, to time.Time) (int, error)
}
