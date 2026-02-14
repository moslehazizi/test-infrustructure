package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
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
	// ResetOrphanedScenarios recovers scenarios that were in running state
	// when the service crashed and ensures they're added to the in-memory executor box.
	ResetOrphanedScenarios(ctx context.Context) error
	DeprovisionAllPods(ctx context.Context) error
}

func NewTestScenarioUsecase(
	db database.Database,
	testScenarioRepository repository.TestScenarioRepository,
	testCategoryRepository repository.TestCategory,
	testServiceConfigRepository repository.TestServiceConfigRepository,
	motherService repository.MotherServiceRepository,
	scenarioExecutorBox interfaces.ScenarioExecutorBox,
	testServiceRepo repository.TestServiceRepository,
	provisioningService provider.ProvisioningService,
) TestScenario {
	return &testScenario{
		db:                          db,
		testScenarioRepository:      testScenarioRepository,
		testCategoryRepository:      testCategoryRepository,
		testServiceConfigRepository: testServiceConfigRepository,
		motherService:               motherService,
		scenarioExecutorBox:         scenarioExecutorBox,
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
	scenarioExecutorBox         interfaces.ScenarioExecutorBox
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

	testScenario.Status = entity.ScenarioStatusPending

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

	tx := service.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)
	defer func() {
		if e != nil {
			_ = tx.Rollback()
			span.SetAttributes(attribute.String("transaction.status", "rolled_back"))
			zap.L().Error("test scenario creation failed, transaction rolled back",
				zap.String(logger.FieldRequestID, requestID),
				zap.String("name", testScenario.Name),
				zap.Error(e),
			)
		}
	}()

	testSciID, err := service.testScenarioRepository.Create(dbCtx, testScenario)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "create_error"), attribute.String("error.message", err.Error()))
		zap.L().Error("failed to create test scenario in database",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", testScenario.Name),
			zap.Error(err),
		)

		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testScenario.TestServiceConfig.TestScenarioID = testSciID

	err = service.testServiceConfigRepository.Create(dbCtx, testScenario.TestServiceConfig)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "create_test_service_config_error"), attribute.String("error.message", err.Error()))
		zap.L().Error("failed to create test service config",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("test_scenario_id", testSciID),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	_ = tx.Commit()

	span.SetAttributes(attribute.String("transaction.status", "committed"))
	zap.L().Info("test scenario created successfully",
		zap.String(logger.FieldRequestID, requestID),
		zap.Uint64("id", testSciID),
		zap.String("name", testScenario.Name),
	)

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
	if scenario.Status != entity.ScenarioStatusPending {
		span.SetAttributes(attribute.String("error.type", "invalid_status"))

		return pkg.ErrOnlyPendingScenariosCanBeStarted
	}

	// mark scenario as running
	err = service.testScenarioRepository.SetStatus(ctx, id, entity.ScenarioStatusRunning)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "set_status_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatus, err)
	}

	// add scenario to executor.
	service.scenarioExecutorBox.Add(
		NewScenarioExecutor(
			*scenario,
			NewScenarioTypeRunnerGroupA(
				service.testServiceRepo,
				service.provisioningService,
				service.testScenarioRepository,
			),
		))

	return nil
}

func (service *testScenario) ResetOrphanedScenarios(ctx context.Context) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "reset_orphaned_scenarios")
	defer span.End()

	runningScenarios, err := service.testScenarioRepository.GetByStatus(ctx, entity.ScenarioStatusRunning)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_by_status_error"), attribute.String("error.message", err.Error()))
		zap.L().Error("failed to get running test scenarios", zap.Error(err))

		return fmt.Errorf("failed to get running test scenarios: %w", err)
	}

	for _, sc := range runningScenarios {
		if service.scenarioExecutorBox.HasExecutor(sc.ID) {
			continue
		}

		service.scenarioExecutorBox.Add(NewScenarioExecutor(
			*sc,
			NewScenarioTypeRunnerGroupA(
				service.testServiceRepo,
				service.provisioningService,
				service.testScenarioRepository,
			),
		))
		zap.L().Info("recovered running scenario and added to executor box", zap.Uint64("id", sc.ID))
	}

	return nil
}

func (service *testScenario) DeprovisionAllPods(ctx context.Context) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "deprovision_test_scenario")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	scenarios, err := service.testScenarioRepository.GetByStatus(ctx, entity.ScenarioStatusRunning)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenariosByStatus, err)
	}

	for _, scenario := range scenarios {
		err := service.provisioningService.DeprovisionTestService(ctx, scenario, scenario.DeploymentNumber)

		if err != nil {
			span.SetAttributes(attribute.String("error.type", "deprovision_error"), attribute.String("error.message", err.Error()))
			zap.L().Error("failed to deprovision running test scenarios", zap.Error(err))

			continue
		}

		err = service.testScenarioRepository.SetStatus(ctx, scenario.ID, entity.ScenarioStatusAborted)

		if err != nil {
			span.SetAttributes(attribute.String("error.type", "set_status"), attribute.String("error.message", err.Error()))
			zap.L().Error("failed to set status running test scenarios as aborted", zap.Error(err))

			continue
		}
	}

	return nil
}
