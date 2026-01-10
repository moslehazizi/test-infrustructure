package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockTestScenario struct {
	mock.Mock
}

func (m *MockTestScenario) Create(ctx context.Context, testSci *entity.TestScenario) error {
	args := m.Called(ctx, testSci)

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

func (m *MockTestScenario) GetPaginated(ctx context.Context, pagRequest entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error) {
	args := m.Called(ctx, pagRequest)

	var result []*entity.TestScenario
	if args.Get(0) != nil {
		result = args.Get(0).([]*entity.TestScenario)
	}

	return result, args.Error(1)
}
