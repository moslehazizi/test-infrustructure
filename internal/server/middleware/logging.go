package middleware

import (
	"control-panel-service/pkg/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	httpStatusServerError = http.StatusInternalServerError // 500
	httpStatusClientError = http.StatusBadRequest          // 400
)

// LoggingMiddleware logs HTTP requests/responses with a request ID.
func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// Generate request ID if not present in headers
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			w.Header().Set("X-Request-ID", requestID)

			// Add request ID to context for handlers and downstream services
			ctx := logger.WithRequestID(r.Context(), requestID)
			r = r.WithContext(ctx)

			// Extract request information
			method := r.Method
			path := r.URL.Path
			originalURL := r.RequestURI
			clientIP := middleware.GetClientIP(ctx) // set by ClientIPFromRemoteAddr earlier
			userAgent := r.UserAgent()

			zap.L().Info("incoming HTTP request",
				zap.String(logger.FieldRequestID, requestID),
				zap.String(logger.FieldMethod, method),
				zap.String(logger.FieldPath, path),
				zap.String("original_path", originalURL),
				zap.String(logger.FieldClientIP, clientIP),
				zap.String("user_agent", userAgent),
			)

			// Wrap the ResponseWriter so we can observe status + size
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process request
			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			statusCode := ww.Status()
			if statusCode == 0 {
				statusCode = http.StatusOK // handler wrote nothing explicit
			}
			responseSize := ww.BytesWritten()

			fields := []zap.Field{
				zap.String(logger.FieldRequestID, requestID),
				zap.String(logger.FieldMethod, method),
				zap.String(logger.FieldPath, path),
				zap.Int(logger.FieldStatusCode, statusCode),
				zap.Int64(logger.FieldDuration, duration.Milliseconds()),
				zap.Int("response_size", responseSize),
				zap.String(logger.FieldClientIP, clientIP),
			}

			switch {
			case statusCode >= httpStatusServerError:
				zap.L().Error("HTTP request completed with server error", fields...)
			case statusCode >= httpStatusClientError:
				zap.L().Warn("HTTP request completed with client error", fields...)
			default:
				zap.L().Info("HTTP request completed successfully", fields...)
			}
		})
	}
}
