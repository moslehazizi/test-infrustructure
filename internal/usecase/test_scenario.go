package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type TestScenario interface {
	Create(ctx context.Context, testScenario *entity.TestScenario) error
	GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error)
	GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error)
	Start(ctx context.Context, id uint64) error
}

func NewTestScenarioUsecase(
	db database.Database,
	testScenarioRepository repository.TestScenarioRepository,
	testCategoryRepository repository.TestCategory,
	testServiceConfigRepository repository.TestServiceConfigRepository,
	motherService repository.MotherServiceRepository,
	scenarioExecutorEngine ScenarioExecutorEngine,
) TestScenario {
	return &testScenario{
		db:                          db,
		testScenarioRepository:      testScenarioRepository,
		testCategoryRepository:      testCategoryRepository,
		testServiceConfigRepository: testServiceConfigRepository,
		motherService:               motherService,
		scenarioExecutorEngine:      scenarioExecutorEngine,
	}
}

type testScenario struct {
	db                          database.Database
	testScenarioRepository      repository.TestScenarioRepository
	testCategoryRepository      repository.TestCategory
	testServiceConfigRepository repository.TestServiceConfigRepository
	motherService               repository.MotherServiceRepository
	scenarioExecutorEngine      ScenarioExecutorEngine
}

func (service *testScenario) Create(
	ctx context.Context,
	testScenario *entity.TestScenario,
) (e error) {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "create_test_scenario")
	defer span.End()

	span.SetAttributes(attribute.String("test_scenario.name", testScenario.Name), attribute.String("test_category.id", fmt.Sprintf("%d", testScenario.TestCategoryID)), attribute.String("mother_service.id", fmt.Sprintf("%d", testScenario.MotherServiceID)))

	_, err := service.motherService.GetByID(ctx, testScenario.MotherServiceID)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "mother_service_not_found"), attribute.String("error.message", err.Error()))
		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherService, err)
	}

	testCat, err := service.testCategoryRepository.GetByID(ctx, testScenario.TestCategoryID)
	if err != nil {
		if errors.Is(err, pkg.ErrTestCategoryNotFound) {
			span.SetAttributes(attribute.String("error.type", "test_category_not_found"))
			return pkg.ErrTestCategoryNotFound
		}

		span.SetAttributes(attribute.String("error.type", "get_test_category_error"), attribute.String("error.message", err.Error()))
		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategory, err)
	}

	err = testScenario.Validate(testCat)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "validation_error"), attribute.String("error.message", err.Error()))
		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testScenario.Status = entity.ScenarioStatusPending

	if testScenario.TestServiceConfig == nil {
		span.SetAttributes(attribute.String("error.type", "missing_test_service_config"))
		return pkg.ErrTestServiceConfigIsRequired
	}
	err = testScenario.TestServiceConfig.Validate()
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "test_service_config_validation_error"), attribute.String("error.message", err.Error()))
		return fmt.Errorf("%w: %w", pkg.ErrFailedToValidateTestSvcCfg, err)
	}

	tx := service.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)
	defer func() {
		if e != nil {
			_ = tx.Rollback()
			span.SetAttributes(attribute.String("transaction.status", "rolled_back"))
		}
	}()

	testSciID, err := service.testScenarioRepository.Create(dbCtx, testScenario)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "create_error"), attribute.String("error.message", err.Error()))
		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testScenario.TestServiceConfig.TestScenarioID = testSciID

	err = service.testServiceConfigRepository.Create(dbCtx, testScenario.TestServiceConfig)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "create_test_service_config_error"), attribute.String("error.message", err.Error()))
		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	_ = tx.Commit()
	span.SetAttributes(attribute.String("transaction.status", "committed"))

	return nil
}

func (service *testScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "get_test_scenario_by_id")
	defer span.End()

	span.SetAttributes(attribute.String("test_scenario.id", fmt.Sprintf("%d", id)))

	result, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))
			return nil, pkg.ErrTestScenarioNotFound
		}

		span.SetAttributes(attribute.String("error.type", "get_error"), attribute.String("error.message", err.Error()))
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	return result, nil
}

func (service *testScenario) GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error) {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "get_paginated_test_scenarios")
	defer span.End()

	span.SetAttributes(attribute.String("pagination.page", fmt.Sprintf("%d", pagReq.Page)), attribute.String("pagination.per_page", fmt.Sprintf("%d", pagReq.PerPage)))

	result, err := service.testScenarioRepository.GetPaginated(ctx, pagReq)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "get_paginated_error"), attribute.String("error.message", err.Error()))
		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetTestScenarios, err)
	}

	return result, nil
}

func (service *testScenario) Start(ctx context.Context, id uint64) error {
	tracer := otel.Tracer("test-scenario-usecase")
	_, span := tracer.Start(ctx, "start_test_scenario")
	defer span.End()

	span.SetAttributes(attribute.String("test_scenario.id", fmt.Sprintf("%d", id)))

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
		return fmt.Errorf("%w: %w", pkg.ErrFailedToSetScenarioStatusAsRunning, err)
	}

	// add scenario to executor.
	service.scenarioExecutorEngine.Add(scenario)

	return nil
}
