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
)

func NewTestScenarioOperationUsecase(
	db database.Database,
	testScenarioRepository repository.TestScenarioRepository,
	stressTestExecutionManager interfaces.ExecutionManager,
	provisioningService provider.ProvisioningService,
) *testScenarioOperation {
	return &testScenarioOperation{
		db:                         db,
		testScenarioRepository:     testScenarioRepository,
		stressTestExecutionManager: stressTestExecutionManager,
		provisioningService:        provisioningService,
	}
}

type testScenarioOperation struct {
	db                         database.Database
	testScenarioRepository     repository.TestScenarioRepository
	stressTestExecutionManager interfaces.ExecutionManager
	provisioningService        provider.ProvisioningService
}

func (service *testScenarioOperation) Start(ctx context.Context, id uint64) error {
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

func (service *testScenarioOperation) Pause(ctx context.Context, id uint64) error {
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

func (service *testScenarioOperation) Resume(ctx context.Context, id uint64) error {
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

func (service *testScenarioOperation) Stop(ctx context.Context, id uint64) error {
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

func (service *testScenarioOperation) Delete(ctx context.Context, id uint64) error {
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
