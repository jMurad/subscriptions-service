package mocks

import (
	"context"

	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (m *Repository) Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error) {
	args := m.Called(ctx, sub)

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *Repository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *Repository) List(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	args := m.Called(ctx, limit, offset)

	return args.Get(0).([]model.Subscription), args.Error(1)
}

func (m *Repository) Update(ctx context.Context, id uuid.UUID, update model.SubscriptionUpdate) error {
	args := m.Called(ctx, id, update)

	return args.Error(0)
}

func (m *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}
