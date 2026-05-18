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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Success test
func TestHandler_Create_Success(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	userID := uuid.New()
	subID := uuid.New()

	payload, _ := json.Marshal(dto.CreateRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      userID.String(),
		StartDate:   "2025-01-01",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewReader(payload),
	)

	w := httptest.NewRecorder()

	expected := model.Subscription{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      userID,
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	service.On(
		"Create",
		mock.Anything,
		expected,
	).Return(subID, nil)

	handler.Create(w, req)

	assert.Equal(
		t,
		http.StatusCreated,
		w.Code,
	)

	service.AssertExpectations(t)
}

// Invalid JSON
func TestHandler_Create_InvalidJSON(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewBufferString("{invalid"),
	)

	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "Create")
}

// Invalid UUID
func TestHandler_Create_InvalidUUID(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	payload, _ := json.Marshal(dto.CreateRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      "invalid",
		StartDate:   "2025-01-01",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewReader(payload),
	)

	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "Create")
}

// Invalid start_date
func TestHandler_Create_InvalidStartDate(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	payload, _ := json.Marshal(dto.CreateRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New().String(),
		StartDate:   "invalid",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewReader(payload),
	)

	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)

	service.AssertNotCalled(t, "Create")
}

// Service error propagation
func TestHandler_Create_ServiceError(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	payload, _ := json.Marshal(dto.CreateRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New().String(),
		StartDate:   "2025-01-01",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewReader(payload),
	)

	w := httptest.NewRecorder()

	service.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Return(
		uuid.Nil,
		apperrors.ErrConflict,
	)

	handler.Create(w, req)

	assert.Equal(
		t,
		http.StatusConflict,
		w.Code,
	)

	service.AssertExpectations(t)
}

// Canceled propagation
func TestHandler_Create_ContextCanceled(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	payload, _ := json.Marshal(dto.CreateRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New().String(),
		StartDate:   "2025-01-01",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewReader(payload),
	)

	ctx, cancel := context.WithCancel(
		req.Context(),
	)

	cancel()

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	service.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Return(
		uuid.Nil,
		apperrors.ErrCanceled,
	)

	handler.Create(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}

// Timeout propagation
func TestHandler_Create_ContextDeadlineExceeded(t *testing.T) {
	service := new(mocks.Service)

	handler := subscriptions.NewHandler(service)

	payload, _ := json.Marshal(dto.CreateRequest{
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New().String(),
		StartDate:   "2025-01-01",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewReader(payload),
	)

	ctx, cancel := context.WithTimeout(
		req.Context(),
		time.Millisecond,
	)
	defer cancel()

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	service.On(
		"Create",
		mock.Anything,
		mock.Anything,
	).Return(
		uuid.Nil,
		apperrors.ErrTimeout,
	)

	handler.Create(w, req)

	assert.Equal(
		t,
		http.StatusRequestTimeout,
		w.Code,
	)

	service.AssertExpectations(t)
}
