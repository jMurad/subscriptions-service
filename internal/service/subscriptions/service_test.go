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
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	expected := &model.Subscription{
		ID:          uuid.New(),
		ServiceName: "Netflix",
		Price:       500,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	mockRepo.
		On(
			"GetByID",
			context.Background(),
			expected.ID,
		).
		Return(expected, nil)

	result, err := svc.GetByID(
		context.Background(),
		expected.ID,
	)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	id := uuid.New()

	expectedErr := errors.New("not found")

	mockRepo.
		On(
			"GetByID",
			context.Background(),
			id,
		).
		Return(nil, expectedErr)

	result, err := svc.GetByID(
		context.Background(),
		id,
	)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

// List tests
func TestList_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	expected := []model.Subscription{
		{
			ID:          uuid.New(),
			ServiceName: "Netflix",
			Price:       500,
		},
	}

	mockRepo.
		On("List", context.Background(), 10, 0).
		Return(expected, nil)

	result, err := svc.List(
		context.Background(),
		10,
		0,
	)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestList_Empty(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	mockRepo.
		On("List", context.Background(), 10, 0).
		Return([]model.Subscription{}, nil)

	result, err := svc.List(
		context.Background(),
		10,
		0,
	)

	assert.NoError(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

// Update tests
func TestUpdate_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	id := uuid.New()

	price := 400

	update := model.SubscriptionUpdate{
		Price: &price,
	}

	mockRepo.
		On("Update", context.Background(), id, update).
		Return(nil)

	err := svc.Update(
		context.Background(),
		id,
		update,
	)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdate_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	id := uuid.New()

	price := 400

	update := model.SubscriptionUpdate{
		Price: &price,
	}

	expectedErr := errors.New("db error")

	mockRepo.
		On("Update", context.Background(), id, update).
		Return(expectedErr)

	err := svc.Update(
		context.Background(),
		id,
		update,
	)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

// Delete tests
func TestDelete_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	id := uuid.New()

	mockRepo.
		On("Delete", context.Background(), id).
		Return(nil)

	err := svc.Delete(
		context.Background(),
		id,
	)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDelete_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	id := uuid.New()

	expectedErr := errors.New("delete failed")

	mockRepo.
		On("Delete", context.Background(), id).
		Return(expectedErr)

	err := svc.Delete(
		context.Background(),
		id,
	)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

// Total tests
func TestTotal_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	from := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	to := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)

	mockRepo.
		On("Total", context.Background(), (*uuid.UUID)(nil), (*string)(nil), from, to).
		Return(2000, nil)

	total, err := svc.Total(
		context.Background(),
		nil,
		nil,
		from,
		to,
	)

	assert.NoError(t, err)

	assert.Equal(t, 2000, total)

	mockRepo.AssertExpectations(t)
}

func TestTotal_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := subscriptions.NewSubscriptionService(mockRepo)

	from := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	to := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)

	expectedErr := errors.New("repository error")

	mockRepo.
		On("Total", context.Background(), (*uuid.UUID)(nil), (*string)(nil), from, to).
		Return(0, expectedErr)

	total, err := svc.Total(
		context.Background(),
		nil,
		nil,
		from,
		to,
	)

	assert.Error(t, err)

	assert.Equal(t, 0, total)

	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}
