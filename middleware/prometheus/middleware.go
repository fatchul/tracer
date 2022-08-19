package prometheus

import (
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prometheusHandler() echo.HandlerFunc {
	h := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		),
	)
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func Middleware(e *echo.Echo, opts ...Option) echo.MiddlewareFunc {
	cfg := &config{
		skipper:        middleware.DefaultSkipper,
		getTraceIDFunc: GetTraceIDContext,
	}

	parseConfigFromEnv(cfg)

	for _, opt := range opts {
		opt.Apply(cfg)
	}

	metric := initMetric(cfg)

	// register metric path to server
	e.GET(cfg.MetricPath, prometheusHandler())

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == cfg.MetricPath {
				return next(c)
			}

			if cfg.skipper(c) {
				return next(c)
			}

			start := time.Now()

			err := next(c)

			metric.IncRequestTotal(c)
			metric.RecordRequestDuration(c, time.Since(start))

			return err
		}
	}
}
