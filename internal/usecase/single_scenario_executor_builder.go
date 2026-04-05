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
	allAgentsHealthy bool,
	allAgentsReadyForTesting bool,
	scenario *entity.TestScenario,
	executionID uuid.UUID,
	running bool,
) interfaces.SingleScenarioExecutor {
	return &singleScenarioExecutor{
		agents,
		allAgentsHealthy,
		allAgentsReadyForTesting,
		scenario,
		executionID,
		running,
	}
}
