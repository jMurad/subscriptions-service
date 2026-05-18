package subscriptions_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/transport/http/handler/subscriptions"
	"subscriptions-service/internal/transport/http/handler/subscriptions/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetByID_Success(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	sub := &model.Subscription{
		ID:          id,
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/"+id.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add("id", id.String())

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	w := httptest.NewRecorder()

	service.On(
		"GetByID",
		mock.Anything,
		id,
	).Return(sub, nil)

	handler.GetByID(w, req)

	assert.Equal(
		t,
		http.StatusOK,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestHandler_GetByID_InvalidUUID(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/invalid",
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"id",
		"invalid",
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "GetByID")
}

func TestHandler_GetByID_NotFound(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/"+id.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"id",
		id.String(),
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	w := httptest.NewRecorder()

	service.On(
		"GetByID",
		mock.Anything,
		id,
	).Return(
		nil,
		apperrors.ErrNotFound,
	)

	handler.GetByID(w, req)

	assert.Equal(
		t,
		http.StatusNotFound,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestHandler_GetByID_InternalError(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/"+id.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"id",
		id.String(),
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	w := httptest.NewRecorder()

	service.On(
		"GetByID",
		mock.Anything,
		id,
	).Return(
		nil,
		apperrors.ErrInternal,
	)

	handler.GetByID(w, req)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestHandler_GetByID_ContextCanceled(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/"+id.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"id",
		id.String(),
	)

	ctx, cancel := context.WithCancel(
		req.Context(),
	)

	cancel()

	req = req.WithContext(
		context.WithValue(
			ctx,
			chi.RouteCtxKey,
			rctx,
		),
	)

	w := httptest.NewRecorder()

	service.On(
		"GetByID",
		mock.Anything,
		id,
	).Return(
		nil,
		apperrors.ErrCanceled,
	)

	handler.GetByID(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestHandler_GetByID_ContextDeadlineExceeded(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/"+id.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"id",
		id.String(),
	)

	ctx, cancel := context.WithTimeout(
		req.Context(),
		time.Millisecond,
	)

	defer cancel()

	req = req.WithContext(
		context.WithValue(
			ctx,
			chi.RouteCtxKey,
			rctx,
		),
	)

	w := httptest.NewRecorder()

	service.On(
		"GetByID",
		mock.Anything,
		id,
	).Return(
		nil,
		apperrors.ErrTimeout,
	)

	handler.GetByID(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}
