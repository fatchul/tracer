package opentelemetry

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/nonrecording"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	otelTrace "go.opentelemetry.io/otel/trace"
)

func setNoopTracerMeter() {
	global.SetMeterProvider(nonrecording.NewNoopMeterProvider())
	otel.SetTracerProvider(otelTrace.NewNoopTracerProvider())
	globalMeter = global.MeterProvider().Meter(meterLibName)
	globalTracer = otel.GetTracerProvider().Tracer(tracerLibName)
}

func handleErr(err error, message string) {
	if err != nil {
		log.Printf("%s: %v\n", message, err)
		cfg.IsActive = false
		setNoopTracerMeter()
	}
}

func getSampler(percentage int) sdkTrace.Sampler {
	switch {
	case percentage >= 100:
		return sdkTrace.AlwaysSample()

	case percentage <= 0:
		return sdkTrace.NeverSample()

	default:
		traceRatio := float64(percentage) / 100
		return sdkTrace.ParentBased(sdkTrace.TraceIDRatioBased(traceRatio))
	}
}
