package opentelemetry

import (
	"context"

	"github.com/labstack/echo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	otelTrace "go.opentelemetry.io/otel/trace"

	"github.com/fatchul/tracer/apm/opentelemetry"
	attrMomolog "github.com/fatchul/tracer/apm/opentelemetry/attribute"
)

func Middleware(opts ...opentelemetry.Option) echo.MiddlewareFunc {
	opentelemetry.Initialize(opts...)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if getSkipper()(c) || !opentelemetry.IsActive() {
				return next(c)
			}

			req := c.Request()
			ctx := req.Context()

			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))
			spanStartOpts := []otelTrace.SpanStartOption{
				otelTrace.WithAttributes(attrMomolog.Header(req)...),
				otelTrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", req)...),
				otelTrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(req)...),
				otelTrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(opentelemetry.ServiceName(), c.Path(), req)...),
				otelTrace.WithSpanKind(otelTrace.SpanKindServer),
			}

			var span otelTrace.Span
			ctx, span = opentelemetry.Tracer().Start(ctx, c.Path(), spanStartOpts...)
			defer span.End()

			// add trace ID to context and response header
			traceID := span.SpanContext().TraceID().String()
			c.Response().Header().Set(opentelemetry.HeaderTraceID, traceID)
			ctx = context.WithValue(ctx, opentelemetry.HeaderTraceID, traceID)

			c.SetRequest(req.WithContext(ctx))

			err := next(c)
			if err != nil {
				span.SetAttributes(attribute.String("echo.error", err.Error()))
			}

			attrs := semconv.HTTPAttributesFromHTTPStatusCode(c.Response().Status)
			spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(c.Response().Status)
			span.SetAttributes(attrs...)
			span.SetStatus(spanStatus, spanMessage)

			return err
		}
	}
}
