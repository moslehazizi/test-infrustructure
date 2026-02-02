package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
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
	requestID := logger.GetRequestID(ctx)
	err := motherService.Validate()
	if err != nil {
		zap.L().Warn("failed to validate mother service",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", motherService.Name),
			zap.Error(err),
		)

		return fmt.Errorf("failed to validate request: %w", err)
	}

	motherService.Status = entity.MotherServiceStatusPending

	tx := service.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)
	defer func() {
		if e != nil {
			_ = tx.Rollback()
			zap.L().Error("mother service creation failed, transaction rolled back",
				zap.String(logger.FieldRequestID, requestID),
				zap.String("name", motherService.Name),
				zap.Error(e),
			)
		}
	}()

	err = service.motherServiceRepo.Create(dbCtx, motherService)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceAlreadyExist) {
			zap.L().Warn("mother service already exists",
				zap.String(logger.FieldRequestID, requestID),
				zap.String("name", motherService.Name),
			)

			return pkg.ErrMotherServiceAlreadyExist
		}

		zap.L().Error("failed to create mother service in database",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", motherService.Name),
			zap.Error(err),
		)

		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateMotherService, err)
	}

	// send kafka event for provisioning purpose
	bts, err := json.Marshal(&motherService)
	if err != nil {
		zap.L().Error("failed to marshal mother service for event",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", motherService.ID),
			zap.Error(err),
		)

		return fmt.Errorf("failed to marshal mother service data to send event: %w", err)
	}
	err = service.eventProducer.SendEvent(ctx, bts, "provisioning")
	if err != nil {
		zap.L().Error("failed to send provisioning event",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", motherService.ID),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToSendProvisioningEvent, err)
	}

	_ = tx.Commit()
	zap.L().Info("mother service created successfully",
		zap.String(logger.FieldRequestID, requestID),
		zap.Uint64("id", motherService.ID),
		zap.String("name", motherService.Name),
	)

	return nil
}

func (service *motherService) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	requestID := logger.GetRequestID(ctx)
	result, err := service.motherServiceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceNotFound) {
			zap.L().Warn("mother service not found",
				zap.String(logger.FieldRequestID, requestID),
				zap.Uint64("id", id),
			)

			return nil, pkg.ErrMotherServiceNotFound
		}

		zap.L().Error("failed to get mother service by ID",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", id),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetMotherService, err)
	}

	return result, nil
}

func (service *motherService) GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error) {
	requestID := logger.GetRequestID(ctx)
	result, err := service.motherServiceRepo.GetPaginated(ctx, paginationRequest)
	if err != nil {
		zap.L().Error("failed to get paginated mother services",
			zap.String(logger.FieldRequestID, requestID),
			zap.Int("page", paginationRequest.Page),
			zap.Int("per_page", paginationRequest.PerPage),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherServices, err)
	}

	zap.L().Debug("retrieved paginated mother services",
		zap.String(logger.FieldRequestID, requestID),
		zap.Int("count", len(result)),
		zap.Int("page", paginationRequest.Page),
	)

	return result, nil
}
