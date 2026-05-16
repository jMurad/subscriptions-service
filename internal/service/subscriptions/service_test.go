package subscriptions_test

import (
	"context"
	"testing"
	"time"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service/subscriptions"
	"subscriptions-service/internal/service/subscriptions/mocks"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Create tests
func TestCreate_Success(t *testing.T) {
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

func TestCreate_EmptyServiceName(t *testing.T) {
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

func TestCreate_InvalidPrice(t *testing.T) {
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

func TestCreate_InvalidDateRange(t *testing.T) {
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

func TestCreate_TrimServiceName(t *testing.T) {
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

func TestCreate_RepositoryError(t *testing.T) {
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

func TestCreate_ContextCanceled(t *testing.T) {
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

func TestCreate_ContextDeadlineExceeded(t *testing.T) {
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

// GetByID tests
func TestGetByID_Success(t *testing.T) {
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

func TestGetByID_NotFound(t *testing.T) {
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

func TestGetByID_RepositoryError(t *testing.T) {
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

func TestGetByID_ContextCanceled(t *testing.T) {
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

func TestGetByID_ContextDeadlineExceeded(t *testing.T) {
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

// GetByUserID tests
func TestGetByUserID_Success(t *testing.T) {
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

func TestGetByUserID_DefaultLimit(t *testing.T) {
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

func TestGetByUserID_MaxLimit(t *testing.T) {
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

func TestGetByUserID_InvalidOffset(t *testing.T) {
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

func TestGetByUserID_RepositoryError(t *testing.T) {
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

func TestGetByUserID_ContextCanceled(t *testing.T) {
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

func TestGetByUserID_ContextDeadlineExceeded(t *testing.T) {
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

// List tests
func TestList_Success(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	expected := []model.Subscription{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			ServiceName: "Netflix",
			Price:       999,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
		},
	}

	repo.On("List",
		mock.Anything,
		10,
		0,
	).Return(expected, nil)

	subs, err := svc.List(
		context.Background(),
		10,
		0,
	)

	assert.NoError(t, err)

	assert.Equal(t, expected, subs)

	repo.AssertExpectations(t)
}

func TestList_DefaultLimit(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	expected := []model.Subscription{}

	repo.On("List",
		mock.Anything,
		10,
		0,
	).Return(expected, nil)

	subs, err := svc.List(
		context.Background(),
		0,
		0,
	)

	assert.NoError(t, err)

	assert.Equal(t, expected, subs)

	repo.AssertExpectations(t)
}

func TestList_MaxLimit(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	expected := []model.Subscription{}

	repo.On(
		"List",
		mock.Anything,
		100,
		0,
	).Return(expected, nil)

	subs, err := svc.List(
		context.Background(),
		1000,
		0,
	)

	assert.NoError(t, err)

	assert.Equal(t, expected, subs)

	repo.AssertExpectations(t)
}

func TestList_InvalidOffset(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	subs, err := svc.List(
		context.Background(),
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

	repo.AssertNotCalled(t, "List")
}

func TestList_RepositoryError(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	expectedErr := apperrors.ErrInternal

	repo.On("List",
		mock.Anything,
		10,
		0,
	).Return(nil, expectedErr)

	subs, err := svc.List(
		context.Background(),
		10,
		0,
	)

	assert.Error(t, err)

	assert.Nil(t, subs)

	assert.Equal(t, expectedErr, err)

	repo.AssertExpectations(t)
}

func TestList_ContextCanceled(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	repo.On("List",
		mock.Anything,
		10,
		0,
	).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		<-ctx.Done()
	}).Return(
		nil,
		apperrors.MapPostgresError(context.Canceled),
	)

	subs, err := svc.List(
		ctx,
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

func TestList_ContextDeadlineExceeded(t *testing.T) {
	repo := new(mocks.Repository)

	svc := subscriptions.NewService(repo)

	repo.On("List",
		mock.Anything,
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

	subs, err := svc.List(
		ctx,
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

// Update tests
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

// Delete tests
func TestDelete_Success(t *testing.T) {
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

func TestDelete_NotFound(t *testing.T) {
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

func TestDelete_RepositoryError(t *testing.T) {
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

func TestDelete_ContextCanceled(t *testing.T) {
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

func TestDelete_ContextDeadlineExceeded(t *testing.T) {
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

// Total tests
func TestTotal_SuccessWithoutFilters(t *testing.T) {
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

func TestTotal_SuccessWithAllFilters(t *testing.T) {
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

func TestTotal_InvalidDateRange(t *testing.T) {
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

func TestTotal_RepositoryError(t *testing.T) {
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

func TestTotal_ContextCanceled(t *testing.T) {
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

func TestTotal_ContextDeadlineExceeded(t *testing.T) {
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
