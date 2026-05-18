package subscriptions_test

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/service/subscriptions"
	"subscriptions-service/internal/service/subscriptions/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteServiceSubscriptions_Success(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	repo.On("Delete",
		mock.Anything,
		id,
	).Return(nil)

	err := svc.Delete(
		context.Background(),
		id,
	)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestDeleteServiceSubscriptions_NotFound(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	expectedErr := apperrors.Wrap(
		pgx.ErrNoRows,
		apperrors.ErrNotFound,
	)

	repo.On("Delete",
		mock.Anything,
		id,
	).Return(expectedErr)

	err := svc.Delete(
		context.Background(),
		id,
	)

	assert.Error(t, err)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrNotFound.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}

func TestDeleteServiceSubscriptions_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	expectedErr := apperrors.ErrInternal

	repo.On("Delete",
		mock.Anything,
		id,
	).Return(expectedErr)

	err := svc.Delete(
		context.Background(),
		id,
	)

	assert.Error(t, err)

	assert.Equal(t, expectedErr, err)

	repo.AssertExpectations(t)
}

func TestDeleteServiceSubscriptions_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	repo.On("Delete",
		mock.Anything,
		id,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		apperrors.MapPostgresError(context.Canceled),
	)

	err := svc.Delete(
		ctx,
		id,
	)

	assert.Error(t, err)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrCanceled.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}

func TestDeleteServiceSubscriptions_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	repo.On("Delete",
		mock.Anything,
		id,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		apperrors.MapPostgresError(context.DeadlineExceeded),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond,
	)
	defer cancel()

	err := svc.Delete(
		ctx,
		id,
	)

	assert.Error(t, err)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrTimeout.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}
