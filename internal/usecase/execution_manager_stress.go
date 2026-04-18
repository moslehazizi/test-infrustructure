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
		scenarios:                  make(map[uint64]interfaces.ScenarioExecutor),
		scenarioRepo:               scenarioRepo,
		scenarioExecutorBuilder:    scenarioExecutorBuilder,
	}
}

// #region StressTestExecutionManager
type StressTestExecutionManager struct {
	scenarios                  map[uint64]interfaces.ScenarioExecutor
	testAgentControllerToolBox interfaces.TestAgentControllerToolBox
	scenarioRepo               repository.TestScenarioRepository
	mx                         sync.Mutex
	scenarioExecutorBuilder    interfaces.ScenarioExecutorBuilder
}

func (ex *StressTestExecutionManager) RunScenario(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("StressTestExecutionManager")
	_, span := tracer.Start(ctx, "AddScenario")
	defer span.End()

	ex.mx.Lock()
	defer ex.mx.Unlock()
	ex.scenarios[scenario.ID] = ex.scenarioExecutorBuilder.Build(scenario, ex.scenarioRepo, ex.testAgentControllerToolBox, &singleScenarioExecutorBuilder{})
	ex.scenarios[scenario.ID].SetRunning(true)

	//nolint
	go func() {
		_ = ex.scenarios[scenario.ID].Run(context.Background())
		// delete scenario from memory
		delete(ex.scenarios, scenario.ID)
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

	ex.mx.Lock()
	testScenario, ok := ex.scenarios[scenario.ID]
	if !ok {
		zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
		ex.mx.Unlock()

		return pkg.ErrTestScenarioNotFound
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
			ex.mx.Unlock()

			return fmt.Errorf("%w - %w", pkg.ErrFailedToPauseTestService, err)
		}
	}

	ex.mx.Unlock()

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

	ex.mx.Lock()
	testScenario, ok := ex.scenarios[scenario.ID]
	if !ok {
		zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
		ex.mx.Unlock()

		return fmt.Errorf("%w", pkg.ErrTestScenarioNotFound)
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
			ex.mx.Unlock()

			return fmt.Errorf("%w - %w", pkg.ErrFailedToResumeTestService, err)
		}
	}

	ex.mx.Unlock()

	return nil
}

func (ex *StressTestExecutionManager) StopScenario(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("StressTestExecutionManager")
	_, span := tracer.Start(ctx, "StopScenario")
	defer span.End()

	if scenario.MaxTestServiceCount == nil {
		zap.L().Error("max test service count value is null but required", zap.Uint64("scenarioID", scenario.ID))

		return pkg.ErrMaxTestServiceCountNotSet
	}

	ex.mx.Lock()
	sci, ok := ex.scenarios[scenario.ID]
	if !ok {
		zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
		ex.mx.Unlock()

		return pkg.ErrTestScenarioNotFound
	}
	sci.SetRunning(false)

	agents := sci.GetAgents()
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
		err := agent.StopTesting(ctx)
		if err != nil {
			zap.L().Error("test agent controller couldn't stop test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))
			ex.mx.Unlock()

			return fmt.Errorf("%w - %w", pkg.ErrFailedToStopTestService, err)
		}
	}

	ex.mx.Unlock()

	return nil
}

func (ex *StressTestExecutionManager) AbortScenario(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("StressTestExecutionManager")
	_, span := tracer.Start(ctx, "AbortScenario")
	defer span.End()

	if scenario.MaxTestServiceCount == nil {
		zap.L().Error("max test service count value is null but required", zap.Uint64("scenarioID", scenario.ID))

		return pkg.ErrMaxTestServiceCountNotSet
	}

	ex.mx.Lock()
	sci, ok := ex.scenarios[scenario.ID]
	if !ok {
		zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
		ex.mx.Unlock()

		return pkg.ErrTestScenarioNotFound
	}
	sci.SetRunning(false)

	agents := sci.GetAgents()
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
		err := agent.AbortTesting(ctx)
		if err != nil {
			zap.L().Error("test agent controller couldn't abort test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))
			ex.mx.Unlock()

			return fmt.Errorf("%w - %w", pkg.ErrFailedToAbortTestService, err)
		}
	}

	ex.mx.Unlock()

	return nil
}

//#endregion StressTestExecutionManager
