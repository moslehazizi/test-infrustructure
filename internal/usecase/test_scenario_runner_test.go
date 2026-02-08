package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewScenarioExecutor(t *testing.T) {
	ex := NewScenarioExecutor(entity.TestScenario{ID: 1})
	assert.NotNil(t, ex)

	e, ok := ex.(*scenarioExecutor)
	assert.True(t, ok)
	assert.NotNil(t, e.scenario)
}

func Test_scenarioExecutor_GetID(t *testing.T) {
	ex := NewScenarioExecutor(entity.TestScenario{ID: 1})
	assert.NotNil(t, ex)

	assert.Equal(t, uint64(1), ex.GetID())
}

func Test_scenarioExecutor_GetScenario(t *testing.T) {
	sc := entity.TestScenario{ID: 1}
	ex := NewScenarioExecutor(sc)
	assert.NotNil(t, ex)

	assert.NotNil(t, ex.GetScenario())
	assert.Equal(t, sc.ID, ex.GetScenario().ID)
}

func TestNewInMemoryScenarioExecutorEngine(t *testing.T) {
	eng := NewInMemoryScenarioExecutorBox()
	assert.NotNil(t, eng)
}

func Test_inMemoryScenarioExecutorBox_Add(t *testing.T) {
	t.Parallel()

	scenarioID := uint64(1)
	mockExe := new(mocks.MockScenarioExecutorWithWG)
	mockExe.WG.Add(1)

	mockExe.On("GetID").Return(scenarioID)
	mockExe.On("GetScenario").Return(&entity.TestScenario{ID: scenarioID})
	mockExe.On("ResumeOrStart", mock.Anything)

	eng := NewInMemoryScenarioExecutorBox()
	eng.Add(mockExe)

	mockExe.WG.Wait()

	assert.True(t, eng.HasExecutor(scenarioID))
	mockExe.AssertCalled(t, "GetID")
	mockExe.AssertCalled(t, "GetScenario")
	mockExe.AssertCalled(t, "ResumeOrStart", mock.Anything)
}

func Test_inMemoryScenarioExecutorBox_HasExecutor(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		mockExe := new(mocks.MockScenarioExecutor)
		mockExe.On("GetID").Return(uint64(1))
		mockExe.On("GetScenario").Return(&entity.TestScenario{ID: 1})
		mockExe.On("ResumeOrStart", mock.Anything)
		eng := NewInMemoryScenarioExecutorBox()
		eng.Add(mockExe)

		assert.True(t, eng.HasExecutor(1))
	})
	t.Run("false", func(t *testing.T) {
		eng := NewInMemoryScenarioExecutorBox()

		assert.False(t, eng.HasExecutor(1))
	})
}
