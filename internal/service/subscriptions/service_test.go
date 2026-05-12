package subscriptions_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service/mocks"
	"subscriptions-service/internal/service/subscriptions"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreate_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	expectedID := uuid.New()

	sub := model.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	mockRepo.
		On("Create", context.Background(), sub).
		Return(expectedID, nil)

	id, err := svc.Create(
		context.Background(),
		sub,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	mockRepo.AssertExpectations(t)
}

func TestCreate_RepositoryError(t *testing.T) {

	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	sub := model.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	expectedErr := errors.New("db error")

	mockRepo.
		On("Create", context.Background(), sub).
		Return(uuid.Nil, expectedErr)

	id, err := svc.Create(
		context.Background(),
		sub,
	)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}
