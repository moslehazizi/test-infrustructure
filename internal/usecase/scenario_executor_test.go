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
func Test_ScenarioExecutor_assignExecutionID(t *testing.T) {
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
	ex.assignExecutionID()

	assert.NotEqual(t, ex.executionID, "00000000-0000-0000-0000-000000000000")
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
	t.Run("make_sure_execution_id_is_set", func(t *testing.T) {
		repo := new(repoMocks.MockTestScenario)
		ex := &scenarioExecutor{
			scenario:         &entity.TestScenario{ID: 1},
			agents:           []interfaces.TestAgentController{},
			allAgentsHealthy: false,
			running:          true,
			scenarioRepo:     repo,
		}

		err := ex.Run(context.Background())
		assert.NoError(t, err)
		assert.NotEqual(t, ex.executionID.String(), "00000000-0000-0000-0000-000000000000")
	})

	t.Run("scenario_is_not_running", func(t *testing.T) {
		repo := new(repoMocks.MockTestScenario)
		ex := &scenarioExecutor{
			scenario:         &entity.TestScenario{ID: 1},
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{},
			allAgentsHealthy: false,
			running:          false,
			scenarioRepo:     repo,
		}

		err := ex.Run(context.Background())
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrScenarioIsNotRunning)
	})
	t.Run("scenario_max_test_service_count_is_null", func(t *testing.T) {
		repo := new(repoMocks.MockTestScenario)
		ex := &scenarioExecutor{
			scenario:         &entity.TestScenario{ID: 1},
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{},
			allAgentsHealthy: false,
			running:          true,
			scenarioRepo:     repo,
		}

		err := ex.Run(context.Background())
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})
	t.Run("scenario_max_test_service_count_is_less_than_one", func(t *testing.T) {
		repo := new(repoMocks.MockTestScenario)
		ex := &scenarioExecutor{
			scenario:         &entity.TestScenario{ID: 1, MaxTestServiceCount: new(int64(-1))},
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{},
			allAgentsHealthy: false,
			running:          true,
			scenarioRepo:     repo,
		}

		err := ex.Run(context.Background())
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountLessThanOne)
	})

	t.Run("1_step_6_dynamic_agents_no_factorial_increment", func(t *testing.T) {
		t.Parallel()

		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            1,
			IncreaseAgentNumber: 2,
			ExecNumMultiAgent:   3,
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                     1,
				TestScenarioID:         1,
				MaxRequests:            100,
				MaxDuration:            1000,
				RequestDelayDuration:   new(10),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(43),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           2,
				NegativeValueRate:      10,
				RealValueRate:          10,
				ZeroValueRate:          0,
				StringValueRate:        80,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "dbname",
				DatabaseTableName:      "tbl",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		repo := new(repoMocks.MockTestScenario)
		ac := new(mocks.MockTestAgentControllerToolBox)
		ex := &scenarioExecutor{
			scenario:                   scenario,
			executionID:                uuid.New(),
			agents:                     []interfaces.TestAgentController{},
			allAgentsHealthy:           false,
			running:                    true,
			scenarioRepo:               repo,
			testAgentControllerToolBox: ac,
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)
		agent3 := new(mocks.MockTestAgentController)
		agent4 := new(mocks.MockTestAgentController)
		agent5 := new(mocks.MockTestAgentController)
		agent6 := new(mocks.MockTestAgentController)
		agent7 := new(mocks.MockTestAgentController)
		agent8 := new(mocks.MockTestAgentController)
		agent9 := new(mocks.MockTestAgentController)
		ac.On("Build", mock.Anything).Times(1).Return(agent1)
		ac.On("Build", mock.Anything).Times(1).Return(agent2)
		ac.On("Build", mock.Anything).Times(1).Return(agent3)
		ac.On("Build", mock.Anything).Times(1).Return(agent4)
		ac.On("Build", mock.Anything).Times(1).Return(agent5)
		ac.On("Build", mock.Anything).Times(1).Return(agent6)
		ac.On("Build", mock.Anything).Times(1).Return(agent7)
		ac.On("Build", mock.Anything).Times(1).Return(agent8)
		ac.On("Build", mock.Anything).Times(1).Return(agent9)

		agent1.On("Run").Times(1).Return(nil)
		agent2.On("Run").Times(1).Return(nil)
		agent3.On("Run").Times(1).Return(nil)
		agent4.On("Run").Times(1).Return(nil)
		agent5.On("Run").Times(1).Return(nil)
		agent6.On("Run").Times(1).Return(nil)
		agent7.On("Run").Times(1).Return(nil)
		agent8.On("Run").Times(1).Return(nil)
		agent9.On("Run").Times(1).Return(nil)
		// agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)

		err := ex.Run(context.Background())
		assert.NoError(t, err)

		ac.AssertCalled(t, "Build", mock.Anything)
		ac.AssertExpectations(t)
		ac.AssertNumberOfCalls(t, "Build", 9)

		// make sure all goroutines are called
		time.Sleep(time.Millisecond)
		agent1.AssertCalled(t, "Run")
		agent1.AssertNumberOfCalls(t, "Run", 1)
		agent2.AssertCalled(t, "Run")
		agent2.AssertNumberOfCalls(t, "Run", 1)
		agent3.AssertCalled(t, "Run")
		agent3.AssertNumberOfCalls(t, "Run", 1)
		agent4.AssertCalled(t, "Run")
		agent4.AssertNumberOfCalls(t, "Run", 1)
		agent5.AssertCalled(t, "Run")
		agent5.AssertNumberOfCalls(t, "Run", 1)
		agent6.AssertCalled(t, "Run")
		agent6.AssertNumberOfCalls(t, "Run", 1)
		agent7.AssertCalled(t, "Run")
		agent7.AssertNumberOfCalls(t, "Run", 1)
		agent8.AssertCalled(t, "Run")
		agent8.AssertNumberOfCalls(t, "Run", 1)
		agent9.AssertCalled(t, "Run")
		agent9.AssertNumberOfCalls(t, "Run", 1)
	})
}

