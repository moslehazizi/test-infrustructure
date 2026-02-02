package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type TestCategoryService interface {
	GetAll(ctx context.Context) ([]entity.TestCategory, error)
	GetByID(ctx context.Context, id uint64) (*entity.TestCategory, error)
}

func NewTestCategoryService(testCategoryRepo repository.TestCategory) TestCategoryService {
	return &testCategoryService{testCategoryRepo}
}

type testCategoryService struct {
	testCategoryRepo repository.TestCategory
}

func (srv *testCategoryService) GetAll(ctx context.Context) ([]entity.TestCategory, error) {
	tracer := otel.Tracer("test-category-usecase")
	_, span := tracer.Start(ctx, "get_test_categories")
	defer span.End()

	items, err := srv.testCategoryRepo.GetAll(ctx)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_all_error"), attribute.String("error.message", err.Error()))
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategoriesFromRepository, err)
	}

	return items, nil
}

func (srv *testCategoryService) GetByID(ctx context.Context, id uint64) (*entity.TestCategory, error) {
	tracer := otel.Tracer("test-category-usecase")
	_, span := tracer.Start(ctx, "get_test_category_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", fmt.Sprintf("%d", id)))

	item, err := srv.testCategoryRepo.GetByID(ctx, id)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_by_id_error"), attribute.String("error.message", err.Error()))
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategoryFromRepository, err)
	}

	return item, nil
}
