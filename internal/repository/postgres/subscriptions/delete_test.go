package subscriptions_test

import (
	"context"
	"errors"
	"testing"

	apperrors "subscriptions-service/internal/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteRepositorySubscriptions_Success(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	err := repo.Delete(context.Background(), sub.ID)

	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), sub.ID)

	require.Error(t, err)

	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}

func TestDeleteRepositorySubscriptions_NotFound(t *testing.T) {
	setupTest(t)

	err := repo.Delete(context.Background(), uuid.New())

	require.Error(t, err)

	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}

func TestDeleteRepositorySubscriptions_AlreadyDeleted(t *testing.T) {
	setupTest(t)

	sub := createTestSubscription(t, repo)

	err := repo.Delete(context.Background(), sub.ID)

	require.NoError(t, err)

	err = repo.Delete(context.Background(), sub.ID)

	require.Error(t, err)

	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}