func Test_scenarioExecutor_executeScenarioSteps(t *testing.T) {
	t.Run("scenario_is_not_running", func(t *testing.T) {
		ctx := context.Background()

		// setting up agents
		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
			NumSteps:            1,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                     1,
				TestScenarioID:         1,
				MaxRequests:            100,
				MaxDuration:            1000,
				RequestDelayDuration:   new(10),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(43),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           2,
				NegativeValueRate:      10,
				RealValueRate:          10,
				ZeroValueRate:          0,
				StringValueRate:        80,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "dbname",
				DatabaseTableName:      "tbl",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
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
		ex.SetRunning(false)

		err := ex.executeScenarioSteps(ctx)
		assert.ErrorIs(t, err, pkg.ErrScenarioIsNotRunning)
	})
	t.Run("success_single_step", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		// setting up agents
		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
			NumSteps:            1,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                     1,
				TestScenarioID:         1,
				MaxRequests:            100,
				MaxDuration:            1000,
				RequestDelayDuration:   new(10),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(43),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           2,
				NegativeValueRate:      10,
				RealValueRate:          10,
				ZeroValueRate:          0,
				StringValueRate:        80,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "dbname",
				DatabaseTableName:      "tbl",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: true,
			running:          false,
			scenarioRepo:     repo,
		}

		healthyCheckSleep = time.Millisecond * 10

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		ex.SetRunning(true)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		repo.On("SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true).Return(nil)
		err := ex.executeScenarioSteps(ctx)
		assert.NoError(t, err)

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		repo.AssertCalled(t, "SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true)
	})
	t.Run("success_2_steps", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		// setting up agents
		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
			NumSteps:            2,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                     1,
				TestScenarioID:         1,
				MaxRequests:            100,
				MaxDuration:            1000,
				RequestDelayDuration:   new(10),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(43),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           2,
				NegativeValueRate:      10,
				RealValueRate:          10,
				ZeroValueRate:          0,
				StringValueRate:        80,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "dbname",
				DatabaseTableName:      "tbl",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: true,
			running:          false,
			scenarioRepo:     repo,
		}

		healthyCheckSleep = time.Millisecond * 10

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		ex.SetRunning(true)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Times(2).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Times(2).Return(nil)

		agent1.On("ReadyForTesting").Times(1).Return(true)
		agent2.On("ReadyForTesting").Times(1).Return(true)

		repo.On("SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true).
			Times(1).
			Return(nil)
		err := ex.executeScenarioSteps(ctx)
		assert.NoError(t, err)

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)

		agent1.AssertCalled(t, "ReadyForTesting")
		agent2.AssertCalled(t, "ReadyForTesting")

		repo.AssertCalled(t, "SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true)
	})
	t.Run("failed_to_set_scenario_status", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		// setting up agents
		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
			NumSteps:            1,
			IncreaseAgentNumber: 0,
			ExecNumMultiAgent:   1,
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                     1,
				TestScenarioID:         1,
				MaxRequests:            100,
				MaxDuration:            1000,
				RequestDelayDuration:   new(10),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(43),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           2,
				NegativeValueRate:      10,
				RealValueRate:          10,
				ZeroValueRate:          0,
				StringValueRate:        80,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "dbname",
				DatabaseTableName:      "tbl",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: true,
			running:          false,
			scenarioRepo:     repo,
		}

		healthyCheckSleep = time.Millisecond * 10

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		ex.SetRunning(true)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		repo.On("SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true).Return(errors.New("something went wrong"))
		err := ex.executeScenarioSteps(ctx)
		assert.ErrorIs(t, err, pkg.ErrFailedToSetScenarioStatus)

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		repo.AssertCalled(t, "SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true)
	})
	t.Run("success_single_step_dynamic_agents", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		// setting up agents
		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
			NumSteps:            1,
			IncreaseAgentNumber: 1,
			ExecNumMultiAgent:   5,
			TestServiceConfig: &entity.TestServiceConfig{
				ID:                     1,
				TestScenarioID:         1,
				MaxRequests:            100,
				MaxDuration:            1000,
				RequestDelayDuration:   new(10),
				RandomRequestDelayMin:  nil,
				RandomRequestDelayMax:  nil,
				FixedTestNumber:        new(43),
				RandomTestNumberMin:    nil,
				RandomTestNumberMax:    nil,
				BadValueRate:           2,
				NegativeValueRate:      10,
				RealValueRate:          10,
				ZeroValueRate:          0,
				StringValueRate:        80,
				LongStringValueRate:    0,
				NullValueRate:          0,
				DatabaseName:           "dbname",
				DatabaseTableName:      "tbl",
				IncreaseFixedInput:     0,
				ExecNumMultiFixedInput: 1,
			},
		}

		healthyCheckSleep = time.Millisecond * 10

		// #region 1st
		// 1 agent + 1 step
		agent1 := new(mocks.MockTestAgentController)
		repo := new(repoMocks.MockTestScenario)
		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1},
			allAgentsHealthy: true,
			running:          false,
			scenarioRepo:     repo,
		}

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)
		ex.SetRunning(true)

		agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		repo.On("SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true).Return(nil)
		err := ex.executeScenarioSteps(ctx)
		assert.NoError(t, err)
		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		repo.AssertCalled(t, "SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true)

		// #endregion 1st

		// // #region 1st
		// // 1 agent + 1 step
		// agent1 := new(mocks.MockTestAgentController)
		// repo := new(repoMocks.MockTestScenario)
		// ex := &scenarioExecutor{
		// 	scenario:         scenario,
		// 	executionID:      uuid.New(),
		// 	agents:           []interfaces.TestAgentController{agent1},
		// 	allAgentsHealthy: true,
		// 	running:          false,
		// 	scenarioRepo:     repo,
		// }

		// time.Sleep(time.Millisecond)
		// time.Sleep(time.Millisecond * 10)
		// ex.SetRunning(true)

		// agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		// repo.On("SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true).Return(nil)
		// err := ex.executeScenarioSteps(ctx)
		// assert.NoError(t, err)
		// agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		// repo.AssertCalled(t, "SetStatus", mock.Anything, scenario.ID, entity.ScenarioStatusPending, true)

		// // #endregion 1st

	})
}

//#endregion scenarioExecutor
