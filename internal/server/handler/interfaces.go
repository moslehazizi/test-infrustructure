package handler

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/server/dto/request"
)

type MotherService interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
	GetByID(ctx context.Context, id uint64) (*entity.MotherService, error)
	GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, int64, error)
	Delete(ctx context.Context, id uint64) error
}

type TestScenario interface {
	Create(ctx context.Context, testScenario *entity.TestScenario) error
	GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error)
	GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, int64, error)
	Update(ctx context.Context, testScenarioUpdateRequest *request.TestScenarioUpdateRequest) error
}

type TestScenarioOperation interface {
	Start(ctx context.Context, id uint64) error
	Pause(ctx context.Context, id uint64) error
	Resume(ctx context.Context, id uint64) error
	Stop(ctx context.Context, id uint64) error
	Delete(ctx context.Context, id uint64) error
}

type TestCategoryService interface {
	GetAll(ctx context.Context) ([]entity.TestCategory, error)
	GetByID(ctx context.Context, id uint64) (*entity.TestCategory, error)
}

type DatabaseMetadata interface {
	GetAll(ctx context.Context) ([]string, error)
	GetTablesByDBName(ctx context.Context, dbName string) (*entity.TablesByType, error)
}
