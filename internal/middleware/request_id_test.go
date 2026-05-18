package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"subscriptions-service/internal/middleware"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID_HeaderSet(t *testing.T) {
	handler := middleware.RequestID(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	requestID := w.Header().Get(
		"X-Request-ID",
	)

	assert.NotEmpty(t, requestID)
}
