package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/logger"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type TestScenario interface {
	Create(ctx context.Context, testScenario *entity.TestScenario) error
	GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error)
	GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error)
}

func NewTestScenarioUsecase(
	db database.Database,
	testScenarioRepository repository.TestScenarioRepository,
	testCategoryRepository repository.TestCategory,
	testServiceConfigRepository repository.TestServiceConfigRepository,
	motherService repository.MotherServiceRepository) TestScenario {
	return &testScenario{
		db:                          db,
		testScenarioRepository:      testScenarioRepository,
		testCategoryRepository:      testCategoryRepository,
		testServiceConfigRepository: testServiceConfigRepository,
		motherService:               motherService,
	}
}

type testScenario struct {
	db                          database.Database
	testScenarioRepository      repository.TestScenarioRepository
	testCategoryRepository      repository.TestCategory
	testServiceConfigRepository repository.TestServiceConfigRepository
	motherService               repository.MotherServiceRepository
}

func (service *testScenario) Create(
	ctx context.Context,
	testScenario *entity.TestScenario,
) (e error) {
	requestID := logger.GetRequestID(ctx)
	_, err := service.motherService.GetByID(ctx, testScenario.MotherServiceID)
	if err != nil {
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
			zap.L().Warn("test category not found",
				zap.String(logger.FieldRequestID, requestID),
				zap.Uint64("test_category_id", testScenario.TestCategoryID),
			)

			return pkg.ErrTestCategoryNotFound
		}

		zap.L().Error("failed to get test category",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("test_category_id", testScenario.TestCategoryID),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategory, err)
	}

	err = testScenario.Validate(testCat)
	if err != nil {
		zap.L().Warn("failed to validate test scenario",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", testScenario.Name),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testScenario.Status = entity.ScenarioStatusPending

	if testScenario.TestServiceConfig == nil {
		zap.L().Warn("test service config is required",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("name", testScenario.Name),
		)

		return pkg.ErrTestServiceConfigIsRequired
	}
	err = testScenario.TestServiceConfig.Validate()
	if err != nil {
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
			zap.L().Error("test scenario creation failed, transaction rolled back",
				zap.String(logger.FieldRequestID, requestID),
				zap.String("name", testScenario.Name),
				zap.Error(e),
			)
		}
	}()

	testSciID, err := service.testScenarioRepository.Create(dbCtx, testScenario)
	if err != nil {
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
		zap.L().Error("failed to create test service config",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("test_scenario_id", testSciID),
			zap.Error(err),
		)

		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	_ = tx.Commit()
	zap.L().Info("test scenario created successfully",
		zap.String(logger.FieldRequestID, requestID),
		zap.Uint64("id", testSciID),
		zap.String("name", testScenario.Name),
	)

	return nil
}

func (service *testScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	requestID := logger.GetRequestID(ctx)
	result, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			zap.L().Warn("test scenario not found",
				zap.String(logger.FieldRequestID, requestID),
				zap.Uint64("id", id),
			)

			return nil, pkg.ErrTestScenarioNotFound
		}

		zap.L().Error("failed to get test scenario by ID",
			zap.String(logger.FieldRequestID, requestID),
			zap.Uint64("id", id),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	return result, nil
}

func (service *testScenario) GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error) {
	requestID := logger.GetRequestID(ctx)
	result, err := service.testScenarioRepository.GetPaginated(ctx, pagReq)
	if err != nil {
		zap.L().Error("failed to get paginated test scenarios",
			zap.String(logger.FieldRequestID, requestID),
			zap.Int("page", pagReq.Page),
			zap.Int("per_page", pagReq.PerPage),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetTestScenarios, err)
	}

	zap.L().Debug("retrieved paginated test scenarios",
		zap.String(logger.FieldRequestID, requestID),
		zap.Int("count", len(result)),
		zap.Int("page", pagReq.Page),
	)

	return result, nil
}
