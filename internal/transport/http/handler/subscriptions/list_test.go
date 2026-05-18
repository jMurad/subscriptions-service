package subscriptions_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	dto "subscriptions-service/internal/transport/http/dto/subscriptions"
	"subscriptions-service/internal/transport/http/handler/subscriptions"
	"subscriptions-service/internal/transport/http/handler/subscriptions/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListHandlerSubscriptions_Success(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions?limit=10&offset=0",
		nil,
	)

	w := httptest.NewRecorder()

	expected := []model.Subscription{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			ServiceName: "Netflix",
			Price:       999,
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
		},
	}

	service.On(
		"List",
		mock.Anything,
		10,
		0,
	).Return(expected, nil)

	handler.List(w, req)

	assert.Equal(
		t,
		http.StatusOK,
		w.Code,
	)

	var response []dto.SubscriptionResponse

	err := json.Unmarshal(
		w.Body.Bytes(),
		&response,
	)

	assert.NoError(t, err)

	assert.Len(t, response, 1)

	service.AssertExpectations(t)
}

func TestListHandlerSubscriptions_DefaultPagination(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions",
		nil,
	)

	w := httptest.NewRecorder()

	service.On(
		"List",
		mock.Anything,
		10,
		0,
	).Return([]model.Subscription{}, nil)

	handler.List(w, req)

	assert.Equal(
		t,
		http.StatusOK,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestListHandlerSubscriptions_InvalidLimit(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions?limit=invalid",
		nil,
	)

	w := httptest.NewRecorder()

	handler.List(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "List")
}

func TestListHandlerSubscriptions_InvalidOffset(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions?offset=invalid",
		nil,
	)

	w := httptest.NewRecorder()

	handler.List(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "List")
}

func TestListHandlerSubscriptions_NegativeOffset(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions?offset=-1",
		nil,
	)

	w := httptest.NewRecorder()

	handler.List(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "List")
}

func TestListHandlerSubscriptions_InternalError(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions",
		nil,
	)

	w := httptest.NewRecorder()

	service.On(
		"List",
		mock.Anything,
		10,
		0,
	).Return(
		nil,
		apperrors.ErrInternal,
	)

	handler.List(w, req)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestListHandlerSubscriptions_ContextCanceled(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions",
		nil,
	).WithContext(ctx)

	w := httptest.NewRecorder()

	service.On(
		"List",
		mock.Anything,
		10,
		0,
	).Run(func(args mock.Arguments) {
		<-args.Get(0).(context.Context).Done()
	}).Return(
		nil,
		apperrors.ErrCanceled,
	)

	handler.List(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestListHandlerSubscriptions_ContextDeadlineExceeded(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond,
	)
	defer cancel()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions",
		nil,
	).WithContext(ctx)

	w := httptest.NewRecorder()

	service.On(
		"List",
		mock.Anything,
		10,
		0,
	).Run(func(args mock.Arguments) {
		<-args.Get(0).(context.Context).Done()
	}).Return(
		nil,
		apperrors.ErrTimeout,
	)

	handler.List(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}
