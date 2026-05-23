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

func TestTotalRepositorySubscriptions_Success(t *testing.T) {
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

	effectiveFrom := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	effectiveTo := time.Now()

	total, err := repo.Total(context.Background(), &userID, nil, effectiveFrom, effectiveTo)

	require.NoError(t, err)

	assert.Equal(t, sub1.Price+sub2.Price, total)
}

func TestTotalRepositorySubscriptions_Empty(t *testing.T) {
	setupTest(t)

	effectiveFrom := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	effectiveTo := time.Now()

	total, err := repo.Total(context.Background(), nil, nil, effectiveFrom, effectiveTo)

	require.NoError(t, err)

	assert.Equal(t, 0, total)
}

func TestTotalRepositorySubscriptions_FilterByService(t *testing.T) {
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

	effectiveFrom := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	effectiveTo := time.Now()

	total, err := repo.Total(context.Background(), nil, &seviceName, effectiveFrom, effectiveTo)

	require.NoError(t, err)

	assert.Equal(t, sub1.Price+sub3.Price, total)
}

func TestTotalRepositorySubscriptions_FilterByUser(t *testing.T) {
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

	effectiveFrom := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	effectiveTo := time.Now()

	total, err := repo.Total(context.Background(), &userID, nil, effectiveFrom, effectiveTo)

	require.NoError(t, err)

	assert.Equal(t, sub1.Price+sub3.Price, total)
}

func TestTotalRepositorySubscriptions_OverlapPeriod(t *testing.T) {
	setupTest(t)

	userID := uuid.New()

	startDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	endDate := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)

	sub := model.Subscription{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     &endDate,
	}

	_, err := repo.Create(context.Background(), sub)

	require.NoError(t, err)

	effectiveFrom := time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC)

	effectiveTo := time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC)

	total, err := repo.Total(context.Background(), &userID, nil, effectiveFrom, effectiveTo)

	require.NoError(t, err)

	assert.Equal(t, 2*sub.Price, total)
}

func TestTotalRepositorySubscriptions_DeletedExcluded(t *testing.T) {
	setupTest(t)

	var (
		total_before, total_after int
		err                       error
	)

	effectiveFrom := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	effectiveTo := time.Now()

	sub := createTestSubscription(t, repo)

	total_before, err = repo.Total(context.Background(), nil, nil, effectiveFrom, effectiveTo)

	require.NoError(t, err)

	err = repo.Delete(context.Background(), sub.ID)

	require.NoError(t, err)

	total_after, err = repo.Total(context.Background(), nil, nil, effectiveFrom, effectiveTo)

	require.NoError(t, err)

	assert.Equal(t, total_before, total_after)
}
