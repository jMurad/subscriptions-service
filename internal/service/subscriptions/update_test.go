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

func TestUpdate_Success(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	price := 999

	update := model.SubscriptionUpdate{
		Price: &price,
	}

	repo.On("Update",
		mock.Anything,
		id,
		update,
	).Return(nil)

	err := svc.Update(
		context.Background(),
		id,
		update,
	)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUpdate_EmptyUpdate(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	err := svc.Update(
		context.Background(),
		uuid.New(),
		model.SubscriptionUpdate{},
	)

	assert.Error(t, err)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrValidation.Code,
		appErr.Code,
	)

	repo.AssertNotCalled(t, "Update")
}

func TestUpdate_InvalidPrice(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	price := 0

	update := model.SubscriptionUpdate{
		Price: &price,
	}

	err := svc.Update(
		context.Background(),
		uuid.New(),
		update,
	)

	assert.Error(t, err)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrValidation.Code,
		appErr.Code,
	)

	repo.AssertNotCalled(t, "Update")
}

func TestUpdate_EmptyServiceName(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	name := ""

	update := model.SubscriptionUpdate{
		ServiceName: &name,
	}

	err := svc.Update(
		context.Background(),
		uuid.New(),
		update,
	)

	assert.Error(t, err)

	appErr, ok := err.(*apperrors.Error)

	assert.True(t, ok)

	assert.Equal(
		t,
		apperrors.ErrValidation.Code,
		appErr.Code,
	)

	repo.AssertNotCalled(t, "Update")
}

func TestUpdate_TrimServiceName(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	name := "   Netflix   "

	expectedName := "Netflix"

	update := model.SubscriptionUpdate{
		ServiceName: &name,
	}

	expected := model.SubscriptionUpdate{
		ServiceName: &expectedName,
	}

	repo.On("Update",
		mock.Anything,
		id,
		expected,
	).Return(nil)

	err := svc.Update(
		context.Background(),
		id,
		update,
	)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUpdate_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	price := 999

	update := model.SubscriptionUpdate{
		Price: &price,
	}

	expectedErr := apperrors.ErrConflict

	repo.On("Update",
		mock.Anything,
		id,
		update,
	).Return(expectedErr)

	err := svc.Update(
		context.Background(),
		id,
		update,
	)

	assert.Error(t, err)

	assert.Equal(t, expectedErr, err)

	repo.AssertExpectations(t)
}
func TestUpdate_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	price := 999

	update := model.SubscriptionUpdate{
		Price: &price,
	}

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	repo.On("Update",
		mock.Anything,
		id,
		update,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		apperrors.MapPostgresError(context.Canceled),
	)

	err := svc.Update(
		ctx,
		id,
		update,
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

func TestUpdate_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	id := uuid.New()

	price := 999

	update := model.SubscriptionUpdate{
		Price: &price,
	}

	repo.On("Update",
		mock.Anything,
		id,
		update,
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

	err := svc.Update(
		ctx,
		id,
		update,
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
