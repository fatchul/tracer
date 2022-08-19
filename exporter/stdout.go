package exporter

import (
	"log"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Stdout struct {
}

func NewStdout() *Stdout {
	return &Stdout{}
}

func (stdout *Stdout) Export() sdktrace.SpanExporter {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalln(err)
	}

	return exporter
}
