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
