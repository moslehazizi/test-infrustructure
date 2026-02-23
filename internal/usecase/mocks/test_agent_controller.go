package mocks

import (
	"control-panel-service/internal/domain/entity"
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

type MockTestAgentControllerBuilder struct {
	mock.Mock
}

func (m *MockTestAgentControllerBuilder) Build(scenario *entity.TestScenario) interfaces.TestAgentController {
	args := m.Called(scenario)

	return args.Get(0).(interfaces.TestAgentController)
}
