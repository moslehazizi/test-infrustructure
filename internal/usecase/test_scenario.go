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
}

func NewTestScenarioUsecase(testScenarioRepository repository.TestScenarioRepository) TestScenario {
	return &testScenario{
		testScenarioRepository: testScenarioRepository,
	}
}

type testScenario struct {
	testScenarioRepository repository.TestScenarioRepository
}

func (service *testScenario) Create(ctx context.Context, testScenario *entity.TestScenario) error {
	testScenario.Status = entity.ScenarioStatusPending

	err := service.testScenarioRepository.Create(ctx, testScenario)
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
