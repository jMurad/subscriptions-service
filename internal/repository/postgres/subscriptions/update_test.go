package subscriptions_test

import (
	"context"
	"errors"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_Success(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	newServiceName := "YouTube"
	newPrice := 1999

	upd := model.SubscriptionUpdate{
		ServiceName: &newServiceName,
		Price:       &newPrice,
	}

	err := repo.Update(context.Background(), sub.ID, upd)

	require.NoError(t, err)

	updated, err := repo.GetByID(context.Background(), sub.ID)

	require.NoError(t, err)

	assert.Equal(t, newServiceName, updated.ServiceName)
	assert.Equal(t, newPrice, updated.Price)
}

func TestUpdate_NotFound(t *testing.T) {
	setupTest(t)

	createTestSubscription(t, repo)

	nonExistID := uuid.New()

	newServiceName := "YouTube"

	upd := model.SubscriptionUpdate{
		ServiceName: &newServiceName,
	}

	err := repo.Update(context.Background(), nonExistID, upd)

	require.Error(t, err)

	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}

func TestUpdate_Deleted(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	err := repo.Delete(context.Background(), sub.ID)

	require.NoError(t, err)

	newPrice := 1999

	upd := model.SubscriptionUpdate{
		Price: &newPrice,
	}

	err = repo.Update(context.Background(), sub.ID, upd)

	require.Error(t, err)

	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}

func TestUpdate_InvalidDates(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	endDate := sub.StartDate.AddDate(-1, 0, 0)

	upd := model.SubscriptionUpdate{
		EndDate: &endDate,
	}

	err := repo.Update(context.Background(), sub.ID, upd)

	require.Error(t, err)

	assert.True(t, errors.Is(err, apperrors.ErrValidation))
}
