package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type ExecutorRepository struct {
	mock.Mock
}

func (m *ExecutorRepository) Create(ctx context.Context, executor *entity.Executor) error {
	args := m.Called(ctx, executor)
	return args.Error(0)
}
