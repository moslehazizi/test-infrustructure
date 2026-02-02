package logger

// Standard field names for consistent logging across the application.
// Use these constants to avoid typos and ensure consistency.
const (
	// Service identification.
	FieldService = "service"
	FieldVersion = "version"
	FieldEnv     = "environment"
	FieldHost    = "host"

	// Request tracking.
	FieldRequestID     = "request_id"
	FieldCorrelationID = "correlation_id"
	FieldUserID        = "user_id"

	// HTTP fields.
	FieldMethod     = "method"
	FieldPath       = "path"
	FieldEndpoint   = "endpoint"
	FieldStatusCode = "status_code"
	FieldDuration   = "duration_ms"
	FieldIP         = "ip"
	FieldClientIP   = "client_ip"

	// Tracing.
	FieldTraceID = "trace_id"
	FieldSpanID  = "span_id"

	// Error fields.
	FieldError      = "error"
	FieldErrorType  = "error_type"
	FieldStackTrace = "stack_trace"

	// Operation fields.
	FieldOperation = "operation"
	FieldResult    = "result"

	// Kafka fields.
	FieldTopic     = "topic"
	FieldPartition = "partition"
	FieldOffset    = "offset"
	FieldMessageID = "message_id"

	// Database fields.
	FieldQuery        = "query"
	FieldTable        = "table"
	FieldRowsAffected = "rows_affected"
)
