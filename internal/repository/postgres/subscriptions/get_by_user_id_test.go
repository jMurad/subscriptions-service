package subscriptions_test

import (
	"context"
	"testing"
	"time"

	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByUserID_Success(t *testing.T) {
	setupTest(t)

	userID := uuid.New()

	sub1 := model.Subscription{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      userID,
		StartDate:   time.Now(),
	}

	sub2 := model.Subscription{
		ServiceName: "Spotify",
		Price:       499,
		UserID:      userID,
		StartDate:   time.Now(),
	}

	_, err := repo.Create(context.Background(), sub1)
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), sub2)
	require.NoError(t, err)

	result, err := repo.GetByUserID(context.Background(), userID, 10, 0)

	require.NoError(t, err)

	assert.Len(t, result, 2)

	assert.Equal(t, userID, result[0].UserID)
	assert.Equal(t, userID, result[1].UserID)
}

func TestGetByUserID_Empty(t *testing.T) {
	setupTest(t)

	result, err := repo.GetByUserID(context.Background(), uuid.New(), 10, 0)

	require.NoError(t, err)

	assert.Len(t, result, 0)

	assert.Equal(t, []model.Subscription{}, result)
}

func TestGetByUserID_DeletedSubscriptions(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	err := repo.Delete(context.Background(), sub.ID)

	require.NoError(t, err)

	result, err := repo.GetByUserID(context.Background(), sub.UserID, 10, 0)

	require.NoError(t, err)

	assert.Len(t, result, 0)

	assert.Equal(t, []model.Subscription{}, result)
}
