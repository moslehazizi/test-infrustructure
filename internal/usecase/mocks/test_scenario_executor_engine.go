package mocks

import (
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockScenarioExecutorEngine struct {
	mock.Mock
}

func (m *MockScenarioExecutorEngine) Add(scenario *entity.TestScenario) {
	_ = m.Called(scenario)
}
