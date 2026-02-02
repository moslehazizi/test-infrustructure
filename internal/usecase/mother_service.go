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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type MotherService interface {
	Create(ctx context.Context, motherService *entity.MotherService) error
	GetByID(ctx context.Context, id uint64) (*entity.MotherService, error)
	GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, int64, error)
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
	tracer := otel.Tracer("mother-service-usecase")
	_, span := tracer.Start(ctx, "create_mother_service")
	defer span.End()

	requestID := logger.GetRequestID(ctx)

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

	motherService.Status = entity.MotherServiceStatusPending
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

	err = service.motherServiceRepo.Create(dbCtx, motherService)
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

	// send kafka event for provisioning purpose
	bts, err := json.Marshal(&motherService)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "marshal_error"))
		zap.L().Error("failed to marshal mother service for event",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", motherService.ID),
			zap.Error(err),
		)

		return fmt.Errorf("failed to marshal mother service data to send event: %w", err)
	}
	err = service.eventProducer.SendEvent(ctx, bts, "provisioning")
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "event_send_error"))
		zap.L().Error("failed to send provisioning event",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", motherService.ID),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToSendProvisioningEvent, err)
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

	span.SetAttributes(attribute.String("service.id", fmt.Sprintf("%d", id)))

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

	span.SetAttributes(attribute.String("pagination.page", fmt.Sprintf("%d", paginationRequest.Page)), attribute.String("pagination.per_page", fmt.Sprintf("%d", paginationRequest.PerPage)))

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
