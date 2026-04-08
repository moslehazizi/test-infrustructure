package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository/mocks"

	"testing"

	usecasemocks "control-panel-service/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
)

func TestScenarioExecutorBuilder(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		scenarioId := uint64(1)

		seb := new(scenarioExecutorBuilder)
		repo := new(mocks.MockTestScenario)
		toolbox := new(usecasemocks.MockTestAgentControllerToolBox)
		sseb := new(usecasemocks.MockSingleScenarioExecutorBuilder)

		sse := seb.Build(&entity.TestScenario{ID: scenarioId}, repo, toolbox, sseb)
		e, ok := sse.(*scenarioExecutor)
		assert.True(t, ok)
		assert.Equal(t, e.scenario.ID, scenarioId)
		assert.NotNil(t, e.scenarioRepo)
		assert.NotNil(t, e.testAgentControllerToolBox)
		assert.NotNil(t, e.scenarioExecutorBuilder)
	})
}
