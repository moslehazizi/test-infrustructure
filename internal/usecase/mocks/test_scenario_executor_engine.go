package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"
	"sync"

	"github.com/stretchr/testify/mock"
)

type MockScenarioExecutorBox struct {
	mock.Mock
}

func (m *MockScenarioExecutorBox) Add(exe interfaces.ScenarioExecutor) {
	_ = m.Called(exe)
}

func (m *MockScenarioExecutorBox) HasExecutor(id uint64) bool {
	args := m.Called(id)

	return args.Get(0).(bool)
}

type MockScenarioExecutor struct {
	mock.Mock
}

func (m *MockScenarioExecutor) GetID() uint64 {
	args := m.Called()

	return args.Get(0).(uint64)
}

func (m *MockScenarioExecutor) GetScenario() *entity.TestScenario {
	args := m.Called()

	return args.Get(0).(*entity.TestScenario)
}

func (m *MockScenarioExecutor) ResumeOrStart(ctx context.Context) {
	m.Called()
}

type MockScenarioExecutorWithWG struct {
	mock.Mock

	WG sync.WaitGroup
}

func (m *MockScenarioExecutorWithWG) GetID() uint64 {
	args := m.Called()

	return args.Get(0).(uint64)
}

func (m *MockScenarioExecutorWithWG) GetScenario() *entity.TestScenario {
	args := m.Called()

	return args.Get(0).(*entity.TestScenario)
}

func (m *MockScenarioExecutorWithWG) ResumeOrStart(ctx context.Context) {
	m.Called()
	m.WG.Done()
}
