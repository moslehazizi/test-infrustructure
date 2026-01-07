package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"errors"
	"fmt"
)

type MotherService interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
	GetByID(ctx context.Context, id uint64) (*entity.MotherService, error)
	GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error)
}

func NewMotherService(motherServiceRepo repository.MotherServiceRepository) MotherService {
	return &motherService{
		motherServiceRepo: motherServiceRepo,
	}
}

type motherService struct {
	motherServiceRepo repository.MotherServiceRepository
}

func (service *motherService) Create(ctx context.Context, motherService *entity.MotherService) error {
	err := service.motherServiceRepo.Create(ctx, motherService)
	if err != nil {
		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateMotherService, err)
	}

	// TODO: create instance in k8s

	return nil
}

func (service *motherService) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	result, err := service.motherServiceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceNotFound) {
			return nil, pkg.ErrMotherServiceNotFound
		}

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetMotherService, err)
	}

	return result, nil
}

func (service *motherService) GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error) {
	result, err := service.motherServiceRepo.GetPaginated(ctx, paginationRequest)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherServices, err)
	}

	return result, nil
}
