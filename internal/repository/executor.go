package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type ExecutorRepository interface {
	Create(ctx context.Context, executor *entity.Executor) error
}
