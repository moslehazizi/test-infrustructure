package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ScenarioExecutor_execute(t *testing.T) {
	t.Run("failed_case_scenario_is_not_running", func(t *testing.T) {
		agent := new(mocks.MockTestAgentController)

		ctx := context.Background()

		ex := &singleScenarioExecutor{
			scenario: &entity.TestScenario{
				ID: 1,
				TestCategory: &entity.TestCategory{
					ID:   7,
					Name: entity.STRESS,
				},
				MaxTestServiceCount: new(int64(1)),
				NumSteps:            int64(1),
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
			},
			agents:                   []interfaces.TestAgentController{agent},
			allAgentsHealthy:         true,
			running:                  false,
			allAgentsReadyForTesting: true,
		}

		agent.On("Healthy").Return(true)
		agent.On("ReadyForTesting").Return(true)

		err := ex.Execute(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrScenarioIsNotRunning)
		agent.AssertCalled(t, "Healthy")
		agent.AssertCalled(t, "ReadyForTesting")
	})
	t.Run("success_case", func(t *testing.T) {
		agent := new(mocks.MockTestAgentController)
		ctx := context.Background()

		ss := &singleScenarioExecutor{
			scenario: &entity.TestScenario{
				ID: 1,
				TestCategory: &entity.TestCategory{
					ID:   7,
					Name: entity.STRESS,
				},
				MaxTestServiceCount: new(int64(1)),
				NumSteps:            int64(1),
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
			},
			agents:                   []interfaces.TestAgentController{agent},
			allAgentsHealthy:         true,
			running:                  true,
			allAgentsReadyForTesting: true,
		}

		req := *request.NewRunRequestFromTestServiceConfig(
			1,
			ss.executionID,
			ss.scenario.TestServiceConfig,
		)

		agent.On("Healthy").Return(true)
		agent.On("ReadyForTesting").Return(true)
		agent.On("StartTesting", ctx, req).Return(nil)

		err := ss.Execute(ctx)

		assert.NoError(t, err)
		agent.AssertCalled(t, "Healthy")
		agent.AssertCalled(t, "ReadyForTesting")
		agent.AssertCalled(t, "StartTesting", ctx, req)
	})
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

		ss := &singleScenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent},
			allAgentsHealthy: false,
			running:          false,
		}

		ss.awaitAgentsToBeHealthy()

		agent.AssertCalled(t, "Healthy")

		assert.True(t, ss.allAgentsHealthy)
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

		ss := &singleScenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: false,
			running:          false,
		}

		healthyCheckSleep = time.Millisecond * 10
		agent1.On("Healthy").Return(true) // may be called multiple times in the health-check loop
		agent2.On("Healthy").Return(true)

		go ss.awaitAgentsToBeHealthy()
		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)
		assert.True(t, ss.allAgentsHealthy)

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

		ss := &singleScenarioExecutor{
			scenario:                 scenario,
			executionID:              uuid.New(),
			agents:                   []interfaces.TestAgentController{agent},
			allAgentsHealthy:         false,
			allAgentsReadyForTesting: false,
			running:                  false,
		}

		ss.awaitAgentsToBeReadyToStartTesting()

		agent.AssertCalled(t, "ReadyForTesting")

		assert.True(t, ss.allAgentsReadyForTesting)
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

		ss := &singleScenarioExecutor{
			scenario:                 scenario,
			executionID:              uuid.New(),
			agents:                   []interfaces.TestAgentController{agent},
			allAgentsHealthy:         false,
			allAgentsReadyForTesting: false,
			running:                  false,
		}

		go ss.awaitAgentsToBeReadyToStartTesting()
		time.Sleep(time.Millisecond)

		agent.AssertCalled(t, "ReadyForTesting")

		assert.False(t, ss.allAgentsReadyForTesting)

		agent.On("ReadyForTesting").Times(1).Return(true)
		time.Sleep(readyForTestingCheckSleep)

		time.Sleep(time.Millisecond)

		agent.AssertCalled(t, "ReadyForTesting")
		assert.True(t, ss.allAgentsReadyForTesting)
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

		ss := &singleScenarioExecutor{scenario: scenario,
			executionID:              uuid.New(),
			agents:                   []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy:         false,
			allAgentsReadyForTesting: false,
			running:                  false,
		}

		ss.awaitAgentsToBeReadyToStartTesting()

		agent1.AssertCalled(t, "ReadyForTesting")
		agent2.AssertCalled(t, "ReadyForTesting")

		assert.True(t, ss.allAgentsReadyForTesting)
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

		ss := &singleScenarioExecutor{scenario: scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: false,
			running:          false,
		}

		healthyCheckSleep = time.Millisecond * 10

		err := ss.executeScenarioSteps(ctx)
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

		ss := &singleScenarioExecutor{scenario: scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: true,
			running:          true,
		}

		healthyCheckSleep = time.Millisecond * 10

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		err := ss.executeScenarioSteps(ctx)
		assert.NoError(t, err)

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
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

		ss := &singleScenarioExecutor{scenario: scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: true,
			running:          true,
		}

		healthyCheckSleep = time.Millisecond * 10

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Times(2).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Times(2).Return(nil)

		agent1.On("ReadyForTesting").Times(1).Return(true)
		agent2.On("ReadyForTesting").Times(1).Return(true)

		err := ss.executeScenarioSteps(ctx)
		assert.NoError(t, err)

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)

		agent1.AssertCalled(t, "ReadyForTesting")
		agent2.AssertCalled(t, "ReadyForTesting")
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
		ss := &singleScenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1},
			allAgentsHealthy: true,
			running:          true,
		}

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		agent1.On("StartTesting", mock.Anything, mock.Anything).Return(nil)
		err := ss.executeScenarioSteps(ctx)
		assert.NoError(t, err)
		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)

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
