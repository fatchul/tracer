package momolog

import (
	"context"
	"runtime"

	"github.com/fatchul/tracer/internal"
	"github.com/rs/zerolog"
)

type writer struct {
	event *zerolog.Event
	ctx   context.Context
}

func newWriter(ctx context.Context, event *zerolog.Event) *writer {
	return &writer{
		event: event,
		ctx:   ctx,
	}
}

// Field: add additional parameter in the log print
func (w *writer) Field(key string, value interface{}) *writer {
	w.event.Interface(key, value)
	return w
}

// Err: print error
func (w *writer) Err(err error) *writer {
	w.event.Err(err)
	return w
}

// Stack: enable stack trace printing for the error purpose
func (w *writer) Stack() *writer {
	w.event.Stack()
	return w
}

// Layer: purpose print the log for the domain layer
func (w *writer) Layer(msg string) {
	w.event.Msgf(LayerTextFormatter.String(), internal.RequestId(w.ctx), FuncName(), msg)
}

// Msg: print the log with message
func (w *writer) Msg(msg string) {
	w.event.Msg(msg)
}

// Msg: print the log message with format
func (w *writer) Msgf(format string, v ...interface{}) {
	w.event.Msgf(format, v...)
}

// Send: print the log without an empty msg. Equivalent with Msg("")
func (w *writer) Send() {
	w.event.Send()
}

// FuncName: given the function name where the function is called
// Example use this function in the function name called Stack()
// result should be [package_name.Stack]
func FuncName() string {
	pc, _, _, _ := runtime.Caller(2)

	return runtime.FuncForPC(pc).Name()
}
