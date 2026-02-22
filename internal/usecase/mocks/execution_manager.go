package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockExecutionManage struct {
	mock.Mock
}

func (m *MockExecutionManage) AddScenario(ctx context.Context, scenario *entity.TestScenario) error {
	args := m.Called(ctx, scenario)

	return args.Error(0)
}
