package main

import (
	"context"
	"fmt"
	"github.com/fatchul/tracer"
	"github.com/fatchul/tracer/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	hello := http.HandlerFunc(helloWorld)
	s := NewStructuredLog(context.Background(), momolog.New())
	mux.Handle("/", middleware.RequestLogger(s)(hello))

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Hello world\n")
}

type StructuredLog struct {
	log *momolog.Log
	ctx context.Context
}

func NewStructuredLog(ctx context.Context, log *momolog.Log) *StructuredLog {
	return &StructuredLog{
		log: log,
		ctx: ctx,
	}
}

func (l *StructuredLog) Write(r *http.Request, w *middleware.Writer, elapsed time.Duration, requestID string) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	l.log.Info(l.ctx).Msgf(momolog.HttpTextFormatter.String(), requestID, r.Method, scheme, r.Host, r.URL.Path, w.Code, elapsed)
}
