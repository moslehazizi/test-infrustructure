package response

import "time"

type TestCategory struct {
	ID                     uint64    `json:"id"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	Name                   string    `json:"name"`
	Label                  string    `json:"label"`
	HasMaxTestServiceCount bool      `json:"has_max_test_service_count"`
	HasNumSteps            bool      `json:"has_num_steps"`
	Active                 bool      `json:"active"`
}
