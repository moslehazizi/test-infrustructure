package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"

	"github.com/stretchr/testify/mock"
)

type MockMotherService struct {
	mock.Mock
}

func (m *MockMotherService) Begin() repository.MotherServiceRepository {
	_ = m.Called()

	return m
}

func (m *MockMotherService) Commit() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockMotherService) Rollback() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockMotherService) Create(ctx context.Context, motherService *entity.MotherService) error {
	args := m.Called(ctx, motherService)

	return args.Error(0)
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
