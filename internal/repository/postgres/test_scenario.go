package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

func NewTestScenarioRepository(db *gorm.DB) repository.TestScenarioRepository {
	return &testScenario{
		db: db,
	}
}

type testScenario struct {
	db *gorm.DB
}

func (repo *testScenario) Create(ctx context.Context, testSci *entity.TestScenario) error {
	err := repo.db.WithContext(ctx).Create(testSci).Error
	if err != nil {
		return fmt.Errorf("failed to create test scenario record: %w", err)
	}

	return nil
}
