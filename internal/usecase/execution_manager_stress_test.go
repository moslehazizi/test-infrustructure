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
func TestStressTestExecutionManager_RunScenario(t *testing.T) {
	t.Run("success_case_delete_scenario_from_memory_when_scenario_added_to_manager_for_first_time", func(t *testing.T) {
		scenario := &entity.TestScenario{
			ID: 1,
		}
		toolbox := new(mocks.MockTestAgentControllerToolBox)
		sseb := new(mocks.MockSingleScenarioExecutorBuilder)
		repo := new(repoMocks.MockTestScenario)
		seb := new(mocks.MockScenarioExecutorBuilderWithWait)
		seb.Wait = time.Second

		sampleSE := &scenarioExecutor{
			scenario:                   scenario,
			scenarioRepo:               repo,
			testAgentControllerToolBox: toolbox,
			scenarioExecutorBuilder:    sseb,
		}

		ex := NewStressTestExecutionManager(toolbox, repo, seb)
		repo.On("SetStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		seb.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sampleSE)

		mng, ok := ex.(*StressTestExecutionManager)
		assert.True(t, ok)
		assert.NotContains(t, mng.scenarios, scenario.ID)

		err := ex.RunScenario(context.Background(), scenario)
		assert.NoError(t, err)

		assert.Contains(t, mng.scenarios, scenario.ID)

		time.Sleep(1100 * time.Millisecond)

		assert.NotContains(t, mng.scenarios, scenario.ID)

		seb.AssertCalled(t, "Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("success_case_scenario_does_not_exist_in_execution_manager", func(t *testing.T) {
		scenario := &entity.TestScenario{
			ID: 1,
		}
		toolbox := new(mocks.MockTestAgentControllerToolBox)
		sseb := new(mocks.MockSingleScenarioExecutorBuilder)
		repo := new(repoMocks.MockTestScenario)
		seb := new(mocks.MockScenarioExecutorBuilder)

		sampleSE := &scenarioExecutor{
			scenario:                   scenario,
			scenarioRepo:               repo,
			testAgentControllerToolBox: toolbox,
			scenarioExecutorBuilder:    sseb,
		}

		ex := NewStressTestExecutionManager(toolbox, repo, seb)
		repo.On("SetStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		seb.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sampleSE)

		err := ex.RunScenario(context.Background(), scenario)
		assert.NoError(t, err)

		_, ok := ex.(*StressTestExecutionManager)
		assert.True(t, ok)

		seb.AssertCalled(t, "Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
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
		seb := new(mocks.MockScenarioExecutorBuilder)
		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)
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
		seb := new(mocks.MockScenarioExecutorBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Healthy").Return(true)
		agent.On("PauseTesting", mock.Anything).Return(errors.New("something went wrong"))

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		agent.AssertCalled(t, "Healthy")

	})

	t.Run("success_case", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		seb := new(mocks.MockScenarioExecutorBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Healthy").Return(true)
		agent.On("PauseTesting", mock.Anything).Return(nil)

		ex := NewStressTestExecutionManager(builder, repo, seb)

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

		agent.AssertCalled(t, "Healthy")
		agent.AssertCalled(t, "PauseTesting", mock.Anything)
	})

	t.Run("failed_case", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		seb := new(mocks.MockScenarioExecutorBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Healthy").Return(true)
		agent.On("PauseTesting", mock.Anything).Return(nil)

		ex := NewStressTestExecutionManager(builder, repo, seb)

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

		agent.AssertCalled(t, "Healthy")
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
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)
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
		seb := new(mocks.MockScenarioExecutorBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Healthy").Return(true)
		agent.On("ResumeTesting", mock.Anything).Return(errors.New("something went wrong"))

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		agent.AssertCalled(t, "Healthy")
	})

	t.Run("success_case", func(t *testing.T) {
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)
		seb := new(mocks.MockScenarioExecutorBuilder)

		agent := new(mocks.MockTestAgentController)
		builder.On("Get", scenario).Times(3).Return(agent)
		agent.On("Healthy").Return(true)
		agent.On("ResumeTesting", mock.Anything).Return(nil)

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		agent.AssertCalled(t, "Healthy")
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
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)
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
		seb := new(mocks.MockScenarioExecutorBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Healthy").Return(true)
		agent.On("StopTesting", mock.Anything).Return(errors.New("something went wrong"))

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		agent.AssertCalled(t, "Healthy", mock.Anything)
	})

	t.Run("success_case", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}
		builder := new(mocks.MockTestAgentControllerToolBox)
		agent := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)

		executionID := uuid.New()
		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
			running:     false,
		}

		builder.On("Get", scenario).Times(3).Return(agent)
		agent.On("Healthy", mock.Anything).Return(true)
		agent.On("StopTesting", mock.Anything).Return(nil)

		err := ex.StopScenario(context.Background(), scenario)

		assert.Nil(t, err)

		agent.AssertCalled(t, "StopTesting", mock.Anything)
		agent.AssertCalled(t, "Healthy")
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
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)
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
		seb := new(mocks.MockScenarioExecutorBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("AbortTesting", mock.Anything).Return(errors.New("something went wrong"))

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
		seb := new(mocks.MockScenarioExecutorBuilder)

		ex := NewStressTestExecutionManager(builder, repo, seb)

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
