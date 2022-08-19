package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fatchul/tracer/exporter"
	"github.com/fatchul/tracer/middleware"
	"github.com/fatchul/tracer/trace"
)

func main() {
	// ctx := context.Background()
	tr := trace.New(
		trace.SetServiceName("github.com/momolog/tracertest"),
		trace.SetVersionNumber("1.0.0"),
	).TraceProvider(exporter.NewStdout())
	defer tr.Shutdown(context.Background())

	svc := NewSvc(tr)

	mux := http.NewServeMux()
	mux.Handle("/hello-world", middleware.Tracer(http.HandlerFunc(svc.HelloWorld)))
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalln(err)
	}
}

type Svc struct {
	Tracer *trace.Trace
}

func NewSvc(tracer *trace.Trace) *Svc {
	return &Svc{
		Tracer: tracer,
	}
}

func (s *Svc) HelloWorld(w http.ResponseWriter, r *http.Request) {
	ctx, span := s.Tracer.Start(r.Context(), "hello-world")
	defer span.End()

	span.AddEvent("Start Multiply")
	result := s.Multiply(ctx, 10, 5)

	log.Println(result)
}

func (s *Svc) Multiply(ctx context.Context, x, y int) int {
	_, span := s.Tracer.Start(ctx, "multiple")
	defer span.End()

	return x + y
}
