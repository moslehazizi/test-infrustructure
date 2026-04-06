package mocks

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/internal/usecase/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockScenarioExecutorBuilder struct {
	mock.Mock
}

func (m *MockScenarioExecutorBuilder) Build(
	scenario *entity.TestScenario,
	scenarioRepo repository.TestScenarioRepository,
	testAgentControllerToolBox interfaces.TestAgentControllerToolBox,
	scenarioExecutorBuilder interfaces.SingleScenarioExecutorBuilder,
) interfaces.ScenarioExecutor {
	args := m.Called(scenario, scenarioRepo, testAgentControllerToolBox, scenarioExecutorBuilder)

	return args.Get(0).(interfaces.ScenarioExecutor)
}
