package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/internal/server/dto/request"
	"control-panel-service/internal/usecase/interfaces"
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

type TestScenario interface {
	Create(ctx context.Context, testScenario *entity.TestScenario) error
	GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error)
	GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, int64, error)
	Start(ctx context.Context, id uint64) error
	Pause(ctx context.Context, id uint64) error
	Resume(ctx context.Context, id uint64) error
	Stop(ctx context.Context, id uint64) error
	Delete(ctx context.Context, id uint64) error
	// ResetOrphanedScenarios recovers scenarios that were in running state
	// when the service crashed and ensures they're added to the in-memory executor box.
	// ResetOrphanedScenarios(ctx context.Context) error
	Update(ctx context.Context, testScenarioUpdateRequest *request.TestScenarioUpdateRequest) (e error)
}

func NewTestScenarioUsecase(
	db database.Database,
	testScenarioRepository repository.TestScenarioRepository,
	testCategoryRepository repository.TestCategory,
	testServiceConfigRepository repository.TestServiceConfigRepository,
	motherService repository.MotherServiceRepository,
	stressTestExecutionManager interfaces.ExecutionManager,
	testServiceRepo repository.TestServiceRepository,
	provisioningService provider.ProvisioningService,
) TestScenario {
	return &testScenario{
		db:                          db,
		testScenarioRepository:      testScenarioRepository,
		testCategoryRepository:      testCategoryRepository,
		testServiceConfigRepository: testServiceConfigRepository,
		motherService:               motherService,
		stressTestExecutionManager:  stressTestExecutionManager,
		testServiceRepo:             testServiceRepo,
		provisioningService:         provisioningService,
	}
}

type testScenario struct {
	db                          database.Database
	testScenarioRepository      repository.TestScenarioRepository
	testCategoryRepository      repository.TestCategory
	testServiceConfigRepository repository.TestServiceConfigRepository
	motherService               repository.MotherServiceRepository
	stressTestExecutionManager  interfaces.ExecutionManager
	testServiceRepo             repository.TestServiceRepository
	provisioningService         provider.ProvisioningService
}

func (service *testScenario) Create(ctx context.Context, testScenario *entity.TestScenario) (e error) {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "create_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("test_scenario.name", testScenario.Name), attribute.String("test_category.id", strconv.FormatUint(testScenario.TestCategoryID, 10)), attribute.String("mother_service.id", strconv.FormatUint(testScenario.MotherServiceID, 10)))

	_, err := service.motherService.GetByID(ctx, testScenario.MotherServiceID)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "mother_service_not_found"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get mother service for test scenario",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("mother_service_id", testScenario.MotherServiceID),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherService, err)
	}

	testCat, err := service.testCategoryRepository.GetByID(ctx, testScenario.TestCategoryID)
	if err != nil {
		if errors.Is(err, pkg.ErrTestCategoryNotFound) {
			span.SetAttributes(attribute.String("error.type", "test_category_not_found"))
			zap.L().Warn("test category not found",
				zap.String(logger.FieldRequestID, requestID),
				zap.Uint64("test_category_id", testScenario.TestCategoryID),
			)

			return pkg.ErrTestCategoryNotFound
		}

		span.SetAttributes(attribute.String("error.type", "get_test_category_error"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get test category",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("test_category_id", testScenario.TestCategoryID),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategory, err)
	}

	if !testCat.HasNumSteps {
		testScenario.NumSteps = int64(1)
	}

	err = testScenario.Validate(testCat)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "validation_error"), attribute.String("error.message", err.Error()))
		zap.L().Warn("failed to validate test scenario",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", testScenario.Name),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testScenario.Status = entity.ScenarioStatusReady

	if testScenario.TestServiceConfig == nil {
		span.SetAttributes(attribute.String("error.type", "missing_test_service_config"))
		zap.L().Warn("test service config is required",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", testScenario.Name),
		)

		return pkg.ErrTestServiceConfigIsRequired
	}
	err = testScenario.TestServiceConfig.Validate()
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "test_service_config_validation_error"), attribute.String("error.message", err.Error()))
		zap.L().Warn("failed to validate test service config",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", testScenario.Name),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToValidateTestSvcCfg, err)
	}

	err = service.testScenarioRepository.CreateScenarioAndConfig(ctx, testScenario)
	if err != nil {
		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	return nil
}

