package service_test

import (
	"context"
	"testing"
	"time"

	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service"
	"subscriptions-service/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubscription(t *testing.T) {

	mockRepo := new(mocks.Repository)

	svc := service.NewSubscriptionService(mockRepo)

	expectedID := uuid.New()

	sub := model.Subscription{
		ServiceName: "Yandex Plus",
		Price:       500,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	mockRepo.
		On("Create", context.Background(), sub).
		Return(expectedID, nil)

	id, err := svc.Create(context.Background(), sub)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	mockRepo.AssertExpectations(t)
}
