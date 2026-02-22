package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"
)

func NewStressTestExecutionManager() interfaces.ExecutionManager {
	return &StressTestExecutionManager{}
}

type StressTestExecutionManager struct {
}

func (ex *StressTestExecutionManager) AddScenario(ctx context.Context, scenario *entity.TestScenario) error {
	return nil
}
