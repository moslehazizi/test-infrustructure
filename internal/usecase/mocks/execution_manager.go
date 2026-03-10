package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockExecutionManage struct {
	mock.Mock
}

func (m *MockExecutionManage) Run() {
	m.Called()
}

func (m *MockExecutionManage) AddScenario(ctx context.Context, scenario *entity.TestScenario, executionID uuid.UUID) error {
	args := m.Called(ctx, scenario, executionID)

	return args.Error(0)
}

func (m *MockExecutionManage) RunScenario(ctx context.Context, scenario *entity.TestScenario) error {
	args := m.Called(ctx, scenario)

	return args.Error(0)
}

func (m *MockExecutionManage) PauseScenario(ctx context.Context, scenario *entity.TestScenario) error {
	args := m.Called(ctx, scenario)

	return args.Error(0)
}

func (m *MockExecutionManage) ResumeScenario(ctx context.Context, scenario *entity.TestScenario) error {
	args := m.Called(ctx, scenario)

	return args.Error(0)
}

func (m *MockExecutionManage) RestartScenario(ctx context.Context, scenario *entity.TestScenario) error {
	args := m.Called(ctx, scenario)

	return args.Error(0)
}
