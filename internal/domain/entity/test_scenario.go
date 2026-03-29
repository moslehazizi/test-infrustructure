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
	NumSteps            int64              `gorm:"column:num_steps"`
	MaxTestServiceCount *int64             `gorm:"column:max_test_service_count"`
	TestServiceConfig   *TestServiceConfig `gorm:"ForeignKey:TestScenarioID"`
	DeploymentNumber    int32              `gorm:"column:deployment_number"`
	StartedAt           *time.Time         `gorm:"column:started_at"`
	Editable            bool               `gorm:"column:editable;default:true;not null"`
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

	if ts.NumSteps < 1 {
		return pkg.ErrNumStepsNotSet
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

	return nil
}
