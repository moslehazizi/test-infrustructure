package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	repoMocks "control-panel-service/internal/repository/mocks"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #region StressTestExecutionManager
func TestStressTestExecutionManager_AddScenario(t *testing.T) {
	t.Run("failed_case_scenario_max_service_count_is_null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)

		err := ex.AddScenario(context.Background(), &scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		agent := new(mocks.MockTestAgentController)
		agent.On("Run").Return(nil).Times(3)
		builder.On("Build", scenario).Times(3).Return(agent)

		ex := NewStressTestExecutionManager(builder, repo)

		err := ex.AddScenario(context.Background(), scenario)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.NoError(t, err)

		agent.AssertCalled(t, "Run")
		builder.AssertCalled(t, "Build", scenario)

	})
}

func TestStressTestExecutionManager_RunScenario(t *testing.T) {
	t.Run("success_case_does_not_exist", func(t *testing.T) {
		scenario := &entity.TestScenario{
			ID: 1,
		}
		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)
		err := ex.RunScenario(context.Background(), scenario)
		assert.NoError(t, err)

		mng, ok := ex.(*StressTestExecutionManager)
		assert.True(t, ok)
		assert.Contains(t, mng.scenarios, scenario.ID)
	})
	t.Run("success_case_scenario_already_exists_in_manager", func(t *testing.T) {
		scenario := &entity.TestScenario{
			ID: 1,
		}
		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		// add for first time
		ex := NewStressTestExecutionManager(builder, repo)
		err := ex.RunScenario(context.Background(), scenario)
		assert.NoError(t, err)

		mng, ok := ex.(*StressTestExecutionManager)
		assert.True(t, ok)
		assert.Contains(t, mng.scenarios, scenario.ID)

		// make sure scenario is not running
		mng.scenarios[scenario.ID].SetRunning(false)

		// try to run it again
		err = ex.RunScenario(context.Background(), scenario)
		assert.NoError(t, err)
		assert.True(t, mng.scenarios[scenario.ID].IsRunning())
	})
}

func TestStressTestExecutionManager_PauseScenario(t *testing.T) {
	t.Run("failed_case_scenario_max_service_count_is_null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)

		err := ex.PauseScenario(context.Background(), &scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)
		err := ex.PauseScenario(context.Background(), scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
	})

	t.Run("failed_case_pause_testing_error", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		agent := new(mocks.MockTestAgentController)
		agent.On("PauseTesting", mock.Anything).Return(errors.New("something went wrong"))

		ex := NewStressTestExecutionManager(builder, repo)

		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     true,
		}

		err := ex.PauseScenario(context.Background(), scenario)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToPauseTestService)

		agent.AssertCalled(t, "PauseTesting", mock.Anything)
	})

	t.Run("successـcase", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		agent := new(mocks.MockTestAgentController)
		agent.On("PauseTesting", mock.Anything).Return(nil)

		ex := NewStressTestExecutionManager(builder, repo)

		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     true,
		}

		err := ex.PauseScenario(context.Background(), scenario)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.NoError(t, err)

		agent.AssertCalled(t, "PauseTesting", mock.Anything)
	})
}

func TestStressTestExecutionManager_ResumeScenario(t *testing.T) {
	t.Run("failed_case_scenario_max_service_count_is_null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)

		err := ex.ResumeScenario(context.Background(), &scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)
		err := ex.ResumeScenario(context.Background(), scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
	})

	t.Run("failed_case_resume_testing_error", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		agent := new(mocks.MockTestAgentController)
		agent.On("ResumeTesting", mock.Anything).Return(errors.New("something went wrong"))

		ex := NewStressTestExecutionManager(builder, repo)

		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     true,
		}

		err := ex.ResumeScenario(context.Background(), scenario)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToResumeTestService)

		agent.AssertCalled(t, "ResumeTesting", mock.Anything)
	})

	t.Run("successـcase", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		agent := new(mocks.MockTestAgentController)
		builder.On("Get", scenario).Times(3).Return(agent)
		agent.On("ResumeTesting", mock.Anything).Return(nil)

		ex := NewStressTestExecutionManager(builder, repo)

		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     true,
		}

		err := ex.ResumeScenario(context.Background(), scenario)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.NoError(t, err)

		agent.AssertCalled(t, "ResumeTesting", mock.Anything)
	})
}

func TestStressTestExecutionManager_Stop(t *testing.T) {
	t.Run("failed_case_scenario_max_service_count_is_null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)

		err := ex.StopScenario(context.Background(), &scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)
		err := ex.StopScenario(context.Background(), scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
	})

	t.Run("failed_case_failed_to_abort", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		agent := new(mocks.MockTestAgentController)
		agent.On("StopTesting", mock.Anything).Return(errors.New("something went wrong"))

		ex := NewStressTestExecutionManager(builder, repo)

		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     true,
		}

		err := ex.StopScenario(context.Background(), scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToStopTestService)

		agent.AssertCalled(t, "StopTesting", mock.Anything)
	})

	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}
		builder := new(mocks.MockTestAgentControllerToolBox)
		agent := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)

		executionID := uuid.New()
		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     false,
		}

		builder.On("Get", scenario).Times(3).Return(agent)
		agent.On("StopTesting", mock.Anything).Return(nil)

		err := ex.StopScenario(context.Background(), scenario)

		assert.Nil(t, err)

		agent.AssertCalled(t, "StopTesting", mock.Anything)
	})
}

func TestStressTestExecutionManager_Abort(t *testing.T) {
	t.Run("failed_case_scenario_max_service_count_is_null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)

		err := ex.AbortScenario(context.Background(), &scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)
		err := ex.AbortScenario(context.Background(), scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
	})

	t.Run("failed_case_failed_to_abort", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		agent := new(mocks.MockTestAgentController)
		agent.On("AbortTesting", mock.Anything).Return(errors.New("something went wrong"))

		ex := NewStressTestExecutionManager(builder, repo)

		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     true,
		}

		err := ex.AbortScenario(context.Background(), scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToAbortTestService)

		agent.AssertCalled(t, "AbortTesting", mock.Anything)
	})

	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}
		builder := new(mocks.MockTestAgentControllerToolBox)
		agent := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)

		executionID := uuid.New()
		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     false,
		}

		builder.On("Get", scenario).Times(3).Return(agent)
		agent.On("AbortTesting", mock.Anything).Return(nil)

		err := ex.AbortScenario(context.Background(), scenario)

		assert.Nil(t, err)

		agent.AssertCalled(t, "AbortTesting", mock.Anything)
	})
}

// #endregion StressTestExecutionManager
