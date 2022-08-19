package middleware

import (
	"log"
	"net/http"
	"time"
)

var (
	DefaultLogger func(next http.Handler) http.Handler
)

type LogFormatter interface {
	Write(r *http.Request, w *Writer, requestTime time.Duration, requestID string)
}

type defaultLogger struct {
}

func (d *defaultLogger) Write(r *http.Request, w *Writer, elapsed time.Duration, requestID string) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	log.Printf("%s %s %s://%s%s - %d in %s", requestID, r.Method, scheme, r.Host, r.URL.Path, w.Code, elapsed)
}

func RequestLogger(log LogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			id := GetRequestID(r)
			ww := newWriter(w)

			defer func() {
				log.Write(r, ww, time.Since(t1), id)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func Logger(next http.Handler) http.Handler {
	return DefaultLogger(next)
}

func init() {
	l := &defaultLogger{}
	DefaultLogger = RequestLogger(l)
}
