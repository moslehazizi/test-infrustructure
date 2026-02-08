package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	requestIDKey     contextKey = "request_id"
	correlationIDKey contextKey = "correlation_id"
	userIDKey        contextKey = "user_id"
	traceIDKey       contextKey = "trace_id"
	spanIDKey        contextKey = "span_id"
)

// WithRequestID adds a request_id to the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// WithCorrelationID adds a correlation_id to the context.
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

// WithUserID adds a user_id to the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithTraceID adds a trace_id to the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// WithSpanID adds a span_id to the context.
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDKey, spanID)
}

// FieldsFromContext extracts logging fields from context.
// Returns zap fields for request_id, correlation_id, user_id, trace_id, and span_id if present.
func FieldsFromContext(ctx context.Context) []zap.Field {
	var fields []zap.Field

	if requestID, ok := ctx.Value(requestIDKey).(string); ok && requestID != "" {
		fields = append(fields, zap.String(FieldRequestID, requestID))
	}

	if correlationID, ok := ctx.Value(correlationIDKey).(string); ok && correlationID != "" {
		fields = append(fields, zap.String(FieldCorrelationID, correlationID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(FieldUserID, userID))
	}

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(FieldTraceID, traceID))
	}

	if spanID, ok := ctx.Value(spanIDKey).(string); ok && spanID != "" {
		fields = append(fields, zap.String(FieldSpanID, spanID))
	}

	return fields
}

// GetRequestID extracts request_id from context.
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}

	return ""
}

// GetCorrelationID extracts correlation_id from context.
func GetCorrelationID(ctx context.Context) string {
	if correlationID, ok := ctx.Value(correlationIDKey).(string); ok {
		return correlationID
	}

	return ""
}

// GetUserID extracts user_id from context.
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}

	return ""
}

// GetTraceID extracts trace_id from context.
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}

	return ""
}

// GetSpanID extracts span_id from context.
func GetSpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value(spanIDKey).(string); ok {
		return spanID
	}

	return ""
}

// WithContext returns a logger with fields extracted from context.
// This is a convenience function that combines zap.L() with FieldsFromContext().
// Usage: logger.WithContext(ctx).Info("message", zap.String("key", "value")).
func WithContext(ctx context.Context) *zap.Logger {
	fields := FieldsFromContext(ctx)
	if len(fields) == 0 {
		return zap.L()
	}

	return zap.L().With(fields...)
}
