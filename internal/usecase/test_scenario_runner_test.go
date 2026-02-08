package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewScenarioExecutor(t *testing.T) {
	ex := NewScenarioExecutor(
		entity.TestScenario{ID: 1},

		new(mocks.MockScenarioTypeRunner),
	)
	assert.NotNil(t, ex)

	e, ok := ex.(*scenarioExecutor)
	assert.True(t, ok)
	assert.NotNil(t, e.scenario)
	assert.NotNil(t, e.runnerGroupA)
}

func Test_scenarioExecutor_GetID(t *testing.T) {
	ex := NewScenarioExecutor(
		entity.TestScenario{ID: 1},
		new(mocks.MockScenarioTypeRunner),
	)
	assert.NotNil(t, ex)

	assert.Equal(t, uint64(1), ex.GetID())
}

func Test_scenarioExecutor_GetScenario(t *testing.T) {
	sc := entity.TestScenario{ID: 1}
	ex := NewScenarioExecutor(
		sc,
		new(mocks.MockScenarioTypeRunner),
	)
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
	mockExe.On("ResumeOrStart", mock.Anything).Return(nil)

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
		mockExe.On("ResumeOrStart", mock.Anything).Return(nil)

		eng := NewInMemoryScenarioExecutorBox()
		eng.Add(mockExe)

		assert.True(t, eng.HasExecutor(1))
	})
	t.Run("false", func(t *testing.T) {
		eng := NewInMemoryScenarioExecutorBox()

		assert.False(t, eng.HasExecutor(1))
	})
}

func Test_scenarioExecutor_ResumeOrStart(t *testing.T) {
	t.Run("success case: group A", func(t *testing.T) {
		serviceCnt := 3
		dur := time.Millisecond * 2000
		scenario := &entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   &dur,
			StartedAt:           nil,
			TestCategoryID:      1,
			TestCategory: &entity.TestCategory{
				ID:                     1,
				Name:                   "load",
				Label:                  "load",
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
		}

		runner := new(mocks.MockScenarioTypeRunner)
		runner.
			On("Run", mock.Anything, scenario).
			Return(nil)

		ex := NewScenarioExecutor(*scenario, runner)
		err := ex.ResumeOrStart(context.Background())
		assert.NoError(t, err)

		runner.AssertCalled(t, "Run", mock.Anything, scenario)
	})
	t.Run("failed case: group A", func(t *testing.T) {
		serviceCnt := 3
		dur := time.Millisecond * 2000
		scenario := &entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   &dur,
			StartedAt:           nil,
			TestCategoryID:      1,
			TestCategory: &entity.TestCategory{
				ID:                     1,
				Name:                   "load",
				Label:                  "load",
				HasMaxTestServiceCount: true,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  false,
			},
		}

		runner := new(mocks.MockScenarioTypeRunner)
		runner.
			On("Run", mock.Anything, scenario).
			Return(errors.New("something went wrong"))

		ex := NewScenarioExecutor(*scenario, runner)
		err := ex.ResumeOrStart(context.Background())
		assert.Error(t, err)

		runner.AssertCalled(t, "Run", mock.Anything, scenario)
	})
	t.Run("failed case: group not implemented yet", func(t *testing.T) {
		serviceCnt := 3
		dur := time.Millisecond * 2000
		scenario := &entity.TestScenario{
			MaxTestServiceCount: &serviceCnt,
			ExecutionDuration:   &dur,
			StartedAt:           nil,
			TestCategoryID:      1,
			TestCategory: &entity.TestCategory{
				ID:                     1,
				Name:                   "load",
				Label:                  "load",
				HasMaxTestServiceCount: false,
				HasExecutionDuration:   true,
				HasAutoStepChangeRate:  true,
			},
		}

		runner := new(mocks.MockScenarioTypeRunner)

		ex := NewScenarioExecutor(*scenario, runner)
		err := ex.ResumeOrStart(context.Background())
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrNotImplemented)
	})
}
