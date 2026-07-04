package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	provision "control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/internal/usecase/interfaces"
	"control-panel-service/pkg"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func NewMotherService(
	motherServiceRepo repository.MotherServiceRepository,
	testScenarioRepo repository.TestScenarioRepository,
	provisioningService provision.ProvisioningService,
	stressTestExecutionManager interfaces.ExecutionManager,
	outboxMaxAttempts int,
) *motherService {
	return &motherService{
		motherServiceRepo,
		testScenarioRepo,
		provisioningService,
		stressTestExecutionManager,
		outboxMaxAttempts,
	}
}

type motherService struct {
	motherServiceRepo          repository.MotherServiceRepository
	testScenarioRepo           repository.TestScenarioRepository
	provisioningService        provision.ProvisioningService
	stressTestExecutionManager interfaces.ExecutionManager
	outboxMaxAttempts          int
}

func (service *motherService) Create(ctx context.Context, motherService *entity.MotherService) (e error) {
	tracer := otel.Tracer("mother-service-usecase")
	useCaseCTX, span := tracer.Start(ctx, "create-mother-service-usecase")
	defer span.End()

	err := motherService.Validate()
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "validation_error"))

		zap.L().Warn("failed to validate mother service",
			zap.String("name", motherService.Name),
			zap.Error(err),
		)

		return fmt.Errorf("failed to validate request: %w", err)
	}

	// Mother service is created in "ready" status. Provisioning it in
	// Kubernetes is slow, network-bound, and not itself transactional, so it
	// is not performed here: instead the repository records an outbox item
	// in the same local transaction as the insert, and a background worker
	// (see usecase.outboxProcessor) claims it afterwards and calls
	// ProvisioningService.ProvisionMotherService outside of any DB
	// transaction, retrying on failure and updating status when it settles.
	motherService.Status = entity.MotherServiceStatusReady
	span.SetAttributes(attribute.String("service.name", motherService.Name))

	outboxItem := &entity.Outbox{
		OperationType: entity.OutboxOperationProvisionMotherService,
		Status:        entity.OutboxStatusPending,
		MaxAttempts:   service.outboxMaxAttempts,
		AvailableAt:   time.Now(),
	}

	id, err := service.motherServiceRepo.CreateWithOutboxItem(useCaseCTX, motherService, outboxItem)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceAlreadyExist) {
			span.SetAttributes(attribute.String("error.type", "already_exists"))
			zap.L().Warn("mother service already exists",
				zap.String("name", motherService.Name),
			)

			return pkg.ErrMotherServiceAlreadyExist
		}

		if errors.Is(err, pkg.ErrFailedToCreateOutboxItem) {
			span.SetAttributes(attribute.String("error.type", "outbox_create_error"), attribute.String("error.message", err.Error()))
			zap.L().Error("failed to record mother service provisioning outbox item",
				zap.String("name", motherService.Name),
				zap.Error(err),
			)

			return err
		}

		span.SetAttributes(attribute.String("error.type", "create_error"), attribute.String("error.message", err.Error()))
		zap.L().Error("failed to create mother service in database",
			zap.String("name", motherService.Name),
			zap.Error(err),
		)

		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateMotherService, err)
	}

	motherService.ID = id

	zap.L().Info("mother service created successfully, provisioning queued",
		zap.Uint64("id", motherService.ID),
		zap.String("name", motherService.Name),
	)

	return nil
}

func (service *motherService) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	tracer := otel.Tracer("mother-service-usecase")
	useCaseCTX, span := tracer.Start(ctx, "get-mother-service-by-id-usecase")
	defer span.End()

	span.SetAttributes(attribute.String("service.id", strconv.FormatUint(id, 10)))

	result, err := service.motherServiceRepo.GetByID(useCaseCTX, id)
	if err != nil {
		if errors.Is(err, pkg.ErrMotherServiceNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))
			zap.L().Warn("mother service not found",
				zap.Uint64("id", id),
			)

			return nil, pkg.ErrMotherServiceNotFound
		}

		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get mother service by ID",
			zap.Uint64("id", id),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetMotherService, err)
	}

	return result, nil
}

func (service *motherService) GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, int64, error) {
	tracer := otel.Tracer("mother-service-usecase")
	useCaseCTX, span := tracer.Start(ctx, "get-paginated-mother-services-usecase")
	defer span.End()

	span.SetAttributes(attribute.String("pagination.page", strconv.Itoa(paginationRequest.Page)), attribute.String("pagination.per_page", strconv.Itoa(paginationRequest.PerPage)))

	result, count, err := service.motherServiceRepo.GetPaginated(useCaseCTX, paginationRequest)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_paginated_error"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get paginated mother services",
			zap.Int("page", paginationRequest.Page),
			zap.Int("per_page", paginationRequest.PerPage),
			zap.Error(err),
		)

		return nil, 0, fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherServices, err)
	}

	zap.L().Debug("retrieved paginated mother services",
		zap.Int64("count", count),
		zap.Int("page", paginationRequest.Page),
	)

	return result, count, nil
}

func (service *motherService) Delete(ctx context.Context, id uint64) error {
	// get mother service by its id
	motherService, err := service.motherServiceRepo.GetByID(ctx, id)
	if err != nil {
		zap.L().Error("failed to get mother service", zap.Uint64("id", motherService.ID), zap.String("name", motherService.Name), zap.Error(err))

		return err
	}

	// deprovision mother service
	err = service.provisioningService.DeprovisionMotherService(ctx, motherService)
	if err != nil {
		zap.L().Error("failed to deprovision mother service", zap.Uint64("id", motherService.ID), zap.String("name", motherService.Name), zap.Error(err))

		return err
	}

	// update status
	err = service.motherServiceRepo.SetStatus(ctx, motherService.ID, entity.MotherServiceStatusDeleted)
	if err != nil {
		zap.L().Error("failed to update mother service status", zap.Uint64("id", motherService.ID), zap.String("name", motherService.Name), zap.Error(err))

		return err
	}

	// get test scenarios by mother service id
	testScenarios, err := service.testScenarioRepo.GetByMotherServiceId(ctx, motherService.ID)
	if err != nil {
		zap.L().Error("failed to get test scenarios by mother service id", zap.Uint64("id", motherService.ID), zap.String("name", motherService.Name), zap.Error(err))

		return err
	}

	for _, testScenario := range testScenarios {
		// delete  scenarios
		err := service.stressTestExecutionManager.DeleteScenario(ctx, testScenario)
		if err != nil {
			zap.L().Error("failed to deprovision test scenario", zap.Uint64("mother_service_id", motherService.ID), zap.Uint64("test_scenario_id", testScenario.ID), zap.String("name", motherService.Name), zap.Error(err))
		}

		// update status
		err = service.testScenarioRepo.SetStatus(ctx, testScenario.ID, entity.ScenarioStatusPending, true)
		if err != nil {
			zap.L().Error("failed to update scenario test status", zap.Uint64("mother_service_id", motherService.ID), zap.Uint64("test_scenario_id", testScenario.ID), zap.String("name", motherService.Name), zap.Error(err))
		}
	}

	return nil
}
