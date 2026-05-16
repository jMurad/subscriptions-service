package subscriptions

import (
	"context"
	"subscriptions-service/internal/model"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(context.Context, model.Subscription) (uuid.UUID, error)
	GetByID(context.Context, uuid.UUID) (*model.Subscription, error)
	GetByUserID(context.Context, uuid.UUID, int, int) ([]model.Subscription, error)
	List(context.Context, int, int) ([]model.Subscription, error)
	Update(context.Context, uuid.UUID, model.SubscriptionUpdate) error
	Delete(context.Context, uuid.UUID) error
	Total(context.Context, *uuid.UUID, *string, *time.Time, *time.Time) (int, error)
}
