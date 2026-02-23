package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type TestScenarioRepository interface {
	Create(ctx context.Context, testSci *entity.TestScenario) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error)
	GetPaginated(ctx context.Context, pagRequest entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, int64, error)
	SetStatus(ctx context.Context, id uint64, status entity.ScenarioStatus, editable bool) error
	GetByStatus(ctx context.Context, status entity.ScenarioStatus) ([]*entity.TestScenario, error)
	GetDeploymentNumberByScenarioID(ctx context.Context, id uint64) (int32, error)
	UpdateDeploymentNumber(ctx context.Context, id uint64, newDeploymentNumber int32) error
	Update(ctx context.Context, scenario *entity.TestScenario) error
}
