package response

import "time"

type HealthResponse struct {
	OK bool `json:"ok"`
}

type PauseResponse struct {
	Message string `json:"message"`
}

type ResumeResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MetricsSnapshot struct {
	Requests    int64         `json:"requests"`
	Success     int64         `json:"success"`
	Failed      int64         `json:"failed"`
	MinDuration time.Duration `json:"min_duration"`
	MaxDuration time.Duration `json:"max_duration"`
	AvgDuration time.Duration `json:"avg_duration"`
}

type FactorialExecutionResult struct {
	SuccessCount int64 `json:"success_count"`
	FailedCount  int64 `json:"failed_count"`
	MaxDuration  int64 `json:"max_duration"`
	MinDuration  int64 `json:"min_duration"`
	AveDuration  int64 `json:"ave_duration"`
}
