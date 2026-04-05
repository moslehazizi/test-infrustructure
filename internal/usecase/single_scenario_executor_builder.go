package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"

	"github.com/google/uuid"
)

type singleScenarioExecutorBuilder struct {
}

func (sseb *singleScenarioExecutorBuilder) Build(
	agents []interfaces.TestAgentController,
	scenario *entity.TestScenario,
	executionID uuid.UUID,
) interfaces.SingleScenarioExecutor {
	return &singleScenarioExecutor{
		agents:      agents,
		scenario:    scenario,
		executionID: executionID,
	}
}