func (service *testScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "get_test_scenario_by_id")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	result, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))
			zap.L().Warn("test scenario not found",
				zap.String(logger.FieldRequestID, requestID),
				zap.Uint64("id", id),
			)

			return nil, pkg.ErrTestScenarioNotFound
		}

		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))
		zap.L().Error("failed to get test scenario by ID",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", id),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	return result, nil
}

func (service *testScenario) GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, int64, error) {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "get_paginated_test_scenarios")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("pagination.page", strconv.Itoa(pagReq.Page)), attribute.String("pagination.per_page", strconv.Itoa(pagReq.PerPage)))

	result, count, err := service.testScenarioRepository.GetPaginated(ctx, pagReq)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_paginated_error"), attribute.String("error.message", err.Error()))

		zap.L().Error("failed to get paginated test scenarios",
			zap.String(logger.FieldRequestID, requestID),
			zap.Int("page", pagReq.Page),
			zap.Int("per_page", pagReq.PerPage),
			zap.Error(err),
		)

		return nil, 0, fmt.Errorf("%w, %w", pkg.ErrFailedToGetTestScenarios, err)
	}

	zap.L().Debug("retrieved paginated test scenarios",
		zap.String(logger.FieldRequestID, requestID),
		zap.Int("count", len(result)),
		zap.Int("page", pagReq.Page),
	)

	return result, count, nil
}

func (service *testScenario) Start(ctx context.Context, id uint64) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "start_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	scenario, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return pkg.ErrTestScenarioNotFound
		}
		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	// make sure scenario has correct status
	if scenario.Status != entity.ScenarioStatusReady {
		span.SetAttributes(attribute.String("error.type", "invalid_status"))

		return pkg.ErrOnlyReadyScenariosCanBeStarted
	}

	// mark scenario as running
	err = service.testScenarioRepository.SetStatus(ctx, id, entity.ScenarioStatusRunning, false)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "set_status_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
	}

	switch scenario.TestCategory.Name {
	case entity.STRESS:
		err := service.stressTestExecutionManager.RunScenario(ctx, scenario)
		if err != nil {
			return fmt.Errorf("%w: %w", pkg.ErrFailedToRunScenarioInExecutionManager, err)
		}
	default:
		return pkg.ErrStartingTestNotImplemented
	}

	return nil
}

func (service *testScenario) Pause(ctx context.Context, id uint64) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "pause_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	scenario, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return pkg.ErrTestScenarioNotFound
		}
		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	// make sure scenario has correct status
	if scenario.Status != entity.ScenarioStatusRunning {
		span.SetAttributes(attribute.String("error.type", "invalid_status"))

		return pkg.ErrOnlyRunningScenariosCanBePaused
	}

	switch scenario.TestCategory.Name {
	case entity.STRESS:
		err := service.stressTestExecutionManager.PauseScenario(ctx, scenario)
		if err != nil {
			return fmt.Errorf("%w: %w", pkg.ErrFailedToPauseScenarioToExecutionManager, err)
		}
	default:
		return pkg.ErrPausingTestNotImplemented
	}

	// mark scenario as paused
	err = service.testScenarioRepository.SetStatus(ctx, id, entity.ScenarioStatusPaused, false)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "set_status_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
	}

	return nil
}

func (service *testScenario) Resume(ctx context.Context, id uint64) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "resume_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	scenario, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return pkg.ErrTestScenarioNotFound
		}
		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	// make sure scenario has correct status
	if scenario.Status != entity.ScenarioStatusPaused {
		span.SetAttributes(attribute.String("error.type", "invalid_status"))

		return pkg.ErrOnlyPausedScenariosCanBeResume
	}

	switch scenario.TestCategory.Name {
	case entity.STRESS:
		err := service.stressTestExecutionManager.ResumeScenario(ctx, scenario)
		if err != nil {
			return fmt.Errorf("%w: %w", pkg.ErrFailedToResumeScenarioToExecutionManager, err)
		}
	default:
		return pkg.ErrResumingTestNotImplemented
	}

	// mark scenario as running
	err = service.testScenarioRepository.SetStatus(ctx, id, entity.ScenarioStatusRunning, false)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "set_status_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
	}

	return nil
}

