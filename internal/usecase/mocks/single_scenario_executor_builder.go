package mocks

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSingleScenarioExecutorBuilder struct {
	mock.Mock
}

func (m *MockSingleScenarioExecutorBuilder) Build(
	agents []interfaces.TestAgentController,
	allAgentsHealthy bool,
	allAgentsReadyForTesting bool,
	scenario *entity.TestScenario,
	executionID uuid.UUID,
	running bool,
) interfaces.SingleScenarioExecutor {
	args := m.Called(agents, allAgentsHealthy, allAgentsReadyForTesting, scenario, executionID, running)

	return args.Get(0).(interfaces.SingleScenarioExecutor)
}
