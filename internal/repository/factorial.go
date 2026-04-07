package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type FactorialRepository interface {
	Create(ctx context.Context, factorial *entity.Factorial) error
}
