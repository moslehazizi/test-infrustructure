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

func (se *scenarioExecutor) IsRunning() bool {
	return se.running
}

func (se *scenarioExecutor) SetRunning(status bool) {
	se.running = status
}

func (se *scenarioExecutor) assignExecutionID() {
	se.executionID = uuid.New()
}

func (se *scenarioExecutor) AddAgent(agent interfaces.TestAgentController) {
	se.agents = append(se.agents, agent)
}

func (se *scenarioExecutor) GetAgents() []interfaces.TestAgentController {
	return se.agents
}

func (se *scenarioExecutor) AllAgentsAreHealthy() bool {
	return se.allAgentsHealthy
}

func (se *scenarioExecutor) Run(ctx context.Context) (e error) {
	zap.L().Info("scenarioExecutor.Run Called")

	defer func() {
		for _, agent := range se.agents {
			err := agent.AbortTesting(ctx)
			if err != nil {
				// returning err is not required.
				zap.L().Error("failed to deprovision test agent",
					zap.Uint64("scenarioID", se.scenario.ID),
					zap.String("executionID", se.executionID.String()),
				)
			}
		}

		// set running false
		se.SetRunning(false)
		// set status
		err := se.scenarioRepo.SetStatus(ctx, se.scenario.ID, entity.ScenarioStatusPending, true)
		if err != nil && e == nil {
			e = fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
		}
	}()

	se.assignExecutionID()
	// check scenario is running
	if !se.running {
		zap.L().Error("scenarioExecutor Run called but scenario is not running",
			zap.Uint64("scenarioID", se.scenario.ID),
			zap.String("executionID", se.executionID.String()),
		)

		return pkg.ErrScenarioIsNotRunning
	}

	if se.scenario.MaxTestServiceCount == nil {
		return pkg.ErrMaxTestServiceCountNotSet
	}

	if *se.scenario.MaxTestServiceCount < 1 {
		return pkg.ErrMaxTestServiceCountLessThanOne
	}

	// As we need stage zero in all scenarios , this addition is necessary.
	se.scenario.TestServiceConfig.ExecNumMultiFixedInput++
	se.scenario.ExecNumMultiAgent++

	// iteration of agent count increment.
	for s := range se.scenario.ExecNumMultiAgent {
		if s == 0 {
			for range *se.scenario.MaxTestServiceCount {
				agent := se.testAgentControllerToolBox.Build(se.scenario)
				go func() {
					_ = agent.Run()
				}()
				se.agents = append(se.agents, agent)
			}
		} else {
			for range se.scenario.IncreaseAgentNumber {
				agent := se.testAgentControllerToolBox.Build(se.scenario)
				go func() {
					_ = agent.Run()
				}()
				se.agents = append(se.agents, agent)
			}
		}

		if se.scenario.TestServiceConfig.FixedTestNumber != nil {
			number := *se.scenario.TestServiceConfig.FixedTestNumber
			for range se.scenario.TestServiceConfig.ExecNumMultiFixedInput {
				sse := se.scenarioExecutorBuilder.Build(
					se.agents,
					se.scenario,
					se.executionID,
				)

				err := sse.Execute(ctx)
				if err != nil {
					return fmt.Errorf("%w: %w", pkg.ErrFailedToExecuteSingleScenario, err)
				}

				number += se.scenario.TestServiceConfig.IncreaseFixedInput

			}
		} else {
			sse := se.scenarioExecutorBuilder.Build(
				se.agents,
				se.scenario,
				se.executionID,
			)

			err := sse.Execute(ctx)
			if err != nil {
				return fmt.Errorf("%w: %w", pkg.ErrFailedToExecuteSingleScenario, err)
			}
		}
	}

	return nil
}
