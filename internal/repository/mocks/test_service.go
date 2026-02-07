package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTestServiceRepository struct {
	mock.Mock
}

func (m *MockTestServiceRepository) GetCountAllRunningByScenario(ctx context.Context, scenarioID uint64) (int, error) {
	args := m.Called(ctx, scenarioID)

	return args.Get(0).(int), args.Error(1)
}
