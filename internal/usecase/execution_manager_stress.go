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

var checkLoopSleep = time.Second

// See: https://github.com/farbodan/challenge-control-panel-service/blob/main/internal/usecase/test_scenario_runner.md.
func NewStressTestExecutionManager(testAgentControllerToolBox interfaces.TestAgentControllerToolBox, scenarioRepo repository.TestScenarioRepository) interfaces.ExecutionManager {
	return &StressTestExecutionManager{
		testAgentControllerToolBox: testAgentControllerToolBox,
		scenarios:                  make(map[uint64]interfaces.ScenarioExecutor),
		scenarioRepo:               scenarioRepo,
		running:                    true,
	}
}

// #region StressTestExecutionManager
type StressTestExecutionManager struct {
	running                    bool
	scenarios                  map[uint64]interfaces.ScenarioExecutor
	testAgentControllerToolBox interfaces.TestAgentControllerToolBox
	scenarioRepo               repository.TestScenarioRepository
	mx                         sync.Mutex
}

func (ex *StressTestExecutionManager) Run() {
	ctx := context.Background()
	for ex.running {
		// check all scenarios
		// for each scenario, make sure the executor is running.
		for _, sc := range ex.scenarios {
			go func() {
				_ = sc.Run(ctx)
			}()
		}
		time.Sleep(checkLoopSleep)
	}

	zap.L().Info("exiting from StressTestExecutionManager.Run")
}

// AddScenario
// @Deprecated no longer needed.
func (ex *StressTestExecutionManager) AddScenario(ctx context.Context, scenario *entity.TestScenario) error {
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
		agent := ex.testAgentControllerToolBox.Build(scenario)
		go func() {
			_ = agent.Run()
		}()

		ex.mx.Lock()
		_, ok := ex.scenarios[scenario.ID]
		if !ok {
			ex.scenarios[scenario.ID] = &scenarioExecutor{
				scenario:     scenario,
				agents:       []interfaces.TestAgentController{agent},
				running:      false,
				scenarioRepo: ex.scenarioRepo,
			}
		} else {
			ex.scenarios[scenario.ID].AddAgent(agent)
		}

		ex.mx.Unlock()
	}

	return nil
}

func (ex *StressTestExecutionManager) RunScenario(ctx context.Context, scenario *entity.TestScenario) error {
	tracer := otel.Tracer("StressTestExecutionManager")
	_, span := tracer.Start(ctx, "AddScenario")
	defer span.End()

	ex.mx.Lock()
	defer ex.mx.Unlock()
	testScenario, ok := ex.scenarios[scenario.ID]
	if ok {
		// set scenario as running
		testScenario.SetRunning(true)

		return nil
	}

	ex.scenarios[scenario.ID] = &scenarioExecutor{
		scenario:     scenario,
		agents:       []interfaces.TestAgentController{},
		running:      true,
		scenarioRepo: ex.scenarioRepo,
	}

	return nil
}

// func (ex *StressTestExecutionManager) RunScenario(ctx context.Context, scenario *entity.TestScenario) error {
// 	tracer := otel.Tracer("StressTestExecutionManager")
// 	_, span := tracer.Start(ctx, "RunScenario")
// 	defer span.End()

// 	if scenario == nil {
// 		zap.L().Error("test scenario is nil")

// 		return pkg.ErrTestScenarioServiceIsNil
// 	}

// 	ex.mx.Lock()
// 	defer ex.mx.Unlock()

// 	testScenario, ok := ex.scenarios[scenario.ID]
// 	if !ok {
// 		zap.L().Error("test scenario not found")

// 		return fmt.Errorf("%w", pkg.ErrTestScenarioNotFound)
// 	}

// 	testScenario.SetExecutionID(uuid.New())
// 	testScenario.SetRunning(true)

// 	return nil
// }

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
