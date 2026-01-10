package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"fmt"
)

type TestScenario interface {
	Create(ctx context.Context, testScenario *entity.TestScenario) error
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
