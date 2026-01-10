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
