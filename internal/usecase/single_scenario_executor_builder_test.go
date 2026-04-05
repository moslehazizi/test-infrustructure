package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/internal/usecase/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSingleScenarioExecutorBuilder(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)
		agents := []interfaces.TestAgentController{agent1, agent2}
		executionID := uuid.New()
		scenarioId := uint64(1)

		sseb := new(singleScenarioExecutorBuilder)

		sse := sseb.Build(agents, &entity.TestScenario{ID: scenarioId}, executionID)
		e, ok := sse.(*singleScenarioExecutor)
		assert.True(t, ok)
		assert.Equal(t, e.agents, agents)
		assert.Equal(t, e.executionID.String(), executionID.String())
		assert.Equal(t, e.scenario.ID, scenarioId)
	})
}
