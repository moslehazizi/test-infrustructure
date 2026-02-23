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
)

func Test_scenarioExecutor_AddAgent(t *testing.T) {
	ex := new(scenarioExecutor)
	assert.Len(t, ex.agents, 0)

	ex.AddAgent(&mocks.MockTestAgentController{})
	assert.Len(t, ex.agents, 1)
}
func Test_scenarioExecutor_Run(t *testing.T) {
	t.Run("all agents' test services are healthy - single agent", func(t *testing.T) {
		t.Parallel()

		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		}

		agent := new(mocks.MockTestAgentController)
		agent.On("Healthy").Times(1).Return(true)

		ex := &scenarioExecutor{
			scenarioID:       scenario.ID,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent},
			allAgentsHealthy: false,
			running:          false,
		}

		ex.Run()

		agent.AssertCalled(t, "Healthy")

		assert.True(t, ex.running)
		assert.True(t, ex.allAgentsHealthy)

	})

	t.Run("all agents' test services are healthy - 2 agents", func(t *testing.T) {
		t.Parallel()

		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(1)),
		}

		agent1 := new(mocks.MockTestAgentController)
		agent2 := new(mocks.MockTestAgentController)

		ex := &scenarioExecutor{
			scenarioID:       scenario.ID,
			executionID:      uuid.New(),
			agents:           []interfaces.TestAgentController{agent1, agent2},
			allAgentsHealthy: false,
			running:          false,
		}

		healthyCheckSleep = time.Millisecond * 10
		agent1.On("Healthy").Times(1).Return(true)
		agent2.On("Healthy").Times(1).Return(true)

		go ex.Run()
		time.Sleep(time.Millisecond)
		assert.True(t, ex.running)
		time.Sleep(time.Millisecond * 10)
		assert.True(t, ex.allAgentsHealthy)

		agent1.AssertCalled(t, "Healthy")
		agent2.AssertCalled(t, "Healthy")
	})
}

// #region StressTestExecutionManager
func TestStressTestExecutionManager_AddScenario(t *testing.T) {
	t.Run("failed case: scenario max service count is null", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: nil,
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
		}

		builder := new(mocks.MockTestAgentControllerBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Run").Return(nil).Times(1)
		agent.On("Healthy").Times(1).Return(false)

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
