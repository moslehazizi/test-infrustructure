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

func (m *MockTestCategory) GetByID(ctx context.Context, id int) (*entity.TestCategory, error) {
	args := m.Called(ctx, id)

	var result *entity.TestCategory
	if args.Get(0) != nil {
		result = args.Get(0).(*entity.TestCategory)
	}

	return result, args.Error(1)
}
