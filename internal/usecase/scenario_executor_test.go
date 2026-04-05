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
		// t.Parallel()

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
		msse := new(mocks.MockSingleScenarioExecutor)
		msseb := new(mocks.MockSingleScenarioExecutorBuilder)
		se := &scenarioExecutor{
			scenario:                   scenario,
			executionID:                uuid.New(),
			agents:                     []interfaces.TestAgentController{},
			allAgentsHealthy:           false,
			running:                    true,
			scenarioRepo:               repo,
			testAgentControllerToolBox: ac,
			scenarioExecutorBuilder:    msseb,
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

		msseb.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(msse).
			Times(1)
		msse.On("Execute", mock.Anything).Return(nil).Times(1)

		err := se.Run(context.Background())
		assert.NoError(t, err)

		ac.AssertCalled(t, "Build", mock.Anything)
		ac.AssertExpectations(t)
		ac.AssertNumberOfCalls(t, "Build", 9) // 3 in stage zero + 6 in multi stage

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

		msseb.AssertCalled(t, "Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		msseb.AssertNumberOfCalls(t, "Build", 4) // 1 in stage zero + 3 in mult stage

		msse.AssertCalled(t, "Execute", mock.Anything)
		msse.AssertNumberOfCalls(t, "Execute", 4) // 1 in stage zero + 3 in multi stage
	})

	t.Run("1_step_6_dynamic_agents_no_factorial_increment_failed_execute", func(t *testing.T) {
		// t.Parallel()

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
		msse := new(mocks.MockSingleScenarioExecutor)
		msseb := new(mocks.MockSingleScenarioExecutorBuilder)
		se := &scenarioExecutor{
			scenario:                   scenario,
			executionID:                uuid.New(),
			agents:                     []interfaces.TestAgentController{},
			allAgentsHealthy:           false,
			running:                    true,
			scenarioRepo:               repo,
			testAgentControllerToolBox: ac,
			scenarioExecutorBuilder:    msseb,
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)
		agent3 := new(mocks.MockTestAgentController)
		ac.On("Build", mock.Anything).Times(1).Return(agent1)
		ac.On("Build", mock.Anything).Times(1).Return(agent2)
		ac.On("Build", mock.Anything).Times(1).Return(agent3)

		agent1.On("Run").Times(1).Return(nil)
		agent2.On("Run").Times(1).Return(nil)
		agent3.On("Run").Times(1).Return(nil)

		msseb.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(msse)
		msse.On("Execute", mock.Anything).Return(errors.New("something went wrong"))

		err := se.Run(context.Background())
		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrFailedToExecuteSingleScenario)
	})
}

//#endregion scenarioExecutor
