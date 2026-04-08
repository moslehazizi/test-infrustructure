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
	scenario *entity.TestScenario,
	executionID uuid.UUID,
) interfaces.SingleScenarioExecutor {
	args := m.Called(agents, scenario, executionID)

	return args.Get(0).(interfaces.SingleScenarioExecutor)
}
