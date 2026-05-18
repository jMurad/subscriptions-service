package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"subscriptions-service/internal/middleware"

	"github.com/stretchr/testify/assert"
)

func TestTimeout_ContextCanceled(t *testing.T) {
	done := make(chan error, 1)
	handler := middleware.Timeout(10 * time.Millisecond)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-r.Context().Done()
			done <- r.Context().Err()
		}),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	err := <-done

	assert.ErrorIs(
		t,
		err,
		context.DeadlineExceeded,
	)
}
