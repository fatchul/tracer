# Prometheus Middleware

Echo middleware to automatically expose some metrics to `/metrics`:

Exposed metrics:

- Total incoming request, partitioned by HTTP code
- Request latency per endpoint

## How to Use

Minimal use:

```go
func main() {
    e := echo.New()

    ...

    // Optional
    // For the best use, use and place OpenTelemetry middleware before Prometheus middleware
    e.Use(opentelemetry.Middleware())

    e.Use(prometheus.Middleware(e))

    ...
}
```

You can also configure this middleware. See [Configuration](#configuration)

## TLDR;

```go
func main() {
    e := echo.New()

    ...

    // Optional
    // For the best use, use and place OpenTelemetry middleware before Prometheus middleware
    e.Use(opentelemetry.Middleware())

    // Default skipper func: middleware.DefaultSkipper (from echo middleware)
    skipperFunc := func (c echo.Context) bool {
        // complete this func
        return false
    }

    // Default trace ID func: prometheus.DefaultGetTraceIDFunc
    tracerIDFunc := func (c echo.Context) string {
        // complete this func
        return "example"
    }
    
    // Default value: `SERVICE` env var
    svcName := "example-service"

    // Default value: `ENV` env var
    env := "example-env"
	
    // Default value: `PROMETHEUS_METRIC_PATH` env var. If not set, will fallback to `/metrics`
    path := "/prom-metrics"
    
    e.Use(prometheus.Middleware(e,
        prometheus.WithSkipper(skipperFunc),
        prometheus.WithGetTraceIDFunc(tracerIDFunc),
        prometheus.WithServiceName(svcName),
		prometheus.WithEnvironment(env),
        prometheus.WithMetricPath(path),
    ))
	
    ...
}
```

## Configuration

We can control some config with env variable or by code. This middleware first will get the
configuration set by code, if not exist, then it will get from environment variable.

### Environment Variable Config

- `SERVICE`. For set the metrics service name param. Default: `evm-service`
- `ENV`. For set the metrics env param. Default: `env`
- `PROMETHEUS_METRIC_PATH`. Path to expose metrics. Default: `/metrics`

### Code Config

You can override default config (or env var config) using code config. You just need to pass
the options below after the first param when initializing the middleware.

- **Skipper**

  Set the skipper function.

  **Default**: echo middleware default skipper.

  **How to override**:
  ```go
  skipperFunc := func(c echo.Context) bool {
    // complete this func
    return false
  }
  
  e.Use(prometheus.Middleware(e, prometheus.WithSkipper(skipperFunc)))
  ```

- **Get Trace ID Func**

  For [Exemplars](https://grafana.com/docs/grafana/latest/basics/exemplars/) feature,
  this middleware need to get the trace ID from the tracer used in the service.
  
  **Default**: Get the tracer ID from the response header with key `X-Tracer-ID`. This header
  sets by OpenTelemetry middleware.

  **How to override**:
  ```go
  tracerIDFunc := func(c echo.Context) string {
    // complete this func
    return "example"
  }
  
  e.Use(prometheus.Middleware(e, prometheus.WithGetTraceIDFunc(tracerIDFunc)))
  ```

- **Service Name**

  This middleware will differentiate each metric by service name, so the metric will have 
  service name param.

  **Default**: use `SERVICE` environment variable

  **How to override**:
  ```go
  svcName := "example-service"
  e.Use(prometheus.Middleware(e, prometheus.WithServiceName(svcName)))
  ```

- **Environment**

  This middleware will differentiate each metric by service name, so the metric will have
  env param.

  **Default**: use `ENV` environment variable

  **How to override**:
  ```go
  env := "example-env"
  e.Use(prometheus.Middleware(e, prometheus.WithEnvironment(env)))
  ```

- **Metric Path**

  This middleware will automatically create the metric endpoint. This endpoint will contain all
  metrics.

  **Default**: use `PROMETHEUS_METRIC_PATH` environment variable

  **How to override**:
  ```go
  path := "/prom-metrics"
  e.Use(prometheus.Middleware(e, prometheus.WithMetricPath(path)))
  ```