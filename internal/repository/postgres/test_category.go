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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

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
	tracer := otel.Tracer("test-category-repository")
	_, span := tracer.Start(ctx, "get_test_categories")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "select"))

	var items []entity.TestCategory
	err := postgres.QueryBuilder(ctx, repo.db).Order("id ASC").Find(&items).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))
		return nil, fmt.Errorf("failed to load test categories from db: %w", err)
	}

	return items, nil
}

func (repo *testCategory) GetByID(ctx context.Context, id uint64) (*entity.TestCategory, error) {
	tracer := otel.Tracer("test-category-repository")
	_, span := tracer.Start(ctx, "get_test_category_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("service.id", fmt.Sprintf("%d", id)))

	var testCategory entity.TestCategory
	err := postgres.QueryBuilder(ctx, repo.db).First(&testCategory, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))
			return nil, pkg.ErrTestCategoryNotFound
		}

		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))
		return nil, fmt.Errorf("failed to get test category record: %w", err)
	}

	span.SetAttributes(attribute.String("service.id", fmt.Sprintf("%d", testCategory.ID)))
	return &testCategory, nil
}
