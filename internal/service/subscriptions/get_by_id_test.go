package subscriptions_test

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service/subscriptions"
	"subscriptions-service/internal/service/subscriptions/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetByIDServiceSubscriptions_Success(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	expected := &model.Subscription{
		ID:          id,
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
	}

	repo.On("GetByID",
		mock.Anything,
		id,
	).Return(expected, nil)

	sub, err := svc.GetByID(
		context.Background(),
		id,
	)

	assert.NoError(t, err)

	assert.Equal(t, expected, sub)

	repo.AssertExpectations(t)
}

func TestGetByIDServiceSubscriptions_NotFound(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	expectedErr := apperrors.Wrap(
		pgx.ErrNoRows,
		apperrors.ErrNotFound,
	)

	repo.On("GetByID",
		mock.Anything,
		id,
	).Return(nil, expectedErr)

	sub, err := svc.GetByID(
		context.Background(),
		id,
	)

	assert.Error(t, err)

	assert.Nil(t, sub)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrNotFound.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}

func TestGetByIDServiceSubscriptions_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	expectedErr := apperrors.ErrInternal

	repo.On("GetByID",
		mock.Anything,
		id,
	).Return(nil, expectedErr)

	sub, err := svc.GetByID(
		context.Background(),
		id,
	)

	assert.Error(t, err)

	assert.Nil(t, sub)

	assert.Equal(t, expectedErr, err)

	repo.AssertExpectations(t)
}

func TestGetByIDServiceSubscriptions_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	repo.On("GetByID",
		mock.Anything,
		id,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()

	}).Return(
		nil,
		apperrors.MapPostgresError(context.Canceled),
	)

	sub, err := svc.GetByID(
		ctx,
		id,
	)

	assert.Error(t, err)

	assert.Nil(t, sub)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrCanceled.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}

func TestGetByIDServiceSubscriptions_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	repo.On("GetByID",
		mock.Anything,
		id,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		nil,
		apperrors.MapPostgresError(context.DeadlineExceeded),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond,
	)

	defer cancel()

	sub, err := svc.GetByID(
		ctx,
		id,
	)

	assert.Error(t, err)

	assert.Nil(t, sub)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrTimeout.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}
