package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"go.uber.org/zap"
)

func NewTestCategoryService(testCategoryRepo repository.TestCategory) *testCategoryService {
	return &testCategoryService{testCategoryRepo}
}

type testCategoryService struct {
	testCategoryRepo repository.TestCategory
}

func (srv *testCategoryService) GetAll(ctx context.Context) ([]entity.TestCategory, error) {
	tracer := otel.Tracer("test-category-usecase")
	useCaseCTX, span := tracer.Start(ctx, "get-test-categories-usecase")
	defer span.End()

	items, err := srv.testCategoryRepo.GetAll(useCaseCTX)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_all_error"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get all test categories",
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategoriesFromRepository, err)
	}

	zap.L().Debug("retrieved all test categories",
		zap.Int("count", len(items)),
	)

	return items, nil
}

func (srv *testCategoryService) GetByID(ctx context.Context, id uint64) (*entity.TestCategory, error) {
	tracer := otel.Tracer("test-category-usecase")
	useCaseCTX, span := tracer.Start(ctx, "get-test-category-by-id-usecase")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", strconv.FormatUint(id, 10)))

	item, err := srv.testCategoryRepo.GetByID(useCaseCTX, id)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_by_id_error"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get test category by ID",
			zap.Uint64("id", id),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategoryFromRepository, err)
	}

	return item, nil
}
