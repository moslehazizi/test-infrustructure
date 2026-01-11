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
	Create(ctx context.Context, testScenario *entity.TestScenario) error
	GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error)
	GetPaginated(ctx context.Context, pagReq entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error)
}

func NewTestScenarioUsecase(testScenarioRepository repository.TestScenarioRepository, testCategoryRepository repository.TestCategory) TestScenario {
	return &testScenario{
		testScenarioRepository: testScenarioRepository,
		testCategoryRepository: testCategoryRepository,
	}
}

type testScenario struct {
	testScenarioRepository repository.TestScenarioRepository
	testCategoryRepository repository.TestCategory
}

func (service *testScenario) Create(ctx context.Context, testScenario *entity.TestScenario) error {
	_, err := service.testCategoryRepository.GetByID(ctx, int(testScenario.TestCategoryID))

	err = testScenario.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", pkg.ErrFailedToCreateTestScenario, err)
	}

	testScenario.Status = entity.ScenarioStatusPending

	err = service.testScenarioRepository.Create(ctx, testScenario)
	if err != nil {
		return fmt.Errorf("%w, %w", pkg.ErrFailedToCreateTestScenario, err)
	}

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
