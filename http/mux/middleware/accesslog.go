package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/zbum/mantyboot/http/mux"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	length     int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK, 0}
}

func (w *LoggingResponseWriter) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.length += n
	return
}

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *LoggingResponseWriter) Length() int {
	return w.length
}

func (w *LoggingResponseWriter) StatusCode() int {
	return w.statusCode
}

func AccessLogger(logger *log.Logger) mux.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			nw := NewLoggingResponseWriter(w)
			start := uint64(time.Now().UnixMicro())
			defer logAccess(logger, nw, r, start)
			next(nw, r)
		}
	}
}

func logAccess(logger *log.Logger, w *LoggingResponseWriter, r *http.Request, start uint64) {
	logger.Printf("%dμs ua:%10.10s %s %s %d %d\n", uint64(time.Now().UnixMicro())-start, r.UserAgent(), r.Method, r.URL.Path, w.StatusCode(), w.Length())
}
