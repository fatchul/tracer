package opentelemetry

import (
	otelMetric "go.opentelemetry.io/otel/metric"
	otelTrace "go.opentelemetry.io/otel/trace"
)

var (
	cfg *config

	globalMeter  otelMetric.Meter
	globalTracer otelTrace.Tracer
)

const (
	tracerLibName = "MomologOpenTelemetryTracer"
	meterLibName  = "MomologOpenTelemetryMeter"

	HeaderTraceID = "X-Trace-ID"
)
