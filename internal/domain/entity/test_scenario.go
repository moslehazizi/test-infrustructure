package entity

import (
	"time"

	"gorm.io/gorm"
)

type ScenarioStatus string

const (
	ScenarioStatusPending ScenarioStatus = "pending"
	ScenarioStatusRunning ScenarioStatus = "running"
	ScenarioStatusPaused  ScenarioStatus = "paused"
	ScenarioStatusStopped ScenarioStatus = "stopped"
	ScenarioStatusAborted ScenarioStatus = "aborted"
	ScenarioStatusSucceed ScenarioStatus = "succeed"
)

type TestScenario struct {
	ID                   uint64          `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt            time.Time       `gorm:"column:created_at"`
	UpdatedAt            time.Time       `gorm:"column:updated_at"`
	DeletedAt            *gorm.DeletedAt `gorm:"column:deleted_at"`
	Name                 string          `gorm:"column:name"`
	TestCategoryID       uint64          `gorm:"column:test_category_id"`
	MotherServiceID      uint64          `gorm:"column:mother_service_id"`
	Status               ScenarioStatus  `gorm:"column:status"`
	MaxTestServiceCount  *int            `gorm:"column:max_test_service_count"`
	ExecutionDuration    *int            `gorm:"column:execution_duration"`
	AutoStepIncreaseRate *int            `gorm:"column:auto_step_increase_rate"`
}

func (TestScenario) TableName() string {
	return "test_scenarios"
}

type TestScenarioPaginationRequest struct {
	Page    int
	PerPage int
}
