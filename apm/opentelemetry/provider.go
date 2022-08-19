package opentelemetry

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func getHTTPEndpoint() (endpoint string, secure bool) {
	secure = strings.HasPrefix(endpoint, "https://")
	if secure {
		endpoint = strings.TrimLeft(cfg.HTTPEndpoint, "https://")
	} else {
		endpoint = strings.TrimLeft(cfg.HTTPEndpoint, "http://")
	}

	return
}

func InitMetricProvider(ctx context.Context) {
	endpoint, secure := getHTTPEndpoint()

	httpOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint),
	}

	if !secure {
		httpOpts = append(httpOpts, otlpmetrichttp.WithInsecure())
	}

	metricClient := otlpmetrichttp.NewClient(httpOpts...)

	// use GRPC if available
	if cfg.GRPCEndpoint != "" {
		metricClient = otlpmetricgrpc.NewClient(
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithEndpoint(cfg.GRPCEndpoint),
		)
	}

	metricExp, err := otlpmetric.New(ctx, metricClient)
	handleErr(err, "Failed to create the metric collector exporter")

	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			metricExp,
		),
		controller.WithExporter(metricExp),
		controller.WithCollectPeriod(2*time.Second),
	)
	err = pusher.Start(ctx)
	handleErr(err, "Failed to start metric pusher")

	SetGlobalMeter(pusher)
}

func InitTraceProvider(ctx context.Context) {
	endpoint, secure := getHTTPEndpoint()

	httpOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}

	if !secure {
		httpOpts = append(httpOpts, otlptracehttp.WithInsecure())
	}

	traceClient := otlptracehttp.NewClient(httpOpts...)

	if cfg.GRPCEndpoint != "" {
		traceClient = otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(cfg.GRPCEndpoint),
		)
	}

	traceExp, err := otlptrace.New(ctx, traceClient)
	handleErr(err, "Failed to create the trace collector exporter")

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	handleErr(err, "failed to create trace resource")

	bsp := sdkTrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(getSampler(cfg.SamplerPercentage)),
		sdkTrace.WithResource(res),
		sdkTrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	SetGlobalTracer(tracerProvider)
}
