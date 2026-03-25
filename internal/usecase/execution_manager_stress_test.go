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

// #region scenarioExecutor
func Test_ScenarioExecutor_SetExecutionID(t *testing.T) {
	agent := new(mocks.MockTestAgentController)

	ex := &scenarioExecutor{
		scenario: &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		},
		agents:           []interfaces.TestAgentController{agent},
		allAgentsHealthy: false,
		running:          false,
	}
	sampleUUID := uuid.New()
	ex.SetExecutionID(sampleUUID)

	assert.Equal(t, sampleUUID, ex.executionID)
}
func Test_ScenarioExecutor_GetAgents(t *testing.T) {
	agent := new(mocks.MockTestAgentController)

	ex := &scenarioExecutor{
		scenario: &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		},
		executionID:      uuid.New(),
		agents:           []interfaces.TestAgentController{agent},
		allAgentsHealthy: false,
		running:          false,
	}
	agents := ex.GetAgents()

	assert.NotNil(t, agents)
	assert.Equal(t, agents[0], agent)
}
func Test_scenarioExecutor_AddAgent(t *testing.T) {
	ex := new(scenarioExecutor)
	assert.Len(t, ex.agents, 0)

	ex.AddAgent(&mocks.MockTestAgentController{})
	assert.Len(t, ex.agents, 1)
}
func Test_scenarioExecutor_AllAgentsAreHealthy(t *testing.T) {
	agent := new(mocks.MockTestAgentController)

	ex := &scenarioExecutor{
		scenario: &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		},
		executionID:      uuid.New(),
		agents:           []interfaces.TestAgentController{agent},
		allAgentsHealthy: false,
		running:          false,
	}

	assert.False(t, ex.AllAgentsAreHealthy())

	ex.allAgentsHealthy = true
	assert.True(t, ex.AllAgentsAreHealthy())
}
func Test_scenarioExecutor_IsRunning(t *testing.T) {
	agent := new(mocks.MockTestAgentController)

	ex := &scenarioExecutor{
		scenario: &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		},
		executionID:      uuid.New(),
		agents:           []interfaces.TestAgentController{agent},
		allAgentsHealthy: false,
		running:          false,
	}

	assert.False(t, ex.IsRunning())
}
func Test_scenarioExecutor_SetRunning(t *testing.T) {
	agent := new(mocks.MockTestAgentController)

	ex := &scenarioExecutor{
		scenario: &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID: 7,

				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		},
		executionID:      uuid.New(),
		agents:           []interfaces.TestAgentController{agent},
		allAgentsHealthy: false,
		running:          false,
	}

	assert.False(t, ex.IsRunning())
	ex.SetRunning(true)
	assert.True(t, ex.IsRunning())
	ex.SetRunning(false)
	assert.False(t, ex.IsRunning())

}
func Test_scenarioExecutor_awaitAgentsToBeHealthy(t *testing.T) {
	t.Run("all_agents_test_services_are_healthy_single_agent", func(t *testing.T) {
		t.Parallel()

		scenario := &entity.TestScenario{
			ID:       1,
			NumSteps: 2,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		}

		agent := new(mocks.MockTestAgentController)
		agent.On("Healthy").Times(1).Return(true)

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent},
			allAgentsHealthy: false,
			running:          false,
		}

		ex.awaitAgentsToBeHealthy()

		agent.AssertCalled(t, "Healthy")

		assert.True(t, ex.allAgentsHealthy)
	})

	t.Run("all_agents_test_services_are_healthy_2_agents", func(t *testing.T) {
		t.Parallel()

		scenario := &entity.TestScenario{
			ID:       1,
			NumSteps: 2,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: false,
			running:          false,
		}

		healthyCheckSleep = time.Millisecond * 10
		agent1.On("Healthy").Return(true) // may be called multiple times in the health-check loop
		agent2.On("Healthy").Return(true)

		go ex.awaitAgentsToBeHealthy()
		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)
		assert.True(t, ex.allAgentsHealthy)

		agent1.AssertCalled(t, "Healthy")
		agent2.AssertCalled(t, "Healthy")
	})
}
func Test_scenarioExecutor_awaitAgentsToBeReadyToStartTesting(t *testing.T) {
	t.Run("all_agents_test_services_are_ready_single_agent", func(t *testing.T) {
		scenario := &entity.TestScenario{
			ID:       1,
			NumSteps: 2,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		}

		agent := new(mocks.MockTestAgentController)
		agent.On("ReadyForTesting").Times(1).Return(true)

		ex := &scenarioExecutor{
			scenario:                 scenario,
			executionID:              uuid.New(),
			agents:                   []interfaces.TestAgentController{agent},
			allAgentsHealthy:         false,
			allAgentsReadyForTesting: false,
			running:                  false,
		}

		ex.awaitAgentsToBeReadyToStartTesting()

		agent.AssertCalled(t, "ReadyForTesting")

		assert.True(t, ex.allAgentsReadyForTesting)
	})

	t.Run("all_agents_test_services_are_ready_single_agent_first_time_is_not_ready", func(t *testing.T) {
		t.Parallel()

		scenario := &entity.TestScenario{
			ID:       1,
			NumSteps: 2,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		}

		agent := new(mocks.MockTestAgentController)
		agent.On("ReadyForTesting").Times(1).Return(false)

		ex := &scenarioExecutor{
			scenario:                 scenario,
			executionID:              uuid.New(),
			agents:                   []interfaces.TestAgentController{agent},
			allAgentsHealthy:         false,
			allAgentsReadyForTesting: false,
			running:                  false,
		}

		go ex.awaitAgentsToBeReadyToStartTesting()
		time.Sleep(time.Millisecond)

		agent.AssertCalled(t, "ReadyForTesting")

		assert.False(t, ex.allAgentsReadyForTesting)

		agent.On("ReadyForTesting").Times(1).Return(true)
		time.Sleep(readyForTestingCheckSleep)

		time.Sleep(time.Millisecond)

		agent.AssertCalled(t, "ReadyForTesting")
		assert.True(t, ex.allAgentsReadyForTesting)
	})

	t.Run("all_agents_test_services_are_ready_two_agents", func(t *testing.T) {
		scenario := &entity.TestScenario{
			ID:       1,
			NumSteps: 2,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
		}

		agent1 := new(mocks.MockTestAgentController)
		agent1.On("ReadyForTesting").Times(1).Return(true)

		agent2 := new(mocks.MockTestAgentController)
		agent2.On("ReadyForTesting").Times(1).Return(true)

		ex := &scenarioExecutor{
			scenario:                 scenario,
			executionID:              uuid.New(),
			agents:                   []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy:         false,
			allAgentsReadyForTesting: false,
			running:                  false,
		}

		ex.awaitAgentsToBeReadyToStartTesting()

		agent1.AssertCalled(t, "ReadyForTesting")
		agent2.AssertCalled(t, "ReadyForTesting")

		assert.True(t, ex.allAgentsReadyForTesting)
	})
}
func Test_scenarioExecutor_Run(t *testing.T) {
	t.Run("success_single_step", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()

		// setting up agents
		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
			NumSteps:            1,
			ExecutionDuration:   new(int64(100)),
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                    1,
				TestScenarioID:        1,
				MaxRequests:           100,
				MaxDuration:           1000,
				RequestDelayDuration:  new(10),
				RandomRequestDelayMin: nil,
				RandomRequestDelayMax: nil,
				FixedTestNumber:       new(43),
				RandomTestNumberMin:   nil,
				RandomTestNumberMax:   nil,
				BadValueRate:          2,
				NegativeValueRate:     10,
				RealValueRate:         10,
				ZeroValueRate:         0,
				StringValueRate:       80,
				LongStringValueRate:   0,
				NullValueRate:         0,
				DatabaseName:          "dbname",
				DatabaseTableName:     "tbl",
			},
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: false,
			running:          false,
			scenarioRepo:     repo,
		}

		healthyCheckSleep = time.Millisecond * 10
		agent1.On("Healthy").Return(true)
		agent2.On("Healthy").Return(true)

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		ex.SetRunning(true)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		repo.On("SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true).Return(nil)
		err := ex.Run(ctx)

		agent1.AssertCalled(t, "Healthy")
		agent2.AssertCalled(t, "Healthy")

		assert.True(t, ex.allAgentsHealthy)

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		repo.AssertCalled(t, "SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true)

		assert.NoError(t, err)
	})

	t.Run("success_2_steps", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()

		// setting up agents
		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
			NumSteps:            2,
			ExecutionDuration:   new(int64(100)),
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                    1,
				TestScenarioID:        1,
				MaxRequests:           100,
				MaxDuration:           1000,
				RequestDelayDuration:  new(10),
				RandomRequestDelayMin: nil,
				RandomRequestDelayMax: nil,
				FixedTestNumber:       new(43),
				RandomTestNumberMin:   nil,
				RandomTestNumberMax:   nil,
				BadValueRate:          2,
				NegativeValueRate:     10,
				RealValueRate:         10,
				ZeroValueRate:         0,
				StringValueRate:       80,
				LongStringValueRate:   0,
				NullValueRate:         0,
				DatabaseName:          "dbname",
				DatabaseTableName:     "tbl",
			},
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: false,
			running:          true,
			scenarioRepo:     repo,
		}

		healthyCheckSleep = time.Millisecond * 10
		agent1.On("Healthy").Return(true)
		agent2.On("Healthy").Return(true)

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		repo.On("SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true).Return(nil)

		agent1.On("ReadyForTesting").Return(true)
		agent2.On("ReadyForTesting").Return(true)

		err := ex.Run(ctx)

		agent1.AssertCalled(t, "Healthy")
		agent2.AssertCalled(t, "Healthy")

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)

		agent1.AssertCalled(t, "ReadyForTesting")
		agent2.AssertCalled(t, "ReadyForTesting")

		repo.AssertCalled(t, "SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true)

		assert.NoError(t, err)
	})
}

