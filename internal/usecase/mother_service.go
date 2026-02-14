package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	provision "control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/logger"
	"errors"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type MotherService interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
	GetByID(ctx context.Context, id uint64) (*entity.MotherService, error)
	GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, int64, error)
	DeprovisionAllPods(ctx context.Context) error
}

func NewMotherService(
	db database.Database,
	motherServiceRepo repository.MotherServiceRepository,
	provisioningService provision.ProvisioningService,
) MotherService {
	return &motherService{
		db,
		motherServiceRepo,
		provisioningService,
	}
}

type motherService struct {
	db                  database.Database
	motherServiceRepo   repository.MotherServiceRepository
	provisioningService provision.ProvisioningService
}

func (service *motherService) Create(ctx context.Context, motherService *entity.MotherService) (e error) {
	tracer := otel.Tracer("mother-service-usecase")
	_, span := tracer.Start(ctx, "create_mother_service")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	err := motherService.Validate()
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "validation_error"))

		zap.L().Warn("failed to validate mother service",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", motherService.Name),
			zap.Error(err),
		)

		return fmt.Errorf("failed to validate request: %w", err)
	}

	// for now when we create a mother service in database in the same time we send it to provision.
	// this operation done in a transaction so if mother service created and provisioned it is in running status.
	motherService.Status = entity.MotherServiceStatusRunning
	span.SetAttributes(attribute.String("service.name", motherService.Name))

	tx := service.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)
	defer func() {
		if e != nil {
			_ = tx.Rollback()
			span.SetAttributes(attribute.String("transaction.status", "rolled_back"))
			zap.L().Error("mother service creation failed, transaction rolled back",
				zap.String(logger.FieldRequestID, requestID),
				zap.String("name", motherService.Name),
				zap.Error(e),
			)
		}
	}()

	id, err := service.motherServiceRepo.Create(dbCtx, motherService)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceAlreadyExist) {
			span.SetAttributes(attribute.String("error.type", "already_exists"))
			zap.L().Warn("mother service already exists",
				zap.String(logger.FieldRequestID, requestID),
				zap.String("name", motherService.Name),
			)

			return pkg.ErrMotherServiceAlreadyExist
		}

		span.SetAttributes(attribute.String("error.type", "create_error"), attribute.String("error.message", err.Error()))
		zap.L().Error("failed to create mother service in database",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", motherService.Name),
			zap.Error(err),
		)

		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateMotherService, err)
	}

	motherService.ID = id
	err = service.provisioningService.ProvisionMotherService(dbCtx, motherService)
	if err != nil {
		return fmt.Errorf("%w, %w", pkg.ErrFailedToDeployMotherService, err)
	}

	_ = tx.Commit()

	span.SetAttributes(attribute.String("transaction.status", "committed"))
	zap.L().Info("mother service created successfully",
		zap.String(logger.FieldRequestID, requestID),
		zap.Uint64("id", motherService.ID),
		zap.String("name", motherService.Name),
	)

	return nil
}

func (service *motherService) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	tracer := otel.Tracer("mother-service-usecase")
	_, span := tracer.Start(ctx, "get_mother_service_by_id")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("service.id", strconv.FormatUint(id, 10)))

	result, err := service.motherServiceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))
			zap.L().Warn("mother service not found",
				zap.String(logger.FieldRequestID, requestID),
				zap.Uint64("id", id),
			)

			return nil, pkg.ErrMotherServiceNotFound
		}

		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get mother service by ID",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", id),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetMotherService, err)
	}

	return result, nil
}

func (service *motherService) GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, int64, error) {
	tracer := otel.Tracer("mother-service-usecase")
	_, span := tracer.Start(ctx, "get_paginated_mother_services")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("pagination.page", strconv.Itoa(paginationRequest.Page)), attribute.String("pagination.per_page", strconv.Itoa(paginationRequest.PerPage)))

	result, count, err := service.motherServiceRepo.GetPaginated(ctx, paginationRequest)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_paginated_error"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get paginated mother services",
			zap.String(logger.FieldRequestID, requestID),
			zap.Int("page", paginationRequest.Page),
			zap.Int("per_page", paginationRequest.PerPage),
			zap.Error(err),
		)

		return nil, 0, fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherServices, err)
	}

	zap.L().Debug("retrieved paginated mother services",
		zap.String(logger.FieldRequestID, requestID),
		zap.Int64("count", count),
		zap.Int("page", paginationRequest.Page),
	)

	return result, count, nil
}

func (service *motherService) DeprovisionAllPods(ctx context.Context) error {
	motherservices, count, err := service.motherServiceRepo.GetPaginated(ctx, entity.PaginationRequest{Page: 0, PerPage: 0})
	if err != nil {
		zap.L().Error("failed to get mother services for deprovisioning", zap.Error(err))

		return fmt.Errorf("failed to get mother services for deprovisioning: %w", err)
	}

	if count == 0 {
		zap.L().Info("no mother services found for deprovisioning")

		return nil
	}

	for _, ms := range motherservices {
		err := service.provisioningService.DeprovisionMotherService(ctx, ms)
		if err != nil {
			zap.L().Error("failed to deprovision mother service", zap.Uint64("id", ms.ID), zap.String("name", ms.Name), zap.Error(err))

			continue
		}

		if ms.Status != entity.MotherServiceStatusAborted {
			err = service.motherServiceRepo.SetStatus(ctx, ms.ID, entity.MotherServiceStatusAborted)
			if err != nil {
				zap.L().Error("failed to update mother service status", zap.Uint64("id", ms.ID), zap.String("name", ms.Name), zap.Error(err))

				continue
			}
		}

	}

	return nil
}
