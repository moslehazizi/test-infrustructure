package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"
)

type TestServiceExecutorResultRepository interface {
	Create(ctx context.Context,
		executor *entity.Executor,
		dbInitializer database.DBInitializerFn,
	) error
}
