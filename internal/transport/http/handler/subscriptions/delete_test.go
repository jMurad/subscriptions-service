package subscriptions_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/transport/http/handler/subscriptions"
	"subscriptions-service/internal/transport/http/handler/subscriptions/mocks"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Delete_Success(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodDelete,
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
		"Delete",
		mock.Anything,
		id,
	).Return(nil)

	handler.Delete(w, req)

	assert.Equal(
		t,
		http.StatusOK,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/subscriptions/invalid",
		nil,
	)

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add("id", "invalid")

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			rctx,
		),
	)

	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "Delete")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodDelete,
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
		"Delete",
		mock.Anything,
		id,
	).Return(
		apperrors.ErrNotFound,
	)

	handler.Delete(w, req)

	assert.Equal(
		t,
		http.StatusNotFound,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestHandler_Delete_InternalError(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodDelete,
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
		"Delete",
		mock.Anything,
		id,
	).Return(
		apperrors.ErrInternal,
	)

	handler.Delete(w, req)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestHandler_Delete_ContextCanceled(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/subscriptions/"+id.String(),
		nil,
	).WithContext(ctx)

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
		"Delete",
		mock.Anything,
		id,
	).Run(func(args mock.Arguments) {
		<-args.Get(0).(context.Context).Done()
	}).Return(
		apperrors.ErrCanceled,
	)

	handler.Delete(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestHandler_Delete_ContextDeadlineExceeded(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond,
	)

	defer cancel()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/subscriptions/"+id.String(),
		nil,
	).WithContext(ctx)

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
		"Delete",
		mock.Anything,
		id,
	).Run(func(args mock.Arguments) {
		<-args.Get(0).(context.Context).Done()
	}).Return(
		apperrors.ErrTimeout,
	)

	handler.Delete(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}
