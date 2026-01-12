package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"fmt"
)

func NewTestServiceConfigRepository(db database.Database) repository.TestServiceConfigRepository {
	return &testServiceConfig{
		db: db,
	}
}

type testServiceConfig struct {
	db database.Database
}

func (repo *testServiceConfig) Create(ctx context.Context, testSvcCfg *entity.TestServiceConfig) error {
	err := postgres.QueryBuilder(ctx, repo.db).Create(testSvcCfg).Error
	if err != nil {
		return fmt.Errorf("failed to create test service config record: %w", err)
	}

	return nil
}
