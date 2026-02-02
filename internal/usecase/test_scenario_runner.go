package usecase

import (
	"control-panel-service/internal/domain/entity"
)

type ScenarioExecutorEngine interface {
	Add(scenario *entity.TestScenario)
}

func NewInMemoryScenarioExecutorEngine() ScenarioExecutorEngine {
	return &inMemoryScenarioExecutorEngine{}
}

type inMemoryScenarioExecutorEngine struct{}

func (e *inMemoryScenarioExecutorEngine) Add(scenario *entity.TestScenario) {}

// GetID() uuid.UUID
// GetScenario() *entity.TestScenario
// // ResumeOrStart tries to resume execution from last step
// // and if never started, it starts from the beginning.
// ResumeOrStart(ctx context.Context)
