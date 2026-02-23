package interfaces

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"

	"github.com/google/uuid"
)

type ExecutionManager interface {
	// Run will run ExecutionManager as a background job.
	Run()
	// Add will initialize new TestAgentControllers based on scenario config.
	AddScenario(ctx context.Context, scenario *entity.TestScenario, executionID uuid.UUID) error
}

type TestAgentController interface {
	// Run runs agent.
	Run() error

	StartTesting(ctx context.Context, req request.RunRequest) error

	// Healthy checks if related test service is up and running.
	Healthy() bool
}

type TestAgentControllerBuilder interface {
	// Build will define new TestAgentController based on given config.
	Build(scenario *entity.TestScenario) TestAgentController
}
