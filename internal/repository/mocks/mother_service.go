package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockMotherService struct {
	mock.Mock
}

func (m *MockMotherService) Create(ctx context.Context, motherService *entity.MotherService) (uint64, error) {
	args := m.Called(ctx, motherService)

	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockMotherService) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	args := m.Called(ctx, id)

	var result *entity.MotherService
	if args.Get(0) != nil {
		result = args.Get(0).(*entity.MotherService)
	}

	return result, args.Error(1)
}

func (m *MockMotherService) GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error) {
	args := m.Called(ctx, paginationRequest)

	var result []*entity.MotherService
	if args.Get(0) != nil {
		result = args.Get(0).([]*entity.MotherService)
	}

	return result, args.Error(1)
}
