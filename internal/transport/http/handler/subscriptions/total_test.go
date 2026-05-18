package subscriptions_test

import (
	"net/http"
	"net/http/httptest"
	apperrors "subscriptions-service/internal/errors"
	"subscriptions-service/internal/transport/http/handler/subscriptions"
	"subscriptions-service/internal/transport/http/handler/subscriptions/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTotalHandlerSubscriptions_SuccessWithoutFilters(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total",
		nil,
	)

	rec := httptest.NewRecorder()

	service.On(
		"Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Return(999, nil)

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusOK,
		rec.Code,
	)

	service.AssertExpectations(t)
}

func TestTotalHandlerSubscriptions_SuccessWithFilters(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	userID := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total?user_id="+userID.String()+
			"&service_name=Netflix"+
			"&from=2025-01-01"+
			"&to=2025-12-31",
		nil,
	)

	rec := httptest.NewRecorder()

	serviceName := "Netflix"

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	to := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	service.On(
		"Total",
		mock.Anything,
		&userID,
		&serviceName,
		&from,
		&to,
	).Return(1999, nil)

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusOK,
		rec.Code,
	)

	service.AssertExpectations(t)
}

func TestTotalHandlerSubscriptions_InvalidUserID(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total?user_id=invalid",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		rec.Code,
	)

	service.AssertNotCalled(t, "Total")
}

func TestTotalHandlerSubscriptions_InvalidFromDate(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total?from=invalid",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		rec.Code,
	)

	service.AssertNotCalled(t, "Total")
}

func TestTotalHandlerSubscriptions_InvalidToDate(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total?to=invalid",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		rec.Code,
	)

	service.AssertNotCalled(t, "Total")
}

func TestTotalHandlerSubscriptions_ValidationError(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total",
		nil,
	)

	rec := httptest.NewRecorder()

	service.On(
		"Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Return(
		0,
		apperrors.WrapMessage(
			nil,
			apperrors.ErrValidation,
			"invalid date range",
		),
	)

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		rec.Code,
	)

	service.AssertExpectations(t)
}

func TestTotalHandlerSubscriptions_InternalError(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total",
		nil,
	)

	rec := httptest.NewRecorder()

	service.On(
		"Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Return(
		0,
		apperrors.ErrInternal,
	)

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		rec.Code,
	)

	service.AssertExpectations(t)
}

func TestTotalHandlerSubscriptions_Timeout(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total",
		nil,
	)

	rec := httptest.NewRecorder()

	service.On(
		"Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Return(
		0,
		apperrors.ErrTimeout,
	)

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		rec.Code,
	)

	service.AssertExpectations(t)
}

func TestTotalHandlerSubscriptions_Canceled(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total",
		nil,
	)

	rec := httptest.NewRecorder()

	service.On(
		"Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Return(
		0,
		apperrors.ErrCanceled,
	)

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		rec.Code,
	)

	service.AssertExpectations(t)
}

func TestTotalHandlerSubscriptions_ConflictError(t *testing.T) {

	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions/total",
		nil,
	)

	rec := httptest.NewRecorder()

	service.On(
		"Total",
		mock.Anything,
		(*uuid.UUID)(nil),
		(*string)(nil),
		(*time.Time)(nil),
		(*time.Time)(nil),
	).Return(
		0,
		apperrors.ErrConflict,
	)

	handler.Total(rec, req)

	assert.Equal(
		t,
		http.StatusConflict,
		rec.Code,
	)

	service.AssertExpectations(t)
}
