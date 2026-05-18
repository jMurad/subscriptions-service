package subscriptions_test

import (
	"context"
	"errors"
	apperrors "subscriptions-service/internal/errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByID_Success(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	got, err := repo.GetByID(context.Background(), sub.ID)

	require.NoError(t, err)

	assert.Equal(t, sub.ID, got.ID)
	assert.Equal(t, sub.ServiceName, got.ServiceName)
	assert.Equal(t, sub.Price, got.Price)
	assert.Equal(t, sub.UserID, got.UserID)
	assert.Equal(t, sub.StartDate, got.StartDate)
}

func TestGetByID_NotFound(t *testing.T) {
	setupTest(t)

	got, err := repo.GetByID(context.Background(), uuid.New())

	require.Error(t, err)

	require.Nil(t, got)

	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}

func TestGetByID_DeletedSubscription(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	err := repo.Delete(context.Background(), sub.ID)

	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), sub.ID)

	require.Error(t, err)

	require.Nil(t, got)

	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}
