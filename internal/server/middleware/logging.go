package middleware

import (
	"control-panel-service/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	httpStatusServerError = http.StatusInternalServerError // 500
	httpStatusClientError = http.StatusBadRequest          // 400
)

// LoggingMiddleware creates a middleware that logs HTTP requests and responses.
// It generates a request ID, adds it to context, and logs structured information.
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Generate request ID if not present in headers
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}

		// Add request ID to context for use in handlers and downstream services
		ctx := logger.WithRequestID(c.Context(), requestID)
		c.SetUserContext(ctx)

		// Extract request information
		method := c.Method()
		path := c.Path()
		originalPath := c.OriginalURL()
		clientIP := c.IP()
		userAgent := c.Get("User-Agent")

		// Log incoming request
		zap.L().Info("incoming HTTP request",
			zap.String(logger.FieldRequestID, requestID),
			zap.String(logger.FieldMethod, method),
			zap.String(logger.FieldPath, path),
			zap.String("original_path", originalPath),
			zap.String(logger.FieldClientIP, clientIP),
			zap.String("user_agent", userAgent),
		)

		// Process request
		errNext := c.Next()

		// Calculate duration
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()
		responseSize := len(c.Response().Body())

		// Prepare log fields
		fields := []zap.Field{
			zap.String(logger.FieldRequestID, requestID),
			zap.String(logger.FieldMethod, method),
			zap.String(logger.FieldPath, path),
			zap.Int(logger.FieldStatusCode, statusCode),
			zap.Int64(logger.FieldDuration, duration.Milliseconds()),
			zap.Int("response_size", responseSize),
			zap.String(logger.FieldClientIP, clientIP),
		}

		// Add error if present
		if errNext != nil {
			fields = append(fields, zap.Error(errNext))
			zap.L().Error("HTTP request completed with error", fields...)
		} else {
			// Log successful request at appropriate level based on status code
			switch {
			case statusCode >= httpStatusServerError:
				zap.L().Error("HTTP request completed with server error", fields...)
			case statusCode >= httpStatusClientError:
				zap.L().Warn("HTTP request completed with client error", fields...)
			default:
				zap.L().Info("HTTP request completed successfully", fields...)
			}
		}

		if errNext != nil {
			return fmt.Errorf("request processing failed: %w", errNext)
		}

		return nil
	}
}
