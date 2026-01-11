package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockTestCategoryService struct {
	mock.Mock
}

func (mock *MockTestCategoryService) GetAll(ctx context.Context) ([]entity.TestCategory, error) {
	args := mock.Called(ctx)

	return args.Get(0).([]entity.TestCategory), args.Error(1)
}

func (mock *MockTestCategoryService) GetByID(ctx context.Context, id uint64) (*entity.TestCategory, error) {
	args := mock.Called(ctx)

	return args.Get(0).(*entity.TestCategory), args.Error(1)
}
