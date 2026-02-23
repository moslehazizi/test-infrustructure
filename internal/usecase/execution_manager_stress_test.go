package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/mocks"
	"control-panel-service/pkg"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

		builder := new(mocks.MockTestAgentControllerBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Run").Return(nil).Times(1)
		agent.On("Healthy").Times(1).Return(false)

		builder.On("Build", scenario).Times(1).Return(agent)

		ex := NewStressTestExecutionManager(builder)

		exeID := uuid.New()
		err := ex.AddScenario(context.Background(), scenario, exeID)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.NoError(t, err)

		agent.AssertCalled(t, "Run")
		builder.AssertCalled(t, "Build", scenario)

		// now, we should watch on manager Run function
		go ex.Run()

		stEx := ex.(*StressTestExecutionManager)
		assert.False(t, stEx.scenarios[1].allAgentsHealthy)

		agent.On("Healthy").Return(true)

		time.Sleep(time.Second + time.Millisecond*100)

		assert.True(t, stEx.scenarios[1].allAgentsHealthy)
	})

	t.Run("all agents' test services are healthy - 2 agents", func(t *testing.T) {
		t.Parallel()

		scenario := &entity.TestScenario{
			ID: 1,
			TestCategory: &entity.TestCategory{
				ID:   7,
				Name: entity.STRESS,
			},
			MaxTestServiceCount: new(int64(2)),
		}

		builder := new(mocks.MockTestAgentControllerBuilder)

		agent1 := new(mocks.MockTestAgentController)
		agent1.On("Run").Return(nil).Times(1)
		agent1.On("Healthy").Times(1).Return(false)

		agent2 := new(mocks.MockTestAgentController)
		agent2.On("Run").Return(nil).Times(1)
		agent2.On("Healthy").Times(1).Return(false)

		builder.On("Build", scenario).Times(1).Return(agent1)
		builder.On("Build", scenario).Times(1).Return(agent2)

		ex := NewStressTestExecutionManager(builder)

		exeID := uuid.New()
		err := ex.AddScenario(context.Background(), scenario, exeID)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.NoError(t, err)

		agent1.AssertCalled(t, "Run")
		agent2.AssertCalled(t, "Run")

		builder.AssertCalled(t, "Build", scenario)

		checkLoopSleep = time.Millisecond * 10

		// now, we should watch on manager Run function
		go ex.Run()

		stEx := ex.(*StressTestExecutionManager)
		assert.False(t, stEx.scenarios[1].allAgentsHealthy)

		agent1.On("Healthy").Return(true)

		time.Sleep(time.Millisecond * 12)

		// agent2 is still not ready
		assert.False(t, stEx.scenarios[1].allAgentsHealthy)

		agent2.On("Healthy").Return(true)
		time.Sleep(time.Millisecond * 12)

		assert.True(t, stEx.scenarios[1].allAgentsHealthy)
	})
}
