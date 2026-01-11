package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

func NewTestServiceConfigRepository(db *gorm.DB) repository.TestServiceConfigRepository {
	return &testServiceConfig{
		db: db,
	}
}

type testServiceConfig struct {
	db *gorm.DB
}

func (repo *testServiceConfig) Create(ctx context.Context, testSvcCfg *entity.TestServiceConfig) error {
	err := repo.db.WithContext(ctx).Create(testSvcCfg).Error
	if err != nil {
		return fmt.Errorf("failed to create test service config record: %w", err)
	}

	return nil
}
