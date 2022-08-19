package opentelemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	otelMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/nonrecording"
	otelTrace "go.opentelemetry.io/otel/trace"
)

var (
	syncInit sync.Once
)

func init() {
	setNoopTracerMeter()
}

func Initialize(opts ...Option) {
	syncInit.Do(func() {
		InitConfig(opts...)
		if IsActive() {
			InitTraceProvider(context.Background())
			InitMetricProvider(context.Background())
		}
	})

}

func IsActive() bool {
	if cfg == nil {
		return false
	}
	return cfg.IsActive
}

func ServiceName() string {
	if cfg == nil {
		return ""
	}
	return cfg.ServiceName
}

func TraceID(ctx context.Context) string {
	tid := otelTrace.SpanContextFromContext(ctx).TraceID()
	if tid.IsValid() {
		return tid.String()
	}
	return ""
}

func SetGlobalMeter(mp otelMetric.MeterProvider) {
	if mp == nil {
		mp = nonrecording.NewNoopMeterProvider()
	}
	global.SetMeterProvider(mp)
	globalMeter = global.MeterProvider().Meter(meterLibName)
}

func Meter() otelMetric.Meter {
	if global.MeterProvider() == nil {
		SetGlobalMeter(nonrecording.NewNoopMeterProvider())
	}
	return globalMeter
}

func SetGlobalTracer(tp otelTrace.TracerProvider) {
	if tp == nil {
		tp = otelTrace.NewNoopTracerProvider()
	}
	otel.SetTracerProvider(tp)
	globalTracer = otel.GetTracerProvider().Tracer(tracerLibName)
}

func Tracer() otelTrace.Tracer {
	if globalTracer == nil {
		SetGlobalTracer(otelTrace.NewNoopTracerProvider())
	}
	return globalTracer
}
