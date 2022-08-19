package middleware

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	requestID = func() string {
		return "4129bc5f-a934-4eec-b4d6-894016896d3f"
	}
	type args struct {
		httpMock http.HandlerFunc
		method   string
		url      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "log GET Request",
			args: args{
				httpMock: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("lorem"))
				}),
				method: "GET",
				url:    "http://localhost:1234/",
			},
			want: fmt.Sprintf("%s GET http://localhost:1234/ - 200 in", requestID()),
		},
		{
			name: "log POST Request",
			args: args{
				httpMock: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte("lorem"))
				}),
				method: "POST",
				url:    "http://localhost:1234/",
			},
			want: fmt.Sprintf("%s POST http://localhost:1234/ - 400 in", requestID()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock log
			var buf bytes.Buffer
			log.SetOutput(&buf)

			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()

			handler := Logger(tt.args.httpMock)
			handler.ServeHTTP(w, r)
			assert.Contains(t, buf.String(), tt.want)
		})
	}
}
