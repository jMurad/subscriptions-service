package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"subscriptions-service/internal/middleware"

	"github.com/stretchr/testify/assert"
)

func TestRecovery_PanicRecovered(t *testing.T) {
	handler := middleware.Recovery(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		w.Code,
	)
}
