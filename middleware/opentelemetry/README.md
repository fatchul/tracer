# OpenTelemetry Middleware

Echo middleware to trace incoming requests and metrics using OpenTelemetry and forward it
to the collector

## How to Use

Minimal use:

```go
func main() {
    e := echo.New()

    ...

    e.Use(opentelemetry.Middleware())

    ...
}
```

You can also configure this middleware. See [Configuration](#configuration)

## TLDR;

```go
import (
    otelAPM "github.com/fatchul/tracer/apm/opentelemetry"
    "github.com/fatchul/tracer/middleware/opentelemetry"
)

func main() {
    e := echo.New()

    ...

    // Default value: true
    active := true

    // Default value: `SERVICE` env var
    serviceName := "example-service"

    // Default value: `OTEL_COLLECTOR_GRPC_ENDPOINT` env var
    grpcEndpoint := "localhost:4137"

    // Default value: `OTEL_COLLECTOR_HTTP_ENDPOINT` env var
    httpEndpoint := "localhost:4138"

    // Default value: `OTEL_SAMPLER_PERCENTAGE` env var
    samplerPercentage := 100

    // Default skipper func: middleware.DefaultSkipper (from echo middleware)
    skipperFunc := func (c echo.Context) bool {
        // complete this func
        return false
    }
	
    e.Use(opentelemetry.Middleware(
        otelAPM.Activated(active),
        otelAPM.WithServiceName(serviceName),
        otelAPM.WithGRPCEndpoint(grpcEndpoint),
        otelAPM.WithHTTPEndpoint(httpEndpoint),
        otelAPM.WithSamplerPercentage(samplerPercentage),
        otelAPM.WithMiddlewareOpts(
            opentelemetry.WithSkipper(skipperFunc),
        ),
    ))
	
    ...
}
```

## Configuration

We can control some config with env variable or by code. This middleware first will get the
configuration set by code, if not exist, then it will get from environment variable.

### Environment Variable Config

- `SERVICE`. For set the parent span name. If not set, this middleware is mark as inactive
- `OTEL_ACTIVE`. Mark this endpoint to active. **Default**: `true`
- `OTEL_COLLECTOR_GRPC_ENDPOINT`. GRPC Endpoint where to send traces and metrics.
- `OTEL_COLLECTOR_HTTP_ENDPOINT`. HTTP Endpoint where to send traces and metrics. Will use
  GRPC if both HTTP and GRPC endpoint is set
- `OTEL_SAMPLER_PERCENTAGE`. Sampler percentage. `100` means trace all the incoming request.
  `0` means trace nothing (but the middleware still active and still generate trace ID for 
  each incoming request). **Default**: `30`

### Code Config

You can override default config (or env var config) using code config. You just need to pass
the options below as param when initializing the middleware.

**Note:**

Use this import to differentiate `opentelemetry` package from middleware and apm

```go
import (
    otelAPM "github.com/fatchul/tracer/apm/opentelemetry"
    "github.com/fatchul/tracer/middleware/opentelemetry"
)
```

- **Activate**

  Set the middleware active or not.

  **Default**: `true`

  **How to override**:
  ```go
  e.Use(opentelemetry.Middleware(otelAPM.Activated(false)))
  ```

- **Service Name**

  Set the service name.

  **Default**: use `SERVICE` environment variable

  **How to override**:
  ```go
  serviceName := "example-service"
  e.Use(opentelemetry.Middleware(otelAPM.WithServiceName(serviceName)))
  ```

- **GRPC Endpoint**

  Set the collector GRPC endpoint.
  
  **Default**: use `OTEL_COLLECTOR_GRPC_ENDPOINT` environment variable

  **How to override**:
  ```go
  grpcEndpoint := "localhost:4137"
  e.Use(opentelemetry.Middleware(otelAPM.WithGRPCEndpoint(grpcEndpoint)))
  ```

- **HTTP Endpoint**

  Set the collector HTTP endpoint

  **Default**: use `OTEL_COLLECTOR_HTTP_ENDPOINT` environment variable

  **How to override**:
  ```go
  httpEndpoint := "localhost:4138"
  e.Use(opentelemetry.Middleware(otelAPM.WithHTTPEndpoint(httpEndpoint)))
  ```

- **Sampler Percentage**

  Set the sampler percentage

  **Default**: use `OTEL_SAMPLER_PERCENTAGE` environment variable

  **How to override**:
  ```go
  samplerPercentage := 100
  e.Use(opentelemetry.Middleware(otelAPM.WithSamplerPercentage(samplerPercentage)))
  ```

- **Skipper**

  Set the skipper function.

  **Default**: echo middleware default skipper.

  **How to override**:
  ```go
  skipperFunc := func(c echo.Context) bool {
    // complete this func
    return false
  }
  
  e.Use(opentelemetry.Middleware(
    otelAPM.WithMiddlewareOpts(
        opentelemetry.WithSkipper(skipperFunc),
    ),
  )
  ```