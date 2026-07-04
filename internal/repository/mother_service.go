package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type MotherServiceRepository interface {
	// CreateWithOutboxItem creates a mother service record and its
	// associated outbox item atomically, in a single local transaction.
	CreateWithOutboxItem(ctx context.Context, motherService *entity.MotherService, outboxItem *entity.Outbox) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*entity.MotherService, error)
	GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, int64, error)
	SetStatus(ctx context.Context, id uint64, status entity.MotherServiceStatus) error
}
