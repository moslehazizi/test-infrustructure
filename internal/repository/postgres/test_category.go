package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func NewTestCategoryRepository(db database.Database) repository.TestCategory {
	return &testCategory{
		db,
	}
}

type testCategory struct {
	db database.Database
}

func (repo *testCategory) GetAll(ctx context.Context) ([]entity.TestCategory, error) {
	var items []entity.TestCategory
	err := postgres.QueryBuilder(ctx, repo.db).Order("id ASC").Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load test categories from db: %w", err)
	}

	return items, nil
}

func (repo *testCategory) GetByID(ctx context.Context, id uint64) (*entity.TestCategory, error) {
	var testCategory entity.TestCategory
	err := postgres.QueryBuilder(ctx, repo.db).First(&testCategory, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrTestCategoryNotFound
		}

		return nil, fmt.Errorf("failed to get test category record: %w", err)
	}

	return &testCategory, nil
}
