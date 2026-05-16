package subscriptions_test

import (
	"context"
	"testing"
	"time"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service/mocks"
	"subscriptions-service/internal/service/subscriptions"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Create tests
func TestCreate_Success(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

	subID := uuid.New()

	sub := model.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   time.Now(),
	}

	repo.On("Create",
		mock.Anything,
		sub,
	).Return(subID, nil)

	id, err := svc.Create(
		context.Background(),
		sub,
	)

	assert.NoError(t, err)
	assert.Equal(t, subID, id)

	repo.AssertExpectations(t)
}

func TestCreate_EmptyServiceName(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

	sub := model.Subscription{
		UserID:      uuid.New(),
		ServiceName: "",
		Price:       999,
		StartDate:   time.Now(),
	}

	id, err := svc.Create(
		context.Background(),
		sub,
	)

	assert.Error(t, err)

	assert.Equal(t, uuid.Nil, id)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrValidation.Code,
		appErr.Code,
	)

	repo.AssertNotCalled(t, "Create")
}

func TestCreate_InvalidPrice(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

	sub := model.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       0,
		StartDate:   time.Now(),
	}

	id, err := svc.Create(
		context.Background(),
		sub,
	)

	assert.Error(t, err)

	assert.Equal(t, uuid.Nil, id)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrValidation.Code,
		appErr.Code,
	)

	repo.AssertNotCalled(t, "Create")
}

func TestCreate_InvalidDateRange(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

	start := time.Now()
	end := start.Add(-24 * time.Hour)

	sub := model.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   start,
		EndDate:     &end,
	}

	id, err := svc.Create(
		context.Background(),
		sub,
	)

	assert.Error(t, err)

	assert.Equal(t, uuid.Nil, id)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrValidation.Code,
		appErr.Code,
	)

	repo.AssertNotCalled(t, "Create")
}

func TestCreate_TrimServiceName(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

	subID := uuid.New()

	sub := model.Subscription{
		UserID:      uuid.New(),
		ServiceName: "   Netflix   ",
		Price:       999,
		StartDate:   time.Now(),
	}

	expected := sub
	expected.ServiceName = "Netflix"

	repo.On("Create",
		mock.Anything,
		expected,
	).Return(subID, nil)

	id, err := svc.Create(
		context.Background(),
		sub,
	)

	assert.NoError(t, err)

	assert.Equal(t, subID, id)

	repo.AssertExpectations(t)
}

func TestCreate_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

	sub := model.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   time.Now(),
	}

	expectedErr := apperrors.ErrConflict

	repo.On("Create",
		mock.Anything,
		sub,
	).Return(uuid.Nil, expectedErr)

	id, err := svc.Create(
		context.Background(),
		sub,
	)

	assert.Error(t, err)

	assert.Equal(t, uuid.Nil, id)

	assert.Equal(t, expectedErr, err)

	repo.AssertExpectations(t)
}

func TestCreate_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

	sub := model.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   time.Now(),
	}

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	repo.On("Create",
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		uuid.Nil,
		apperrors.MapPostgresError(context.Canceled),
	)

	id, err := svc.Create(ctx, sub)

	assert.Error(t, err)

	assert.Equal(t, uuid.Nil, id)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrCanceled.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}

func TestCreate_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

	repo.On("Create",
		mock.Anything,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		uuid.Nil,
		apperrors.MapPostgresError(context.DeadlineExceeded),
	)

	sub := model.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   time.Now(),
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond,
	)

	defer cancel()

	id, err := svc.Create(ctx, sub)

	assert.Error(t, err)

	assert.Equal(t, uuid.Nil, id)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrTimeout.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}

// GetByID tests
func TestGetByID_Success(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

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

func TestGetByID_NotFound(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

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

func TestGetByID_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

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

func TestGetByID_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

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

func TestGetByID_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(repo)

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
