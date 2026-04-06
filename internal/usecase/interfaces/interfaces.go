package interfaces

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"

	"github.com/google/uuid"
)

type ExecutionManager interface {
	// RunScenario is for run scenario.
	RunScenario(ctx context.Context, scenario *entity.TestScenario) error
	// PauseScenario is for pause ran scenario.
	PauseScenario(ctx context.Context, scenario *entity.TestScenario) error
	// ResumeScenario is for resume paused scenario.
	ResumeScenario(ctx context.Context, scenario *entity.TestScenario) error
	// StopScenario is for stop running scenario.
	StopScenario(ctx context.Context, scenario *entity.TestScenario) error
	// AbortScenario is for stop running scenario.
	AbortScenario(ctx context.Context, scenario *entity.TestScenario) error
}

type ScenarioExecutor interface {
	Run(ctx context.Context) error
	IsRunning() bool
	SetRunning(status bool)
	AllAgentsAreHealthy() bool
	AddAgent(agent TestAgentController)
	GetAgents() []TestAgentController
}

type SingleScenarioExecutor interface {
	Execute(ctx context.Context) error
	// AwaitAgentsToBeHealthy()
	// AwaitAgentsToBeReadyToStartTesting()
}

type SingleScenarioExecutorBuilder interface {
	Build(
		agents []TestAgentController,
		scenario *entity.TestScenario,
		executionID uuid.UUID,
	) SingleScenarioExecutor
}

type TestAgentController interface {
	// Run runs agent.
	Run() error
	// StartTesting is responsible for sending start command.
	StartTesting(ctx context.Context, req request.RunRequest) error
	// Healthy checks if related test service is up and running.
	Healthy() bool
	// ReadyForTesting tests that test service is not executing
	// any test and is ready to get execution command.
	ReadyForTesting() bool
	// PauseTesting pause agent.
	PauseTesting(ctx context.Context) error
	// ResumeTesting resume agent.
	ResumeTesting(ctx context.Context) error
	// StopTesting stop agent.
	StopTesting(ctx context.Context) error
	// AbortTesting is responsible for sending abort command.
	AbortTesting(ctx context.Context) error
}

type TestAgentControllerToolBox interface {
	// Build will define new TestAgentController based on given config.
	Build(scenario *entity.TestScenario) TestAgentController
}
