package subscriptions_test

import (
	"context"
	"errors"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList_Success(t *testing.T) {
	setupTest(t)

	createTestSubscription(t, repo)
	createTestSubscription(t, repo)

	result, err := repo.List(context.Background(), 10, 0)

	require.NoError(t, err)

	assert.Len(t, result, 2)
}

func TestList_Empty(t *testing.T) {
	setupTest(t)

	result, err := repo.List(context.Background(), 10, 0)

	require.NoError(t, err)

	assert.Len(t, result, 0)

	assert.Equal(t, []model.Subscription{}, result)
}

func TestList_Pagination(t *testing.T) {
	setupTest(t)

	createTestSubscription(t, repo)

	result, err := repo.List(context.Background(), -10, -10)

	assert.Error(t, err)

	assert.True(t, errors.Is(err, apperrors.ErrValidation))

	assert.Nil(t, result)
}
