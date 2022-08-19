package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

var (
	requestIDKey = "x-request-id"

	requestID = func() string {
		return uuid.NewString()
	}
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r)
		ctx := context.WithValue(r.Context(), requestIDKey, id)

		w.Header().Set(requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(r *http.Request) string {
	id := r.Header.Get(requestIDKey)
	if id != "" {
		return id
	}

	if id, ok := r.Context().Value(requestIDKey).(string); ok {
		return id
	}

	return requestID()
}
