package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type TestCategory interface {
	GetAll(ctx context.Context) ([]entity.TestCategory, error)
}
