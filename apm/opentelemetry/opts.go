package opentelemetry

type (
	Option func()

	MDOption func() (key string, value interface{})
)

func (o Option) Apply() {
	o()
}

// Activated sets the OpenTelemetry status, if false, it won't send any trace/metric
func Activated(b bool) Option {
	return func() {
		cfg.IsActive = b
	}
}

// WithServiceName sets the service name
func WithServiceName(name string) Option {
	return func() {
		cfg.ServiceName = name
	}
}

// WithGRPCEndpoint sets the opentelemetry collector endpoint
func WithGRPCEndpoint(endpoint string) Option {
	return func() {
		cfg.GRPCEndpoint = endpoint
	}
}

// WithHTTPEndpoint sets the opentelemetry collector http endpoint
// still use GRPC endpoint if any
func WithHTTPEndpoint(endpoint string) Option {
	return func() {
		cfg.HTTPEndpoint = endpoint
	}
}

// WithSamplerPercentage sets the sampler percentage
// value <= 0 means Never sample
// value >= 100 means Always sample
// value 1 - 99 means Percentage sample value
func WithSamplerPercentage(percentage int) Option {
	return func() {
		cfg.SamplerPercentage = percentage
	}
}

// WithMiddlewareOpts sets the middleware specific options
func WithMiddlewareOpts(fs ...MDOption) Option {
	return func() {
		for _, f := range fs {
			k, v := f()
			cfg.MDOpts[k] = v
		}
	}
}
