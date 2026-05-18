package subscriptions_test

import (
	"bytes"
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

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateHandlerSubscriptions_Success(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	price := 1499

	payload, _ := json.Marshal(dto.UpdateRequest{
		Price: &price,
	})

	req := httptest.NewRequest(
		http.MethodPatch,
		"/subscriptions/"+id.String(),
		bytes.NewReader(payload),
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

	expected := model.SubscriptionUpdate{
		Price: &price,
	}

	service.On(
		"Update",
		mock.Anything,
		id,
		expected,
	).Return(nil)

	handler.Update(w, req)

	assert.Equal(
		t,
		http.StatusOK,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestUpdateHandlerSubscriptions_InvalidJSON(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodPatch,
		"/subscriptions/"+id.String(),
		bytes.NewBufferString("{invalid"),
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

	handler.Update(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "Update")
}

func TestUpdateHandlerSubscriptions_InvalidEndDate(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	endDate := "invalid"

	payload, _ := json.Marshal(dto.UpdateRequest{
		EndDate: &endDate,
	})

	req := httptest.NewRequest(
		http.MethodPatch,
		"/subscriptions/"+id.String(),
		bytes.NewReader(payload),
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

	handler.Update(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "Update")
}

func TestUpdateHandlerSubscriptions_InvalidID(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	price := 999

	payload, _ := json.Marshal(dto.UpdateRequest{
		Price: &price,
	})

	req := httptest.NewRequest(
		http.MethodPatch,
		"/subscriptions/invalid",
		bytes.NewReader(payload),
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

	handler.Update(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "Update")
}

func TestUpdateHandlerSubscriptions_NotFound(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	price := 999

	payload, _ := json.Marshal(dto.UpdateRequest{
		Price: &price,
	})

	req := httptest.NewRequest(
		http.MethodPatch,
		"/subscriptions/"+id.String(),
		bytes.NewReader(payload),
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
		"Update",
		mock.Anything,
		id,
		mock.Anything,
	).Return(
		apperrors.ErrNotFound,
	)

	handler.Update(w, req)

	assert.Equal(
		t,
		http.StatusNotFound,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestUpdateHandlerSubscriptions_InternalError(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	price := 1499

	payload, _ := json.Marshal(dto.UpdateRequest{
		Price: &price,
	})

	req := httptest.NewRequest(
		http.MethodPatch,
		"/subscriptions/"+id.String(),
		bytes.NewReader(payload),
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
		"Update",
		mock.Anything,
		id,
		mock.Anything,
	).Return(
		apperrors.ErrInternal,
	)

	handler.Update(w, req)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestUpdateHandlerSubscriptions_ContextCanceled(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	price := 999

	payload, _ := json.Marshal(dto.UpdateRequest{
		Price: &price,
	})

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add("id", id.String())

	ctx = context.WithValue(
		ctx,
		chi.RouteCtxKey,
		rctx,
	)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/subscriptions/"+id.String(),
		bytes.NewReader(payload),
	).WithContext(ctx)

	w := httptest.NewRecorder()

	service.On(
		"Update",
		mock.Anything,
		id,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		<-args.Get(0).(context.Context).Done()
	}).Return(
		apperrors.ErrCanceled,
	)

	handler.Update(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}

func TestUpdateHandlerSubscriptions_ContextDeadlineExceeded(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	id := uuid.New()

	price := 999

	payload, _ := json.Marshal(dto.UpdateRequest{
		Price: &price,
	})

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond,
	)
	defer cancel()

	rctx := chi.NewRouteContext()

	rctx.URLParams.Add("id", id.String())

	ctx = context.WithValue(
		ctx,
		chi.RouteCtxKey,
		rctx,
	)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/subscriptions/"+id.String(),
		bytes.NewReader(payload),
	).WithContext(ctx)

	w := httptest.NewRecorder()

	service.On(
		"Update",
		mock.Anything,
		id,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		<-args.Get(0).(context.Context).Done()
	}).Return(
		apperrors.ErrTimeout,
	)

	handler.Update(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}
