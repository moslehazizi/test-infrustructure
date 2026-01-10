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
