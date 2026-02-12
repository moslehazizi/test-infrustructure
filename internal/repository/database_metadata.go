package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type DatabaseMetadata interface {
	GetAll(ctx context.Context) ([]string, error)
	GetTablesByDBName(ctx context.Context, dbName string) (*entity.TablesByType, error)
}
