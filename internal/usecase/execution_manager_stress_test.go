package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #region scenarioExecutor
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

	ex.running = true
	assert.True(t, ex.IsRunning())
}

func Test_scenarioExecutor_awaitAgentsToBeHealthy(t *testing.T) {
	t.Run("all agents' test services are healthy - single agent", func(t *testing.T) {
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

	t.Run("all agents' test services are healthy - 2 agents", func(t *testing.T) {
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
	t.Run("all agents' test services are ready - single agent", func(t *testing.T) {
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

	t.Run("all agents' test services are ready - single agent - first time is not ready", func(t *testing.T) {
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

	t.Run("all agents' test services are ready - two agents", func(t *testing.T) {
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
	t.Run("success - single step", func(t *testing.T) {
		t.Parallel()

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

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: false,
			running:          false,
		}

		healthyCheckSleep = time.Millisecond * 10
		agent1.On("Healthy").Return(true)
		agent2.On("Healthy").Return(true)

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Times(1).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Times(1).Return(nil)
		err := ex.Run()

		agent1.AssertCalled(t, "Healthy")
		agent2.AssertCalled(t, "Healthy")

		assert.True(t, ex.allAgentsHealthy)

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)

		assert.NoError(t, err)
	})

	t.Run("success - 2 steps", func(t *testing.T) {
		t.Parallel()

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

		ex := &scenarioExecutor{
			scenario:         scenario,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: false,
			running:          false,
		}

		healthyCheckSleep = time.Millisecond * 10
		agent1.On("Healthy").Return(true)
		agent2.On("Healthy").Return(true)

		time.Sleep(time.Millisecond)
		time.Sleep(time.Millisecond * 10)

		// we have 2 agents ready

		agent1.On("StartTesting", mock.Anything, mock.Anything).Times(2).Return(nil)
		agent2.On("StartTesting", mock.Anything, mock.Anything).Times(2).Return(nil)

		agent1.On("ReadyForTesting").Times(1).Return(true)
		agent2.On("ReadyForTesting").Times(1).Return(true)

		err := ex.Run()

		agent1.AssertCalled(t, "Healthy")
		agent2.AssertCalled(t, "Healthy")

		agent1.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)
		agent2.AssertCalled(t, "StartTesting", mock.Anything, mock.Anything)

		agent1.AssertCalled(t, "ReadyForTesting")
		agent2.AssertCalled(t, "ReadyForTesting")

		assert.NoError(t, err)
	})
}

//#endregion scenarioExecutor

// #region StressTestExecutionManager
func TestStressTestExecutionManager_AddScenario(t *testing.T) {
	t.Run("failed case: scenario max service count is null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerBuilder)

		ex := NewStressTestExecutionManager(builder)

		exeID := uuid.New()
		err := ex.AddScenario(context.Background(), &scenario, exeID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("success case", func(t *testing.T) {
		scenario := &entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
			NumSteps:            2,
		}

		builder := new(mocks.MockTestAgentControllerBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Run").Return(nil).Times(3)
		builder.On("Build", scenario).Times(3).Return(agent)

		ex := NewStressTestExecutionManager(builder)

		exeID := uuid.New()
		err := ex.AddScenario(context.Background(), scenario, exeID)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.NoError(t, err)

		agent.AssertCalled(t, "Run")
		builder.AssertCalled(t, "Build", scenario)
	})
}

func TestStressTestExecutionManager_Run(t *testing.T) {
	t.Run("all executors are running", func(t *testing.T) {
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

		builder := new(mocks.MockTestAgentControllerBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Run").Return(nil).Times(1)
		agent.On("Healthy").Return(false)

		builder.On("Build", scenario).Times(1).Return(agent)

		ex := NewStressTestExecutionManager(builder)

		exeID := uuid.New()
		err := ex.AddScenario(context.Background(), scenario, exeID)
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
		assert.True(t, stEx.scenarios[1].IsRunning())

		// make sure if already is running, it hit continue.
		time.Sleep(time.Millisecond * 11)
		assert.True(t, stEx.scenarios[1].IsRunning())
	})
}

// #endregion StressTestExecutionManager
