package subscriptions_test

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate_Success(t *testing.T) {
	setupTest(t)

	sub := model.Subscription{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	id, err := repo.Create(context.Background(), sub)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	created, err := repo.GetByID(context.Background(), id)

	require.NoError(t, err)

	assert.Equal(t, sub.ServiceName, created.ServiceName)
	assert.Equal(t, sub.Price, created.Price)
	assert.Equal(t, sub.UserID, created.UserID)
}

func TestCreate_DuplicateSubscription(t *testing.T) {
	setupTest(t)

	sub := model.Subscription{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	id, err := repo.Create(context.Background(), sub)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	id, err = repo.Create(context.Background(), sub)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	assert.True(t, errors.Is(err, apperrors.ErrConflict))
}

func TestCreate_EndDateBeforeStartDate(t *testing.T) {
	setupTest(t)

	endDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	sub := model.Subscription{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     &endDate,
	}

	id, err := repo.Create(context.Background(), sub)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	assert.True(t, errors.Is(err, apperrors.ErrValidation))
}

func TestCreate_EmptyServiceName(t *testing.T) {
	setupTest(t)

	sub := model.Subscription{
		ServiceName: "",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	id, err := repo.Create(context.Background(), sub)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	assert.True(t, errors.Is(err, apperrors.ErrValidation))
}
