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
