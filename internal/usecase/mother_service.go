package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
)

type MotherService interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
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
	return nil
}
