package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockTestServiceRepository struct {
	mock.Mock
}

func (m *MockTestServiceRepository) GetCountAllRunningByScenario(ctx context.Context, scenarioID uint64) (int, error) {
	args := m.Called(ctx, scenarioID)

	return args.Get(0).(int), args.Error(1)
}

func (m *MockTestServiceRepository) GetRunningByScenario(ctx context.Context, scenarioID uint64, limit int) ([]entity.TestService, error) {
	args := m.Called(ctx, scenarioID, limit)

	var result []entity.TestService
	if args.Get(0) != nil {
		result = args.Get(0).([]entity.TestService)
	}

	return result, args.Error(1)
}
