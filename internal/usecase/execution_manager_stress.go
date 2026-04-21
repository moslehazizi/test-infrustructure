package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/pkg"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// See: https://github.com/farbodan/challenge-control-panel-service/blob/main/internal/usecase/test_scenario_runner.md.
func NewStressTestExecutionManager(testAgentControllerToolBox interfaces.TestAgentControllerToolBox, scenarioRepo repository.TestScenarioRepository, scenarioExecutorBuilder interfaces.ScenarioExecutorBuilder) interfaces.ExecutionManager {
	return &StressTestExecutionManager{
		testAgentControllerToolBox: testAgentControllerToolBox,
		scenarios:                  sync.Map{},
		scenarioRepo:               scenarioRepo,
		scenarioExecutorBuilder:    scenarioExecutorBuilder,
	}
}

// #region StressTestExecutionManager
type StressTestExecutionManager struct {
	// scenarios                  map[uint64]interfaces.ScenarioExecutor
	// mx                         sync.Mutex
	scenarios                  sync.Map
	testAgentControllerToolBox interfaces.TestAgentControllerToolBox
	scenarioRepo               repository.TestScenarioRepository
	scenarioExecutorBuilder    interfaces.ScenarioExecutorBuilder
}

func (ex *StressTestExecutionManager) RunScenario(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("StressTestExecutionManager")
	_, span := tracer.Start(ctx, "AddScenario")
	defer span.End()

	// scenario also running
	rawScenarioExec, ok := ex.scenarios.Load(scenario.ID)
	if ok && rawScenarioExec != nil {
		scenarioExec, _ := rawScenarioExec.(interfaces.ScenarioExecutor)
		if scenarioExec.IsRunning() {
			zap.L().Error("scenario is also running", zap.Uint64("scenarioID", scenario.ID))

			return pkg.ErrScenarioIsRunning
		}
	}

	scenarioExecutor := ex.scenarioExecutorBuilder.Build(scenario, ex.scenarioRepo, ex.testAgentControllerToolBox, &singleScenarioExecutorBuilder{})
	scenarioExecutor.SetRunning(true)
	ex.scenarios.Store(scenario.ID, scenarioExecutor)

	//nolint
	go func() {
		_ = scenarioExecutor.Run(context.Background())
		// delete scenario from memory
		ex.scenarios.Delete(scenario.ID)
	}()

	return nil
}

func (ex *StressTestExecutionManager) PauseScenario(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("StressTestExecutionManager")
	_, span := tracer.Start(ctx, "PauseScenario")
	defer span.End()

	if scenario.MaxTestServiceCount == nil {
		zap.L().Error("max test service count value is null but required", zap.Uint64("scenarioID", scenario.ID))

		return pkg.ErrMaxTestServiceCountNotSet
	}

	rawTestScenario, ok := ex.scenarios.Load(scenario.ID)
	if !ok {
		zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))

		return pkg.ErrTestScenarioNotFound
	}

	testScenario, ok := rawTestScenario.(interfaces.ScenarioExecutor)
	if !ok {
		zap.L().Error("failed to type assertion", zap.Error(pkg.ErrTypeAssertionAnyToScenarioExecutor), zap.Uint64("scenarioID", scenario.ID))

		return pkg.ErrTypeAssertionAnyToScenarioExecutor
	}

	agents := testScenario.GetAgents()

	// We add health check agent so that if user immediately click on pause right after run first wait to all agents be ready.
	for {
		allHealthy := true
		for _, agent := range agents {
			if !agent.Healthy() {
				allHealthy = false

				break
			}
		}

		if allHealthy {
			break
		}

		time.Sleep(healthyCheckSleep)
	}

	for _, agent := range agents {
		err := agent.PauseTesting(ctx)
		if err != nil {
			zap.L().Error("test agent controller can not pause test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))

			return fmt.Errorf("%w - %w", pkg.ErrFailedToPauseTestService, err)
		}
	}

	return nil
}

