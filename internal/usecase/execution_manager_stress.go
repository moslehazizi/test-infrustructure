package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/pkg"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var checkLoopSleep = time.Second

// See: https://github.com/farbodan/challenge-control-panel-service/blob/main/internal/usecase/test_scenario_runner.md.
func NewStressTestExecutionManager(testAgentControllerBuilder interfaces.TestAgentControllerBuilder) interfaces.ExecutionManager {
	mng := &StressTestExecutionManager{
		testAgentControllerBuilder: testAgentControllerBuilder,
	}
	mng.scenarios = make(map[uint64]*scenarioExecution)
	mng.running = true

	return mng
}

type scenarioExecution struct {
	scenarioID       uint64
	executionID      uuid.UUID
	agents           []interfaces.TestAgentController
	allAgentsHealthy bool
}

type StressTestExecutionManager struct {
	running                    bool
	scenarios                  map[uint64]*scenarioExecution
	testAgentControllerBuilder interfaces.TestAgentControllerBuilder
	mx                         sync.Mutex
}

func (ex *StressTestExecutionManager) Run() {
	for ex.running {
		// check all scenarios
		// for each scenario, all test services should be healthy.
		for k, sc := range ex.scenarios {
			_ = k

			allHealthy := true
			for i, agent := range sc.agents {
				_ = i
				h := agent.Healthy()
				if !h {
					allHealthy = false

					break
				}
			}

			sc.allAgentsHealthy = allHealthy
		}

		time.Sleep(checkLoopSleep)
	}

	zap.L().Info("exiting from StressTestExecutionManager.Run")
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
		_, ok := ex.scenarios[scenario.ID]
		if !ok {
			ex.scenarios[scenario.ID] = &scenarioExecution{
				scenarioID:  scenario.ID,
				executionID: executionID,
				agents:      []interfaces.TestAgentController{agent},
			}
		} else {
			ex.scenarios[scenario.ID].agents = append(ex.scenarios[scenario.ID].agents, agent)
		}

		ex.mx.Unlock()
	}

	return nil
}
