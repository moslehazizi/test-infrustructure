package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/pkg"
	"sync"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// See: https://github.com/farbodan/challenge-control-panel-service/blob/main/internal/usecase/test_scenario_runner.md.
func NewStressTestExecutionManager(testAgentControllerBuilder interfaces.TestAgentControllerBuilder) interfaces.ExecutionManager {
	mng := &StressTestExecutionManager{
		testAgentControllerBuilder: testAgentControllerBuilder,
	}
	mng.scenarios = make(map[uint64][]interfaces.TestAgentController)

	return mng
}

type StressTestExecutionManager struct {
	scenarios                  map[uint64][]interfaces.TestAgentController
	testAgentControllerBuilder interfaces.TestAgentControllerBuilder
	mx                         sync.Mutex
}

func (ex *StressTestExecutionManager) Run() {

}

func (ex *StressTestExecutionManager) AddScenario(ctx context.Context, scenario *entity.TestScenario, executionID uuid.UUID) error {
	tracer := otel.Tracer("StressTestExecutionManager")
	_, span := tracer.Start(ctx, "AddScenario")
	defer span.End()

	if scenario.MaxTestServiceCount == nil {
		zap.L().Error("max test service count value is null but required",
			zap.Uint64("scenarioID", scenario.ID),
		)

		return pkg.ErrMaxTestServiceCountNotSet
	}

	for i := int64(1); i <= *scenario.MaxTestServiceCount; i++ {
		agent := ex.testAgentControllerBuilder.Build()
		go agent.Run()

		ex.mx.Lock()
		ex.scenarios[scenario.ID] = append(ex.scenarios[scenario.ID], agent)
		ex.mx.Unlock()
	}

	return nil
}
