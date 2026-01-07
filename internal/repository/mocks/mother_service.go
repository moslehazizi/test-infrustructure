package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockMotherService struct {
	mock.Mock
}

func (m *MockMotherService) Create(ctx context.Context, motherService *entity.MotherService) error {
	args := m.Called(ctx, motherService)

	return args.Error(1)
}

func (m *MockMotherService) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*entity.MotherService), args.Error(1)
}

func (m *MockMotherService) GetAll(ctx context.Context) ([]*entity.MotherService, error) {
	args := m.Called(ctx)

	return args.Get(0).([]*entity.MotherService), args.Error(1)
}
