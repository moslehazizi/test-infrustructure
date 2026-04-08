package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/internal/usecase/interfaces"
)

func NewScenarioExecutorBuilder() *scenarioExecutorBuilder {
	return &scenarioExecutorBuilder{}
}

type scenarioExecutorBuilder struct {
}

func (seb *scenarioExecutorBuilder) Build(
	scenario *entity.TestScenario,
	scenarioRepo repository.TestScenarioRepository,
	testAgentControllerToolBox interfaces.TestAgentControllerToolBox,
	scenarioExecutorBuilder interfaces.SingleScenarioExecutorBuilder,
) interfaces.ScenarioExecutor {
	return &scenarioExecutor{
		scenario:                   scenario,
		scenarioRepo:               scenarioRepo,
		testAgentControllerToolBox: testAgentControllerToolBox,
		scenarioExecutorBuilder:    scenarioExecutorBuilder,
	}
}
