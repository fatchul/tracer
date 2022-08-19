package opentelemetry

import (
	"errors"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type (
	mdOpts map[string]interface{}

	config struct {
		ServiceName       string `envconfig:"SERVICE" default:""`
		IsActive          bool   `envconfig:"OTEL_ACTIVE" default:"true"`
		GRPCEndpoint      string `envconfig:"OTEL_COLLECTOR_GRPC_ENDPOINT" default:""`
		HTTPEndpoint      string `envconfig:"OTEL_COLLECTOR_HTTP_ENDPOINT" default:""`
		SamplerPercentage int    `envconfig:"OTEL_SAMPLER_PERCENTAGE" default:"30"`
		MDOpts            mdOpts `ignored:"true"`
	}
)

func InitConfig(opts ...Option) {
	InitConfigFromEnv()

	for _, opt := range opts {
		opt.Apply()
	}

	MustValidateConfig()
}

func InitConfigFromEnv() {
	cfg = &config{
		MDOpts: make(mdOpts),
	}

	_ = envconfig.Process("", cfg)
}

func GetMDOpts(key string) interface{} {
	return cfg.MDOpts[key]
}

func MustValidateConfig() {
	handleErr(ValidateConfig(), "some config is missing")
}

func ValidateConfig() error {
	if !IsActive() {
		log.Println("OpenTelemetry is disabled")
		return nil
	}

	if cfg.ServiceName == "" {
		return errors.New("service name is empty")
	}

	if cfg.GRPCEndpoint == "" && cfg.HTTPEndpoint == "" {
		return errors.New("collector endpoint is empty, you must set one of HTTP or GRPC endpoint")
	}

	return nil
}
