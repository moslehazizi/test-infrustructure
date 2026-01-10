package entity

import (
	"time"

	"gorm.io/gorm"
)

type TestScenario struct {
	ID                   uint64          `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt            time.Time       `gorm:"column:created_at"`
	UpdatedAt            time.Time       `gorm:"column:updated_at"`
	DeletedAt            *gorm.DeletedAt `gorm:"column:deleted_at"`
	Name                 string          `gorm:"column:name"`
	TestCategoryID       uint64          `gorm:"column:test_category_id"`
	MotherServiceID      uint64          `gorm:"column:mother_service_id"`
	MaxTestServiceCount  *int            `gorm:"coulmn:max_test_service_count"`
	ExecutionDuration    *int            `gorm:"column:execution_duration"`
	AutoStepIncreaseRate *int            `gorm:"column:auto_step_increase_rate"`
	StoppedAt            *time.Time      `gorm:"column:stopped_at"`
	RestartedAt          *time.Time      `gorm:"column:restarted_at"`
	StartedAt            *time.Time      `gorm:"column:started_at"`
}

func (TestScenario) TableName() string {
	return "test_scenarios"
}
