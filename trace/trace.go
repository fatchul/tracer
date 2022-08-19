package trace

import (
	"context"
	"log"

	"github.com/fatchul/tracer/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
	config   *Config
}

type Config struct {
	DisabledGlobalProvider bool
	ServiceName            string
	VersionNumber          string
}

func SetServiceName(name string) Apply {
	return func(c *Config) {
		c.ServiceName = name
	}
}

func SetVersionNumber(number string) Apply {
	return func(c *Config) {
		c.VersionNumber = number
	}
}

func SetDisabledGlobalProvider(disabled bool) Apply {
	return func(c *Config) {
		c.DisabledGlobalProvider = true
	}
}

func defaultConfig() *Config {
	return &Config{
		DisabledGlobalProvider: false,
		ServiceName:            "github.com/fatchul/tracer/exporter",
		VersionNumber:          "0.0.1",
	}
}

type Apply func(*Config)

func New(applies ...Apply) *Trace {
	cfg := defaultConfig()
	for _, apply := range applies {
		apply(cfg)
	}

	return &Trace{
		config: cfg,
	}
}

func (e *Trace) createResources() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(e.config.ServiceName),
		semconv.ServiceVersionKey.String(e.config.VersionNumber),
	)
}

// TraceProvider: create trace provider and given the exporter
func (e *Trace) TraceProvider(batcher Batcher) *Trace {
	e.provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(batcher.Export()),
		sdktrace.WithResource(e.createResources()),
	)

	e.tracer = e.provider.Tracer(e.config.ServiceName)

	if !e.config.DisabledGlobalProvider {
		otel.SetTracerProvider(e.provider)
	}

	return e
}

// Shutdown: handle shutdown properly so nothing leaks.
func (e *Trace) Shutdown(ctx context.Context) func() {
	return func() {
		if err := e.provider.Shutdown(ctx); err != nil {
			log.Fatalf("stopping tracer provider: %v", err)
		}
	}
}

// Start: creating new span
func (e *Trace) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := e.tracer.Start(ctx, spanName, opts...)
	span.SetAttributes(attribute.KeyValue{
		Key:   "RequestID",
		Value: attribute.StringValue(internal.RequestId(ctx)),
	})

	return ctx, span
}

// Batcher: purpose for attach the exporter in trace provider
type Batcher interface {
	Export() sdktrace.SpanExporter
}
