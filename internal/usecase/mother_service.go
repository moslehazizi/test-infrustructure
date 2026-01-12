package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"encoding/json"
	"errors"
	"fmt"
)

type MotherService interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
	GetByID(ctx context.Context, id uint64) (*entity.MotherService, error)
	GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error)
}

func NewMotherService(db database.Database, motherServiceRepo repository.MotherServiceRepository, eventProducer provider.EventProducer) MotherService {
	return &motherService{
		db,
		motherServiceRepo,
		eventProducer,
	}
}

type motherService struct {
	db                database.Database
	motherServiceRepo repository.MotherServiceRepository
	eventProducer     provider.EventProducer
}

func (service *motherService) Create(ctx context.Context, motherService *entity.MotherService) (e error) {
	err := motherService.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate request: %w", err)
	}

	motherService.Status = entity.MotherServiceStatusPending

	tx := service.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)
	defer func() {
		if e != nil {
			_ = tx.Rollback()
		}
	}()

	err = service.motherServiceRepo.Create(dbCtx, motherService)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceAlreadyExist) {
			return pkg.ErrMotherServiceAlreadyExist
		}

		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateMotherService, err)
	}

	// send kafka event for provisioning purpose
	bts, err := json.Marshal(&motherService)
	if err != nil {
		return fmt.Errorf("failed to marshal mother service data to send event: %w", err)
	}
	err = service.eventProducer.SendEvent(ctx, bts, "provisioning")
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToSendProvisioningEvent, err)
	}

	_ = tx.Commit()

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
