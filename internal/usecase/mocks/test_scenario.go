package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockTestScenario struct {
	mock.Mock
}

func (m *MockTestScenario) Create(ctx context.Context, testScenario *entity.TestScenario) error {
	args := m.Called(ctx, testScenario)

	return args.Error(0)
}

func (m *MockTestScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	args := m.Called(ctx, id)

	var result *entity.TestScenario
	if args.Get(0) != nil {
		result = args.Get(0).(*entity.TestScenario)
	}

	return result, args.Error(1)
}

func (m *MockTestScenario) GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, int64, error) {
	args := m.Called(ctx, pagReq)

	var result []*entity.TestScenario
	if args.Get(0) != nil {
		result = args.Get(0).([]*entity.TestScenario)
	}

	return result, args.Get(1).(int64), args.Error(2)
}

func (m *MockTestScenario) Start(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockTestScenario) DeployTestScenarioService(ctx context.Context, motherService *entity.TestScenario) error {
	args := m.Called(ctx, motherService)
	return args.Error(0)
}

func (m *MockTestScenario) DeprovisionAllPods(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}
