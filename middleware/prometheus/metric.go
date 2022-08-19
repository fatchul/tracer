package prometheus

import (
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

type metricType struct {
	cfg             *config
	reqTotalCounter *prometheus.CounterVec
	reqDurHistogram *prometheus.HistogramVec
}

func initMetric(cfg *config) *metricType {
	reqTotalCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total",
			Help: "How many HTTP requests processed. Partitioned by status code and HTTP method",
		},
		[]string{"service", "env", "path", "code", "method"},
	)

	reqDurHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "The HTTP request latencies in seconds",
		},
		[]string{"service", "env", "path", "code", "method"},
	)

	collectors := []prometheus.Collector{
		reqTotalCounter,
		reqDurHistogram,
	}

	for _, collector := range collectors {
		if err := prometheus.Register(collector); err != nil {
			log.Println("register prometheus metric error.", err)
		}
	}

	return &metricType{
		cfg:             cfg,
		reqTotalCounter: reqTotalCounter,
		reqDurHistogram: reqDurHistogram,
	}
}

func (m *metricType) IncRequestTotal(c echo.Context) {
	if m == nil || m.reqTotalCounter == nil {
		return
	}

	ob, err := m.reqTotalCounter.GetMetricWithLabelValues(
		m.cfg.ServiceName,
		m.cfg.Env,
		c.Path(),
		strconv.Itoa(c.Response().Status),
		c.Request().Method,
	)
	if err != nil {
		return
	}

	ob.Inc()
}

func (m *metricType) RecordRequestDuration(c echo.Context, dur time.Duration) {
	if m == nil || m.reqDurHistogram == nil {
		return
	}

	ob, err := m.reqDurHistogram.GetMetricWithLabelValues(
		m.cfg.ServiceName,
		m.cfg.Env,
		c.Path(),
		strconv.Itoa(c.Response().Status),
		c.Request().Method,
	)
	if err != nil {
		return
	}

	durSec := float64(dur) / float64(time.Second)

	traceID := m.cfg.getTraceIDFunc(c)
	eo, ok := ob.(prometheus.ExemplarObserver)
	if traceID == "" || !ok {
		ob.Observe(durSec)
		return
	}

	eo.ObserveWithExemplar(durSec, prometheus.Labels{
		"traceID": traceID,
	})
}
