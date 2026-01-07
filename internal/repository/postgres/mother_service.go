package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type motherServiceRepository struct {
	db *gorm.DB
}

func NewMotherServiceRepository(db *gorm.DB) repository.MotherServiceRepository {
	return &motherServiceRepository{
		db: db,
	}
}

func (m *motherServiceRepository) Create(ctx context.Context, motherService *entity.MotherService) error {
	err := m.db.WithContext(ctx).Create(motherService).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateMotherService, pkg.ErrMotherServiceAlreadyExist)
		}

		return fmt.Errorf("failed to create mother service record: %w", err)
	}

	return nil
}

func (m *motherServiceRepository) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	var motherService entity.MotherService
	err := m.db.WithContext(ctx).First(&motherService, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, fmt.Errorf("failed to get mother service record: %w", err)
	}

	return &motherService, nil
}

func (m *motherServiceRepository) GetAll(ctx context.Context) ([]*entity.MotherService, error) {
	var motherServices []*entity.MotherService
	err := m.db.WithContext(ctx).Order("id DESC").Find(&motherServices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get mother service records: %w", err)
	}

	return motherServices, nil
}
