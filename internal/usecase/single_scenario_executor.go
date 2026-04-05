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
	"go.uber.org/zap"
)

type singleScenarioExecutor struct {
	agents                   []interfaces.TestAgentController
	allAgentsHealthy         bool
	allAgentsReadyForTesting bool
	scenario                 *entity.TestScenario
	executionID              uuid.UUID
	running                  bool
}

func (ss *singleScenarioExecutor) Execute(ctx context.Context) error {
	// await to be healthy
	ss.awaitAgentsToBeHealthy()

	// await to be ready to start testing
	ss.awaitAgentsToBeReadyToStartTesting()

	// executing tests
	err := ss.executeScenarioSteps(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute scenario: %w", err)
	}

	return nil
}

func (ss *singleScenarioExecutor) awaitAgentsToBeHealthy() {
	for {
		allHealthy := true
		for _, agent := range ss.agents {
			if !agent.Healthy() {
				allHealthy = false

				break
			}
		}

		ss.allAgentsHealthy = allHealthy

		if allHealthy {
			break
		}

		time.Sleep(healthyCheckSleep)
	}
}

func (ss *singleScenarioExecutor) awaitAgentsToBeReadyToStartTesting() {
	for {
		allReady := true
		for _, agent := range ss.agents {
			if !agent.ReadyForTesting() {
				allReady = false

				break
			}
		}

		ss.allAgentsReadyForTesting = allReady

		if allReady {
			break
		}

		time.Sleep(readyForTestingCheckSleep)
	}
}
func (ss *singleScenarioExecutor) executeScenarioSteps(ctx context.Context) error {
	// sc.scenario.IncreaseAgentNumber
	// sc.scenario.ExecNumMultiAgent
	// sc.scenario.TestServiceConfig.IncreaseFixedInput
	// sc.scenario.TestServiceConfig.ExecNumMultiFixedInput

	for i := int64(1); i <= ss.scenario.NumSteps; i++ {
		if !ss.running {
			zap.L().Info("scenario executor is not running; exiting scenarioExecutor.Run")

			return pkg.ErrScenarioIsNotRunning
		}
		req := request.NewRunRequestFromTestServiceConfig(
			int(i),
			ss.executionID,
			ss.scenario.TestServiceConfig,
		)

		// send command to all agents.
		wg := sync.WaitGroup{}
		for _, agent := range ss.agents {
			wg.Add(1)
			//nolint
			go func() {
				defer wg.Done()
				_ = agent.StartTesting(context.Background(), *req)
			}()
		}

		wg.Wait()

		if i < ss.scenario.NumSteps {
			// wait for all agents to be ready to execute next step.
			ss.awaitAgentsToBeReadyToStartTesting()
		} else {
			// // set status ScenarioStatusPending
			// err := sc.scenarioRepo.SetStatus(ctx, sc.scenario.ID, entity.ScenarioStatusPending, true)
			// if err != nil {
			// 	zap.L().Error("failed to update scenario executor status", zap.Error(err))

			// 	return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
			// }
			// // set running false , to prevent run again and again
			// sc.SetRunning(false)
		}
	}

	return nil
}
