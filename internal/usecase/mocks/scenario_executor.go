package mocks

import (
	"context"
	"control-panel-service/internal/usecase/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockScenarioExecutor struct {
	mock.Mock
}

func (m *MockScenarioExecutor) Run(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *MockScenarioExecutor) IsRunning() bool {
	args := m.Called()

	return args.Get(0).(bool)
}

func (m *MockScenarioExecutor) SetRunning(status bool) {
	_ = m.Called(status)
}

func (m *MockScenarioExecutor) AllAgentsAreHealthy() bool {
	args := m.Called()

	return args.Get(0).(bool)
}

func (m *MockScenarioExecutor) AddAgent(agent interfaces.TestAgentController) {
	_ = m.Called(agent)
}

func (m *MockScenarioExecutor) GetAgents() []interfaces.TestAgentController {
	args := m.Called()

	return args.Get(0).([]interfaces.TestAgentController)
}
