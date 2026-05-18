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

func TestTotal_Success(t *testing.T) {
	setupTest(t)

	userID := uuid.New()

	sub1 := model.Subscription{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      userID,
		StartDate:   time.Now(),
	}

	sub2 := model.Subscription{
		ServiceName: "Yandex Plus",
		Price:       200,
		UserID:      userID,
		StartDate:   time.Now(),
	}

	_, err := repo.Create(context.Background(), sub1)

	require.NoError(t, err)

	_, err = repo.Create(context.Background(), sub2)

	require.NoError(t, err)

	total, err := repo.Total(context.Background(), &userID, nil, nil, nil)

	require.NoError(t, err)

	assert.Equal(t, sub1.Price+sub2.Price, total)
}

func TestTotal_Empty(t *testing.T) {
	setupTest(t)

	total, err := repo.Total(context.Background(), nil, nil, nil, nil)

	require.NoError(t, err)

	assert.Equal(t, 0, total)
}

func TestTotal_FilterByService(t *testing.T) {
	setupTest(t)

	seviceName := "Netflix"

	sub1 := model.Subscription{
		ServiceName: seviceName,
		Price:       100,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	sub2 := model.Subscription{
		ServiceName: "Spotify",
		Price:       500,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	sub3 := model.Subscription{
		ServiceName: seviceName,
		Price:       200,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	_, err := repo.Create(context.Background(), sub1)

	require.NoError(t, err)

	_, err = repo.Create(context.Background(), sub2)

	require.NoError(t, err)

	_, err = repo.Create(context.Background(), sub3)

	require.NoError(t, err)

	total, err := repo.Total(context.Background(), nil, &seviceName, nil, nil)

	require.NoError(t, err)

	assert.Equal(t, sub1.Price+sub3.Price, total)
}

func TestTotal_FilterByUser(t *testing.T) {
	setupTest(t)

	userID := uuid.New()

	sub1 := model.Subscription{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      userID,
		StartDate:   time.Now(),
	}

	sub2 := model.Subscription{
		ServiceName: "Spotify",
		Price:       500,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	sub3 := model.Subscription{
		ServiceName: "Youtube",
		Price:       200,
		UserID:      userID,
		StartDate:   time.Now(),
	}

	_, err := repo.Create(context.Background(), sub1)
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), sub2)
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), sub3)
	require.NoError(t, err)

	total, err := repo.Total(context.Background(), &userID, nil, nil, nil)

	require.NoError(t, err)

	assert.Equal(t, sub1.Price+sub3.Price, total)
}

func TestTotal_DeletedExcluded(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	err := repo.Delete(context.Background(), sub.ID)

	require.NoError(t, err)

	total, err := repo.Total(context.Background(), nil, nil, nil, nil)

	require.NoError(t, err)

	assert.Equal(t, 0, total)
}
