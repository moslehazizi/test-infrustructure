package interfaces

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

// ScenarioExecutor is responsible to run a single scenario.
type ScenarioExecutor interface {
	GetID() uint64
	GetScenario() *entity.TestScenario
	// ResumeOrStart tries to resume execution from last step
	// and if never started, it starts from the beginning.
	ResumeOrStart(ctx context.Context)
}

type ScenarioExecutorBox interface {
	Add(exe ScenarioExecutor)
	HasExecutor(id uint64) bool
}
