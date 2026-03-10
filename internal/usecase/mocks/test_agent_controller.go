package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"
	"control-panel-service/internal/usecase/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockTestAgentController struct {
	mock.Mock
}

func (m *MockTestAgentController) Run() error {
	arg := m.Called()

	return arg.Error(0)
}

func (m *MockTestAgentController) Healthy() bool {
	args := m.Called()

	return args.Get(0).(bool)
}

func (m *MockTestAgentController) ReadyForTesting() bool {
	args := m.Called()

	return args.Get(0).(bool)
}

func (m *MockTestAgentController) StartTesting(ctx context.Context, req request.RunRequest) error {
	args := m.Called(ctx, req)

	return args.Error(0)
}

func (m *MockTestAgentController) AbortTesting(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *MockTestAgentController) PauseTesting(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *MockTestAgentController) ResumeTesting(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *MockTestAgentController) StopTesting(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

type MockTestAgentControllerToolBox struct {
	mock.Mock
}

func (m *MockTestAgentControllerToolBox) Build(scenario *entity.TestScenario) interfaces.TestAgentController {
	args := m.Called(scenario)

	return args.Get(0).(interfaces.TestAgentController)
}

func (m *MockTestAgentControllerToolBox) Get(scenario *entity.TestScenario) interfaces.TestAgentController {
	args := m.Called(scenario)

	return args.Get(0).(interfaces.TestAgentController)
}
