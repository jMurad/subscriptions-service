package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"subscriptions-service/internal/handler"
	"subscriptions-service/internal/service"
	"subscriptions-service/internal/service/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateSubscription_InvalidJSON(t *testing.T) {
	mockRepo := new(mocks.Repository)

	svc := service.NewSubscriptionService(mockRepo)

	h := handler.NewSubscriptionHandler(svc)

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

	svc := service.NewSubscriptionService(mockRepo)

	h := handler.NewSubscriptionHandler(svc)

	body := `{
		"service_name":"Yandex Plus",
		"price":500,
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

	svc := service.NewSubscriptionService(mockRepo)

	h := handler.NewSubscriptionHandler(svc)

	body := `{
		"service_name":"Yandex Plus",
		"price":500,
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
