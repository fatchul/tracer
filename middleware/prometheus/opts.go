package prometheus

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type config struct {
	skipper        middleware.Skipper `ignored:"true"`
	getTraceIDFunc GetTraceIDFunc     `ignored:"true"`

	Env         string `envconfig:"ENV" default:"env"`
	ServiceName string `envconfig:"SERVICE" default:"evm-service"`
	MetricPath  string `envconfig:"PROMETHEUS_METRIC_PATH" default:"/metrics"`
}

type Option func(cfg *config)

type GetTraceIDFunc func(c echo.Context) string

func (o Option) Apply(cfg *config) {
	o(cfg)
}

func parseConfigFromEnv(cfg *config) {
	_ = envconfig.Process("", cfg)
}

func WithSkipper(skipper middleware.Skipper) Option {
	return func(cfg *config) {
		cfg.skipper = skipper
	}
}

func WithGetTraceIDFunc(f GetTraceIDFunc) Option {
	return func(cfg *config) {
		if f == nil {
			// set noop func to avoid crash
			f = func(c echo.Context) string {
				return ""
			}
		}
		cfg.getTraceIDFunc = f
	}
}

func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.ServiceName = serviceName
	}
}

func WithEnvironment(env string) Option {
	return func(cfg *config) {
		cfg.Env = env
	}
}

func WithMetricPath(path string) Option {
	return func(cfg *config) {
		cfg.MetricPath = path
	}
}
