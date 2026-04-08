package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"
)

type MotherServiceFactorialResultRepository interface {
	Create(
		ctx context.Context,
		factorial *entity.Factorial,
		dbInitializer database.DBInitializerFn,
	) error
}
