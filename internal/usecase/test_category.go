package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/logger"
	"fmt"

	"go.uber.org/zap"
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
	requestID := logger.GetRequestID(ctx)
	items, err := srv.testCategoryRepo.GetAll(ctx)
	if err != nil {
		zap.L().Error("failed to get all test categories",
			zap.String(logger.FieldRequestID, requestID),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategoriesFromRepository, err)
	}

	zap.L().Debug("retrieved all test categories",
		zap.String(logger.FieldRequestID, requestID),
		zap.Int("count", len(items)),
	)

	return items, nil
}

func (srv *testCategoryService) GetByID(ctx context.Context, id uint64) (*entity.TestCategory, error) {
	requestID := logger.GetRequestID(ctx)
	item, err := srv.testCategoryRepo.GetByID(ctx, id)
	if err != nil {
		zap.L().Error("failed to get test category by ID",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", id),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategoryFromRepository, err)
	}

	return item, nil
}
