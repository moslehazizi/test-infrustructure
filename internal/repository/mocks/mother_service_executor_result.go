package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"

	"github.com/stretchr/testify/mock"
)

type MockMotherServiceExecutorResultRepository struct {
	mock.Mock
}

func (m *MockMotherServiceExecutorResultRepository) Create(ctx context.Context, executor *entity.Executor, dbInitializer database.DBInitializerFn) error {
	args := m.Called(ctx, executor, dbInitializer)

	return args.Error(0)
}
