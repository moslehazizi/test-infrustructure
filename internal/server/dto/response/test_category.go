package response

import (
	"control-panel-service/internal/domain/entity"
	"time"
)

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

func (result *TestCategory) FromTestCategoryEntity(entityCategory *entity.TestCategory) {
	if entityCategory == nil {
		return
	}
	result.ID = entityCategory.ID
	result.CreatedAt = entityCategory.CreatedAt
	result.UpdatedAt = entityCategory.UpdatedAt
	result.Name = entityCategory.Name
	result.Label = entityCategory.Label
	result.HasMaxTestServiceCount = entityCategory.HasMaxTestServiceCount
	result.HasNumSteps = entityCategory.HasNumSteps
	result.Active = entityCategory.Active
}
