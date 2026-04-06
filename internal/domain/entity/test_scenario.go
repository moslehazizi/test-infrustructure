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
	// The scenario can be automatically re-executed by increasing the number of test agents.
	// This field specifies how many agents should be added at the beginning of each new run
	// compared to the previous run.
	// For example, in every full execution of the scenario, 5 agents are added.
	// This means that when the scenario reaches the end,
	// it starts again from the beginning, and this time the number of agents
	// will be 5 more than in the previous execution.
	// An important note is that the execution_id field will be different
	// for each run that starts from the beginning, allowing them to be distinguished.
	IncreaseAgentNumber int64 `gorm:"column:increase_agent_number"`
	// This field becomes meaningful when used together with the field above.
	// It specifies how many times the agent increment should occur.
	// For example, if this value is 10, it means the entire scenario will be repeated 10 times,
	// and in each run, the number of agents will increase by the amount defined in the field above.
	// Assume:
	// IncreaseAgentNumber = 5 & ExecNumMultiAgent = 10 & MaxTestServiceCount = 5
	// In this case, the scenario first runs with 5 agents.
	// Then, with an increase of 5 agents, the entire scenario is repeated with 10 agents.
	// This process repeats 10 times, and the final run will be tested with 55 agents.
	ExecNumMultiAgent int64 `gorm:"column:execution_number_multi_agent"`
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

	if testCat.HasNumSteps {
		if ts.NumSteps < 1 {
			return pkg.ErrNumStepsNotSet
		}
	} else {
		// if this field be false it means that we have a kind of test that has one step
		if ts.NumSteps != 1 {
			return pkg.ErrNumStepsShouldBeOne
		}
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

	if ts.IncreaseAgentNumber < 0 {
		return pkg.ErrIncreaseAgentNumNotBeNegative
	}
	if ts.ExecNumMultiAgent < 1 {
		return pkg.ErrExecNumMultiAgentShouldBePositive
	}

	return nil
}
