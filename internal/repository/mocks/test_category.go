package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockTestCategory struct {
	mock.Mock
}

func (m *MockTestCategory) GetAll(ctx context.Context) ([]entity.TestCategory, error) {
	args := m.Called(ctx)

	return args.Get(0).([]entity.TestCategory), args.Error(1)
}
