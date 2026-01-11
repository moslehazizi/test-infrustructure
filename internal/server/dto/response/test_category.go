package response

import "time"

type TestCategory struct {
	ID                     uint64    `json:"id"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	Name                   string    `json:"name"`
	Label                  string    `json:"label"`
	HasMaxTestServiceCount bool      `json:"has_max_test_service_count"`
	HasExecutionDuration   bool      `json:"has_execution_duration"`
	HasAutoStepChangeRate  bool      `json:"has_auto_step_increase_rate"`
}