func (service *testScenario) Stop(ctx context.Context, id uint64) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "stop_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	scenario, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return pkg.ErrTestScenarioNotFound
		}
		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	// make sure scenario has correct status
	if scenario.Status != entity.ScenarioStatusRunning && scenario.Status != entity.ScenarioStatusPaused {
		span.SetAttributes(attribute.String("error.type", "invalid_status"))

		return pkg.ErrOnlyRunAndPauseScenariosCanBeStop
	}

	switch scenario.TestCategory.Name {
	case entity.STRESS:
		err := service.stressTestExecutionManager.StopScenario(ctx, scenario)
		if err != nil {
			return fmt.Errorf("%w: %w", pkg.ErrFailedToStopScenarioToExecutionManager, err)
		}
	default:
		return pkg.ErrStoppingTestNotImplemented
	}

	// mark scenario as running
	err = service.testScenarioRepository.SetStatus(ctx, id, entity.ScenarioStatusReady, false)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "set_status_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
	}

	return nil
}

func (service *testScenario) Delete(ctx context.Context, id uint64) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "delete_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("test_scenario.id", strconv.FormatUint(id, 10)))

	scenario, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return pkg.ErrTestScenarioNotFound
		}
		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	// make sure scenario has correct status
	if scenario.Status != entity.ScenarioStatusReady && scenario.Status != entity.ScenarioStatusPending {
		span.SetAttributes(attribute.String("error.type", "invalid_status"))

		return pkg.ErrScenariosCanNotBeDelete
	}

	// mark scenario as running
	err = service.testScenarioRepository.SetStatus(ctx, id, entity.ScenarioStatusDeleted, false)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "set_status_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
	}

	switch scenario.TestCategory.Name {
	case entity.STRESS:
		err := service.stressTestExecutionManager.DeleteScenario(ctx, scenario)
		if err != nil {
			return fmt.Errorf("%w: %w", pkg.ErrFailedToDeleteScenarioToExecutionManager, err)
		}
	default:
		return pkg.ErrDeleteTestNotImplemented
	}

	return nil
}

func (service *testScenario) Update(ctx context.Context, testScenarioUpdateRequest *request.TestScenarioUpdateRequest) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "update_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(
		attribute.String("request_id", requestID),
		attribute.String("test_scenario.id", strconv.FormatUint(testScenarioUpdateRequest.ID, 10)),
	)

	if testScenarioUpdateRequest.Config == nil {
		return pkg.ErrTestServiceConfigIsRequired
	}

	existing, err := service.testScenarioRepository.GetByID(ctx, testScenarioUpdateRequest.ID)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrTestScenarioNotFound, err)
	}

	if existing.Status != entity.ScenarioStatusReady && existing.Status != entity.ScenarioStatusPending {
		span.SetAttributes(attribute.String("error.type", "invalid_status"))

		return pkg.ErrScenariosCanNotBeUpdated
	}

	motherService, err := service.motherService.GetByID(ctx, testScenarioUpdateRequest.MotherServiceID)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherService, err)
	}

	if existing.TestServiceConfig == nil || existing.TestServiceConfig.ID == 0 {
		return fmt.Errorf("%w: test service config not loaded for scenario %d", pkg.ErrFailedToGetTestServiceConfig, existing.ID)
	}

	// TODO: convert request to entity instead
	existing.ApplyUpdateRequest(testScenarioUpdateRequest, motherService)
	if err = existing.Validate(existing.TestCategory); err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToUpdateTestScenario, err)
	}

	// TODO: convert request to entity instead
	existing.TestServiceConfig.ApplyUpdateFromRequest(testScenarioUpdateRequest.Config)
	if err = existing.TestServiceConfig.Validate(); err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToValidateTestSvcCfg, err)
	}

	err = service.testScenarioRepository.UpdateScenarioAndConfig(ctx, existing)
	if err != nil {
		span.SetAttributes(attribute.String("update.status", "failed"))
		zap.L().Error("test scenario and config update failed",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", existing.Name),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToUpdateTestScenario, err)
	}

	return nil
}
