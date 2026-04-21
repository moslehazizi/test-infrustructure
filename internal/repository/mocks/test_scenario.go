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

func (m *MockTestScenario) GetPaginated(ctx context.Context, pagRequest entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, int64, error) {
	args := m.Called(ctx, pagRequest)

	var result []*entity.TestScenario
	if args.Get(0) != nil {
		result = args.Get(0).([]*entity.TestScenario)
	}

	return result, args.Get(1).(int64), args.Error(2)
}

func (m *MockTestScenario) SetStatus(ctx context.Context, id uint64, status entity.ScenarioStatus, editable bool) error {
	args := m.Called(ctx, id, status, editable)

	return args.Error(0)
}

func (m *MockTestScenario) GetByStatus(ctx context.Context, status entity.ScenarioStatus) ([]*entity.TestScenario, error) {
	args := m.Called(ctx, status)

	var result []*entity.TestScenario
	if args.Get(0) != nil {
		result = args.Get(0).([]*entity.TestScenario)
	}

	return result, args.Error(1)
}

func (m *MockTestScenario) GetDeploymentNumberByScenarioID(ctx context.Context, id uint64) (int32, error) {
	args := m.Called(ctx, id)

	if v, ok := args.Get(0).(int32); ok {
		return v, args.Error(1)
	}

	return 0, args.Error(1)
}

func (m *MockTestScenario) UpdateDeploymentNumber(ctx context.Context, id uint64, newDeploymentNumber int32) error {
	args := m.Called(ctx, id, newDeploymentNumber)

	return args.Error(0)
}

func (m *MockTestScenario) Update(ctx context.Context, scenario *entity.TestScenario) error {
	args := m.Called(ctx, scenario)

	return args.Error(0)
}

func (m *MockTestScenario) GetByMotherServiceId(ctx context.Context, motherServiceId uint64) ([]*entity.TestScenario, error) {
	args := m.Called(ctx, motherServiceId)

	var result []*entity.TestScenario
	if args.Get(0) != nil {
		result = args.Get(0).([]*entity.TestScenario)
	}

	return result, args.Error(1)
}

func (m *MockTestScenario) UpdateScenarioAndConfig(ctx context.Context, scenario *entity.TestScenario) error {
	args := m.Called(ctx, scenario)

	return args.Error(0)
}
