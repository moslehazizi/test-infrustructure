package mocks

import (
	"control-panel-service/internal/usecase/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockTestAgentController struct {
	mock.Mock
}

func (m *MockTestAgentController) Run() {
	m.Called()
}

type MockTestAgentControllerBuilder struct {
	mock.Mock
}

func (m *MockTestAgentControllerBuilder) Build() interfaces.TestAgentController {
	args := m.Called()

	return args.Get(0).(interfaces.TestAgentController)
}
