package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type MotherServiceRepository interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
}
