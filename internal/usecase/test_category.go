package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"fmt"
)

type TestCategoryService interface {
	GetAll(ctx context.Context) ([]entity.TestCategory, error)
}

func NewTestCategoryService(testCategoryRepo repository.TestCategory) TestCategoryService {
	return &testCategoryService{testCategoryRepo}
}

type testCategoryService struct {
	testCategoryRepo repository.TestCategory
}

func (srv *testCategoryService) GetAll(ctx context.Context) ([]entity.TestCategory, error) {
	items, err := srv.testCategoryRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategoriesFromRepository, err)
	}

	return items, nil
}
