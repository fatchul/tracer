package opentelemetry

import (
	"github.com/labstack/echo/middleware"

	"github.com/fatchul/tracer/apm/opentelemetry"
)

const (
	mdKeySkipper = "echo.skipper"
)

// WithSkipper sets the skipper to skip the middleware
func WithSkipper(skipper middleware.Skipper) opentelemetry.MDOption {
	return func() (key string, value interface{}) {
		return mdKeySkipper, skipper
	}
}

func getSkipper() middleware.Skipper {
	o := opentelemetry.GetMDOpts(mdKeySkipper)
	if o == nil {
		return middleware.DefaultSkipper
	}

	skipper, ok := o.(middleware.Skipper)
	if !ok {
		return middleware.DefaultSkipper
	}

	return skipper
}
