package middleware

import "net/http"

// Writer Capture ResponseWriter data for log after request done
type Writer struct {
	http.ResponseWriter
	Code int
}

func newWriter(w http.ResponseWriter) *Writer {
	return &Writer{ResponseWriter: w}
}

func (w *Writer) WriteHeader(code int) {
	w.Code = code
	w.ResponseWriter.WriteHeader(code)
}
