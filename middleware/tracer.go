package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Tracer create a new tracing span every request made it
// full documentation https://opentelemetry.io/docs/instrumentation/go/libraries/
func Tracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		otelhttp.NewHandler(next, r.URL.Path).ServeHTTP(rw, r)
	})
}
