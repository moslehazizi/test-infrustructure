package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/pkg"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// What is ScenarioExecutorBox?
// ScenarioExecutorBox is a collection of scenarios which are executing in the background.
// Each scenario that is added to the BOX is of type ScenarioExecutor.
//
// The inMemoryScenarioExecutorBox is an in-memory implementation of ScenarioExecutorBox.
//
// The most important function of ScenarioExecutor is ResumeOrStart. This function
// is responsible to run scenario in background. For more details about this method, please see its doc.

// #region ScenarioExecutorBox

func NewInMemoryScenarioExecutorBox() interfaces.ScenarioExecutorBox {
	return &inMemoryScenarioExecutorBox{}
}

type inMemoryScenarioExecutorBox struct {
	runningScenarios sync.Map
}

func (e *inMemoryScenarioExecutorBox) Add(exe interfaces.ScenarioExecutor) {
	e.runningScenarios.Store(exe.GetID(), exe.GetScenario())
	go exe.ResumeOrStart(context.Background())
}

func (e *inMemoryScenarioExecutorBox) HasExecutor(id uint64) bool {
	_, ok := e.runningScenarios.Load(id)

	return ok
}

//#endregion ScenarioExecutorBox

//#region ScenarioExecutor

type scenarioExecutor struct {
	scenario            entity.TestScenario
	testServiceRepo     repository.TestServiceRepository
	provisioningService provider.ProvisioningService
}

func NewScenarioExecutor(
	scenario entity.TestScenario,
	testServiceRepo repository.TestServiceRepository,
	provisioningService provider.ProvisioningService,
) interfaces.ScenarioExecutor {
	return &scenarioExecutor{
		scenario,
		testServiceRepo,
		provisioningService,
	}
}

func (ex *scenarioExecutor) GetID() uint64 {
	return ex.scenario.ID
}

func (ex *scenarioExecutor) GetScenario() *entity.TestScenario {
	return &ex.scenario
}

func (ex *scenarioExecutor) ResumeOrStart(ctx context.Context) {
	// test scenario categorization based on how it should be executed:
	// A (total test service + total execution time): smoke, load, soak, spike.
	// B (step execution + increase rate): stress.
	// C: (total test service + step execution + increase rate): scalability.
	// D: (total test service + step execution + increase rate + decrease rate): recovery.
}

// runGroupA
// See: https://github.com/farbodan/challenge-control-panel-service/blob/main/internal/usecase/test_scenario_runner.md#group-a.
func (ex *scenarioExecutor) runGroupA(ctx context.Context) error {
	tracer := otel.Tracer("scenarioExecutor")
	_, span := tracer.Start(ctx, "runGroupA")
	defer span.End()

	if ex.scenario.MaxTestServiceCount == nil {
		zap.L().Error("max test service count value is null but required",
			zap.Uint64("scenarioID", ex.scenario.ID),
		)

		return pkg.ErrMaxTestServiceCountNotSet
	}

	cnt, err := ex.testServiceRepo.GetCountAllRunningByScenario(ctx, ex.scenario.ID)
	if err != nil {
		zap.L().Error("failed to get running test services count by scenario",
			zap.Uint64("scenarioID", ex.scenario.ID),
			zap.Error(err),
		)

		return fmt.Errorf("failed to get running test services count: %w", err)
	}

	remaining := *ex.scenario.MaxTestServiceCount - cnt
	// no more test service to provision and we are done here.
	if remaining > 0 {
		err = ex.provisioningService.ProvisionTestService(ctx, ex.scenario.TestServiceConfig, remaining)
		if err != nil {
			zap.L().Error("failed to provision remaining test services",
				zap.Uint64("scenarioID", ex.scenario.ID),
				zap.Error(err),
			)

			return fmt.Errorf("failed to provision remaining test services: %w", err)
		}
	}

	zap.L().Debug("all test services already provisioned for scenario",
		zap.Uint64("scenarioID", ex.scenario.ID),
	)

	// if ex.scenario.StartedAt
	// ex.provisioningService.DeprovisionTestService(ctx, )
	return nil
}

//#endregion ScenarioExecutor
