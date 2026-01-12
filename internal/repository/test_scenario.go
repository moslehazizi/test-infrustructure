package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type TestScenarioRepository interface {
	Create(ctx context.Context, testSci *entity.TestScenario) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error)
	GetPaginated(ctx context.Context, pagRequest entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error)
}
