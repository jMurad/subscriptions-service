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

// Success
func TestHandler_GetByUserID_Success(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

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

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/"+userID.String()+"?limit=10&offset=0",
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		userID.String(),
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	rec := httptest.NewRecorder()

	service.On(
		"GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(expected, nil)

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusOK,
		rec.Code,
	)

	service.AssertExpectations(t)
}

// Invalid user_id
func TestHandler_GetByUserID_InvalidUserID(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/invalid",
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		"invalid",
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	rec := httptest.NewRecorder()

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		rec.Code,
	)

	service.AssertNotCalled(
		t,
		"GetByUserID",
	)
}

// Invalid limit
func TestHandler_GetByUserID_InvalidLimit(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

	userID := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/"+userID.String()+"?limit=invalid",
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		userID.String(),
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	rec := httptest.NewRecorder()

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		rec.Code,
	)

	service.AssertNotCalled(
		t,
		"GetByUserID",
	)
}

// Invalid offset
func TestHandler_GetByUserID_InvalidOffset(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

	userID := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/"+userID.String()+"?offset=-1",
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		userID.String(),
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	rec := httptest.NewRecorder()

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		rec.Code,
	)

	service.AssertNotCalled(
		t,
		"GetByUserID",
	)
}

// Not found
func TestHandler_GetByUserID_NotFound(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

	userID := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/"+userID.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		userID.String(),
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	rec := httptest.NewRecorder()

	service.On(
		"GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(
		nil,
		apperrors.ErrNotFound,
	)

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusNotFound,
		rec.Code,
	)

	service.AssertExpectations(t)
}

// Internal error
func TestHandler_GetByUserID_InternalError(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

	userID := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/"+userID.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		userID.String(),
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	rec := httptest.NewRecorder()

	service.On(
		"GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(
		nil,
		apperrors.ErrInternal,
	)

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		rec.Code,
	)

	service.AssertExpectations(t)
}

// Conflict error
func TestHandler_GetByUserID_ConflictError(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

	userID := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/"+userID.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		userID.String(),
	)

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	rec := httptest.NewRecorder()

	service.On(
		"GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(
		nil,
		apperrors.ErrConflict,
	)

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusConflict,
		rec.Code,
	)

	service.AssertExpectations(t)
}

// Context canceled
func TestHandler_GetByUserID_ContextCanceled(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

	userID := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/"+userID.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		userID.String(),
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

	rec := httptest.NewRecorder()

	service.On(
		"GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(
		nil,
		apperrors.ErrCanceled,
	)

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		rec.Code,
	)

	service.AssertExpectations(t)
}

// Context deadline exceeded
func TestHandler_GetByUserID_ContextDeadlineExceeded(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(
		service,
	)

	userID := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/user/"+userID.String(),
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add(
		"user_id",
		userID.String(),
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

	rec := httptest.NewRecorder()

	service.On(
		"GetByUserID",
		mock.Anything,
		userID,
		10,
		0,
	).Return(
		nil,
		apperrors.ErrTimeout,
	)

	handler.GetByUserID(rec, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		rec.Code,
	)

	service.AssertExpectations(t)
}
