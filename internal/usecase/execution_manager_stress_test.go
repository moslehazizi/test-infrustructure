package usecase_test

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase"
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

		ex := usecase.NewStressTestExecutionManager(builder)

		exeID := uuid.New()
		err := ex.AddScenario(context.Background(), &scenario, exeID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, pkg.ErrMaxTestServiceCountNotSet)
	})

	t.Run("success case", func(t *testing.T) {
		scenario := entity.TestScenario{
			MaxTestServiceCount: new(int64(3)),
		}

		builder := new(mocks.MockTestAgentControllerBuilder)

		agent := new(mocks.MockTestAgentController)
		agent.On("Run").Times(3)
		builder.On("Build").Times(3).Return(agent)

		ex := usecase.NewStressTestExecutionManager(builder)

		exeID := uuid.New()
		err := ex.AddScenario(context.Background(), &scenario, exeID)

		// will wait to all goroutines be called.
		time.Sleep(time.Millisecond)

		assert.NoError(t, err)

		agent.AssertCalled(t, "Run")
		builder.AssertCalled(t, "Build")
	})
}
