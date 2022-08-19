package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	requestID = func() string {
		return "4129bc5f-a934-4eec-b4d6-894016896d3f"
	}
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("lorem"))
	}))
	handler.ServeHTTP(w, r)

	expected := w.Header().Get(requestIDKey)
	assert.Equal(t, expected, requestID())
}
