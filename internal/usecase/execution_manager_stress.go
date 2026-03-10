package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/pkg"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var checkLoopSleep = time.Second
var healthyCheckSleep = time.Second
var readyForTestingCheckSleep = time.Second

// See: https://github.com/farbodan/challenge-control-panel-service/blob/main/internal/usecase/test_scenario_runner.md.
func NewStressTestExecutionManager(testAgentControllerToolBox interfaces.TestAgentControllerToolBox) interfaces.ExecutionManager {
	mng := &StressTestExecutionManager{
		testAgentControllerToolBox: testAgentControllerToolBox,
	}
	mng.scenarios = make(map[uint64]interfaces.ScenarioExecutor)
	mng.running = true

	return mng
}

type scenarioExecutor struct {
	scenario                 *entity.TestScenario
	executionID              uuid.UUID
	agents                   []interfaces.TestAgentController
	allAgentsHealthy         bool
	allAgentsReadyForTesting bool
	running                  bool
}

func (sc *scenarioExecutor) IsRunning() bool {
	return sc.running
}

func (sc *scenarioExecutor) SetRunning(status bool) {
	sc.running = status
}

func (sc *scenarioExecutor) AddAgent(agent interfaces.TestAgentController) {
	sc.agents = append(sc.agents, agent)
}

func (sc *scenarioExecutor) AllAgentsAreHealthy() bool {
	return sc.allAgentsHealthy
}

func (sc *scenarioExecutor) Run() error {
	sc.running = false

	// Make sure all agents are healthy.
	sc.awaitAgentsToBeHealthy()

	// NOW: all agents are healthy.

	// we can send scheduled commands.

	// we need a loop based on len of steps:
	for i := int64(1); i <= sc.scenario.NumSteps; i++ {
		req := request.NewRunRequestFromTestServiceConfig(
			int(i),
			sc.executionID,
			sc.scenario.TestServiceConfig,
		)

		// send command to all agents.
		wg := sync.WaitGroup{}
		for _, agent := range sc.agents {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = agent.StartTesting(context.Background(), *req)
			}()
		}

		wg.Wait()

		if i < sc.scenario.NumSteps {
			// wait for all agents to be ready to execute next step.
			sc.awaitAgentsToBeReadyToStartTesting()
		}
	}

	return nil
}

func (sc *scenarioExecutor) awaitAgentsToBeHealthy() {
	for {
		allHealthy := true
		for _, agent := range sc.agents {
			if !agent.Healthy() {
				allHealthy = false

				break
			}
		}

		sc.allAgentsHealthy = allHealthy

		if allHealthy {
			break
		}

		time.Sleep(healthyCheckSleep)
	}
}

func (sc *scenarioExecutor) awaitAgentsToBeReadyToStartTesting() {
	for {
		allReady := true
		for _, agent := range sc.agents {
			if !agent.ReadyForTesting() {
				allReady = false

				break
			}
		}

		sc.allAgentsReadyForTesting = allReady

		if allReady {
			break
		}

		time.Sleep(readyForTestingCheckSleep)
	}
}

// #region StressTestExecutionManager
type StressTestExecutionManager struct {
	running                    bool
	scenarios                  map[uint64]interfaces.ScenarioExecutor
	testAgentControllerToolBox interfaces.TestAgentControllerToolBox
	mx                         sync.Mutex
}

func (ex *StressTestExecutionManager) Run() {
	for ex.running {
		// check all scenarios
		// for each scenario, make sure the executor is running.
		for _, sc := range ex.scenarios {
			if !sc.IsRunning() {
				continue
			}

			go func() {
				_ = sc.Run()
			}()
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
		agent := ex.testAgentControllerToolBox.Build(scenario)
		go func() {
			_ = agent.Run()
		}()

		ex.mx.Lock()
		_, ok := ex.scenarios[scenario.ID]
		if !ok {
			ex.scenarios[scenario.ID] = &scenarioExecutor{
				scenario:    scenario,
				executionID: executionID,
				agents:      []interfaces.TestAgentController{agent},
				running:     false,
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
	_, span := tracer.Start(ctx, "RunScenario")
	defer span.End()

	if scenario == nil {
		zap.L().Error("test scenario is nil")

		return pkg.ErrTestScenarioServiceIsNil
	}

	ex.mx.Lock()
	defer ex.mx.Unlock()

	testScenario, ok := ex.scenarios[scenario.ID]
	if !ok {
		zap.L().Error("test scenario not found")

		return fmt.Errorf("%w", pkg.ErrTestScenarioNotFound)
	}

	testScenario.SetRunning(true)

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

	for i := int64(1); i <= *scenario.MaxTestServiceCount; i++ {
		ex.mx.Lock()
		_, ok := ex.scenarios[scenario.ID]
		if !ok {
			zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
			ex.mx.Unlock()

			return pkg.ErrTestScenarioNotFound
		}

		agent := ex.testAgentControllerToolBox.Get(scenario)
		err := agent.PauseTesting(ctx)
		if err != nil {
			zap.L().Error("test agent controller can not pause test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))
			ex.mx.Unlock()

			return fmt.Errorf("%w - %w", pkg.ErrFailedToPauseTestService, err)
		}

		ex.mx.Unlock()
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

	for i := int64(1); i <= *scenario.MaxTestServiceCount; i++ {
		ex.mx.Lock()
		_, ok := ex.scenarios[scenario.ID]
		if !ok {
			zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
			ex.mx.Unlock()

			return fmt.Errorf("%w", pkg.ErrTestScenarioNotFound)
		}

		agent := ex.testAgentControllerToolBox.Get(scenario)
		err := agent.ResumeTesting(ctx)
		if err != nil {
			zap.L().Error("test agent controller can not resume test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))
			ex.mx.Unlock()

			return fmt.Errorf("%w - %w", pkg.ErrFailedToResumeTestService, err)
		}

		ex.mx.Unlock()
	}

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

	for i := int64(1); i <= *scenario.MaxTestServiceCount; i++ {
		ex.mx.Lock()
		sci, ok := ex.scenarios[scenario.ID]
		if !ok {
			zap.L().Error("scenario not executed", zap.Uint64("scenarioID", scenario.ID))
			ex.mx.Unlock()

			return pkg.ErrTestScenarioNotFound
		}
		sci.SetRunning(false)

		agent := ex.testAgentControllerToolBox.Get(scenario)
		err := agent.StopTesting(ctx)
		if err != nil {
			zap.L().Error("test agent controller couldn't stop test service", zap.Error(err), zap.Uint64("scenarioID", scenario.ID))
			ex.mx.Unlock()

			return fmt.Errorf("%w - %w", pkg.ErrFailedToStopTestService, err)
		}

		ex.mx.Unlock()
	}

	return nil
}

//#endregion StressTestExecutionManager
