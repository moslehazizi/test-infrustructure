package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"

	"github.com/stretchr/testify/mock"
)

type MockTestScenario struct {
	mock.Mock
}

func (m *MockTestScenario) Begin() repository.TestScenarioRepository {
	_ = m.Called()

	return m
}

func (m *MockTestScenario) Commit() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockTestScenario) Rollback() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockTestScenario) Create(ctx context.Context, testSci *entity.TestScenario) (uint64, error) {
	args := m.Called(ctx, testSci)

	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockTestScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	args := m.Called(ctx, id)

	var result *entity.TestScenario
	if args.Get(0) != nil {
		result = args.Get(0).(*entity.TestScenario)
	}

	return result, args.Error(1)
}

func (m *MockTestScenario) GetPaginated(ctx context.Context, pagRequest entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error) {
	args := m.Called(ctx, pagRequest)

	var result []*entity.TestScenario
	if args.Get(0) != nil {
		result = args.Get(0).([]*entity.TestScenario)
	}

	return result, args.Error(1)
}
