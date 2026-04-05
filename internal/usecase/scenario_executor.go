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

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RunOnceDelay = time.Millisecond * 100
)

var (
	healthyCheckSleep         = time.Second
	readyForTestingCheckSleep = time.Second
)

type scenarioExecutor struct {
	scenario                   *entity.TestScenario
	executionID                uuid.UUID
	agents                     []interfaces.TestAgentController
	allAgentsHealthy           bool
	allAgentsReadyForTesting   bool
	running                    bool
	once                       sync.Once
	scenarioRepo               repository.TestScenarioRepository
	testAgentControllerToolBox interfaces.TestAgentControllerToolBox
	scenarioExecutorBuilder    interfaces.SingleScenarioExecutorBuilder
}

func (sc *scenarioExecutor) IsRunning() bool {
	return sc.running
}

func (sc *scenarioExecutor) SetRunning(status bool) {
	sc.running = status
}

func (sc *scenarioExecutor) assignExecutionID() {
	sc.executionID = uuid.New()
}

func (sc *scenarioExecutor) AddAgent(agent interfaces.TestAgentController) {
	sc.agents = append(sc.agents, agent)
}

func (sc *scenarioExecutor) GetAgents() []interfaces.TestAgentController {
	return sc.agents
}

func (sc *scenarioExecutor) AllAgentsAreHealthy() bool {
	return sc.allAgentsHealthy
}

func (sc *scenarioExecutor) Run(ctx context.Context) error {
	zap.L().Info("scenarioExecutor.Run Called")

	sc.assignExecutionID()
	// check scenario is running
	if !sc.running {
		zap.L().Error(
			"scenarioExecutor Run called but scenario is not running",
			zap.Any("scenarioID", sc.scenario.ID),
			zap.Any("executionID", sc.executionID.String()),
		)

		return pkg.ErrScenarioIsNotRunning
	}

	if sc.scenario.MaxTestServiceCount == nil {
		return pkg.ErrMaxTestServiceCountNotSet
	}

	if *sc.scenario.MaxTestServiceCount < 1 {
		return pkg.ErrMaxTestServiceCountLessThanOne
	}

	for range *sc.scenario.MaxTestServiceCount {
		agent := sc.testAgentControllerToolBox.Build(sc.scenario)
		go func() {
			_ = agent.Run()
		}()

		sc.agents = append(sc.agents, agent)
	}

	// run scenario
	sse := sc.scenarioExecutorBuilder.Build(
		sc.agents,
		sc.allAgentsHealthy,
		sc.allAgentsReadyForTesting,
		sc.scenario,
		sc.executionID,
		sc.running,
	)

	err := sse.Execute(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToExecuteSingleScenario, err)
	}

	// iteration of agent count increment.
	for range sc.scenario.ExecNumMultiAgent {
		for range sc.scenario.IncreaseAgentNumber {
			agent := sc.testAgentControllerToolBox.Build(sc.scenario)
			go func() {
				_ = agent.Run()
			}()
			sc.agents = append(sc.agents, agent)
		}

		// run scenario
		sse := sc.scenarioExecutorBuilder.Build(
			sc.agents,
			sc.allAgentsHealthy,
			sc.allAgentsReadyForTesting,
			sc.scenario,
			sc.executionID,
			sc.running,
		)

		err := sse.Execute(ctx)
		if err != nil {
			return err
		}

	}

	// TODO: complete implementation
	// panic("complete implementation")
	// loop
	// run scenario(
	// // await to be healthy
	// // await to be ready to start testing
	// // executing tests
	// )
	// end loop
	// deprovision all agents
	// remove scenario from memory
	// ------------

	// start:
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	default:
	// 		for !sc.running {
	// 			time.Sleep(RunOnceDelay)

	// 			zap.L().Debug("scenarioExecutor.RunOnce waiting to be run")
	// 		}

	// 		zap.L().Info("scenarioExecutor.RunOnce Running")

	// 		// provision agents

	// 		// Make sure all agents are healthy.
	// 		sc.awaitAgentsToBeHealthy()

	// 		// NOW: all agents are healthy.
	// 		// we can send scheduled commands.
	// 		// we need a loop based on len of steps:
	// 		_ = sc.executeScenarioSteps(ctx)

	// 		goto start
	// 	}

	return nil
}
