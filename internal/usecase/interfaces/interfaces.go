package interfaces

import (
	"context"
	"control-panel-service/internal/domain/entity"

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
	Run()
}

type TestAgentControllerBuilder interface {
	// Build will define new TestAgentController based on given config.
	Build() TestAgentController
}
