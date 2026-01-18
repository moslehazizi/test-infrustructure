package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"errors"
	"fmt"
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
	motherService, err := service.motherService.GetByID(ctx, testScenario.MotherServiceID)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetMotherService, err)
	}
	testScenario.MotherService = motherService

	testCat, err := service.testCategoryRepository.GetByID(ctx, testScenario.TestCategoryID)
	if err != nil {
		if errors.Is(err, pkg.ErrTestCategoryNotFound) {
			return pkg.ErrTestCategoryNotFound
		}

		return fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestCategory, err)
	}

	err = testScenario.Validate(testCat)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testScenario.Status = entity.ScenarioStatusPending

	if testScenario.TestServiceConfig == nil {
		return pkg.ErrTestServiceConfigIsRequired
	}
	err = testScenario.TestServiceConfig.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToValidateTestSvcCfg, err)
	}

	tx := service.db.Begin()
	dbCtx := context.WithValue(ctx, database.ContextKeyDBTx, tx)
	defer func() {
		if e != nil {
			_ = tx.Rollback()
		}
	}()

	testSciID, err := service.testScenarioRepository.Create(dbCtx, testScenario)
	if err != nil {
		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testScenario.TestServiceConfig.TestScenarioID = testSciID

	err = service.testServiceConfigRepository.Create(dbCtx, testScenario.TestServiceConfig)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	_ = tx.Commit()

	return nil
}

func (service *testScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	result, err := service.testScenarioRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrTestScenarioNotFound) {
			return nil, pkg.ErrTestScenarioNotFound
		}

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	return result, nil
}

func (service *testScenario) GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error) {
	result, err := service.testScenarioRepository.GetPaginated(ctx, pagReq)
	if err != nil {
		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetTestScenarios, err)
	}

	return result, nil
}
