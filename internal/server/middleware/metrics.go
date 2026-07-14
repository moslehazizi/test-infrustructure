package middleware

import (
	"control-panel-service/pkg/telemetry"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (w *metricsResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *metricsResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(b)
	w.size += int64(n)

	return n, err
}

func MetricsMiddleware(next http.Handler) http.Handler {
	meter := telemetry.GetMeter()

	requestCounter, _ := meter.Int64Counter("http.requests.total")
	requestDuration, _ := meter.Int64Histogram("http.request.duration_ms")
	requestSize, _ := meter.Int64Histogram("http.request.size_bytes")
	responseSize, _ := meter.Int64Histogram("http.response.size_bytes")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		metricRespWriter := &metricsResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		requestSizeBytes := r.ContentLength
		if requestSizeBytes < 0 {
			requestSizeBytes = 0
		}

		attrs := attribute.NewSet(
			attribute.String("method", r.Method),
			attribute.String("route", r.URL.Path),
		)

		requestSize.Record(r.Context(), requestSizeBytes, metric.WithAttributeSet(attrs))

		next.ServeHTTP(metricRespWriter, r)

		durationMs := time.Since(start).Milliseconds()

		route := r.URL.Path
		if rc := chi.RouteContext(r.Context()); rc != nil {
			if pattern := rc.RoutePattern(); pattern != "" {
				route = pattern
			}
		}

		attrs = attribute.NewSet(
			attribute.String("method", r.Method),
			attribute.String("route", route),
			attribute.String("status_code", strconv.Itoa(metricRespWriter.statusCode)),
		)

		requestCounter.Add(r.Context(), 1, metric.WithAttributeSet(attrs))
		requestDuration.Record(r.Context(), durationMs, metric.WithAttributeSet(attrs))
		responseSize.Record(r.Context(), metricRespWriter.size, metric.WithAttributeSet(attrs))
	})
}
