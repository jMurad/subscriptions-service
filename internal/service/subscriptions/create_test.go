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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateServiceSubscriptions_Success(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

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

func TestCreateServiceSubscriptions_EmptyServiceName(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

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

func TestCreateServiceSubscriptions_InvalidPrice(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

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

func TestCreateServiceSubscriptions_InvalidDateRange(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

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

func TestCreateServiceSubscriptions_TrimServiceName(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

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

func TestCreateServiceSubscriptions_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

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

func TestCreateServiceSubscriptions_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

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

func TestCreateServiceSubscriptions_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

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
