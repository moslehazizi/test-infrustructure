package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"
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

	// iteration of agent count increment.
	for range sc.scenario.ExecNumMultiAgent {
		for range sc.scenario.IncreaseAgentNumber {
			agent := sc.testAgentControllerToolBox.Build(sc.scenario)
			go func() {
				_ = agent.Run()
			}()
			sc.agents = append(sc.agents, agent)
		}

		// iteration of factorial number increment.
		// for {
		// }
		// run scenario
	}

	// TODO: complete implementation
	panic("complete implementation")
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

func (sc *scenarioExecutor) execute(ctx context.Context) error {
	// await to be healthy
	// await to be ready to start testing
	// executing tests
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

func (sc *scenarioExecutor) executeScenarioSteps(ctx context.Context) error {
	// sc.scenario.IncreaseAgentNumber
	// sc.scenario.ExecNumMultiAgent
	// sc.scenario.TestServiceConfig.IncreaseFixedInput
	// sc.scenario.TestServiceConfig.ExecNumMultiFixedInput

	for i := int64(1); i <= sc.scenario.NumSteps; i++ {
		if !sc.running {
			zap.L().Info("scenario executor is not running; exiting scenarioExecutor.Run")

			return pkg.ErrScenarioIsNotRunning
		}
		req := request.NewRunRequestFromTestServiceConfig(
			int(i),
			sc.executionID,
			sc.scenario.TestServiceConfig,
		)

		// send command to all agents.
		wg := sync.WaitGroup{}
		for _, agent := range sc.agents {
			wg.Add(1)
			//nolint
			go func() {
				defer wg.Done()
				_ = agent.StartTesting(context.Background(), *req)
			}()
		}

		wg.Wait()

		if i < sc.scenario.NumSteps {
			// wait for all agents to be ready to execute next step.
			sc.awaitAgentsToBeReadyToStartTesting()
		} else {
			// set status ScenarioStatusPending
			err := sc.scenarioRepo.SetStatus(ctx, sc.scenario.ID, entity.ScenarioStatusPending, true)
			if err != nil {
				zap.L().Error("failed to update scenario executor status", zap.Error(err))

				return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
			}
			// set running false , to privent run again and again
			sc.SetRunning(false)
		}
	}

	return nil
}
