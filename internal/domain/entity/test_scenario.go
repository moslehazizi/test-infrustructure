package entity

import (
	"control-panel-service/pkg"
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
	ID                  uint64             `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt           time.Time          `gorm:"column:created_at"`
	UpdatedAt           time.Time          `gorm:"column:updated_at"`
	DeletedAt           *gorm.DeletedAt    `gorm:"column:deleted_at"`
	Name                string             `gorm:"column:name"`
	TestCategoryID      uint64             `gorm:"column:test_category_id"`
	TestCategory        *TestCategory      `gorm:"ForeignKey:TestCategoryID"`
	MotherServiceID     uint64             `gorm:"column:mother_service_id"`
	MotherService       *MotherService     `gorm:"ForeignKey:MotherServiceID"`
	Status              ScenarioStatus     `gorm:"column:status"`
	MaxTestServiceCount *int               `gorm:"column:max_test_service_count"`
	ExecutionDuration   *time.Duration     `gorm:"column:execution_duration"`
	AutoStepChangeRate  *int               `gorm:"column:auto_step_change_rate"`
	TestServiceConfig   *TestServiceConfig `gorm:"ForeignKey:TestScenarioID"`
	StartedAt           *time.Time         `gorm:"column:started_at"`
}

func (TestScenario) TableName() string {
	return "test_scenarios"
}

type TestScenarioPaginationRequest struct {
	Page    int
	PerPage int
}

func (ts *TestScenario) Validate(testCat *TestCategory) error {
	if ts.MaxTestServiceCount != nil && *ts.MaxTestServiceCount < 1 {
		return pkg.ErrMaxTestServiceCountLessThanOne
	}

	if ts.ExecutionDuration != nil && *ts.ExecutionDuration < 1 {
		return pkg.ErrExecutionDurationLessThanOne
	}

	if ts.AutoStepChangeRate != nil && *ts.AutoStepChangeRate < 1 {
		return pkg.ErrAutoStepChangeRateLessThanOne
	}

	if testCat.HasMaxTestServiceCount {
		if ts.MaxTestServiceCount == nil {
			return pkg.ErrMaxTestServiceCountNotSet
		}
	} else {
		if ts.MaxTestServiceCount != nil {
			return pkg.ErrNoNeedMaxTestServiceCount
		}
	}

	if testCat.HasExecutionDuration {
		if ts.ExecutionDuration == nil {
			return pkg.ErrExecutionDurationNotSet
		}
	} else {
		if ts.ExecutionDuration != nil {
			return pkg.ErrNoNeedExecutionDuration
		}
	}

	if testCat.HasAutoStepChangeRate {
		if ts.AutoStepChangeRate == nil {
			return pkg.ErrAutoStepChangeNotSet
		}
	} else {
		if ts.AutoStepChangeRate != nil {
			return pkg.ErrNoNeedAutoStepChange
		}
	}

	return nil
}
