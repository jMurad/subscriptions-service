package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"subscriptions-service/internal/handler"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Create tests
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

// GetByID tests
func TestGetByID_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	id := uuid.New()

	sub := &model.Subscription{
		ID:          id,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   time.Now(),
	}

	mockRepo.
		On("GetByID", mock.Anything, id).
		Return(sub, nil)

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

	h.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestGetByID_InvalidID(t *testing.T) {
	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	req := httptest.NewRequest(
		http.MethodGet,
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

	h.GetByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByID_NotFound(t *testing.T) {

	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	id := uuid.New()

	mockRepo.
		On("GetByID", mock.Anything, id).
		Return(nil, errors.New("not found"))

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

	h.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// List tests
func TestListSubscriptions_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	subscriptions := []model.Subscription{
		{
			ID:          uuid.New(),
			ServiceName: "Yandex Plus",
			Price:       400,
		},
	}

	mockRepo.
		On("List", mock.Anything, 10, 0).
		Return(subscriptions, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions",
		nil,
	)

	w := httptest.NewRecorder()

	h.List(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.Subscription

	err := json.NewDecoder(w.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Len(t, response, 1)
}

func TestListSubscriptions_Pagination(t *testing.T) {

	mockRepo := new(mocks.Repository)

	h := handler.NewSubscriptionHandler(mockRepo)

	mockRepo.
		On("List", mock.Anything, 5, 10).
		Return([]model.Subscription{}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/subscriptions?limit=5&offset=10",
		nil,
	)

	w := httptest.NewRecorder()

	h.List(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
