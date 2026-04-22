package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/usecase/interfaces"
	"sync"
	"time"

	"github.com/google/uuid"
)

type singleScenarioExecutor struct {
	agents                   []interfaces.TestAgentController
	allAgentsHealthy         bool
	allAgentsReadyForTesting bool
	scenario                 *entity.TestScenario
	executionID              uuid.UUID
}

func (ss *singleScenarioExecutor) Execute(ctx context.Context) error {
	// await to be healthy
	ss.awaitAgentsToBeHealthy()

	// await to be ready to start testing
	ss.awaitAgentsToBeReadyToStartTesting()

	// executing tests
	ss.executeScenarioSteps(ctx)

	return nil
}

func (ss *singleScenarioExecutor) awaitAgentsToBeHealthy() {
	awaitAllHealthy(ss.agents)
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
func (ss *singleScenarioExecutor) executeScenarioSteps(_ context.Context) {
	for i := int64(1); i <= ss.scenario.NumSteps; i++ {
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
		}
	}
}
