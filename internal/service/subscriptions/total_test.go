package subscriptions_test

import (
	"context"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/service/subscriptions"
	"subscriptions-service/internal/service/subscriptions/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTotalServiceSubscriptions_SuccessWithoutFilters(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	expectedTotal := 999

	repo.On("Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Return(expectedTotal, nil)

	total, err := svc.Total(
		context.Background(),
		nil,
		nil,
		nil,
		nil,
	)

	assert.NoError(t, err)

	assert.Equal(t, expectedTotal, total)

	repo.AssertExpectations(t)
}

func TestTotalServiceSubscriptions_SuccessWithAllFilters(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	userID := uuid.New()

	serviceName := "Netflix"

	from := time.Now().AddDate(0, -1, 0)
	to := time.Now()

	expectedTotal := 1999

	repo.On("Total",
		mock.Anything,
		&userID,
		&serviceName,
		&from,
		&to,
	).Return(expectedTotal, nil)

	total, err := svc.Total(
		context.Background(),
		&userID,
		&serviceName,
		&from,
		&to,
	)

	assert.NoError(t, err)

	assert.Equal(t, expectedTotal, total)

	repo.AssertExpectations(t)
}

func TestTotalServiceSubscriptions_InvalidDateRange(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	from := time.Now()
	to := from.Add(-24 * time.Hour)

	total, err := svc.Total(
		context.Background(),
		nil,
		nil,
		&from,
		&to,
	)

	assert.Error(t, err)

	assert.Equal(t, 0, total)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrValidation.Code,
		appErr.Code,
	)

	repo.AssertNotCalled(t, "Total")
}

func TestTotalServiceSubscriptions_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	expectedErr := apperrors.ErrInternal

	repo.On("Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Return(0, expectedErr)

	total, err := svc.Total(
		context.Background(),
		nil,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err)

	assert.Equal(t, 0, total)

	assert.Equal(t, expectedErr, err)

	repo.AssertExpectations(t)
}

func TestTotalServiceSubscriptions_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	repo.On("Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		0,
		apperrors.MapPostgresError(context.Canceled),
	)

	total, err := svc.Total(
		ctx,
		nil,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err)

	assert.Equal(t, 0, total)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrCanceled.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}

func TestTotalServiceSubscriptions_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	repo.On("Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		0,
		apperrors.MapPostgresError(context.DeadlineExceeded),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond,
	)

	defer cancel()

	total, err := svc.Total(
		ctx,
		nil,
		nil,
		nil,
		nil,
	)

	assert.Error(t, err)

	assert.Equal(t, 0, total)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrTimeout.Code,
		appErr.Code,
	)

	repo.AssertExpectations(t)
}
