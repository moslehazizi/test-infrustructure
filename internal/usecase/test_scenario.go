package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"errors"
	"fmt"
)

type TestScenario interface {
	Create(ctx context.Context, testScenario *entity.TestScenario, testSvcCfg *entity.TestServiceConfig) error
	GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error)
	GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error)
}

func NewTestScenarioUsecase(
	testScenarioRepository repository.TestScenarioRepository,
	testCategoryRepository repository.TestCategory,
	testServiceConfigRepository repository.TestServiceConfigRepository) TestScenario {
	return &testScenario{
		testScenarioRepository:      testScenarioRepository,
		testCategoryRepository:      testCategoryRepository,
		testServiceConfigRepository: testServiceConfigRepository,
	}
}

type testScenario struct {
	testScenarioRepository      repository.TestScenarioRepository
	testCategoryRepository      repository.TestCategory
	testServiceConfigRepository repository.TestServiceConfigRepository
}

func (service *testScenario) Create(
	ctx context.Context,
	testScenario *entity.TestScenario,
	testSvcCfg *entity.TestServiceConfig,
) (e error) {
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

	repo := service.testScenarioRepository.Begin()
	defer func() {
		if e != nil {
			_ = repo.Rollback()
		}
	}()

	testSciID, err := repo.Create(ctx, testScenario)
	if err != nil {
		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testSvcCfg.ID = testSciID

	err = service.testServiceConfigRepository.Create(ctx, testSvcCfg)
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	_ = repo.Commit()

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
