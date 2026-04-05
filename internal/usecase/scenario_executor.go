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
		zap.L().Error("scenarioExecutor Run called but scenario is not running",
			zap.Uint64("scenarioID", sc.scenario.ID),
			zap.String("executionID", sc.executionID.String()),
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

	sse := sc.scenarioExecutorBuilder.Build(
		sc.agents,
		sc.scenario,
		sc.executionID,
	)

	err := sse.Execute(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToExecuteSingleScenario, err)
	}

	if sc.scenario.TestServiceConfig.ExecNumMultiFixedInput == 0 && sc.scenario.ExecNumMultiAgent == 0 {
		for _, agent := range sc.agents {
			err := agent.AbortTesting(ctx)
			if err != nil {
				zap.L().Error("failed to deprovision test agent",
					zap.Uint64("scenarioID", sc.scenario.ID),
					zap.String("executionID", sc.executionID.String()),
				)
			}
		}

		return nil
	}

	if sc.scenario.TestServiceConfig.ExecNumMultiFixedInput == 0 {
		sc.scenario.TestServiceConfig.ExecNumMultiFixedInput++
	}

	if sc.scenario.ExecNumMultiAgent == 0 {
		sc.scenario.ExecNumMultiAgent++
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

		for range sc.scenario.TestServiceConfig.ExecNumMultiFixedInput {
			if sc.scenario.TestServiceConfig.FixedTestNumber != nil {
				*sc.scenario.TestServiceConfig.FixedTestNumber += sc.scenario.TestServiceConfig.IncreaseFixedInput
			}

			sse := sc.scenarioExecutorBuilder.Build(
				sc.agents,
				sc.scenario,
				sc.executionID,
			)

			err := sse.Execute(ctx)
			if err != nil {
				return err
			}
		}
	}

	for _, agent := range sc.agents {
		err := agent.AbortTesting(ctx)
		if err != nil {
			zap.L().Error("failed to deprovision test agent",
				zap.Uint64("scenarioID", sc.scenario.ID),
				zap.String("executionID", sc.executionID.String()),
			)
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

	return nil
}
