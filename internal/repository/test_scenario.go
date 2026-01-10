package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type TestScenarioRepository interface {
	Create(ctx context.Context, testSci *entity.TestScenario) error
}