func (ex *StressTestExecutionManager) ResumeScenario(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("StressTestExecutionManager")
	_, span := tracer.Start(ctx, "ResumeScenario")
	defer span.End()

	if scenario.MaxTestServiceCount == nil {
		zap.L().Error("max test service count value is null but required", zap.Uint64("scenarioID", scenario.ID))

		return pkg.ErrMaxTestServiceCountNotSet
	}

	rawTestScenario, ok := ex.scenarios.Load(scenario.ID)
	if !ok {
		zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))

		return pkg.ErrTestScenarioNotFound
	}

	testScenario, ok := rawTestScenario.(interfaces.ScenarioExecutor)
	if !ok {
		zap.L().Error("failed to type assertion", zap.Error(pkg.ErrTypeAssertionAnyToScenarioExecutor), zap.Uint64("scenarioID", scenario.ID))

		return pkg.ErrTypeAssertionAnyToScenarioExecutor
	}

	agents := testScenario.GetAgents()

	for {
		allHealthy := true
		for _, agent := range agents {
			if !agent.Healthy() {
				allHealthy = false

				break
			}
		}

		if allHealthy {
			break
		}

		time.Sleep(healthyCheckSleep)
	}

	for _, agent := range agents {
		err := agent.ResumeTesting(ctx)
		if err != nil {
			zap.L().Error("test agent controller can not resume test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))

			return fmt.Errorf("%w - %w", pkg.ErrFailedToResumeTestService, err)
		}
	}

	return nil
}

func (ex *StressTestExecutionManager) StopScenario(ctx context.Context, scenario *entity.TestScenario) error {
	return nil

	// tracer := otel.Tracer("StressTestExecutionManager")
	// _, span := tracer.Start(ctx, "StopScenario")
	// defer span.End()

	// if scenario.MaxTestServiceCount == nil {
	// 	zap.L().Error("max test service count value is null but required", zap.Uint64("scenarioID", scenario.ID))

	// 	return pkg.ErrMaxTestServiceCountNotSet
	// }

	// ex.mx.Lock()
	// sci, ok := ex.scenarios[scenario.ID]
	// if !ok {
	// 	zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
	// 	ex.mx.Unlock()

	// 	return pkg.ErrTestScenarioNotFound
	// }
	// sci.SetRunning(false)

	// agents := sci.GetAgents()
	// for {
	// 	allHealthy := true
	// 	for _, agent := range agents {
	// 		if !agent.Healthy() {
	// 			allHealthy = false

	// 			break
	// 		}
	// 	}

	// 	if allHealthy {
	// 		break
	// 	}

	// 	time.Sleep(healthyCheckSleep)
	// }

	// for _, agent := range agents {
	// 	err := agent.StopTesting(ctx)
	// 	if err != nil {
	// 		zap.L().Error("test agent controller couldn't stop test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))
	// 		ex.mx.Unlock()

	// 		return fmt.Errorf("%w - %w", pkg.ErrFailedToStopTestService, err)
	// 	}
	// }

	// ex.mx.Unlock()

	// return nil
}

func (ex *StressTestExecutionManager) DeleteScenario(ctx context.Context, scenario *entity.TestScenario) error {
	return nil

	// tracer := otel.Tracer("StressTestExecutionManager")
	// _, span := tracer.Start(ctx, "DeleteScenario")
	// defer span.End()

	// if scenario.MaxTestServiceCount == nil {
	// 	zap.L().Error("max test service count value is null but required", zap.Uint64("scenarioID", scenario.ID))

	// 	return pkg.ErrMaxTestServiceCountNotSet
	// }

	// ex.mx.Lock()
	// sci, ok := ex.scenarios[scenario.ID]
	// if !ok {
	// 	zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
	// 	ex.mx.Unlock()

	// 	return pkg.ErrTestScenarioNotFound
	// }
	// sci.SetRunning(false)

	// agents := sci.GetAgents()
	// for _, agent := range agents {
	// 	err := agent.DeleteTesting(ctx)
	// 	if err != nil {
	// 		zap.L().Error("test agent controller couldn't delete  test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))
	// 		ex.mx.Unlock()

	// 		return fmt.Errorf("%w - %w", pkg.ErrFailedToDeleteTestService, err)
	// 	}
	// }

	// ex.mx.Unlock()

	// return nil
}

//#endregion StressTestExecutionManager
