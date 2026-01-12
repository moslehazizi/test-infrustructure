package response

import (
	"control-panel-service/internal/domain/entity"
	"time"
)

type TestScenario struct {
	ID                  uint64                `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt           time.Time             `gorm:"column:created_at"`
	UpdatedAt           time.Time             `gorm:"column:updated_at"`
	Name                string                `gorm:"column:name"`
	TestCategoryID      uint64                `gorm:"column:test_category_id"`
	MotherServiceID     uint64                `gorm:"column:mother_service_id"`
	Status              entity.ScenarioStatus `gorm:"column:status"`
	MaxTestServiceCount *int                  `gorm:"column:max_test_service_count"`
	ExecutionDuration   *int                  `gorm:"column:execution_duration"`
	AutoStepChangeRate  *int                  `gorm:"column:auto_step_change_rate"`
}
