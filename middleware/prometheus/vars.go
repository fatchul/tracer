package prometheus

import (
	"github.com/labstack/echo"
)

// GetTraceIDContext will take Trace ID value from the echo context
// the value sets by OpenTelemetry middleware
// this func is the default func to get the Trace ID
var GetTraceIDContext = func(c echo.Context) string {
	tID, _ := c.Request().Context().Value("X-Trace-ID").(string)
	return tID
}

// GetTraceIDHeader will take Trace ID value from the response header
// the value sets by OpenTelemetry middleware
var GetTraceIDHeader = func(c echo.Context) string {
	return c.Response().Header().Get("X-Trace-ID")
}
