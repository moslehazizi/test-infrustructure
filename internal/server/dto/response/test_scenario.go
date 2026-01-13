package response

import (
	"control-panel-service/internal/domain/entity"
	"time"
)

type TestScenario struct {
	ID                  uint64                `json:"id"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
	Name                string                `json:"name"`
	TestCategory        *TestCategory         `json:"test_category"`
	MotherService       *MotherService        `json:"mother_service"`
	Status              entity.ScenarioStatus `json:"status"`
	MaxTestServiceCount *int                  `json:"max_test_service_count"`
	ExecutionDuration   *int                  `json:"execution_duration"`
	AutoStepChangeRate  *int                  `json:"auto_step_change_rate"`
}
