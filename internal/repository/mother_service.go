package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type MotherServiceRepository interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
	GetByID(ctx context.Context, id uint64) (*entity.MotherService, error)
	GetAll(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error)
}
