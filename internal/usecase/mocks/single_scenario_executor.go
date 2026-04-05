package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockSingleScenarioExecutor struct {
	mock.Mock
}

func (m *MockSingleScenarioExecutor) Execute(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}
