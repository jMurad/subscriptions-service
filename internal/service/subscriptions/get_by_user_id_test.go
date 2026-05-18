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

func TestGetByUserIDServiceSubscriptions_Success(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	userID := uuid.New()

	expected := []model.Subscription{
		{
			ID:          uuid.New(),
			UserID:      userID,
			ServiceName: "Netflix",
			Price:       999,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
		},
	}

	repo.On("GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(expected, nil)

	subs, err := svc.GetByUserID(
		context.Background(),
		userID,
		10,
		0,
	)

	assert.NoError(t, err)

	assert.Equal(t, expected, subs)

	repo.AssertExpectations(t)
}

func TestGetByUserIDServiceSubscriptions_DefaultLimit(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	userID := uuid.New()

	expected := []model.Subscription{}

	repo.On("GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(expected, nil)

	subs, err := svc.GetByUserID(
		context.Background(),
		userID,
		0,
		0,
	)

	assert.NoError(t, err)

	assert.Equal(t, expected, subs)

	repo.AssertExpectations(t)
}

func TestGetByUserIDServiceSubscriptions_MaxLimit(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	userID := uuid.New()

	expected := []model.Subscription{}

	repo.On("GetByUserID",
		mock.Anything,
		userID,
		100,
		0,
	).Return(expected, nil)

	subs, err := svc.GetByUserID(
		context.Background(),
		userID,
		1000,
		0,
	)

	assert.NoError(t, err)

	assert.Equal(t, expected, subs)

	repo.AssertExpectations(t)
}

func TestGetByUserIDServiceSubscriptions_InvalidOffset(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	userID := uuid.New()

	subs, err := svc.GetByUserID(
		context.Background(),
		userID,
		10,
		-1,
	)

	assert.Error(t, err)

	assert.Nil(t, subs)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrValidation.Code,
		appErr.Code,
	)

	repo.AssertNotCalled(t, "GetByUserID")
}

func TestGetByUserIDServiceSubscriptions_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	userID := uuid.New()

	expectedErr := apperrors.ErrInternal

	repo.On("GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(nil, expectedErr)

	subs, err := svc.GetByUserID(
		context.Background(),
		userID,
		10,
		0,
	)

	assert.Error(t, err)

	assert.Nil(t, subs)

	assert.Equal(t, expectedErr, err)

	repo.AssertExpectations(t)
}

func TestGetByUserIDServiceSubscriptions_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	userID := uuid.New()

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	repo.On("GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		nil,
		apperrors.MapPostgresError(context.Canceled),
	)

	subs, err := svc.GetByUserID(
		ctx,
		userID,
		10,
		0,
	)

	assert.Error(t, err)

	assert.Nil(t, subs)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrCanceled.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}

func TestGetByUserIDServiceSubscriptions_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	userID := uuid.New()

	repo.On(
		"GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
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

	subs, err := svc.GetByUserID(
		ctx,
		userID,
		10,
		0,
	)

	assert.Error(t, err)

	assert.Nil(t, subs)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrTimeout.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}
