package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type TestServiceConfigRepository interface {
	Create(ctx context.Context, testSvcCfg *entity.TestServiceConfig) error
	GetByID(ctx context.Context, id uint64) (*entity.TestServiceConfig, error)
	UpdateByScenarioID(ctx context.Context, scenarioID uint64, cfg *entity.TestServiceConfig) error
}
