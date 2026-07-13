package middleware

import (
	"control-panel-service/pkg/telemetry"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func MetricsfMiddleware() fiber.Handler {
	meter := telemetry.GetMeter()

	requestCounter, _ := meter.Int64Counter("http.requests.total")
	requestDuration, _ := meter.Int64Histogram("http.request.duration_ms")
	requestSize, _ := meter.Int64Histogram("http.request.size_bytes")
	responseSize, _ := meter.Int64Histogram("http.response.size_bytes")

	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestSizeBytes := int64(len(c.Request().Header.String()) + len(c.Body()))
		attrs := attribute.NewSet(
			attribute.String("method", c.Method()),
			attribute.String("route", c.Route().Path),
		)
		requestSize.Record(c.Context(), requestSizeBytes, metric.WithAttributeSet(attrs))

		err := c.Next()

		duration := time.Since(start)
		durationMs := duration.Milliseconds()

		statusCode := c.Response().StatusCode()
		attrs = attribute.NewSet(
			attribute.String("method", c.Method()),
			attribute.String("route", c.Route().Path),
			attribute.String("status_code", strconv.Itoa(statusCode)),
		)

		requestCounter.Add(c.Context(), 1, metric.WithAttributeSet(attrs))
		requestDuration.Record(c.Context(), durationMs, metric.WithAttributeSet(attrs))

		responseSizeBytes := int64(len(c.Response().Header.String()) + len(c.Response().Body()))
		responseSize.Record(c.Context(), responseSizeBytes, metric.WithAttributeSet(attrs))

		// nolint
		return err
	}
}
