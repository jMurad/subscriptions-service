package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"subscriptions-service/internal/handler"
	"subscriptions-service/internal/service/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSubscription_InvalidJSON(t *testing.T) {
	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	body := `invalid-json`

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewBufferString(body),
	)

	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)
}

func TestCreateSubscription_InvalidUUID(t *testing.T) {
	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	body := `{
		"service_name":"Yandex Plus",
		"price":400,
		"user_id":"invalid",
		"start_date":"07-2025"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewBufferString(body),
	)

	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)
}

func TestCreateSubscription_InvalidDate(t *testing.T) {
	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	body := `{
		"service_name":"Yandex Plus",
		"price":400,
		"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba",
		"start_date":"invalid-date"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewBufferString(body),
	)

	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(
		t,
		http.StatusBadRequest,
		w.Code,
	)
}

func TestCreateSubscription_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	expectedID := uuid.New()

	mockRepo.
		On("Create", mock.Anything, mock.Anything).
		Return(expectedID, nil)

	body := `{
		"service_name":"Yandex Plus",
		"price":400,
		"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba",
		"start_date":"07-2025"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/subscriptions",
		bytes.NewBufferString(body),
	)

	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockRepo.AssertExpectations(t)
}
