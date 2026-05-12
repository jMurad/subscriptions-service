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

// Create tests
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

// GetByID tests
func TestGetByID_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	expected := &model.Subscription{
		ID:          uuid.New(),
		ServiceName: "Netflix",
		Price:       500,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	mockRepo.
		On(
			"GetByID",
			context.Background(),
			expected.ID,
		).
		Return(expected, nil)

	result, err := svc.GetByID(
		context.Background(),
		expected.ID,
	)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	id := uuid.New()

	expectedErr := errors.New("not found")

	mockRepo.
		On(
			"GetByID",
			context.Background(),
			id,
		).
		Return(nil, expectedErr)

	result, err := svc.GetByID(
		context.Background(),
		id,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}