//#endregion scenarioExecutor

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
	t.Run("failed_case_scenario_is_null", func(t *testing.T) {
		var scenario *entity.TestScenario

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)
		err := ex.RunScenario(context.Background(), scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioServiceIsNil)
	})

	t.Run("failed_case_scenario_not_found", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: nil,
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)
		err := ex.RunScenario(context.Background(), scenario)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrTestScenarioNotFound)
	})

	t.Run("success_case", func(t *testing.T) {
		ctx := context.Background()
		executionID := uuid.New()
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		ex := NewStressTestExecutionManager(builder, repo)
		agent := new(mocks.MockTestAgentController)

		stem, _ := ex.(*StressTestExecutionManager)
		stem.scenarios[scenario.ID] = &scenarioExecutor{
			scenario:    scenario,
			executionID: executionID,
			agents:      []interfaces.TestAgentController{agent},
		}

		err := ex.RunScenario(ctx, scenario)

		assert.NoError(t, err)
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

func TestStressTestExecutionManager_Run(t *testing.T) {
	t.Run("all_executors_are_running", func(t *testing.T) {
		t.Parallel()

		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerToolBox)
		repo := new(repoMocks.MockTestScenario)

		agent := new(mocks.MockTestAgentController)
		agent.On("Run").Return(nil).Times(1)
		agent.On("Healthy").Return(false)

		builder.On("Build", scenario).Times(1).Return(agent)

		ex := NewStressTestExecutionManager(builder, repo)

		err := ex.AddScenario(context.Background(), scenario)
		assert.NoError(t, err)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		agent.AssertCalled(t, "Run")
		builder.AssertCalled(t, "Build", scenario)

		checkLoopSleep = time.Millisecond * 10

		// now, we should watch on manager Run function
		go ex.Run()
		time.Sleep(time.Millisecond)
		stEx := ex.(*StressTestExecutionManager)
		assert.False(t, stEx.scenarios[1].IsRunning())

		// make sure if already is running, it hit continue.
		time.Sleep(time.Millisecond * 11)
		assert.False(t, stEx.scenarios[1].IsRunning())
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
